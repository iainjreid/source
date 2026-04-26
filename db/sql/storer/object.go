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

package storer

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"maps"
	"slices"
	"strings"
	"time"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/iainjreid/source/db/sql/shared"
)

type ObjectStorage struct {
	client *sql.DB
	cache  map[string]plumbing.EncodedObject
	buffer map[string]plumbing.EncodedObject
}

// NewEncodedObject returns a new copy of the default plumbing.MemoryObject.
func (o *ObjectStorage) NewEncodedObject() plumbing.EncodedObject {
	return &plumbing.MemoryObject{}
}

var n int

// EncodedObject returns the object with the specified hash and matching object
// type. If AnyObject is passed as the object type, only the object hash will be
// used in this lookup.
func (o *ObjectStorage) EncodedObject(objType plumbing.ObjectType, objHash plumbing.Hash) (plumbing.EncodedObject, error) {
	if o.buffer != nil {
		if obj := o.buffer[objHash.String()]; obj != nil {
			return obj, nil
		}
	}

	n++
	log.Println("READ", n, objHash.String())

	rows, err := o.client.Query(`SELECT type, cont FROM "source_objects" WHERE hash = $1;`, objHash.String())

	if err != nil {
		return nil, &plumbing.UnexpectedError{
			Err: err,
		}
	}

	defer rows.Close()

	if !rows.Next() {
		return nil, plumbing.ErrObjectNotFound
	}

	obj, err := scanMemoryObject(rows)
	if err != nil {
		return nil, &plumbing.UnexpectedError{
			Err: err,
		}
	}

	// If the object found in the database does not have the correct object type,
	// return an error. If the expected object type is AnyObject, skip this check.
	if objType != plumbing.AnyObject && obj.Type() != objType {
		return nil, plumbing.ErrObjectNotFound
	}

	return obj, nil
}

// IterEncodedObjects returns an iterator that traverses all of the available
// objects of the specified type.
func (o *ObjectStorage) IterEncodedObjects(objType plumbing.ObjectType) (storer.EncodedObjectIter, error) {
	rows, err := o.client.Query(`SELECT type, cont FROM "source_objects" WHERE type = $1;`, objType)

	if err != nil {
		return nil, &plumbing.UnexpectedError{
			Err: err,
		}
	}

	return shared.NewIterator(rows.Next, func() { rows.Close() }, func() (plumbing.EncodedObject, error) {
		return scanMemoryObject(rows)
	})
}

// EncodedObjectSize returns the size of the contents stored against the object
// with the specified hash.
func (o *ObjectStorage) EncodedObjectSize(hash plumbing.Hash) (int64, error) {
	if o.buffer != nil {
		if obj := o.buffer[hash.String()]; obj != nil {
			return obj.Size(), nil
		}
	}

	rows, err := o.client.Query(`SELECT length FROM "source_objects" WHERE hash = $1;`, hash)

	if err != nil {
		return 0, &plumbing.UnexpectedError{
			Err: err,
		}
	}

	defer rows.Close()

	var len int64
	err = rows.Scan(&len)
	if err != nil {
		return 0, plumbing.ErrObjectNotFound
	}
	return len, nil
}

// HasEncodedObject checks if an object with the specified hash exists,
// returning nil if the object does exist, and an error if it does not.
func (o *ObjectStorage) HasEncodedObject(hash plumbing.Hash) error {
	if o.buffer != nil {
		if obj := o.buffer[hash.String()]; obj != nil {
			return nil
		}
	}

	rows, err := o.client.Query(`SELECT hash FROM "source_objects" WHERE hash = '$1';`, hash)

	if err != nil {
		return &plumbing.UnexpectedError{
			Err: err,
		}
	}

	defer rows.Close()

	if _, err = scanMemoryObject(rows); err != nil {
		return plumbing.ErrObjectNotFound
	}

	return nil
}

