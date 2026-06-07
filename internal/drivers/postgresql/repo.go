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

	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/google/uuid"
	"github.com/iainjreid/source/storage/driver"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ (driver.Repo) = &Repo{}

// Repo implements the [storage.Repo] interface.
type Repo struct {
	id          uuid.UUID
	name        string
	description string
	pool        *pgxpool.Pool

	ObjectStorage
	ReferenceStorage

	// The following structs are duplicated from the go-git memory storage
	// implementation, however they are omitted from the declaration below.
	memory.IndexStorage
	memory.ShallowStorage
	memory.ModuleStorage
	memory.ConfigStorage
}

func NewRepo(pool *pgxpool.Pool, name string) *Repo {
	id, err := uuid.NewV7()
	if err != nil {
		panic(err)
	}

	return &Repo{
		id:               id,
		name:             name,
		pool:             pool,
		ReferenceStorage: ReferenceStorage{Name: name, Pool: pool},
		ObjectStorage:    ObjectStorage{ID: id.String(), Name: name, Pool: pool},
	}
}

func (r Repo) Name() string {
	return r.name
}

func (r Repo) Description() string {
	return r.description
}

func (r *Repo) Init() error {
	slog.Debug("initialising repository")
	return r.Create(context.Background())
}

func (r *Repo) Create(ctx context.Context) error {
	if r.id == uuid.Nil {
		fmt.Println("makingnew")
		id, err := uuid.NewV7()
		if err != nil {
			return err
		}
		r.id = id
	}

	_, err := r.pool.Exec(ctx, `
		INSERT INTO repos (
			id,
			name
		) VALUES ($1, $2);`, r.id, r.name)

	return err
}

func (r *Repo) Exists() (bool, error) {
	var exists bool

	err := r.pool.QueryRow(context.Background(), `
		SELECT EXISTS (
			SELECT 1
			FROM repos
			WHERE repos.name = $1
		);`, r.name).Scan(&exists)

	return exists, err
}
