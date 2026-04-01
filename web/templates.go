//go:build standalone
package web

import (
    "embed"
	"html/template"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iainjreid/go-stencil"
)

//go:embed templates
var tmpls embed.FS

func LoadTemplates(r *gin.Engine, funcs template.FuncMap) {
    log.Println("Loading embeded templates")

    stencilFS, _ := stencil.New(tmpls)

	r.StaticFS("/static", http.FS(stencilFS))
    r.SetHTMLTemplate(template.Must(template.New("").Funcs(funcs).ParseFS(stencilFS, "*.tmpl", "*/*.tmpl")))
}
