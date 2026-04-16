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
	"github.com/go-git/go-git/v5/storage"
	"github.com/iainjreid/source/git"
	"github.com/iainjreid/source/view"
)

// Does this achieve anything? Requests for specific blobs should
// be cached client side. Set an ETag?
func cacheMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "public, max-age=604800, immutable")
	}
}

func NewServer(storage storage.Storer) error {
	repoUrl := "https://github.com/iainjreid/source.git"

	repo := git.CloneRepo(storage, repoUrl)
	if err := repo.Error(); err != nil {
		panic(err)
	}

	r := gin.Default()

	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.Use(cacheMiddleware())

	r.RedirectTrailingSlash = true
	r.RedirectFixedPath = true

	r.SetTrustedProxies(nil)

	loadTemplates(r, template.FuncMap{
		"add": func(i1, i2 int) int {
			return i1 + i2
		},
		"sub": func(i1, i2 int) int {
			return i1 - i2
		},
		"mul": func(i1, i2 int) int {
			return i1 * i2
		},
		"htmlSafe": func(html string) template.HTML {
			return template.HTML(html)
		},
	})

	r.GET("/blob/:hash", func(c *gin.Context) {
		hash := c.Param("hash")

		// db.GraphLookup(storage, hash)
		renderFile(c, repo, hash, "/")
	})

	r.GET("/blob/:hash/*path", func(c *gin.Context) {
		hash := c.Param("hash")
		path := c.Param("path")

		// db.GraphLookup(storage, hash)
		renderFile(c, repo, hash, path)
	})

	r.GET("/branches", func(c *gin.Context) {
		branches, err := repo.GetBranches()

		if err != nil {
			renderError(c, err)
			return
		}

		c.JSON(http.StatusOK, branches)
	})

	// r.GET("/clone", func(c *gin.Context) {
	// 	db.HardReset()
	// 	repo := git.CloneRepo(storage, repoUrl)
	// 	branches, err := repo.GetBranches()

	// 	if err != nil {
	// 		renderError(c, err)
	// 		return
	// 	}

	// 	c.HTML(http.StatusOK, "index.tmpl", map[string]interface{}{
	// 		"now":      time.Now(),
	// 		"Branches": branches,
	// 	})
	// })

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", view.New(repo))
	})

	r.GET("/feedback", func(c *gin.Context) {
		c.HTML(http.StatusOK, "feedback.tmpl", view.New(repo))
	})

	return r.Run()
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
