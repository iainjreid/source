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

package handlers

import (
	"io"
	"net/http"
	"path"

	"github.com/iainjreid/source/public"
)

var (
	files = http.FileServerFS(public.Files)
)

var SpaHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	name := path.Clean(r.URL.Path)
	if name == "/" {
		name = "index.html"
	} else {
		name = name[1:]
	}

	// Serve the requested asset if it exists.
	if f, err := public.Files.Open(name); err == nil {
		f.Close()
		files.ServeHTTP(w, r)
		return
	}

	// Otherwise serve the SPA entry point directly.
	f, err := public.Files.Open("index.html")
	if err != nil {
		http.NotFound(w, r)
		return
	}
	defer f.Close()

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.Copy(w, f)
})
