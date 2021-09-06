// +build go1.16

package generator

import (
	"embed"
	"text/template"
)

//go:embed enum.tmpl
var content embed.FS

func (g *Generator) addEmbeddedTemplates() {
	g.t = template.Must(g.t.ParseFS(content, "*.tmpl"))
}
