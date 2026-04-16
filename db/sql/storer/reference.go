package storer

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/go-git/go-git/v5/storage"
	"github.com/iainjreid/source/db/sql/shared"
)

type ReferenceStorage struct {
	client *sql.DB
}

// Reference loads a Git reference from storage.
func (r *ReferenceStorage) Reference(name plumbing.ReferenceName) (*plumbing.Reference, error) {
	rows, err := r.client.Query(`SELECT type, hash, name, target FROM "source_refs" WHERE name = $1;`, strings.TrimSpace(name.String()))

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
	rows, err := r.client.Query(`SELECT type, hash, name, target FROM "source_refs";`)

	if err != nil {
		return nil, &plumbing.UnexpectedError{
			Err: err,
		}
	}

	return shared.NewIterator(rows.Next, func() { rows.Close() }, func() (*plumbing.Reference, error) {
		return scanReference(rows)
	})
}

// SetReference writes a Git reference to storage, replacing it if reference
// with the same name alreayd exists.
func (r *ReferenceStorage) SetReference(ref *plumbing.Reference) error {
	r.RemoveReference(ref.Name())

	if _, err := r.client.Exec(`INSERT INTO source_refs(type, hash, name, target) VALUES($1, $2, $3, $4);`, ref.Type(), ref.Hash(), ref.Name(), ref.Target()); err != nil {
		return &plumbing.UnexpectedError{
			Err: err,
		}
	}

	return nil
}

// RemoveReference deletes a Git reference from storage by its unique name.
func (r *ReferenceStorage) RemoveReference(name plumbing.ReferenceName) error {
	if _, err := r.client.Exec(`DELETE FROM "source_refs" WHERE name = $1;`, strings.TrimSpace(name.String())); err != nil {
		return &plumbing.UnexpectedError{
			Err: err,
		}
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
	query, err := r.client.Query(`SELECT COUNT(*) FROM "source_refs";`)

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

func scanReference(row *sql.Rows) (*plumbing.Reference, error) {
	var t plumbing.ReferenceType

	var hash string
	var name plumbing.ReferenceName
	var target plumbing.ReferenceName

	if err := row.Scan(&t, &hash, &name, &target); err != nil {
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
