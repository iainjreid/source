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

// Package handlers hold the behaviour required to support the HTTP API.
package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/iainjreid/source/internal/utils"
	"github.com/iainjreid/source/storage/driver"
)

type Handlers struct {
	Store driver.Store
}

// clientSupportsGzip checks to see if the client will accept gzip encoded
// responses.
func clientSupportsGzip(r *http.Request) bool {
	return strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")
}

// SendJSON encodes the provided data and writes it do the response stream. If
// the client accepted gzipped encoded responses, it will compress the data.
func (h *Handlers) SendJSON(w http.ResponseWriter, r *http.Request, data interface{}) {
	w.Header().Set("Content-Type", "application/json")

	var writer io.Writer = w

	if clientSupportsGzip(r) {
		w.Header().Set("Content-Encoding", "gzip")

		gz := utils.GetGzipWriter(writer)
		defer utils.PutGzipWriter(gz)

		writer = gz
	}

	if err := json.NewEncoder(writer).Encode(data); err != nil {
		slog.Error("unable to marshal json response", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handlers) SendPlain(w http.ResponseWriter, r *http.Request, body io.Reader) {
	w.Header().Set("Content-Type", "text/plain")

	var writer io.Writer = w

	if clientSupportsGzip(r) {
		w.Header().Set("Content-Encoding", "gzip")

		gz := utils.GetGzipWriter(writer)
		defer utils.PutGzipWriter(gz)

		writer = gz
	}

	if _, err := io.Copy(writer, body); err != nil {
		slog.Error("unable to write plain response", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// SendErr sends an HTTP error response with the provided status code and
// message in the format "<status-code> <message>\n".
func (h *Handlers) SendErr(w http.ResponseWriter, code int, msg string) {
	http.Error(w, fmt.Sprintf("%d %s", code, msg), code)
}
