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

package driver_test

import (
	"io"
	"slices"
	"testing"

	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/iainjreid/source/storage/driver"
)

var (
	seq = []int{1, 2, 3, 4, 5}
)

// A simpleScannableIterator is a simple iterator that provides the required
// behaviours to satisfy the [driver.NewScannableIterator] factory for testing.
type simpleScannableIterator struct {
	i      int
	closed bool
}

func (s *simpleScannableIterator) Next() bool {
	return s.i < len(seq)
}

func (s *simpleScannableIterator) Close() {
	s.closed = true
}

func (s *simpleScannableIterator) Scan() (int, error) {
	val := seq[s.i]
	s.i++
	return val, nil
}

func (s *simpleScannableIterator) ToScannableIterator() driver.Iterator[int] {
	return driver.NewScannableIterator(s.Next, s.Close, s.Scan)
}

func TestIteratorNext(t *testing.T) {
	s := simpleScannableIterator{}
	iter := s.ToScannableIterator()

	for _, want := range seq {
		if got, err := iter.Next(); err != nil || got != want {
			t.Errorf("iter.Next() = (%v, %v), want (%d, nil)", got, err, want)
		}
	}

	if _, err := iter.Next(); err != io.EOF {
		t.Errorf("got %v, want io.EOF", err)
	}

	if s.closed {
		t.Error("expected iterator to NOT be closed")
	}
}

func TestIteratorForEach(t *testing.T) {
	s := simpleScannableIterator{}
	iter := s.ToScannableIterator()

	// Collect the values from the iterator here
	got := []int{}

	err := iter.ForEach(func(i int) error {
		got = append(got, i)
		return nil
	})

	// If an error is generated here this is an issue with the test struct, not
	// with the [driver.Iterator] that we are testing.
	if err != nil {
		t.Fatalf("unexpected error, please refer to sourcecode: %v", err)
	}

	if !slices.Equal(got, seq) {
		t.Errorf("got %v, want %v", got, seq)
	}

	if !s.closed {
		t.Error("expected iterator to be closed")
	}
}

func TestIteratorForEachStop(t *testing.T) {
	s := simpleScannableIterator{}
	iter := s.ToScannableIterator()

	// Collect the values from the iterator here
	got := []int{}

	err := iter.ForEach(func(i int) error {
		got = append(got, i)

		if i == 2 {
			return storer.ErrStop
		}
		return nil
	})

	// If an error is generated here this is an issue with the test struct, not
	// with the [driver.Iterator] that we are testing.
	if err != nil {
		t.Fatalf("unexpected error, please refer to sourcecode: %v", err)
	}

	if want := []int{1, 2}; !slices.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}

	if !s.closed {
		t.Error("expected iterator to be closed")
	}
}

func TestIteratorClose(t *testing.T) {
	s := simpleScannableIterator{}
	iter := s.ToScannableIterator()

	iter.Close()

	if !s.closed {
		t.Error("expected close callback to be called")
	}
}
