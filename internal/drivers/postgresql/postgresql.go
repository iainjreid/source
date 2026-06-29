// Copyright 2026 Iain J. Reid
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package postgresql

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/go-git/go-git/v5/plumbing"
	gogit "github.com/go-git/go-git/v5/storage"
	"github.com/go-git/go-git/v5/storage/memory"
	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/iainjreid/source/internal/cache"
	"github.com/iainjreid/source/internal/utils"
	"github.com/iainjreid/source/storage"
	"github.com/iainjreid/source/storage/driver"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Driver struct{}

func (Driver) Open(ctx context.Context, uri string) (driver.Store, error) {
	pool, err := pgxpool.New(ctx, uri)
	if err != nil {
		return nil, fmt.Errorf("error whilst connecting to '%s': %w", uri, err)
	}

	return &Store{
		Pool:  pool,
		Cache: utils.Must(lru.New[string, driver.Repo](1024)),
	}, nil
}

func init() {
	storage.Register("postgres", Driver{})
	storage.Register("postgresql", Driver{})
}

type Store struct {
	Pool  *pgxpool.Pool
	Cache *lru.Cache[string, driver.Repo]
}

func (*Store) Protocol() string {
	return "postgresql"
}

func (s *Store) RepoExists(ctx context.Context, name string) (bool, error) {
	var exists bool

	err := s.Pool.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM repos
			WHERE repos.name = $1
		);`, name).Scan(&exists)

	return exists, err
}

func (s *Store) CreateRepo(ctx context.Context, repo driver.Repo) error {
	_, err := s.Pool.Exec(ctx, `
		INSERT INTO repos (
			id,
			name
		) VALUES ($1, $2);`, repo.ID, repo.Name)

	return err
}

func (s *Store) GetRepo(ctx context.Context, name string) (driver.Repo, error) {
	if repo, exists := s.Cache.Get(name); exists {
		return repo, nil
	}

	var (
		id          string
		description string
	)

	slog.DebugContext(ctx, "loading repo from durable storage", "name", name)

	err := s.Pool.QueryRow(ctx, `
		SELECT
			repos.id,
			repos.description
		FROM repos
		WHERE repos.name = $1;`, name).Scan(&id, &description)

	if err != nil {
		return driver.Repo{}, err
	}

	repo := driver.Repo{
		ID:          id,
		Name:        name,
		Description: description,
	}

	s.Cache.Add(name, repo)
	return repo, nil
}

func (s *Store) EnsureReady(ctx context.Context) error {
	slog.InfoContext(ctx, "creating 'repos' table")
	if _, err := s.Pool.Exec(ctx, storage.ReposSchema); err != nil {
		return fmt.Errorf("error whilst ensuring 'repos' table exist: %w", err)
	}

	slog.InfoContext(ctx, "creating 'objects' table")
	if _, err := s.Pool.Exec(ctx, storage.ObjectsSchema); err != nil {
		return fmt.Errorf("error whilst ensuring 'objects' table exist: %w", err)
	}

	slog.InfoContext(ctx, "creating 'refs' table")
	if _, err := s.Pool.Exec(ctx, storage.RefsSchema); err != nil {
		return fmt.Errorf("error whilst ensuring 'refs' table exist: %w", err)
	}

	return nil
}

func (s *Store) ListRepos(ctx context.Context) ([]driver.Repo, error) {
	rows, err := s.Pool.Query(ctx, "SELECT id, name, description FROM repos;")
	defer rows.Close()

	if err != nil {
		return nil, fmt.Errorf("error whilst listing repositories: %w", err)
	}

	var repos []driver.Repo
	var id, name, description string

	for rows.Next() {
		err := rows.Scan(&id, &name, &description)
		if err != nil {
			return nil, err
		}

		repos = append(repos, driver.Repo{
			ID:          id,
			Name:        name,
			Description: description,
		})
	}

	return repos, nil
}

func (r *Store) IterateRefs(ctx context.Context, repoId string) (driver.Iterator[*driver.Ref], error) {
	slog.Debug("iterating references")

	rows, err := r.Pool.Query(context.Background(), storage.GetRefsQuery, repoId)

	if err != nil {
		return nil, &plumbing.UnexpectedError{
			Err: err,
		}
	}

	return driver.NewScannableIterator(rows.Next, rows.Close, func() (*driver.Ref, error) {
		return scanReference(rows)
	}), nil
}

// Storage is a struct implementation of the [storage.Storer] interface.
type Storage struct {
	ObjectStorage
	ReferenceStorage

	// The following structs are duplicated from the go-git memory storage
	// implementation, however they are omitted from the declaration below.
	memory.IndexStorage
	memory.ShallowStorage
	memory.ModuleStorage
	memory.ConfigStorage
}

func (s *Store) ToStorer(ctx context.Context, name string) (gogit.Storer, error) {
	repo, err := s.GetRepo(ctx, name)
	if err != nil {
		return nil, err
	}

	ns := cache.Register(name)

	return &Storage{
		ReferenceStorage: ReferenceStorage{ID: repo.ID, Pool: s.Pool, Cache: ns.Refs},
		ObjectStorage:    ObjectStorage{ID: repo.ID, Pool: s.Pool, Cache: ns.Objs},
	}, nil
}
