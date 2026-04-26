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

package postgres

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/iainjreid/source/db"
	"github.com/iainjreid/source/db/sql/shared"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ (db.DB) = &Postgres{}

type Postgres struct {
	Pool *pgxpool.Pool
}

func Connect(ctx context.Context, uri string) (*Postgres, error) {
	pool, err := pgxpool.New(ctx, uri)
	if err != nil {
		return nil, fmt.Errorf("error whilst connecting to '%s': %w", uri, err)
	}

	return &Postgres{pool}, nil
}

func (p *Postgres) HardReset(ctx context.Context) error {
	slog.InfoContext(ctx, "dropping all tables")
	if _, err := p.Pool.Exec(ctx, shared.DropTablesQuery); err != nil {
		return fmt.Errorf("error whilst dropping tables: %w", err)
	}

	return nil
}

func (p *Postgres) EnsureReady(ctx context.Context) error {
	slog.InfoContext(ctx, "creating postgres tables")
	if _, err := p.Pool.Exec(ctx, shared.CreateTablesQuery); err != nil {
		return fmt.Errorf("error whilst ensuring tables exist: %w", err)
	}

	slog.InfoContext(ctx, "creating postgres indexes")
	if _, err := p.Pool.Exec(ctx, shared.CreateIndexesQuery); err != nil {
		return fmt.Errorf("error whilst ensuring indexes exist: %w", err)
	}

	return nil
}
