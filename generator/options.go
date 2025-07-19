package generator

import (
	"sort"
	"text/template"
)

// GeneratorOptions holds all the configuration options for the Generator
type GeneratorOptions struct {
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
	jsonPkg           string
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
	userTemplateNames []string
}

// Option is a function that modifies a Generator
type Option func(*Generator)

// WithNoPrefix is used to change the enum const values generated to not have the enum on them.
func WithNoPrefix() Option {
	return func(g *Generator) {
		g.noPrefix = true
	}
}

// WithLowercaseVariant is used to change the enum const values generated to not have the enum on them.
func WithLowercaseVariant() Option {
	return func(g *Generator) {
		g.lowercaseLookup = true
	}
}

// WithCaseInsensitiveParse is used to change the enum const values generated to not have the enum on them.
func WithCaseInsensitiveParse() Option {
	return func(g *Generator) {
		g.lowercaseLookup = true
		g.caseInsensitive = true
	}
}

// WithMarshal is used to add marshalling to the enum
func WithMarshal() Option {
	return func(g *Generator) {
		g.marshal = true
	}
}

// WithSQLDriver is used to add marshalling to the enum
func WithSQLDriver() Option {
	return func(g *Generator) {
		g.sql = true
	}
}

// WithSQLInt is used to signal a string to be stored as an int.
func WithSQLInt() Option {
	return func(g *Generator) {
		g.sqlint = true
	}
}

// WithFlag is used to add flag methods to the enum
func WithFlag() Option {
	return func(g *Generator) {
		g.flag = true
	}
}

// WithNames is used to add Names methods to the enum
func WithNames() Option {
	return func(g *Generator) {
		g.names = true
	}
}

// WithValues is used to add Values methods to the enum
func WithValues() Option {
	return func(g *Generator) {
		g.values = true
	}
}

// WithoutSnakeToCamel is used to add flag methods to the enum
func WithoutSnakeToCamel() Option {
	return func(g *Generator) {
		g.leaveSnakeCase = true
	}
}

// WithJsonPkg is used to add a custom json package to the imports
func WithJsonPkg(pkg string) Option {
	return func(g *Generator) {
		g.jsonPkg = pkg
	}
}

// WithPrefix is used to add a custom prefix to the enum constants
func WithPrefix(prefix string) Option {
	return func(g *Generator) {
		g.prefix = prefix
	}
}

// WithPtr adds a way to get a pointer value straight from the const value.
func WithPtr() Option {
	return func(g *Generator) {
		g.ptr = true
	}
}

// WithSQLNullInt is used to add a null int option for SQL interactions.
func WithSQLNullInt() Option {
	return func(g *Generator) {
		g.sqlNullInt = true
	}
}

// WithSQLNullStr is used to add a null string option for SQL interactions.
func WithSQLNullStr() Option {
	return func(g *Generator) {
		g.sqlNullStr = true
	}
}

// WithMustParse is used to add a method `MustParse` that will panic on failure.
func WithMustParse() Option {
	return func(g *Generator) {
		g.mustParse = true
	}
}

// WithForceLower is used to force enums names to lower case while keeping variable names the same.
func WithForceLower() Option {
	return func(g *Generator) {
		g.forceLower = true
	}
}

// WithForceUpper is used to force enums names to upper case while keeping variable names the same.
func WithForceUpper() Option {
	return func(g *Generator) {
		g.forceUpper = true
	}
}

// WithNoComments is used to remove auto generated comments from the enum.
func WithNoComments() Option {
	return func(g *Generator) {
		g.noComments = true
	}
}

// WithBuildTags will add build tags to the generated file.
func WithBuildTags(tags ...string) Option {
	return func(g *Generator) {
		g.buildTags = append(g.buildTags, tags...)
	}
}

// WithAliases will set up aliases for the generator.
func WithAliases(aliases map[string]string) Option {
	return func(g *Generator) {
		if aliases == nil {
			return
		}
		g.replacementNames = aliases
	}
}

// WithTemplates is used to provide the filenames of additional templates.
func WithTemplates(filenames ...string) Option {
	return func(g *Generator) {
		for _, ut := range template.Must(g.t.ParseFiles(filenames...)).Templates() {
			if _, ok := g.knownTemplates[ut.Name()]; !ok {
				g.userTemplateNames = append(g.userTemplateNames, ut.Name())
			}
		}
		g.updateTemplates()
		sort.Strings(g.userTemplateNames)
	}
}
