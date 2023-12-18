package generator

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"text/template"
	"unicode"

	"github.com/Masterminds/sprig/v3"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"golang.org/x/tools/imports"
)

const (
	skipHolder         = `_`
	parseCommentPrefix = `//`
)

// Generator is responsible for generating validation files for the given in a go source file.
type Generator struct {
	Version           string
	Revision          string
	BuildDate         string
	BuiltBy           string
	t                 *template.Template
	knownTemplates    map[string]*template.Template
	userTemplateNames []string
	fileSet           *token.FileSet
	noPrefix          bool
	lowercaseLookup   bool
	caseInsensitive   bool
	marshal           bool
	sql               bool
	sqlint            bool
	flag              bool
	names             bool
	values            bool
	leaveSnakeCase    bool
	prefix            string
	sqlNullInt        bool
	sqlNullStr        bool
	ptr               bool
	mustParse         bool
	forceLower        bool
	forceUpper        bool
	noComments        bool
	buildTags         []string
	replacementNames  map[string]string
}

// Enum holds data for a discovered enum in the parsed source
type Enum struct {
	Name    string
	Prefix  string
	Type    string
	Values  []EnumValue
	Comment string
}

// EnumValue holds the individual data for each enum value within the found enum.
type EnumValue struct {
	RawName      string
	Name         string
	PrefixedName string
	ValueStr     string
	ValueInt     interface{}
	Comment      string
}

