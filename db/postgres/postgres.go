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
