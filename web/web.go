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
	"html/template"
	"log/slog"
	"net/http"
	"path"
	"time"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/iainjreid/source/git"
	"github.com/iainjreid/source/storage"
	"github.com/iainjreid/source/view"
)

type Index struct {
	PageName string
	Repos    []storage.Repo
}

type Error struct {
	PageName string
	Error    string
}

// Does this achieve anything? Requests for specific blobs should
// be cached client side. Set an ETag?
func cacheMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "public, max-age=604800, immutable")
	}
}

func NewServer(storage storage.Storage, port int) error {
	r := gin.Default()

	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.Use(cacheMiddleware())

	r.RedirectTrailingSlash = true
	r.RedirectFixedPath = true

	r.SetTrustedProxies(nil)

	loadTemplates(r, template.FuncMap{
		"treeDepthToOffset": func(depth int) int {
			return (depth-1)*2 + 1
		},
		"htmlSafe": func(html string) template.HTML {
			return template.HTML(html)
		},
	})

	r.GET("/:repo/blob/:hash", func(c *gin.Context) {
		repo := git.OpenRepo(storage.Repo(c.Param("repo")))
		hash := c.Param("hash")

		renderFile(c, repo, hash, "/")
	})

	r.GET("/:repo/blob/:hash/*path", func(c *gin.Context) {
		repo := git.OpenRepo(storage.Repo(c.Param("repo")))

		hash := c.Param("hash")
		path := c.Param("path")

		renderFile(c, repo, hash, path)
	})

	r.GET("/:repo/branches", func(c *gin.Context) {
		repo := git.OpenRepo(storage.Repo(c.Param("repo")))
		branches, err := repo.GetBranches()

		if err != nil {
			renderError(c, err)
			return
		}

		c.JSON(http.StatusOK, branches)
	})

	r.GET("/", func(c *gin.Context) {
		repos, err := storage.ListRepos(c)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "error.tmpl", Error{
				PageName: "Error",
				Error:    err.Error(),
			})
			return
		}

		c.HTML(http.StatusOK, "index.tmpl", Index{
			PageName: "Home",
			Repos:    repos,
		})
	})

	r.GET("/favicon.ico", func(c *gin.Context) {
		c.JSON(404, gin.H{"message": "Not found"})
	})

	r.GET("/:repo", func(c *gin.Context) {
		repo := git.OpenRepo(storage.Repo(c.Param("repo")))
		c.HTML(http.StatusOK, "repo.tmpl", view.New(repo))
	})

	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"message": "Not found"})
	})

	return r.Run(fmt.Sprintf(":%d", port))
}

func renderFile(c *gin.Context, repo *git.Repo, revision string, filepath string) {
	slog.DebugContext(c, "rendering file", "filepath", filepath)
	dir, file := path.Split(filepath)

	view := view.New(repo)
	view = view.LoadCommit(revision)
	view = view.LoadDir(dir)

	if _, err := view.LoadBlob(file, false); err != nil {
		slog.WarnContext(c, "file not found", "filepath", file)

		if _, err := view.LoadBlob("/README.md", false); err != nil {
			slog.WarnContext(c, "file not found", "filepath", dir+"/README.md")
		}
	}

	c.HTML(http.StatusOK, "file.tmpl", view)
}

func renderError(c *gin.Context, err error) {
	c.HTML(http.StatusOK, "error.tmpl", map[string]interface{}{
		"now":   time.Now(),
		"Error": err.Error(),
	})
}

func formatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d/%02d/%02d", year, month, day)
}
