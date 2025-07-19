package generator

import (
	"sort"
	"text/template"
)

// GeneratorOptions holds all the configuration options for the Generator
type GeneratorOptions struct {
	NoPrefix          bool              `json:"no_prefix"`
	LowercaseLookup   bool              `json:"lowercase_lookup"`
	CaseInsensitive   bool              `json:"case_insensitive"`
	Marshal           bool              `json:"marshal"`
	SQL               bool              `json:"sql"`
	SQLInt            bool              `json:"sql_int"`
	Flag              bool              `json:"flag"`
	Names             bool              `json:"names"`
	Values            bool              `json:"values"`
	LeaveSnakeCase    bool              `json:"leave_snake_case"`
	JSONPkg           string            `json:"json_pkg"`
	Prefix            string            `json:"prefix"`
	SQLNullInt        bool              `json:"sql_null_int"`
	SQLNullStr        bool              `json:"sql_null_str"`
	Ptr               bool              `json:"ptr"`
	MustParse         bool              `json:"must_parse"`
	ForceLower        bool              `json:"force_lower"`
	ForceUpper        bool              `json:"force_upper"`
	NoComments        bool              `json:"no_comments"`
	BuildTags         []string          `json:"build_tags"`
	ReplacementNames  map[string]string `json:"replacement_names"`
	UserTemplateNames []string          `json:"user_template_names"`
}

// Option is a function that modifies a Generator
type Option func(*Generator)

// WithNoPrefix is used to change the enum const values generated to not have the enum on them.
func WithNoPrefix() Option {
	return func(g *Generator) {
		g.NoPrefix = true
	}
}

// WithLowercaseVariant is used to change the enum const values generated to not have the enum on them.
func WithLowercaseVariant() Option {
	return func(g *Generator) {
		g.LowercaseLookup = true
	}
}

// WithCaseInsensitiveParse is used to change the enum const values generated to not have the enum on them.
func WithCaseInsensitiveParse() Option {
	return func(g *Generator) {
		g.LowercaseLookup = true
		g.CaseInsensitive = true
	}
}

// WithMarshal is used to add marshalling to the enum
func WithMarshal() Option {
	return func(g *Generator) {
		g.Marshal = true
	}
}

// WithSQLDriver is used to add marshalling to the enum
func WithSQLDriver() Option {
	return func(g *Generator) {
		g.SQL = true
	}
}

// WithSQLInt is used to signal a string to be stored as an int.
func WithSQLInt() Option {
	return func(g *Generator) {
		g.SQLInt = true
	}
}

// WithFlag is used to add flag methods to the enum
func WithFlag() Option {
	return func(g *Generator) {
		g.Flag = true
	}
}

// WithNames is used to add Names methods to the enum
func WithNames() Option {
	return func(g *Generator) {
		g.Names = true
	}
}

// WithValues is used to add Values methods to the enum
func WithValues() Option {
	return func(g *Generator) {
		g.Values = true
	}
}

// WithoutSnakeToCamel is used to add flag methods to the enum
func WithoutSnakeToCamel() Option {
	return func(g *Generator) {
		g.LeaveSnakeCase = true
	}
}

// WithJsonPkg is used to add a custom json package to the imports
func WithJsonPkg(pkg string) Option {
	return func(g *Generator) {
		g.JSONPkg = pkg
	}
}

// WithPrefix is used to add a custom prefix to the enum constants
func WithPrefix(prefix string) Option {
	return func(g *Generator) {
		g.Prefix = prefix
	}
}

// WithPtr adds a way to get a pointer value straight from the const value.
func WithPtr() Option {
	return func(g *Generator) {
		g.Ptr = true
	}
}

// WithSQLNullInt is used to add a null int option for SQL interactions.
func WithSQLNullInt() Option {
	return func(g *Generator) {
		g.SQLNullInt = true
	}
}

// WithSQLNullStr is used to add a null string option for SQL interactions.
func WithSQLNullStr() Option {
	return func(g *Generator) {
		g.SQLNullStr = true
	}
}

// WithMustParse is used to add a method `MustParse` that will panic on failure.
func WithMustParse() Option {
	return func(g *Generator) {
		g.MustParse = true
	}
}

// WithForceLower is used to force enums names to lower case while keeping variable names the same.
func WithForceLower() Option {
	return func(g *Generator) {
		g.ForceLower = true
	}
}

// WithForceUpper is used to force enums names to upper case while keeping variable names the same.
func WithForceUpper() Option {
	return func(g *Generator) {
		g.ForceUpper = true
	}
}

// WithNoComments is used to remove auto generated comments from the enum.
func WithNoComments() Option {
	return func(g *Generator) {
		g.NoComments = true
	}
}

// WithBuildTags will add build tags to the generated file.
func WithBuildTags(tags ...string) Option {
	return func(g *Generator) {
		g.BuildTags = append(g.BuildTags, tags...)
	}
}

// WithAliases will set up aliases for the generator.
func WithAliases(aliases map[string]string) Option {
	return func(g *Generator) {
		if aliases == nil {
			return
		}
		g.ReplacementNames = aliases
	}
}

// WithTemplates is used to provide the filenames of additional templates.
func WithTemplates(filenames ...string) Option {
	return func(g *Generator) {
		for _, ut := range template.Must(g.t.ParseFiles(filenames...)).Templates() {
			if _, ok := g.knownTemplates[ut.Name()]; !ok {
				g.UserTemplateNames = append(g.UserTemplateNames, ut.Name())
			}
		}
		g.updateTemplates()
		sort.Strings(g.UserTemplateNames)
	}
}
