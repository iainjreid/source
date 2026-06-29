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

package web

import (
	"fmt"
	"net/http"

	"github.com/iainjreid/source/internal/web/handlers"
	"github.com/iainjreid/source/internal/web/middleware"
	"github.com/iainjreid/source/public"
	"github.com/iainjreid/source/storage/driver"
)

func NewServer(store driver.Store, port int) error {
	mux := http.NewServeMux()

	h := handlers.Handlers{
		Store: store,
	}

	mux.Handle("/", http.FileServer(http.FS(public.Files)))

	mux.HandleFunc("GET /{repo}/refs", h.GetRefs)
	mux.HandleFunc("GET /{repo}/tree/{hash}", h.GetTree)
	mux.HandleFunc("GET /{repo}/blob/{hash}/{filepath...}", h.GetBlob)
	mux.HandleFunc("GET /{repo}/raw/{hash}/{filepath...}", h.GetRawBlob)

	handler := middleware.Logging(
		middleware.Recover(
			middleware.ResponseTimeReporter(mux),
		),
	)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: handler,
	}

	return srv.ListenAndServe()
}