// SetEncodedObject stores the object provided.
//
// If the object contents cannot be read, or the object fails to be written to
// the database, a ZeroHash will be returned with the appropriate wrapped error.
func (o *ObjectStorage) SetEncodedObject(obj plumbing.EncodedObject) (plumbing.Hash, error) {
	cont := make([]byte, obj.Size())
	reader, _ := obj.Reader()

	// Read the provided objects contents in to a local byte array.
	//
	// In the event that this operation fails, return a ZeroHash along with the
	// error returned by the Reader.
	if _, err := reader.Read(cont); err != nil {
		return plumbing.ZeroHash, &plumbing.UnexpectedError{
			Err: err,
		}
	}

	// Write the object to the database.
	//
	// Similar to the above, if an error occurs we will return a ZeroHash along
	// with the error returned by the database driver.
	if _, err := o.client.Exec(`INSERT INTO source_objects(type, hash, cont, length) VALUES($1, $2, $3, $4);`, obj.Type(), obj.Hash(), cont, obj.Size()); err != nil {
		return plumbing.ZeroHash, &plumbing.UnexpectedError{
			Err: err,
		}
	}

	return obj.Hash(), nil
}

// AddAlternate is not currently implemented.
func (o *ObjectStorage) AddAlternate(remote string) error {
	return &plumbing.UnexpectedError{
		Err: fmt.Errorf("not supported"),
	}
}

func scanMemoryObject(row *sql.Rows) (plumbing.EncodedObject, error) {
	var t plumbing.ObjectType
	var cont []byte

	// Attempt to scan the row, and in the event that an error occurs, it's likely
	// we have hit a dud query. If this happens, return nil and the wrapped error.
	if err := row.Scan(&t, &cont); err != nil {
		return nil, &plumbing.UnexpectedError{
			Err: err,
		}
	}

	obj := &plumbing.MemoryObject{}

	obj.SetType(t)
	obj.Write(cont)

	return obj, nil
}

func (o *ObjectStorage) BufferEncodedObject(obj plumbing.EncodedObject) (plumbing.Hash, error) {
	cont := make([]byte, obj.Size())
	reader, err := obj.Reader()

	if err != nil {
		return plumbing.ZeroHash, err
	}

	// Read the provided objects contents into a local byte array.
	//
	// In the event that this operation fails, return a ZeroHash along with the
	// error returned by the Reader (we can safely ignore EOF errors).
	if _, err := reader.Read(cont); err != nil && err != io.EOF {
		return plumbing.ZeroHash, &plumbing.UnexpectedError{
			Err: err,
		}
	}

	if o.buffer == nil {
		o.buffer = map[string]plumbing.EncodedObject{}
	} else if _, has := o.buffer[obj.Hash().String()]; has == true {
		panic("buffer hash collision")
	}
	o.buffer[obj.Hash().String()] = obj
	return obj.Hash(), nil
}

func (o *ObjectStorage) Commit() error {
	if o.buffer == nil {
		return nil
	}

	start := time.Now()
	defer func() { fmt.Println("INSERT", time.Since(start)) }()

	ctx := context.Background()

	// Get a Tx for making transaction requests.
	tx, err := o.client.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Defer a rollback in case anything fails.
	defer tx.Rollback()

	query := `INSERT INTO source_objects(type, hash, parent_hash, cont, length) VALUES %s;`

	for chunk := range slices.Chunk(slices.Collect(maps.Values(o.buffer)), 100) {
		log.Println("chunk: ", len(chunk))

		values := []interface{}{}
		placeholders := []string{}

		for i, obj := range chunk {
			cont := make([]byte, obj.Size())
			reader, _ := obj.Reader()
			reader.Read(cont)

			placeholders = append(placeholders, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)", i*5+1, i*5+2, i*5+3, i*5+4, i*5+5))
			values = append(values, obj.Type(), obj.Hash(), "missing", cont, obj.Size())
		}

		finalQuery := fmt.Sprintf(query, strings.Join(placeholders, ","))

		if _, err := tx.ExecContext(ctx, finalQuery, values...); err != nil {
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	o.buffer = nil
	return nil
}
