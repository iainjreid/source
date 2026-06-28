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

package utils_test

import (
	"bytes"
	"compress/gzip"
	"io"
	"testing"

	"github.com/iainjreid/source/internal/utils"
)

func TestGzipWriter(t *testing.T) {
	want := "some data to compress"

	if got := decompress(t, compress(t, want)); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestGzipWriterReuse(t *testing.T) {
	const (
		first  = "some data to compress (first pass)"
		second = "some data to compress (second pass)"
	)

	// Compress the first string, loading a writer in the process and then
	// returning it to the pool.
	_ = compress(t, first)

	// Compress the second string, ensuring that the first string has been
	// correctly removed from the underlying writer.
	b := compress(t, second)

	if got := decompress(t, b); got != second {
		t.Errorf("got %q, want %q", got, second)
	}
}

// compress writes to the provided [gzip.Writer]. If an error occurs while doing
// this, it's highly likely to be an initalisation issue somewhere either in
// this file or within the pool itself.
func compress(t *testing.T, str string) []byte {
	t.Helper()

	var buf bytes.Buffer
	gz := utils.GetGzipWriter(&buf)

	if _, err := gz.Write([]byte(str)); err != nil {
		t.Fatalf("unexpected error, failed to write to gzip writer: %v", err)
	}

	utils.PutGzipWriter(gz)
	return buf.Bytes()
}

func decompress(t *testing.T, data []byte) string {
	t.Helper()

	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("unexpected error, failed to create gzip reader: %v", err)
	}
	defer r.Close()

	out, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("unexpected error, failed read from the gzip reader: %v", err)
	}

	return string(out)
}
