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
