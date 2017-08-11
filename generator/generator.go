//go:generate go-bindata -o assets.go -pkg=generator enum.tmpl

package generator

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"sort"
	"strings"
	"text/template"

	"golang.org/x/tools/imports"

	"github.com/Masterminds/sprig"
)

// Generator is responsible for generating validation files for the given in a go source file.
type Generator struct {
	t               *template.Template
	knownTemplates  map[string]*template.Template
	fileSet         *token.FileSet
	noPrefix        bool
	lowercaseLookup bool
	marshal         bool
}

type Enum struct {
	Name   string
	Prefix string
	Type   string
	Values []EnumValue
}

type EnumValue struct {
	Name         string
	PrefixedName string
	Value        int
}

// NewGenerator is a constructor method for creating a new Generator with default
// templates loaded.
func NewGenerator() *Generator {
	g := &Generator{
		knownTemplates: make(map[string]*template.Template),
		t:              template.New("generator"),
		fileSet:        token.NewFileSet(),
		noPrefix:       false,
	}

	funcs := sprig.TxtFuncMap()

	funcs["stringify"] = Stringify
	funcs["indexify"] = Indexify

	g.t.Funcs(funcs)

	for _, assets := range AssetNames() {
		g.t = template.Must(g.t.Parse(string(MustAsset(assets))))
	}

	g.updateTemplates()

	return g
}

// WithNoPrefix is used to change the enum const values generated to not have the enum on them.
func (g *Generator) WithNoPrefix() *Generator {
	g.noPrefix = true
	return g
}

// WithLowercaseVariant is used to change the enum const values generated to not have the enum on them.
func (g *Generator) WithLowercaseVariant() *Generator {
	g.lowercaseLookup = true
	return g
}

// WithMarshal is used to add marshalling to the enum
func (g *Generator) WithMarshal() *Generator {
	g.marshal = true
	return g
}

// GenerateFromFile is responsible for orchestrating the Code generation.  It results in a byte array
// that can be written to any file desired.  It has already had goimports run on the code before being returned.
func (g *Generator) GenerateFromFile(inputFile string) ([]byte, error) {
	f, err := g.parseFile(inputFile)
	if err != nil {
		return nil, fmt.Errorf("generate: error parsing input file '%s': %s", inputFile, err)
	}
	return g.Generate(f)

}

type byPosition []*ast.Field

func (a byPosition) Len() int           { return len(a) }
func (a byPosition) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byPosition) Less(i, j int) bool { return a[i].Pos() < a[j].Pos() }

