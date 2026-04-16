package web

import (
	"embed"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iainjreid/stencil"
)

//go:embed templates
var tmpls embed.FS

func loadTemplates(r *gin.Engine, funcs template.FuncMap) {
	slog.Debug("Loading templates", "standalone", stencil.IsStandalone)

	stencilFS, err := stencil.New(tmpls, "./web/templates")
	if err != nil {
		panic(err)
	}

	tmpl, err := stencilFS.LoadTemplates(template.New("").Funcs(funcs))
	if err != nil {
		panic(err)
	}

	r.StaticFS("/static", http.FS(stencilFS))
	r.SetHTMLTemplate(tmpl)

	if !stencil.IsStandalone {
		r.Use(func(c *gin.Context) {
			tmpl, err := stencilFS.LoadTemplates(template.New("").Funcs(funcs))
			if err != nil {
				c.AbortWithError(500, err)
			}
			r.SetHTMLTemplate(tmpl)
		})
	}
}
