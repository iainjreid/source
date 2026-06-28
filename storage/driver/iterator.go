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

package driver

import (
	"io"

	"github.com/go-git/go-git/v5/plumbing/storer"
)

// An Iterator provides sequential access to a collection of values. Most
type Iterator[T any] interface {
	// Next returns the next value in the sequence. It returns [io.EOF] when no
	// more values remain.
	Next() (T, error)

	// ForEach calls the provided function for each value in the sequence until
	// there are no more values to iterate over, or an error occurs. If an error
	// is raised by the provided function, it is also returned by this method.
	//
	// Alternativly, the provided function may return [storer.ErrStop] to halt
	// iteration immediately. In this instance, ForEach returns nil.
	//
	// Any underlying resources, such as child iterators, will be closed before
	// ForEach returns.
	ForEach(func(T) error) error

	// Close releases any underlying resources, such as child iterators,
	// associated with the Iterator.
	Close()
}

type iterator[T any] struct {
	next  func() bool
	close func()
	scan  func() (T, error)
}

// NewScannableIterator constructs an [Iterator] over any underlying data type that
// can be abstracted into a sequentially scanned table of values.
//
// Initially this was developed to read the results of an SQL query row by row
// in an idiomatic manner, however, it could just as well be used to iterate
// over values read in from a file, or any other durable object that can benefit
// from staggered reads, rather than reading all of the data at once.
//
//   - `next` is expected to report whether another value is available for
//     reading.
//   - `close` is expected to release any underlying resources associated with
//     the iterator.
//   - `scan` is expected to read the current value from the underlying sequence
//     and has the oppourtunity to raise an error if required.
//
// It is up to the implementation to decide whether the positional portion of
// the underlying resource is updated when Next or Scan is called.
func NewScannableIterator[T any](next func() bool, close func(), scan func() (T, error)) Iterator[T] {
	return &iterator[T]{
		next:  next,
		close: close,
		scan:  scan,
	}
}

// Next returns the next value from the [Iterator]. If no more values are
// available, it returns [io.EOF].
func (i *iterator[T]) Next() (T, error) {
	if i.next() {
		return i.scan()
	}

	var empty T
	return empty, io.EOF
}

// ForEach calls the provided function for each value produced by the [Iterator]
// until iteration completes or an error occurs.
//
// If the provided function returns [storer.ErrStop], iteration immidiately
// halts and ForEach returns nil.
//
// The [Iterator] will close itself before ForEach returns.
func (i *iterator[T]) ForEach(fn func(T) error) error {
	defer i.Close()

	for {
		obj, err := i.Next()

		if err != nil {
			if err == io.EOF {
				return nil
			}

			return err
		}

		if err := fn(obj); err != nil {
			if err == storer.ErrStop {
				return nil
			}

			return err
		}
	}
}

// Close will try to close the underlying resources that back the [Iterator].
//
// In the future this method will allow errors to be returned by the underlying
// close method, but this is blocked while go-git is still in use.
func (i *iterator[T]) Close() {
	i.close()
}
