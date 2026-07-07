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
