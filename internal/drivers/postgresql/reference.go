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
	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/iainjreid/source/storage"
	"github.com/iainjreid/source/storage/driver"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ReferenceStorage struct {
	ID    string
	Pool  *pgxpool.Pool
	Cache *lru.Cache[string, *plumbing.Reference]
}

// Reference loads a Git reference from storage.
func (r *ReferenceStorage) Reference(name plumbing.ReferenceName) (*plumbing.Reference, error) {
	if ref, exists := r.Cache.Get(name.String()); exists {
		return ref, nil
	} else {
		slog.Debug("not found", "name", name)
	}

	slog.Debug("getting reference", "name", name)

	rows, err := r.Pool.Query(context.Background(), storage.GetRefQuery, r.ID, name)
	if err != nil {
		slog.Debug("err", "err", err)
		return nil, &plumbing.UnexpectedError{
			Err: err,
		}
	}

	defer rows.Close()

	if !rows.Next() {
		slog.Debug("Next", "err", err)
		return nil, plumbing.ErrReferenceNotFound
	}

	obj, err := scanReference(rows)
	if err != nil {
		slog.Debug("scanReference", "err", err)
		return nil, &plumbing.UnexpectedError{
			Err: err,
		}
	}

	evicted := r.Cache.Add(name.String(), obj)
	slog.Debug("added to cache", "evicted", evicted)

	return obj, nil
}

// IterReferences returns an iterator capable of walking through all available
// Git references.
func (r *ReferenceStorage) IterReferences() (storer.ReferenceIter, error) {
	slog.Debug("iterating references")

	rows, err := r.Pool.Query(context.Background(), storage.GetRefsQuery, r.ID)

	if err != nil {
		return nil, &plumbing.UnexpectedError{
			Err: err,
		}
	}

	return driver.NewScannableIterator(rows.Next, rows.Close, func() (*plumbing.Reference, error) {
		return scanReference(rows)
	}), nil
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

	if result, err := r.Pool.Exec(context.Background(), storage.InsertRef, r.ID, ref.Type(), ref.Hash(), ref.Name(), ref.Target()); err != nil {
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

	if result, err := r.Pool.Exec(context.Background(), storage.DeleteRef, r.ID, name); err != nil {
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

// CountLooseRefs is not required to be implemented.
func (r *ReferenceStorage) CountLooseRefs() (int, error) {
	return 0, &plumbing.UnexpectedError{
		Err: fmt.Errorf("not supported"),
	}
}

// PackRefs is not required to be implemented.
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
