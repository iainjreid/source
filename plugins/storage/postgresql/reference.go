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
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/go-git/go-git/v5/storage"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ReferenceStorage struct {
	Name string
	Pool *pgxpool.Pool
}

// Reference loads a Git reference from storage.
func (r *ReferenceStorage) Reference(name plumbing.ReferenceName) (*plumbing.Reference, error) {
	slog.Debug("finding reference", "name", name)

	rows, err := r.Pool.Query(context.Background(), `
		SELECT
			refs.type,
			refs.hash,
			refs.name,
			refs.target
		FROM refs
		JOIN repos
			ON refs.repo_id = repos.id
		WHERE repos.name = $1
    		AND refs.name = $2;`, r.Name, name)

	if err != nil {
		return nil, &plumbing.UnexpectedError{
			Err: err,
		}
	}

	defer rows.Close()

	if !rows.Next() {
		return nil, plumbing.ErrReferenceNotFound
	}

	obj, err := scanReference(rows)
	if err != nil {
		return nil, &plumbing.UnexpectedError{
			Err: err,
		}
	}

	return obj, nil
}

// IterReferences returns an iterator capable of walking through all available
// Git references.
func (r *ReferenceStorage) IterReferences() (storer.ReferenceIter, error) {
	rows, err := r.Pool.Query(context.Background(), `
		SELECT
			refs.type,
			refs.hash,
			refs.name,
			refs.target
		FROM refs
		JOIN repos
			ON refs.repo_id = repos.id
		WHERE repos.name = $1;`, r.Name)

	if err != nil {
		return nil, &plumbing.UnexpectedError{
			Err: err,
		}
	}

	return NewIterator(rows.Next, rows.Close, func() (*plumbing.Reference, error) {
		return scanReference(rows)
	})
}

// SetReference writes a Git reference to storage, replacing it if reference
// with the same name alreayd exists.
func (r *ReferenceStorage) SetReference(ref *plumbing.Reference) error {
	if err := r.RemoveReference(ref.Name()); err != nil {
		return &plumbing.UnexpectedError{
			Err: err,
		}
	}

	slog.Debug("setting reference", "name", ref.Name(), "hash", ref.Hash().String())

	query := `
		INSERT INTO refs (
			repo_id,
			type,
			hash,
			name,
			target
		)
		SELECT
			repos.id,
			$2, $3, $4, $5
		FROM repos
		WHERE repos.name = $1;
	`

	if result, err := r.Pool.Exec(context.Background(), query, r.Name, ref.Type(), ref.Hash(), ref.Name(), ref.Target()); err != nil {
		slog.Debug("error whilst setting reference", "name", ref.Name())
		return &plumbing.UnexpectedError{
			Err: err,
		}
	} else {
		slog.Debug("rows inserted", "count", result.RowsAffected())
	}

	return nil
}

// RemoveReference deletes a Git reference from storage by its unique name.
func (r *ReferenceStorage) RemoveReference(name plumbing.ReferenceName) error {
	slog.Debug("deleting reference", "name", name)

	query := `
		DELETE FROM refs
		WHERE EXISTS (
			SELECT 1
			FROM refs r
			JOIN repos
				ON r.repo_id = repos.id
			WHERE repos.name = $1
		)
			AND refs.name = $2;
	`

	if result, err := r.Pool.Exec(context.Background(), query, r.Name, name); err != nil {
		slog.Debug("error whilst deleting reference", "name", name)
		return &plumbing.UnexpectedError{
			Err: err,
		}
	} else {
		slog.Debug("rows deleted", "count", result.RowsAffected())
	}

	return nil
}

func (r *ReferenceStorage) CheckAndSetReference(new, old *plumbing.Reference) error {
	if new == nil {
		return nil
	}

	if old != nil {
		if tmp, _ := r.Reference(new.Name()); tmp != nil && tmp.Hash() != old.Hash() {
			return storage.ErrReferenceHasChanged
		}
	}

	return r.SetReference(new)
}

func (r *ReferenceStorage) CountLooseRefs() (int, error) {
	query, err := r.Pool.Query(context.Background(), `
		SELECT
			COUNT(*)
		FROM refs
		JOIN repos
			ON refs.repo_id = repos.id
		WHERE repos.name = $1;`, r.Name)

	defer query.Close()
	if err != nil {
		return 0, &plumbing.UnexpectedError{
			Err: err,
		}
	}

	var count int
	if err := query.Scan(&count); err != nil {
		return 0, &plumbing.UnexpectedError{
			Err: err,
		}
	}

	return count, nil
}

// PackRefs is not currently implemented.
func (r *ReferenceStorage) PackRefs() error {
	return &plumbing.UnexpectedError{
		Err: fmt.Errorf("not supported"),
	}
}

func scanReference(rows pgx.Rows) (*plumbing.Reference, error) {
	var t plumbing.ReferenceType

	var hash string
	var name plumbing.ReferenceName
	var target plumbing.ReferenceName

	if err := rows.Scan(&t, &hash, &name, &target); err != nil {
		return nil, &plumbing.UnexpectedError{
			Err: err,
		}
	}

	var obj *plumbing.Reference

	switch t {
	case plumbing.HashReference:
		obj = plumbing.NewHashReference(name, plumbing.NewHash(hash))

	case plumbing.SymbolicReference:
		obj = plumbing.NewSymbolicReference(name, target)

	default:
		return nil, &plumbing.UnexpectedError{
			Err: fmt.Errorf("unhandled ref type: %s", t.String()),
		}
	}

	return obj, nil
}
