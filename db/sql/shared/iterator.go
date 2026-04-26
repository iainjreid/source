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

package shared

import (
	"io"

	"github.com/go-git/go-git/v5/plumbing/storer"
)

type Iterator[T interface{}] struct {
	next  func() bool
	close func()
	scan  func() (T, error)
}

func NewIterator[T interface{}](next func() bool, close func(), scan func() (T, error)) (*Iterator[T], error) {
	return &Iterator[T]{
		next:  next,
		close: close,
		scan:  scan,
	}, nil
}

func (i *Iterator[T]) Next() (T, error) {
	if i.next() {
		return i.scan()
	} else {
		var empty T
		return empty, io.EOF
	}
}

func (i *Iterator[T]) ForEach(fn func(T) error) error {
	defer i.close()

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

func (i *Iterator[T]) Close() {
	i.close()
}
