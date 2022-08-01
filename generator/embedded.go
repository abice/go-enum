//go:generate ../bin/go-bindata -nometadata -o assets/assets.go -pkg=assets enum.tmpl enum_string.tmpl
//go:build !go1.16
// +build !go1.16

package generator

import (
	"text/template"

	"github.com/abice/go-enum/generator/assets"
)

func (g *Generator) addEmbeddedTemplates() {
	for _, asset := range assets.AssetNames() {
		g.t = template.Must(g.t.Parse(string(assets.MustAsset(asset))))
	}
}
