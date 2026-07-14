package generator

import (
	"fmt"
	"strconv"
	"strings"
)

// EnumConfigValue holds a configuration value with its validity flag.
type EnumConfigValue[T ~string | ~bool] struct {
	Value T
	Valid bool
}

// GetBool returns the boolean value if valid, otherwise returns the default value.
func (v *EnumConfigValue[bool]) GetBool(def bool) bool {
	if v.Valid {
		return v.Value
	}
	return def
}

// GetString returns the string value if valid, otherwise returns the default value.
func (v *EnumConfigValue[string]) GetString(def string) string {
	if v.Valid {
		return v.Value
	}
	return def
}

// EnumConfig holds configuration options specific to a single enum.
// These options can be specified inline via annotations and override global GeneratorConfig.
type EnumConfig struct {
	// Bool options
	NoPrefix        EnumConfigValue[bool] `json:"no_prefix"`
	NoIota          EnumConfigValue[bool] `json:"no_iota"`
	LowercaseLookup EnumConfigValue[bool] `json:"lowercase_lookup"`
	CaseInsensitive EnumConfigValue[bool] `json:"case_insensitive"`
	Marshal         EnumConfigValue[bool] `json:"marshal"`
	SQL             EnumConfigValue[bool] `json:"sql"`
	SQLInt          EnumConfigValue[bool] `json:"sql_int"`
	Flag            EnumConfigValue[bool] `json:"flag"`
	Names           EnumConfigValue[bool] `json:"names"`
	Values          EnumConfigValue[bool] `json:"values"`
	LeaveSnakeCase  EnumConfigValue[bool] `json:"leave_snake_case"`
	Ptr             EnumConfigValue[bool] `json:"ptr"`
	SQLNullInt      EnumConfigValue[bool] `json:"sql_null_int"`
	SQLNullStr      EnumConfigValue[bool] `json:"sql_null_str"`
	MustParse       EnumConfigValue[bool] `json:"must_parse"`
	ForceLower      EnumConfigValue[bool] `json:"force_lower"`
	ForceUpper      EnumConfigValue[bool] `json:"force_upper"`
	NoComments      EnumConfigValue[bool] `json:"no_comments"`
	NoParse         EnumConfigValue[bool] `json:"no_parse"`

	// String options
	Prefix EnumConfigValue[string] `json:"prefix"`

	// Slice/map options (not supported inline for simplicity)
	// BuildTags         []string
	// ReplacementNames  map[string]string
	// TemplateFileNames []string
}

// NewEnumConfig creates a new EnumConfig with default values.
func NewEnumConfig() *EnumConfig {
	return &EnumConfig{}
}

// ParseAnnotation parses a single annotation string (e.g., "@marshal", "@marshal:true", "@prefix=\"My\"")
// and updates the EnumConfig accordingly.
func (ec *EnumConfig) ParseAnnotation(annotation string) error {
	annotation = strings.TrimSpace(annotation)
	if annotation == "" {
		return nil
	}

	// Remove @ prefix
	if !strings.HasPrefix(annotation, "@") {
		return fmt.Errorf("annotation must start with @: %s", annotation)
	}
	annotation = annotation[1:]

	// Check for key:value format (e.g., @marshal:true, @marshal:false)
	if strings.Contains(annotation, ":") {
		parts := strings.SplitN(annotation, ":", 2)
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Parse boolean value
		if value == "true" || value == "false" {
			boolValue, _ := strconv.ParseBool(value)
			return ec.setBoolOption(key, boolValue)
		}

		// String value (could be quoted)
		if len(value) >= 2 && ((value[0] == '"' && value[len(value)-1] == '"') ||
			(value[0] == '\'' && value[len(value)-1] == '\'')) {
			value = value[1 : len(value)-1]
		}

		return ec.setStringOption(key, value)
	}

	// Check for key=value format (legacy style, e.g., @prefix="My")
	if strings.Contains(annotation, "=") {
		parts := strings.SplitN(annotation, "=", 2)
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Remove quotes if present
		if len(value) >= 2 && ((value[0] == '"' && value[len(value)-1] == '"') ||
			(value[0] == '\'' && value[len(value)-1] == '\'')) {
			value = value[1 : len(value)-1]
		}

		return ec.setStringOption(key, value)
	}

	// Boolean flag without explicit value (defaults to true)
	return ec.setBoolOption(annotation, true)
}

// setBoolOption sets a boolean option in the EnumConfig.
func (ec *EnumConfig) setBoolOption(key string, value bool) error {
	switch key {
	case "noprefix":
		ec.NoPrefix = EnumConfigValue[bool]{Value: value, Valid: true}
	case "noiota":
		ec.NoIota = EnumConfigValue[bool]{Value: value, Valid: true}
	case "lower":
		ec.LowercaseLookup = EnumConfigValue[bool]{Value: value, Valid: true}
	case "nocase":
		ec.CaseInsensitive = EnumConfigValue[bool]{Value: value, Valid: true}
		if value {
			ec.LowercaseLookup = EnumConfigValue[bool]{Value: true, Valid: true} // nocase forces lower
		}
	case "marshal":
		ec.Marshal = EnumConfigValue[bool]{Value: value, Valid: true}
	case "sql":
		ec.SQL = EnumConfigValue[bool]{Value: value, Valid: true}
	case "sqlint":
		ec.SQLInt = EnumConfigValue[bool]{Value: value, Valid: true}
	case "flag":
		ec.Flag = EnumConfigValue[bool]{Value: value, Valid: true}
	case "names":
		ec.Names = EnumConfigValue[bool]{Value: value, Valid: true}
	case "values":
		ec.Values = EnumConfigValue[bool]{Value: value, Valid: true}
	case "nocamel":
		ec.LeaveSnakeCase = EnumConfigValue[bool]{Value: value, Valid: true}
	case "ptr":
		ec.Ptr = EnumConfigValue[bool]{Value: value, Valid: true}
	case "sqlnullint":
		ec.SQLNullInt = EnumConfigValue[bool]{Value: value, Valid: true}
	case "sqlnullstr":
		ec.SQLNullStr = EnumConfigValue[bool]{Value: value, Valid: true}
	case "mustparse":
		ec.MustParse = EnumConfigValue[bool]{Value: value, Valid: true}
	case "forcelower":
		ec.ForceLower = EnumConfigValue[bool]{Value: value, Valid: true}
	case "forceupper":
		ec.ForceUpper = EnumConfigValue[bool]{Value: value, Valid: true}
	case "nocomments":
		ec.NoComments = EnumConfigValue[bool]{Value: value, Valid: true}
	case "noparse":
		ec.NoParse = EnumConfigValue[bool]{Value: value, Valid: true}
	default:
		return fmt.Errorf("unknown annotation: @%s", key)
	}

	return nil
}

// setStringOption sets a string option in the EnumConfig.
func (ec *EnumConfig) setStringOption(key, value string) error {
	switch key {
	case "prefix":
		ec.Prefix = EnumConfigValue[string]{Value: value, Valid: true}
	default:
		return fmt.Errorf("unknown annotation with value: @%s=%s", key, value)
	}
	return nil
}