// Generate does the heavy lifting for the code generation starting from the parsed AST file.
func (g *Generator) Generate(f *ast.File) ([]byte, error) {
	var err error
	enums := g.inspect(f)
	if len(enums) <= 0 {
		return nil, nil
	}

	pkg := f.Name.Name

	vBuff := bytes.NewBuffer([]byte{})
	g.t.ExecuteTemplate(vBuff, "header", map[string]interface{}{"package": pkg})

	// Make the output more consistent by iterating over sorted keys of map
	var keys []string
	for key := range enums {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, name := range keys {
		ts := enums[name]

		// Parse the enum doc statement
		enum, err := g.parseEnum(ts)
		if err != nil {
			continue
		}

		data := map[string]interface{}{
			"enum":      enum,
			"name":      name,
			"lowercase": g.lowercaseLookup,
			"marshal":   g.marshal,
		}

		err = g.t.ExecuteTemplate(vBuff, "enum", data)
	}

	formatted, err := imports.Process(pkg, vBuff.Bytes(), nil)
	if err != nil {
		err = fmt.Errorf("generate: error formatting code %s\n\n%s\n", err, string(vBuff.Bytes()))
	}
	return formatted, err
}

func (g *Generator) getStringForExpr(f ast.Expr) string {
	typeBuff := bytes.NewBuffer([]byte{})
	pErr := printer.Fprint(typeBuff, g.fileSet, f)
	if pErr != nil {
		fmt.Printf("Error getting Type: %s\n", pErr)
	}
	return typeBuff.String()
}

// AddTemplateFiles will be used during generation when the command line accepts
// user templates to add to the generation.
func (g *Generator) AddTemplateFiles(filenames ...string) (err error) {
	g.t, err = g.t.ParseFiles(filenames...)
	if err == nil {
		g.updateTemplates()
	}
	return
}

// updateTemplates will update the lookup map for validation checks that are
// allowed within the template engine.
func (g *Generator) updateTemplates() {
	for _, template := range g.t.Templates() {
		g.knownTemplates[template.Name()] = template
	}
}

// parseFile simply calls the go/parser ParseFile function with an empty token.FileSet
func (g *Generator) parseFile(fileName string) (*ast.File, error) {
	// Parse the file given in arguments
	return parser.ParseFile(g.fileSet, fileName, nil, parser.ParseComments)
}

// parseEnum looks for the ENUM(x,y,z) formatted documentation from the type definition
func (g *Generator) parseEnum(ts *ast.TypeSpec) (*Enum, error) {

	if ts.Doc == nil {
		return nil, errors.New("No Doc on Enum")
	}

	enum := &Enum{}

	enum.Name = ts.Name.Name
	enum.Type = fmt.Sprintf("%s", ts.Type)
	if !g.noPrefix {
		enum.Prefix = ts.Name.Name
	}

	parts := []string{}
	store := false
	for _, comment := range ts.Doc.List {
		if store {
			trimmed := strings.TrimSuffix(strings.TrimSpace(strings.TrimPrefix(comment.Text, `//`)), `,`)
			parts = append(parts, trimmed)
			if strings.Contains(comment.Text, `)`) {
				// End ENUM Declaration
				break
			}
		}
		if strings.Contains(comment.Text, `ENUM(`) {
			// Start ENUM Declaration
			if !strings.Contains(comment.Text, `)`) {
				// Don't store other lines
				store = true
			}
			startIndex := strings.Index(comment.Text, `ENUM(`)
			text := comment.Text
			if startIndex >= 0 {
				text = text[startIndex:]
			}
			trimmed := strings.TrimSuffix(strings.TrimSpace(strings.TrimPrefix(text, `//`)), `,`)
			parts = append(parts, trimmed)

		}
	}
	enumDecl := strings.Join(parts, `,`)
	values := strings.Split(strings.TrimSuffix(strings.TrimPrefix(enumDecl, `ENUM(`), `)`), `,`)
	data := 0
	for _, value := range values {
		if value != "" {
			name := strings.Title(strings.TrimSpace(value))
			enum.Values = append(enum.Values, EnumValue{Name: strings.TrimSpace(name), PrefixedName: enum.Prefix + strings.TrimSpace(name), Value: data})
			data++
		}
	}

	// fmt.Printf("###\nENUM: %+v\n###\n", enum)

	return enum, nil
}

// inspect will walk the ast and fill a map of names and their struct information
// for use in the generation template.
func (g *Generator) inspect(f *ast.File) map[string]*ast.TypeSpec {
	// structs := make(map[string]*ast.StructType)
	enums := make(map[string]*ast.TypeSpec)
	// Inspect the AST and find all structs.
	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.GenDecl:
			// Copy the doc spec to the type or value spec
			// cause they missed this... whoops
			if x.Doc != nil {
				for _, spec := range x.Specs {
					switch s := spec.(type) {
					case *ast.TypeSpec:
						if s.Doc == nil {
							s.Doc = x.Doc
						}
					case *ast.ValueSpec:
						if s.Doc == nil {
							s.Doc = x.Doc
						}
					}
				}
			}
		case *ast.Ident:
			if x.Obj != nil {
				// fmt.Printf("Node: %#v\n", x.Obj)
				// Make sure it's a Type Identifier
				if x.Obj.Kind == ast.Typ {
					// Make sure it's a spec (Type Identifiers can be throughout the code)
					if ts, ok := x.Obj.Decl.(*ast.TypeSpec); ok {
						// fmt.Printf("Type: %+v\n", ts)
						isEnum := false
						if ts.Doc != nil {
							for _, comment := range ts.Doc.List {
								if strings.Contains(comment.Text, `ENUM(`) {
									isEnum = true
								}
								// fmt.Printf("Doc: %s\n", comment.Text)
							}
						}
						// Only store documented enums
						if isEnum {
							// fmt.Printf("EnumType: %T\n", ts.Type)
							enums[x.Name] = ts
						}
					}
				}
			}
		}
		// Return true to continue through the tree
		return true
	})

	return enums
}
