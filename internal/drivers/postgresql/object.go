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
	"io"
	"log/slog"
	"time"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/storer"
	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/iainjreid/source/storage/driver"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/sync/errgroup"
)

type ObjectStorage struct {
	ID     string
	Pool   *pgxpool.Pool
	Cache  *lru.Cache[string, plumbing.EncodedObject]
	buffer map[string]plumbing.EncodedObject
}

// NewEncodedObject returns a new copy of the default plumbing.MemoryObject.
func (o *ObjectStorage) NewEncodedObject() plumbing.EncodedObject {
	return &plumbing.MemoryObject{}
}

// EncodedObject returns the object with the specified hash and matching object
// type. If AnyObject is passed as the object type, only the object hash will be
// used in this lookup.
func (o *ObjectStorage) EncodedObject(objType plumbing.ObjectType, objHash plumbing.Hash) (plumbing.EncodedObject, error) {
	if o.buffer != nil {
		if obj := o.buffer[objHash.String()]; obj != nil {
			return obj, nil
		}
	}

	if obj, exists := o.Cache.Get(objHash.String()); exists {
		return obj, nil
	}

	rows, err := o.Pool.Query(context.Background(), `
		SELECT
			objects.type,
			objects.contents
		FROM objects
		WHERE objects.repo_id = $1
			AND objects.hash = $2;`, o.ID, objHash.String())

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

	o.Cache.Add(objHash.String(), obj)

	return obj, nil
}

// IterEncodedObjects returns an iterator that traverses all of the available
// objects of the specified type.
func (o *ObjectStorage) IterEncodedObjects(objType plumbing.ObjectType) (storer.EncodedObjectIter, error) {
	rows, err := o.Pool.Query(context.Background(), `
		SELECT
			objects.type,
			objects.contents
		FROM objects
		WHERE objects.repo_id = $1
			AND objects.type = $2;`, o.ID, objType)

	if err != nil {
		return nil, &plumbing.UnexpectedError{
			Err: err,
		}
	}

	return driver.NewScannableIterator(rows.Next, rows.Close, func() (plumbing.EncodedObject, error) {
		return scanMemoryObject(rows)
	}), nil
}

// EncodedObjectSize returns the size of the contents stored against the object
// with the specified hash.
func (o *ObjectStorage) EncodedObjectSize(hash plumbing.Hash) (int64, error) {
	if o.buffer != nil {
		if obj := o.buffer[hash.String()]; obj != nil {
			return obj.Size(), nil
		}
	}

	rows, err := o.Pool.Query(context.Background(), `
		SELECT
			length
		FROM objects
		WHERE objects.repo_id = $1
			AND objects.hash = $2;`, o.ID, hash)

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

	rows, err := o.Pool.Query(context.Background(), `
		SELECT
			objects.hash
		FROM objects
		WHERE objects.repo_id = $1
			AND objects.hash = $2;`, o.ID, hash)

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
	query := `
		INSERT INTO objects (
			repo_id,
			type,
			hash,
			contents,
			length
		) VALUES ($1, $2, $3, $4, $5);
	`
	if _, err := o.Pool.Exec(context.Background(), query, o.ID, obj.Type(), obj.Hash(), cont, obj.Size()); err != nil {
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

func scanMemoryObject(rows pgx.Rows) (plumbing.EncodedObject, error) {
	var t plumbing.ObjectType
	var cont []byte

	// Attempt to scan the row, and in the event that an error occurs, it's likely
	// we have hit a dud query. If this happens, return nil and the wrapped error.
	if err := rows.Scan(&t, &cont); err != nil {
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
	defer func() {
		slog.Debug("commit complete", "time", time.Since(start))
	}()

	return o.CommitCopy()
}

type BufferWriter struct {
	repoId string
	buf    []plumbing.EncodedObject
	pos    int
	err    error
}

func (b *BufferWriter) Next() bool {
	b.pos++
	return b.pos <= len(b.buf)
}

var keyMap = map[string]bool{}

func (b *BufferWriter) Values() ([]any, error) {
	obj := b.buf[b.pos-1]

	cont := make([]byte, obj.Size())
	reader, err := obj.Reader()
	if err != nil {
		b.err = err
	}
	reader.Read(cont)

	return []any{b.repoId, obj.Type(), obj.Hash(), cont, obj.Size()}, nil
}

func (b *BufferWriter) Err() error {
	return b.err
}

func (o *ObjectStorage) CommitCopy() error {
	ctx := context.Background()
	slog.DebugContext(ctx, "commiting object storage", "count", len(o.buffer))

	// Get a Tx for making transaction requests.
	conn, err := o.Pool.Acquire(ctx)
	defer func() {
		conn.Release()
	}()

	if err != nil {
		return err
	}

	wg := new(errgroup.Group)

	n := 3

	writes := make([]plumbing.EncodedObject, 0, len(o.buffer))
	for _, v := range o.buffer {
		writes = append(writes, v)
	}

	chunkSize := len(writes) / n
	slog.DebugContext(ctx, "running parallel commit", "write_count", len(writes), "max_conn", n, "max_chunk_size", chunkSize)

	for i := 0; i < len(writes); i += chunkSize {
		logger := slog.With(
			slog.Group("goroutine_info",
				slog.Int("goid", i),
			),
		)

		// Clamp the last chunk to the slice bound as necessary.
		end := min(chunkSize, len(writes[i:]))

		// Set the capacity of each chunk so that appending to a chunk does
		// not modify the original slice.
		chunk := writes[i : i+end : i+end]

		wg.Go(func() error {
			chunkConn, err := o.Pool.Acquire(ctx)
			logger.DebugContext(ctx, "acquired connection")

			defer chunkConn.Release()
			if err != nil {
				return err
			}

			logger.DebugContext(ctx, "copying  chunk")
			source := &BufferWriter{
				repoId: o.ID,
				buf:    chunk,
			}

			count, err := chunkConn.CopyFrom(ctx, pgx.Identifier{"objects"}, []string{"repo_id", "type", "hash", "contents", "length"}, source)
			if err != nil {
				return fmt.Errorf("chunk failed to be written: %w", err)
			}
			logger.DebugContext(ctx, "successfully copied chunk", "written", count)

			return nil
		})
	}

	err = wg.Wait()
	if err != nil {
		return err
	}
	o.buffer = nil

	return nil
}
