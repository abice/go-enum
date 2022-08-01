//go:build go1.16
// +build go1.16

package generator

import (
	"embed"
	"text/template"
)

//go:embed enum.tmpl enum_string.tmpl
var content embed.FS

func (g *Generator) addEmbeddedTemplates() {
	g.t = template.Must(g.t.ParseFS(content, "*.tmpl"))
}
