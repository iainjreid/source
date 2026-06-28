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

package utils

import (
	"compress/gzip"
	"io"
	"sync"
)

var pool = sync.Pool{
	New: func() any {
		return gzip.NewWriter(io.Discard)
	},
}

// GetGzipWriter returns a [gzip.Writer] that writes compressed data to the
// provided [io.Writer].
//
// The returned writer is retrieved from an internal pool and must be returned
// to the pool by calling [PutGzipWriter] when it is no longer needed.
func GetGzipWriter(w io.Writer) *gzip.Writer {
	gz := pool.Get().(*gzip.Writer)
	gz.Reset(w)

	return gz
}

// PutGzipWriter accepts a [gzip.Writer] to be returned to the internal pool. It
// is important to note that the writer must not be used after calling this
// method.
//
// Before returning the writer to the internal pool, any pending writes to the
// underlying [io.Writer] are flushed, and the [gzip.Writer] is then closed.
func PutGzipWriter(gz *gzip.Writer) {
	gz.Close()
	pool.Put(gz)
}
