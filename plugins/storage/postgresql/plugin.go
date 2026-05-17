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

	"github.com/iainjreid/source/storage"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ (storage.Storage) = &Storage{}

type Storage struct {
	Pool *pgxpool.Pool
}

func (*Storage) Protocol() string {
	return "postgresql"
}

func (s *Storage) Repo(name string) storage.Repo {
	return NewRepo(s.Pool, name)
}

func (s *Storage) EnsureReady(ctx context.Context) error {
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

func (s *Storage) ListRepos(ctx context.Context) ([]storage.Repo, error) {
	rows, err := s.Pool.Query(ctx, "SELECT name, description FROM repos;")
	defer rows.Close()

	if err != nil {
		return nil, fmt.Errorf("error whilst listing repositories: %w", err)
	}

	var repos []storage.Repo
	var name, description string

	for rows.Next() {
		err := rows.Scan(&name, &description)
		if err != nil {
			return nil, err
		}
		repos = append(repos, NewRepo(s.Pool, name))
	}

	return repos, nil
}

func Init(ctx context.Context, uri string) (storage.Storage, error) {
	pool, err := pgxpool.New(ctx, uri)
	if err != nil {
		return nil, fmt.Errorf("error whilst connecting to '%s': %w", uri, err)
	}

	return &Storage{
		Pool: pool,
	}, nil
}
