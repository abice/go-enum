package generator

// GeneratorConfig holds all the configuration options for the Generator
type GeneratorConfig struct {
	NoPrefix          bool              `json:"no_prefix"`
	NoIota            bool              `json:"no_iota"`
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
	TemplateFileNames []string          `json:"template_file_names"`
}

func NewGeneratorConfig() *GeneratorConfig {
	return &GeneratorConfig{
		NoPrefix:         false,
		ReplacementNames: map[string]string{},
		JSONPkg:          "encoding/json",
	}
}

// Option is a function that modifies a Generator
type Option func(*GeneratorConfig)

// WithNoPrefix is used to change the enum const values generated to not have the enum on them.
func WithNoPrefix() Option {
	return func(g *GeneratorConfig) {
		g.NoPrefix = true
	}
}

// WithNoIota is used to generate enum constants with explicit values instead of using iota.
func WithNoIota() Option {
	return func(g *GeneratorConfig) {
		g.NoIota = true
	}
}

// WithLowercaseVariant is used to change the enum const values generated to not have the enum on them.
func WithLowercaseVariant() Option {
	return func(g *GeneratorConfig) {
		g.LowercaseLookup = true
	}
}

// WithCaseInsensitiveParse is used to change the enum const values generated to not have the enum on them.
func WithCaseInsensitiveParse() Option {
	return func(g *GeneratorConfig) {
		g.LowercaseLookup = true
		g.CaseInsensitive = true
	}
}

// WithMarshal is used to add marshalling to the enum
func WithMarshal() Option {
	return func(g *GeneratorConfig) {
		g.Marshal = true
	}
}

// WithSQLDriver is used to add marshalling to the enum
func WithSQLDriver() Option {
	return func(g *GeneratorConfig) {
		g.SQL = true
	}
}

// WithSQLInt is used to signal a string to be stored as an int.
func WithSQLInt() Option {
	return func(g *GeneratorConfig) {
		g.SQLInt = true
	}
}

// WithFlag is used to add flag methods to the enum
func WithFlag() Option {
	return func(g *GeneratorConfig) {
		g.Flag = true
	}
}

// WithNames is used to add Names methods to the enum
func WithNames() Option {
	return func(g *GeneratorConfig) {
		g.Names = true
	}
}

// WithValues is used to add Values methods to the enum
func WithValues() Option {
	return func(g *GeneratorConfig) {
		g.Values = true
	}
}

// WithoutSnakeToCamel is used to add flag methods to the enum
func WithoutSnakeToCamel() Option {
	return func(g *GeneratorConfig) {
		g.LeaveSnakeCase = true
	}
}

// WithJsonPkg is used to add a custom json package to the imports
func WithJsonPkg(pkg string) Option {
	return func(g *GeneratorConfig) {
		g.JSONPkg = pkg
	}
}

// WithPrefix is used to add a custom prefix to the enum constants
func WithPrefix(prefix string) Option {
	return func(g *GeneratorConfig) {
		g.Prefix = prefix
	}
}

// WithPtr adds a way to get a pointer value straight from the const value.
func WithPtr() Option {
	return func(g *GeneratorConfig) {
		g.Ptr = true
	}
}

// WithSQLNullInt is used to add a null int option for SQL interactions.
func WithSQLNullInt() Option {
	return func(g *GeneratorConfig) {
		g.SQLNullInt = true
	}
}

// WithSQLNullStr is used to add a null string option for SQL interactions.
func WithSQLNullStr() Option {
	return func(g *GeneratorConfig) {
		g.SQLNullStr = true
	}
}

// WithMustParse is used to add a method `MustParse` that will panic on failure.
func WithMustParse() Option {
	return func(g *GeneratorConfig) {
		g.MustParse = true
	}
}

// WithForceLower is used to force enums names to lower case while keeping variable names the same.
func WithForceLower() Option {
	return func(g *GeneratorConfig) {
		g.ForceLower = true
	}
}

// WithForceUpper is used to force enums names to upper case while keeping variable names the same.
func WithForceUpper() Option {
	return func(g *GeneratorConfig) {
		g.ForceUpper = true
	}
}

// WithNoComments is used to remove auto generated comments from the enum.
func WithNoComments() Option {
	return func(g *GeneratorConfig) {
		g.NoComments = true
	}
}

// WithBuildTags will add build tags to the generated file.
func WithBuildTags(tags ...string) Option {
	return func(g *GeneratorConfig) {
		g.BuildTags = append(g.BuildTags, tags...)
	}
}

// WithAliases will set up aliases for the generator.
func WithAliases(aliases map[string]string) Option {
	return func(g *GeneratorConfig) {
		if aliases == nil {
			return
		}
		g.ReplacementNames = aliases
	}
}

// WithTemplates is used to provide the filenames of additional templates.
func WithTemplates(filenames ...string) Option {
	return func(g *GeneratorConfig) {
		// Note: Template processing is deferred to the generator constructor
		// because we need access to the template collection and knownTemplates
		g.TemplateFileNames = append(g.TemplateFileNames, filenames...)
	}
}
