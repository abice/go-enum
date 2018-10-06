//go:generate go-bindata -o assets/assets.go -pkg=assets enum.tmpl

package generator

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"sort"
	"strconv"
	"strings"
	"text/template"
	"unicode"

	"golang.org/x/tools/imports"

	"github.com/Masterminds/sprig"
	"github.com/abice/go-enum/generator/assets"
	"github.com/pkg/errors"
)

const (
	skipHolder         = `_`
	parseCommentPrefix = `//`
)

// Generator is responsible for generating validation files for the given in a go source file.
type Generator struct {
	t               *template.Template
	knownTemplates  map[string]*template.Template
	fileSet         *token.FileSet
	noPrefix        bool
	lowercaseLookup bool
	marshal         bool
	sql             bool
	flag            bool
	names           bool
}

// Enum holds data for a discovered enum in the parsed source
type Enum struct {
	Name   string
	Prefix string
	Type   string
	Values []EnumValue
}

// EnumValue holds the individual data for each enum value within the found enum.
type EnumValue struct {
	RawName      string
	Name         string
	PrefixedName string
	Value        int
	Comment      string
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
	funcs["mapify"] = Mapify
	funcs["unmapify"] = Unmapify
	funcs["namify"] = Namify

	g.t.Funcs(funcs)

	for _, asset := range assets.AssetNames() {
		g.t = template.Must(g.t.Parse(string(assets.MustAsset(asset))))
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

// WithSQLDriver is used to add marshalling to the enum
func (g *Generator) WithSQLDriver() *Generator {
	g.sql = true
	return g
}

// WithFlag is used to add flag methods to the enum
func (g *Generator) WithFlag() *Generator {
	g.flag = true
	return g
}

// WithNames is used to add Names methods to the enum
func (g *Generator) WithNames() *Generator {
	g.names = true
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

// Generate does the heavy lifting for the code generation starting from the parsed AST file.
func (g *Generator) Generate(f *ast.File) ([]byte, error) {
	enums := g.inspect(f)
	if len(enums) <= 0 {
		return nil, nil
	}

	pkg := f.Name.Name

	vBuff := bytes.NewBuffer([]byte{})
	err := g.t.ExecuteTemplate(vBuff, "header", map[string]interface{}{"package": pkg})
	if err != nil {
		return nil, errors.WithMessage(err, "Failed writing header")
	}

	// Make the output more consistent by iterating over sorted keys of map
	var keys []string
	for key := range enums {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, name := range keys {
		ts := enums[name]

		// Parse the enum doc statement
		enum, pErr := g.parseEnum(ts)
		if pErr != nil {
			continue
		}

		data := map[string]interface{}{
			"enum":      enum,
			"name":      name,
			"lowercase": g.lowercaseLookup,
			"marshal":   g.marshal,
			"sql":       g.sql,
			"flag":      g.flag,
			"names":     g.names,
		}

		err = g.t.ExecuteTemplate(vBuff, "enum", data)
		if err != nil {
			return vBuff.Bytes(), errors.WithMessage(err, fmt.Sprintf("Failed writing enum data for enum: %q", name))
		}
	}

	formatted, err := imports.Process(pkg, vBuff.Bytes(), nil)
	if err != nil {
		err = fmt.Errorf("generate: error formatting code %s\n\n%s", err, vBuff.String())
	}
	return formatted, err
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

	enumDecl := getEnumDeclFromComments(ts.Doc.List)

	values := strings.Split(strings.TrimSuffix(strings.TrimPrefix(enumDecl, `ENUM(`), `)`), `,`)
	data := 0
	for _, value := range values {
		var comment string

		// Trim and store comments
		if strings.Contains(value, parseCommentPrefix) {
			commentStartIndex := strings.Index(value, parseCommentPrefix)
			comment = value[commentStartIndex+len(parseCommentPrefix):]
			comment = strings.TrimSpace(comment)
			// value without comment
			value = value[:commentStartIndex]
		}

		// Make sure to leave out any empty parts
		if value != "" {
			if strings.Contains(value, `=`) {
				// Get the value specified and set the data to that value.
				equalIndex := strings.Index(value, `=`)
				dataVal := strings.TrimSpace(value[equalIndex+1:])
				if dataVal != "" {
					newData, err := strconv.ParseInt(dataVal, 10, 32)
					if err != nil {
						return nil, errors.Wrapf(err, "failed parsing the data part of enum value '%s'", value)
					}
					data = int(newData)
					value = value[:equalIndex]
				} else {
					value = strings.TrimSuffix(value, `=`)
					fmt.Printf("Ignoring enum with '=' but no value after: %s\n", value)
				}
			}
			rawName := strings.TrimSpace(value)
			name := strings.Title(rawName)
			prefixedName := name
			if name != skipHolder {
				prefixedName = enum.Prefix + name
				prefixedName = sanitizeValue(prefixedName)
			}

			ev := EnumValue{Name: name, RawName: rawName, PrefixedName: prefixedName, Value: data, Comment: comment}
			enum.Values = append(enum.Values, ev)
			data++
		}
	}

	// fmt.Printf("###\nENUM: %+v\n###\n", enum)

	return enum, nil
}

// sanitizeValue will ensure the value name generated adheres to golang's
// identifier syntax as described here: https://golang.org/ref/spec#Identifiers
// identifier = letter { letter | unicode_digit }
// where letter can be unicode_letter or '_'
func sanitizeValue(value string) string {

	// Keep skip value holders
	if value == skipHolder {
		return skipHolder
	}

	name := value

	// If the start character is not a unicode letter (this check includes the case of '_')
	// then we need to add an exported prefix, so tack on a 'X' at the beginning
	if !(unicode.IsLetter(rune(name[0]))) {
		name = `X` + name
	}

	// Loop through all the runes and remove any that aren't valid.
	for i, r := range name {
		if !(unicode.IsLetter(r) || unicode.IsNumber(r) || r == '_') {
			if i < len(name) {
				name = name[:i] + name[i+1:]
			} else {
				// At the end of the string, take off the last char
				name = name[:i-1]
			}
		}
	}

	return name
}

// getEnumDeclFromComments parses the array of comment strings and creates a single Enum Declaration statement
// that is easier to deal with for the remainder of parsing.  It turns multi line declarations and makes a single
// string declaration.
func getEnumDeclFromComments(comments []*ast.Comment) string {
	parts := []string{}
	store := false
	for _, comment := range comments {
		lines := breakCommentIntoLines(comment)

		// Go over all the lines in this comment block
		for _, line := range lines {
			if store {
				trimmed := trimAllTheThings(line)
				if trimmed != "" {
					parts = append(parts, trimmed)
				}
				if strings.Contains(line, `)`) {
					// End ENUM Declaration
					break
				}
			}
			if strings.Contains(line, `ENUM(`) {
				// Start ENUM Declaration
				if !strings.Contains(line, `)`) {
					// Store other lines
					store = true
				}
				startIndex := strings.Index(line, `ENUM(`)
				if startIndex >= 0 {
					line = line[startIndex+len(`ENUM(`):]
				}
				trimmed := trimAllTheThings(line)
				if trimmed != "" {
					parts = append(parts, trimmed)
				}
			}
		}
	}
	joined := fmt.Sprintf("ENUM(%s)", strings.Join(parts, `,`))
	return joined
}

// breakCommentIntoLines takes the comment and since single line comments are already broken into lines
// we break multiline comments into separate lines for processing.
func breakCommentIntoLines(comment *ast.Comment) []string {
	lines := []string{}
	text := comment.Text
	if strings.HasPrefix(text, `/*`) {
		// deal with multi line comment
		multiline := strings.TrimSuffix(strings.TrimPrefix(text, `/*`), `*/`)
		lines = append(lines, strings.Split(multiline, "\n")...)
	} else {
		lines = append(lines, strings.TrimPrefix(text, `//`))
	}
	return lines
}

// trimAllTheThings takes off all the cruft of a line that we don't need.
func trimAllTheThings(thing string) string {
	return strings.TrimSpace(strings.TrimSuffix(strings.TrimSuffix(strings.TrimSpace(thing), `,`), `)`))
}

// inspect will walk the ast and fill a map of names and their struct information
// for use in the generation template.
func (g *Generator) inspect(f ast.Node) map[string]*ast.TypeSpec {
	enums := make(map[string]*ast.TypeSpec)
	// Inspect the AST and find all structs.
	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.GenDecl:
			copyGenDeclCommentsToSpecs(x)
		case *ast.Ident:
			if x.Obj != nil {
				// fmt.Printf("Node: %#v\n", x.Obj)
				// Make sure it's a Type Identifier
				if x.Obj.Kind == ast.Typ {
					// Make sure it's a spec (Type Identifiers can be throughout the code)
					if ts, ok := x.Obj.Decl.(*ast.TypeSpec); ok {
						// fmt.Printf("Type: %+v\n", ts)
						isEnum := isTypeSpecEnum(ts)
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

// copyDocsToSpecs will take the GenDecl level documents and copy them
// to the children Type and Value specs.  I think this is actually working
// around a bug in the AST, but it works for now.
func copyGenDeclCommentsToSpecs(x *ast.GenDecl) {
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

}

// isTypeSpecEnum checks the comments on the type spec to determine if there is an enum
// declaration for the type.
func isTypeSpecEnum(ts *ast.TypeSpec) bool {
	isEnum := false
	if ts.Doc != nil {
		for _, comment := range ts.Doc.List {
			if strings.Contains(comment.Text, `ENUM(`) {
				isEnum = true
			}
		}
	}

	return isEnum
}
