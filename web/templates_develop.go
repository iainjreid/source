//go:build !standalone
package web

import (
    "os"
	"html/template"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iainjreid/go-stencil"
)

func LoadTemplates(r *gin.Engine, funcs template.FuncMap) {
    log.Println("Loading live templates, ignore Gin warnings")

    stencilFS, _ := stencil.New(os.DirFS("web/templates"))

	r.StaticFS("/static", http.FS(stencilFS))
	r.Use(func(*gin.Context) {
        r.SetHTMLTemplate(template.Must(template.New("").Funcs(funcs).ParseFS(stencilFS, "*.tmpl", "*/*.tmpl")))
	})
}
