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

package handlers_test

import (
	"compress/gzip"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/iainjreid/source/internal/web/handlers"
)

func TestSendJSON(t *testing.T) {
	h := &handlers.Handlers{}

	tests := []struct {
		name string
		gzip bool
	}{
		{
			name: "plain",
		},
		{
			name: "gzip",
			gzip: true,
		},
	}

	want := map[string]any{
		"hello": "world",
		"id":    float64(123),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)

			if tt.gzip {
				req.Header.Set("Accept-Encoding", "gzip")
			}

			rec := httptest.NewRecorder()

			h.SendJSON(rec, req, map[string]any{
				"hello": "world",
				"id":    123,
			})

			res := rec.Result()

			if got := res.Header.Get("Content-Type"); got != "application/json" {
				t.Errorf("Content-Type = %q, want application/json", got)
			}

			var r io.Reader = res.Body

			if tt.gzip {
				if got := res.Header.Get("Content-Encoding"); got != "gzip" {
					t.Errorf("Content-Encoding = %q, want gzip", got)
				}

				gz, err := gzip.NewReader(res.Body)
				if err != nil {
					t.Error(err)
				}
				defer gz.Close()

				r = gz
			}

			var got map[string]any
			if err := json.NewDecoder(r).Decode(&got); err != nil {
				t.Error(err)
			}

			if got["hello"] != want["hello"] {
				t.Errorf("hello = %v, want %v", got["hello"], want["hello"])
			}

			if got["id"] != want["id"] {
				t.Errorf("id = %v, want %v", got["id"], want["id"])
			}
		})
	}
}

func TestSendErr(t *testing.T) {
	h := &handlers.Handlers{}

	rec := httptest.NewRecorder()

	h.SendErr(rec, http.StatusNotFound, "Not Found")

	res := rec.Result()

	if res.StatusCode != http.StatusNotFound {
		t.Errorf("StatusCode = %d, want %d", res.StatusCode, http.StatusNotFound)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}

	const want = "404 Not Found\n"

	if string(body) != want {
		t.Errorf("body = %q, want %q", body, want)
	}
}