// NewGenerator is a constructor method for creating a new Generator with default
// templates loaded.
func NewGenerator() *Generator {
	g := &Generator{
		Version:           "-",
		Revision:          "-",
		BuildDate:         "-",
		BuiltBy:           "-",
		knownTemplates:    make(map[string]*template.Template),
		userTemplateNames: make([]string, 0),
		t:                 template.New("generator"),
		fileSet:           token.NewFileSet(),
		noPrefix:          false,
		replacementNames:  map[string]string{},
	}

	funcs := sprig.TxtFuncMap()

	funcs["stringify"] = Stringify
	funcs["mapify"] = Mapify
	funcs["unmapify"] = Unmapify
	funcs["namify"] = Namify
	funcs["offset"] = Offset

	g.t.Funcs(funcs)

	g.addEmbeddedTemplates()

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

// WithLowercaseVariant is used to change the enum const values generated to not have the enum on them.
func (g *Generator) WithCaseInsensitiveParse() *Generator {
	g.lowercaseLookup = true
	g.caseInsensitive = true
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

// WithSQLInt is used to signal a string to be stored as an int.
func (g *Generator) WithSQLInt() *Generator {
	g.sqlint = true
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

// WithValues is used to add Values methods to the enum
func (g *Generator) WithValues() *Generator {
	g.values = true
	return g
}

// WithoutSnakeToCamel is used to add flag methods to the enum
func (g *Generator) WithoutSnakeToCamel() *Generator {
	g.leaveSnakeCase = true
	return g
}

// WithPrefix is used to add a custom prefix to the enum constants
func (g *Generator) WithPrefix(prefix string) *Generator {
	g.prefix = prefix
	return g
}

// WithPtr adds a way to get a pointer value straight from the const value.
func (g *Generator) WithPtr() *Generator {
	g.ptr = true
	return g
}

// WithSQLNullInt is used to add a null int option for SQL interactions.
func (g *Generator) WithSQLNullInt() *Generator {
	g.sqlNullInt = true
	return g
}

// WithSQLNullStr is used to add a null string option for SQL interactions.
func (g *Generator) WithSQLNullStr() *Generator {
	g.sqlNullStr = true
	return g
}

// WithMustParse is used to add a method `MustParse` that will panic on failure.
func (g *Generator) WithMustParse() *Generator {
	g.mustParse = true
	return g
}

// WithForceLower is used to force enums names to lower case while keeping variable names the same.
func (g *Generator) WithForceLower() *Generator {
	g.forceLower = true
	return g
}

// WithForceUpper is used to force enums names to upper case while keeping variable names the same.
func (g *Generator) WithForceUpper() *Generator {
	g.forceUpper = true
	return g
}

// WithNoComments is used to remove auto generated comments from the enum.
func (g *Generator) WithNoComments() *Generator {
	g.noComments = true
	return g
}

// WithBuildTags will add build tags to the generated file.
func (g *Generator) WithBuildTags(tags ...string) *Generator {
	g.buildTags = append(g.buildTags, tags...)
	return g
}

// WithAliases will set up aliases for the generator.
func (g *Generator) WithAliases(aliases map[string]string) *Generator {
	if aliases == nil {
		return g
	}
	g.replacementNames = aliases
	return g
}

func (g *Generator) anySQLEnabled() bool {
	return g.sql || g.sqlNullStr || g.sqlint || g.sqlNullInt
}

// ParseAliases is used to add aliases to replace during name sanitization.
func ParseAliases(aliases []string) (map[string]string, error) {
	aliasMap := map[string]string{}

	for _, str := range aliases {
		kvps := strings.Split(str, ",")
		for _, kvp := range kvps {
			parts := strings.Split(kvp, ":")
			if len(parts) != 2 {
				return nil, fmt.Errorf("invalid formatted alias entry %q, must be in the format \"key:value\"", kvp)
			}
			aliasMap[parts[0]] = parts[1]
		}
	}

	return aliasMap, nil
}

// WithTemplates is used to provide the filenames of additional templates.
func (g *Generator) WithTemplates(filenames ...string) *Generator {
	for _, ut := range template.Must(g.t.ParseFiles(filenames...)).Templates() {
		if _, ok := g.knownTemplates[ut.Name()]; !ok {
			g.userTemplateNames = append(g.userTemplateNames, ut.Name())
		}
	}
	g.updateTemplates()
	sort.Strings(g.userTemplateNames)
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
	err := g.t.ExecuteTemplate(vBuff, "header", map[string]interface{}{
		"package":   pkg,
		"version":   g.Version,
		"revision":  g.Revision,
		"buildDate": g.BuildDate,
		"builtBy":   g.BuiltBy,
		"buildTags": g.buildTags,
	})
	if err != nil {
		return nil, fmt.Errorf("failed writing header: %w", err)
	}

	// Make the output more consistent by iterating over sorted keys of map
	var keys []string
	for key := range enums {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var created int
	for _, name := range keys {
		ts := enums[name]

		// Parse the enum doc statement
		enum, pErr := g.parseEnum(ts)
		if pErr != nil {
			continue
		}

		created++
		data := map[string]interface{}{
			"enum":          enum,
			"name":          name,
			"lowercase":     g.lowercaseLookup,
			"nocase":        g.caseInsensitive,
			"nocomments":    g.noComments,
			"marshal":       g.marshal,
			"sql":           g.sql,
			"sqlint":        g.sqlint,
			"flag":          g.flag,
			"names":         g.names,
			"ptr":           g.ptr,
			"values":        g.values,
			"anySQLEnabled": g.anySQLEnabled(),
			"sqlnullint":    g.sqlNullInt,
			"sqlnullstr":    g.sqlNullStr,
			"mustparse":     g.mustParse,
			"forcelower":    g.forceLower,
			"forceupper":    g.forceUpper,
		}

		templateName := "enum"
		if enum.Type == "string" {
			templateName = "enum_string"
		}

		err = g.t.ExecuteTemplate(vBuff, templateName, data)
		if err != nil {
			return vBuff.Bytes(), fmt.Errorf("failed writing enum data for enum: %q: %w", name, err)
		}

		for _, userTemplateName := range g.userTemplateNames {
			err = g.t.ExecuteTemplate(vBuff, userTemplateName, data)
			if err != nil {
				return vBuff.Bytes(), fmt.Errorf("failed writing enum data for enum: %q, template: %v: %w", name, userTemplateName, err)
			}
		}
	}

	if created < 1 {
		// Don't save anything if we didn't actually generate any successful enums.
		return nil, nil
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
		return nil, errors.New("no doc on enum")
	}

	enum := &Enum{}

	enum.Name = ts.Name.Name
	enum.Type = fmt.Sprintf("%s", ts.Type)
	if !g.noPrefix {
		enum.Prefix = ts.Name.Name
	}
	if g.prefix != "" {
		enum.Prefix = g.prefix + enum.Prefix
	}

	commentPreEnumDecl, _, _ := strings.Cut(ts.Doc.Text(), `ENUM(`)
	enum.Comment = strings.TrimSpace(commentPreEnumDecl)

	enumDecl := getEnumDeclFromComments(ts.Doc.List)
	if enumDecl == "" {
		return nil, errors.New("failed parsing enum")
	}

	values := strings.Split(strings.TrimSuffix(strings.TrimPrefix(enumDecl, `ENUM(`), `)`), `,`)
	var (
		data     interface{}
		unsigned bool
	)
	if strings.HasPrefix(enum.Type, "u") {
		data = uint64(0)
		unsigned = true
	} else {
		data = int64(0)
	}
	for _, value := range values {
		var comment string

		// Trim and store comments
		if strings.Contains(value, parseCommentPrefix) {
			commentStartIndex := strings.Index(value, parseCommentPrefix)
			comment = value[commentStartIndex+len(parseCommentPrefix):]
			comment = strings.TrimSpace(unescapeComment(comment))
			// value without comment
			value = value[:commentStartIndex]
		}

		// Make sure to leave out any empty parts
		if value != "" {
			rawName := value
			valueStr := value

			if strings.Contains(value, `=`) {
				// Get the value specified and set the data to that value.
				equalIndex := strings.Index(value, `=`)
				dataVal := strings.TrimSpace(value[equalIndex+1:])
				if dataVal != "" {
					valueStr = dataVal
					rawName = value[:equalIndex]
					if enum.Type == "string" {
						if parsed, err := strconv.ParseInt(dataVal, 0, 64); err == nil {
							data = parsed
							valueStr = rawName
						}
						if isQuoted(dataVal) {
							valueStr = trimQuotes(dataVal)
						}
					} else if unsigned {
						newData, err := strconv.ParseUint(dataVal, 0, 64)
						if err != nil {
							err = fmt.Errorf("failed parsing the data part of enum value '%s': %w", value, err)
							fmt.Println(err)
							return nil, err
						}
						data = newData
					} else {
						newData, err := strconv.ParseInt(dataVal, 0, 64)
						if err != nil {
							err = fmt.Errorf("failed parsing the data part of enum value '%s': %w", value, err)
							fmt.Println(err)
							return nil, err
						}
						data = newData
					}
				} else {
					rawName = strings.TrimSuffix(rawName, `=`)
					fmt.Printf("Ignoring enum with '=' but no value after: %s\n", rawName)
				}
			}
			rawName = strings.TrimSpace(rawName)
			valueStr = strings.TrimSpace(valueStr)
			name := cases.Title(language.Und, cases.NoLower).String(rawName)
			prefixedName := name
			if name != skipHolder {
				prefixedName = enum.Prefix + name
				prefixedName = g.sanitizeValue(prefixedName)
				if !g.leaveSnakeCase {
					prefixedName = snakeToCamelCase(prefixedName)
				}
			}

			ev := EnumValue{Name: name, RawName: rawName, PrefixedName: prefixedName, ValueStr: valueStr, ValueInt: data, Comment: comment}
			enum.Values = append(enum.Values, ev)
			data = increment(data)
		}
	}

	// fmt.Printf("###\nENUM: %+v\n###\n", enum)

	return enum, nil
}

func isQuoted(s string) bool {
	s = strings.TrimSpace(s)
	return (strings.HasPrefix(s, `"`) && strings.HasSuffix(s, `"`)) || (strings.HasPrefix(s, `'`) && strings.HasSuffix(s, `'`))
}

func trimQuotes(s string) string {
	s = strings.TrimSpace(s)
	for _, quote := range []string{`"`, `'`} {
		s = strings.TrimPrefix(s, quote)
		s = strings.TrimSuffix(s, quote)
	}
	return s
}

func increment(d interface{}) interface{} {
	switch v := d.(type) {
	case uint64:
		return v + 1
	case int64:
		return v + 1
	}
	return d
}

func unescapeComment(comment string) string {
	val, err := url.QueryUnescape(comment)
	if err != nil {
		return comment
	}
	return val
}

// sanitizeValue will ensure the value name generated adheres to golang's
// identifier syntax as described here: https://golang.org/ref/spec#Identifiers
// identifier = letter { letter | unicode_digit }
// where letter can be unicode_letter or '_'
func (g *Generator) sanitizeValue(value string) string {
	// Keep skip value holders
	if value == skipHolder {
		return skipHolder
	}

	replacedValue := value
	for k, v := range g.replacementNames {
		replacedValue = strings.ReplaceAll(replacedValue, k, v)
	}

	nameBuilder := strings.Builder{}
	nameBuilder.Grow(len(replacedValue))

	for i, r := range replacedValue {
		// If the start character is not a unicode letter (this check includes the case of '_')
		// then we need to add an exported prefix, so tack on a 'X' at the beginning
		if i == 0 && !unicode.IsLetter(r) {
			nameBuilder.WriteRune('X')
		}

		if unicode.IsLetter(r) || unicode.IsNumber(r) || r == '_' {
			nameBuilder.WriteRune(r)
		}
	}

	return nameBuilder.String()
}

func snakeToCamelCase(value string) string {
	parts := strings.Split(value, "_")
	title := cases.Title(language.Und, cases.NoLower)

	for i, part := range parts {
		parts[i] = title.String(part)
	}
	value = strings.Join(parts, "")

	return value
}

// getEnumDeclFromComments parses the array of comment strings and creates a single Enum Declaration statement
// that is easier to deal with for the remainder of parsing.  It turns multi line declarations and makes a single
// string declaration.
func getEnumDeclFromComments(comments []*ast.Comment) string {
	const EnumPrefix = "ENUM("
	var (
		parts          []string
		lines          []string
		store          bool
		enumParamLevel int
		filteredLines  []string
	)

	for _, comment := range comments {
		lines = append(lines, breakCommentIntoLines(comment)...)
	}

	filteredLines = make([]string, 0, len(lines))
	for idx := range lines {
		line := lines[idx]
		// If we're not in the enum, and this line doesn't contain the
		// start string, then move along
		if !store && !strings.Contains(line, EnumPrefix) {
			continue
		}
		if !store {
			// We must have had the start value in here
			store = true
			enumParamLevel = 1
			start := strings.Index(line, EnumPrefix)
			line = line[start+len(EnumPrefix):]
		}
		lineParamLevel := strings.Count(line, "(")
		lineParamLevel = lineParamLevel - strings.Count(line, ")")

		if enumParamLevel+lineParamLevel < 1 {
			// We've ended, either with more than we need, or with just enough.  Now we need to find the end.
			for lineIdx, ch := range line {
				if ch == '(' {
					enumParamLevel = enumParamLevel + 1
					continue
				}
				if ch == ')' {
					enumParamLevel = enumParamLevel - 1
					if enumParamLevel == 0 {
						// We've found the end of the ENUM() definition,
						// Cut off the suffix and break out of the loop
						line = line[:lineIdx]
						store = false
						break
					}
				}
			}
		}

		filteredLines = append(filteredLines, line)
	}

	if enumParamLevel > 0 {
		fmt.Println("ENUM Parse error, there is a dangling '(' in your comment.")
		return ""
	}

	// Go over all the lines in this comment block
	for _, line := range filteredLines {
		_, trimmed := parseLinePart(line)
		if trimmed != "" {
			parts = append(parts, trimmed)
		}
	}

	joined := fmt.Sprintf("ENUM(%s)", strings.Join(parts, `,`))
	return joined
}

func parseLinePart(line string) (paramLevel int, trimmed string) {
	trimmed = line
	comment := ""
	if idx := strings.Index(line, parseCommentPrefix); idx >= 0 {
		trimmed = line[:idx]
		comment = "//" + url.QueryEscape(strings.TrimSpace(line[idx+2:]))
	}
	trimmed = trimAllTheThings(trimmed)
	trimmed += comment
	opens := strings.Count(line, `(`)
	closes := strings.Count(line, `)`)
	if opens > 0 {
		paramLevel += opens
	}
	if closes > 0 {
		paramLevel -= closes
	}
	return
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
// These lines should be pre-filtered so that we don't have to worry about
// the `ENUM(` prefix and the `)` suffix... those should already be removed.
func trimAllTheThings(thing string) string {
	preTrimmed := strings.TrimSuffix(strings.TrimSpace(thing), `,`)
	return strings.TrimSpace(preTrimmed)
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
