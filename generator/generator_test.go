package generator

import (
	"errors"
	"fmt"
	"go/parser"
	"strings"
	"testing"

	"github.com/bradleyjkemp/cupaloy/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testExample       = `example_test.go`
	testExampleNoIota = `example_no_iota_test.go`
)

// TestNoStructInputFile
func TestNoStructFile(t *testing.T) {
	input := `package test
	// Behavior
	type SomeInterface interface{

	}
	`
	g := NewGenerator(WithoutSnakeToCamel())
	f, err := parser.ParseFile(g.fileSet, "TestRequiredErrors", input, parser.ParseComments)
	assert.Nil(t, err, "Error parsing no struct input")

	output, err := g.Generate(f)
	assert.Nil(t, err, "Error generating formatted code")
	if false { // Debugging statement
		fmt.Println(string(output))
	}
}

// TestNoFile
func TestNoFile(t *testing.T) {
	g := NewGenerator(WithoutSnakeToCamel())
	// Parse the file given in arguments
	_, err := g.GenerateFromFile("")
	assert.NotNil(t, err, "Error generating formatted code")
}

// TestExampleFile
func TestExampleFile(t *testing.T) {
	g := NewGenerator(
		WithMarshal(),
		WithSQLDriver(),
		WithCaseInsensitiveParse(),
		WithNames(),
		WithoutSnakeToCamel(),
	)
	// Parse the file given in arguments
	imported, err := g.GenerateFromFile(testExample)
	require.Nil(t, err, "Error generating formatted code")

	outputLines := strings.Split(string(imported), "\n")
	cupaloy.SnapshotT(t, outputLines)

	if false {
		fmt.Println(string(imported))
	}
}

// TestExampleFileEmptyHeaders tests the generator with empty header values
func TestExampleFileEmptyHeaders(t *testing.T) {
	g := NewGenerator(
		WithMarshal(),
		WithSQLDriver(),
		WithCaseInsensitiveParse(),
		WithNames(),
		WithoutSnakeToCamel(),
	)
	g.Version = ""
	g.Revision = ""
	g.BuildDate = ""
	g.BuiltBy = ""
	// Parse the file given in arguments
	imported, err := g.GenerateFromFile(testExample)
	require.Nil(t, err, "Error generating formatted code")

	outputLines := strings.Split(string(imported), "\n")
	cupaloy.SnapshotT(t, outputLines)

	const maxLinesToCheck = 10
	// Check that the first 10 lines do not contain version information
	// or build information
	for idx := range min(maxLinesToCheck, len(outputLines)) {
		assert.NotContains(t, outputLines[idx], "Version:")
		assert.NotContains(t, outputLines[idx], "Revision:")
		assert.NotContains(t, outputLines[idx], "BuildDate:")
		assert.NotContains(t, outputLines[idx], "BuiltBy:")
	}

	if false {
		fmt.Println(string(imported))
	}
}

// TestExampleFileMoreOptions
func TestExampleFileMoreOptions(t *testing.T) {
	g := NewGenerator(
		WithMarshal(),
		WithSQLDriver(),
		WithCaseInsensitiveParse(),
		WithNames(),
		WithoutSnakeToCamel(),
		WithMustParse(),
		WithForceLower(),
		WithTemplates(`../example/user_template.tmpl`),
	)
	// Parse the file given in arguments
	imported, err := g.GenerateFromFile(testExample)
	require.Nil(t, err, "Error generating formatted code")

	outputLines := strings.Split(string(imported), "\n")
	cupaloy.SnapshotT(t, outputLines)

	if false {
		fmt.Println(string(imported))
	}
}

// TestExampleFileMoreOptionsWithForceUpper — test with force upper option
func TestExampleFileMoreOptionsWithForceUpper(t *testing.T) {
	g := NewGenerator(
		WithMarshal(),
		WithSQLDriver(),
		WithNames(),
		WithoutSnakeToCamel(),
		WithMustParse(),
		WithForceUpper(),
		WithTemplates(`../example/user_template.tmpl`),
	)
	// Parse the file given in arguments
	imported, err := g.GenerateFromFile(testExample)
	require.Nil(t, err, "Error generating formatted code")

	outputLines := strings.Split(string(imported), "\n")
	cupaloy.SnapshotT(t, outputLines)

	if false {
		fmt.Println(string(imported))
	}
}

// TestExampleFile
func TestNoPrefixExampleFile(t *testing.T) {
	g := NewGenerator(
		WithMarshal(),
		WithLowercaseVariant(),
		WithNoPrefix(),
		WithFlag(),
		WithoutSnakeToCamel(),
	)
	// Parse the file given in arguments
	imported, err := g.GenerateFromFile(testExample)
	require.Nil(t, err, "Error generating formatted code")

	outputLines := strings.Split(string(imported), "\n")
	cupaloy.SnapshotT(t, outputLines)

	if false {
		fmt.Println(string(imported))
	}
}

func TestNoIotaExampleFile(t *testing.T) {
	g := NewGenerator(
		WithMarshal(),
		WithLowercaseVariant(),
		WithNoIota(),
		WithFlag(),
		WithoutSnakeToCamel(),
	)
	// Parse the file given in arguments
	imported, err := g.GenerateFromFile(testExample)
	require.Nil(t, err, "Error generating formatted code")

	outputLines := strings.Split(string(imported), "\n")
	cupaloy.SnapshotT(t, outputLines)

	if false {
		fmt.Println(string(imported))
	}
}

func TestNoIotaOnlyExampleFile(t *testing.T) {
	g := NewGenerator(
		WithMarshal(),
		WithLowercaseVariant(),
		WithNoIota(),
		WithFlag(),
		WithoutSnakeToCamel(),
	)
	// Parse the file given in arguments
	imported, err := g.GenerateFromFile(testExampleNoIota)
	require.Nil(t, err, "Error generating formatted code")

	outputLines := strings.Split(string(imported), "\n")
	cupaloy.SnapshotT(t, outputLines)

	if false {
		fmt.Println(string(imported))
	}
}

// TestExampleFile
func TestReplacePrefixExampleFile(t *testing.T) {
	g := NewGenerator(
		WithMarshal(),
		WithLowercaseVariant(),
		WithNoPrefix(),
		WithPrefix("MyPrefix_"),
		WithFlag(),
		WithoutSnakeToCamel(),
	)
	// Parse the file given in arguments
	imported, err := g.GenerateFromFile(testExample)
	require.Nil(t, err, "Error generating formatted code")

	outputLines := strings.Split(string(imported), "\n")
	cupaloy.SnapshotT(t, outputLines)

	if false {
		fmt.Println(string(imported))
	}
}

// TestExampleFile
func TestNoPrefixExampleFileWithSnakeToCamel(t *testing.T) {
	g := NewGenerator(
		WithMarshal(),
		WithLowercaseVariant(),
		WithNoPrefix(),
		WithFlag(),
	)

	// Parse the file given in arguments
	imported, err := g.GenerateFromFile(testExample)
	require.Nil(t, err, "Error generating formatted code")

	outputLines := strings.Split(string(imported), "\n")
	cupaloy.SnapshotT(t, outputLines)

	if false {
		fmt.Println(string(imported))
	}
}

// TestCustomPrefixExampleFile
func TestCustomPrefixExampleFile(t *testing.T) {
	g := NewGenerator(
		WithMarshal(),
		WithLowercaseVariant(),
		WithNoPrefix(),
		WithFlag(),
		WithoutSnakeToCamel(),
		WithPtr(),
		WithSQLNullInt(),
		WithSQLNullStr(),
		WithPrefix("Custom_prefix_"),
	)
	// Parse the file given in arguments
	imported, err := g.GenerateFromFile(testExample)
	require.Nil(t, err, "Error generating formatted code")

	outputLines := strings.Split(string(imported), "\n")
	cupaloy.SnapshotT(t, outputLines)

	if false {
		fmt.Println(string(imported))
	}
}

func TestAliasParsing(t *testing.T) {
	tests := map[string]struct {
		input        []string
		resultingMap map[string]string
		err          error
	}{
		"no aliases": {
			resultingMap: map[string]string{},
		},
		"bad input": {
			input: []string{"a:b,c"},
			err:   errors.New(`invalid formatted alias entry "c", must be in the format "key:value"`),
		},
		"multiple arrays": {
			input: []string{
				`!:Bang,a:a`,
				`@:AT`,
				`&:AND,|:OR`,
			},
			resultingMap: map[string]string{
				"a": "a",
				"!": "Bang",
				"@": "AT",
				"&": "AND",
				"|": "OR",
			},
		},
		"more types": {
			input: []string{
				`*:star,+:PLUS`,
				`-:less`,
				`#:HASH,!:Bang`,
			},
			resultingMap: map[string]string{
				"*": "star",
				"+": "PLUS",
				"-": "less",
				"#": "HASH",
				"!": "Bang",
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			replacementNames, err := ParseAliases(tc.input)
			if tc.err != nil {
				require.Error(t, err)
				require.EqualError(t, err, tc.err.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.resultingMap, replacementNames)
			}
		})
	}
}

// TestEnumParseFailure
func TestEnumParseFailure(t *testing.T) {
	input := `package test
	// Behavior
	type SomeInterface interface{

	}

	// ENUM(
	//	a,
	//}
	type Animal int
	`
	g := NewGenerator(WithoutSnakeToCamel())
	f, err := parser.ParseFile(g.fileSet, "TestRequiredErrors", input, parser.ParseComments)
	assert.Nil(t, err, "Error parsing no struct input")

	output, err := g.Generate(f)
	assert.Nil(t, err, "Error generating formatted code")
	assert.Empty(t, string(output))
	if false { // Debugging statement
		fmt.Println(string(output))
	}
}

// TestUintInvalidParsing
func TestUintInvalidParsing(t *testing.T) {
	input := `package test
	// ENUM(
	//	a=-1,
	//)
	type Animal uint
	`
	g := NewGenerator(WithoutSnakeToCamel())
	f, err := parser.ParseFile(g.fileSet, "TestRequiredErrors", input, parser.ParseComments)
	assert.Nil(t, err, "Error parsing no struct input")

	output, err := g.Generate(f)
	assert.Nil(t, err, "Error generating formatted code")
	assert.Empty(t, string(output))
	if false { // Debugging statement
		fmt.Println(string(output))
	}
}

// TestIntInvalidParsing
func TestIntInvalidParsing(t *testing.T) {
	input := `package test
	// ENUM(
	//	a=c,
	//)
	type Animal int
	`
	g := NewGenerator(WithoutSnakeToCamel())
	f, err := parser.ParseFile(g.fileSet, "TestRequiredErrors", input, parser.ParseComments)
	assert.Nil(t, err, "Error parsing no struct input")

	output, err := g.Generate(f)
	assert.Nil(t, err, "Error generating formatted code")
	assert.Empty(t, string(output))
	if false { // Debugging statement
		fmt.Println(string(output))
	}
}

// TestAliasing
func TestAliasing(t *testing.T) {
	input := `package test
	// ENUM(a,b,CDEF) with some extra text
	type Animal int
	`
	aliases, err := ParseAliases([]string{"CDEF:C"})
	require.NoError(t, err)
	g := NewGenerator(
		WithoutSnakeToCamel(),
		WithAliases(aliases),
	)
	f, err := parser.ParseFile(g.fileSet, "TestRequiredErrors", input, parser.ParseComments)
	assert.Nil(t, err, "Error parsing no struct input")

	output, err := g.Generate(f)
	assert.Nil(t, err, "Error generating formatted code")
	assert.Contains(t, string(output), "// AnimalC is a Animal of type CDEF.")
	if false { // Debugging statement
		fmt.Println(string(output))
	}
}

// TestParanthesesParsing
func TestParenthesesParsing(t *testing.T) {
	input := `package test
	// This is a pre-enum comment that needs (to be handled properly)
	// ENUM(
	//	abc (x),
	//). This is an extra string comment (With parentheses of it's own)
	// And (another line) with Parentheses
	type Animal string
	`
	g := NewGenerator()
	f, err := parser.ParseFile(g.fileSet, "TestRequiredErrors", input, parser.ParseComments)
	assert.Nil(t, err, "Error parsing no struct input")

	output, err := g.Generate(f)
	assert.Nil(t, err, "Error generating formatted code")
	assert.Contains(t, string(output), "// AnimalAbcX is a Animal of type abc (x).")
	assert.NotContains(t, string(output), "// AnimalAnd")
	if false { // Debugging statement
		fmt.Println(string(output))
	}
}

// TestQuotedStrings
func TestQuotedStrings(t *testing.T) {
	input := `package test
	// This is a pre-enum comment that needs (to be handled properly)
	// ENUM(
	//	abc (x),
	//  ghi = "20",
	//). This is an extra string comment (With parentheses of it's own)
	// And (another line) with Parentheses
	type Animal string
	`
	g := NewGenerator()
	f, err := parser.ParseFile(g.fileSet, "TestRequiredErrors", input, parser.ParseComments)
	assert.Nil(t, err, "Error parsing no struct input")

	output, err := g.Generate(f)
	assert.Nil(t, err, "Error generating formatted code")
	assert.Contains(t, string(output), "// AnimalAbcX is a Animal of type abc (x).")
	assert.Contains(t, string(output), "AnimalGhi Animal = \"20\"")
	assert.NotContains(t, string(output), "// AnimalAnd")
	if false { // Debugging statement
		fmt.Println(string(output))
	}
}

func TestStringWithSingleDoubleQuoteValue(t *testing.T) {
	input := `package test
	// ENUM(DoubleQuote='"')
	type Char string
	`
	g := NewGenerator()
	f, err := parser.ParseFile(g.fileSet, "TestRequiredErrors", input, parser.ParseComments)
	assert.Nil(t, err, "Error parsing no struct input")

	output, err := g.Generate(f)
	assert.Nil(t, err, "Error generating formatted code")
	assert.Contains(t, string(output), "CharDoubleQuote Char = \"\\\"\"")
	if false { // Debugging statement
		fmt.Println(string(output))
	}
}

func TestStringWithSingleSingleQuoteValue(t *testing.T) {
	input := `package test
	// ENUM(SingleQuote="'")
	type Char string
	`
	g := NewGenerator()
	f, err := parser.ParseFile(g.fileSet, "TestRequiredErrors", input, parser.ParseComments)
	assert.Nil(t, err, "Error parsing no struct input")

	output, err := g.Generate(f)
	assert.Nil(t, err, "Error generating formatted code")
	assert.Contains(t, string(output), "CharSingleQuote Char = \"'\"")
	if false { // Debugging statement
		fmt.Println(string(output))
	}
}

func TestStringWithSingleBacktickValue(t *testing.T) {
	input := `package test
	// ENUM(SingleQuote="` + "`" + `")
	type Char string
	`
	g := NewGenerator()
	f, err := parser.ParseFile(g.fileSet, "TestRequiredErrors", input, parser.ParseComments)
	assert.Nil(t, err, "Error parsing no struct input")

	output, err := g.Generate(f)
	assert.Nil(t, err, "Error generating formatted code")
	assert.Contains(t, string(output), "CharSingleQuote Char = \"`\"")
	if false { // Debugging statement
		fmt.Println(string(output))
	}
}

// TestNewGeneratorWithConfig tests the NewGeneratorWithConfig constructor
func TestNewGeneratorWithConfig(t *testing.T) {
	config := GeneratorConfig{
		NoPrefix:        true,
		LowercaseLookup: true,
		Marshal:         true,
		SQL:             true,
		SQLInt:          true,
		JSONPkg:         "custom/json",
		Prefix:          "TestPrefix",
		BuildTags:       []string{"tag1", "tag2"},
		Initialisms:     []string{"HTTP", "URL"},
		NoComments:      true,
		Values:          true,
	}

	g := NewGeneratorWithConfig(config)

	assert.NotNil(t, g)
	assert.Equal(t, config.NoPrefix, g.NoPrefix)
	assert.Equal(t, config.LowercaseLookup, g.LowercaseLookup)
	assert.Equal(t, config.Marshal, g.Marshal)
	assert.Equal(t, config.SQL, g.SQL)
	assert.Equal(t, config.SQLInt, g.SQLInt)
	assert.Equal(t, config.JSONPkg, g.JSONPkg)
	assert.Equal(t, config.Prefix, g.Prefix)
	assert.Equal(t, config.BuildTags, g.BuildTags)
	assert.Equal(t, config.Initialisms, g.Initialisms)
	assert.Equal(t, config.NoComments, g.NoComments)
	assert.Equal(t, config.Values, g.Values)

	// Test default values
	assert.Equal(t, "-", g.Version)
	assert.Equal(t, "-", g.Revision)
	assert.Equal(t, "-", g.BuildDate)
	assert.Equal(t, "-", g.BuiltBy)
	assert.NotNil(t, g.knownTemplates)
	assert.NotNil(t, g.t)
	assert.NotNil(t, g.fileSet)
	assert.NotNil(t, g.userTemplateNames)
}

// TestNewGeneratorConfig tests the NewGeneratorConfig constructor
func TestNewGeneratorConfig(t *testing.T) {
	config := NewGeneratorConfig()

	assert.NotNil(t, config)
	assert.False(t, config.NoPrefix)
	assert.NotNil(t, config.ReplacementNames)
	assert.Equal(t, "encoding/json", config.JSONPkg)
	assert.Empty(t, config.ReplacementNames)
}

// TestWithSQLInt tests the WithSQLInt option
func TestWithSQLInt(t *testing.T) {
	config := &GeneratorConfig{}
	option := WithSQLInt()

	assert.False(t, config.SQLInt)
	option(config)
	assert.True(t, config.SQLInt)
}

// TestWithValues tests the WithValues option
func TestWithValues(t *testing.T) {
	config := &GeneratorConfig{}
	option := WithValues()

	assert.False(t, config.Values)
	option(config)
	assert.True(t, config.Values)
}

// TestWithJsonPkg tests the WithJsonPkg option
func TestWithJsonPkg(t *testing.T) {
	config := &GeneratorConfig{}
	testPkg := "custom/json/package"
	option := WithJsonPkg(testPkg)

	assert.Empty(t, config.JSONPkg)
	option(config)
	assert.Equal(t, testPkg, config.JSONPkg)
}

// TestWithNoComments tests the WithNoComments option
func TestWithNoComments(t *testing.T) {
	config := &GeneratorConfig{}
	option := WithNoComments()

	assert.False(t, config.NoComments)
	option(config)
	assert.True(t, config.NoComments)
}

// TestWithBuildTags tests the WithBuildTags option
func TestWithBuildTags(t *testing.T) {
	config := &GeneratorConfig{}
	testTags := []string{"tag1", "tag2", "tag3"}
	option := WithBuildTags(testTags...)

	assert.Empty(t, config.BuildTags)
	option(config)
	assert.Equal(t, testTags, config.BuildTags)

	// Test appending more tags
	moreTags := []string{"tag4", "tag5"}
	option2 := WithBuildTags(moreTags...)
	option2(config)
	expectedTags := append(testTags, moreTags...)
	assert.Equal(t, expectedTags, config.BuildTags)
}

// TestAllOptionsIntegration tests using multiple options together
func TestAllOptionsIntegration(t *testing.T) {
	g := NewGenerator(
		WithSQLInt(),
		WithValues(),
		WithJsonPkg("custom/json"),
		WithNoComments(),
		WithBuildTags("integration", "test"),
		WithInitialisms("HTTP", "URL"),
	)

	assert.True(t, g.SQLInt)
	assert.True(t, g.Values)
	assert.Equal(t, "custom/json", g.JSONPkg)
	assert.True(t, g.NoComments)
	assert.Equal(t, []string{"integration", "test"}, g.BuildTags)
	assert.Equal(t, []string{"HTTP", "URL"}, g.Initialisms)
}

// TestGeneratorConfigWithTemplates tests NewGeneratorWithConfig with templates
func TestGeneratorConfigWithTemplates(t *testing.T) {
	config := GeneratorConfig{
		// Use empty template file names to avoid file not found errors
		TemplateFileNames: []string{},
	}

	g := NewGeneratorWithConfig(config)
	assert.NotNil(t, g)
	assert.Equal(t, config.TemplateFileNames, g.TemplateFileNames)

	// Test with non-empty but valid scenario - no actual templates needed for coverage
	config2 := GeneratorConfig{
		NoPrefix: true,
		Values:   true,
	}
	g2 := NewGeneratorWithConfig(config2)
	assert.NotNil(t, g2)
	assert.True(t, g2.NoPrefix)
	assert.True(t, g2.Values)
}

// TestNoParseOption tests the WithNoParse option
func TestNoParseOption(t *testing.T) {
	g := NewGenerator(WithNoParse())
	assert.True(t, g.NoParse, "NoParse should be true when WithNoParse is used")
}

// TestNoParseWithIntEnum tests that NoParse flag removes Parse method from int enums
func TestNoParseWithIntEnum(t *testing.T) {
	input := `package test

// ENUM(one, two, three)
type Number int
`
	g := NewGenerator(WithNoParse())
	f, err := parser.ParseFile(g.fileSet, "test.go", input, parser.ParseComments)
	require.NoError(t, err)

	output, err := g.Generate(f)
	require.NoError(t, err)
	require.NotNil(t, output)

	outputStr := string(output)

	// Should not contain Parse method
	assert.NotContains(t, outputStr, "func ParseNumber")
	assert.NotContains(t, outputStr, "func parseNumber")

	// Should still contain other methods
	assert.Contains(t, outputStr, "func (x Number) String()")
	assert.Contains(t, outputStr, "func (x Number) IsValid()")
}

// TestNoParseWithStringEnum tests that NoParse flag removes Parse method from string enums
func TestNoParseWithStringEnum(t *testing.T) {
	input := `package test

// ENUM(alpha, beta, gamma)
type Greek string
`
	g := NewGenerator(WithNoParse())
	f, err := parser.ParseFile(g.fileSet, "test.go", input, parser.ParseComments)
	require.NoError(t, err)

	output, err := g.Generate(f)
	require.NoError(t, err)
	require.NotNil(t, output)

	outputStr := string(output)

	// Should not contain Parse method
	assert.NotContains(t, outputStr, "func ParseGreek")
	assert.NotContains(t, outputStr, "func parseGreek")

	// Should still contain other methods
	assert.Contains(t, outputStr, "func (x Greek) String()")
	assert.Contains(t, outputStr, "func (x Greek) IsValid()")
}

// TestNoParseWithMarshalCreatesUnexported tests that NoParse with Marshal creates unexported parse
func TestNoParseWithMarshalCreatesUnexported(t *testing.T) {
	input := `package test

// ENUM(one, two, three)
type Number int
`
	g := NewGenerator(WithNoParse(), WithMarshal())
	f, err := parser.ParseFile(g.fileSet, "test.go", input, parser.ParseComments)
	require.NoError(t, err)

	output, err := g.Generate(f)
	require.NoError(t, err)
	require.NotNil(t, output)

	outputStr := string(output)

	// Should not contain public Parse method
	assert.NotContains(t, outputStr, "func ParseNumber(")

	// Should contain unexported parse method
	assert.Contains(t, outputStr, "func parseNumber(")

	// Should contain marshal methods that use the unexported parse
	assert.Contains(t, outputStr, "func (x *Number) UnmarshalText(")
	assert.Contains(t, outputStr, "parseNumber(name)")
}

// TestNoParseWithSQLCreatesUnexported tests that NoParse with SQL creates unexported parse
func TestNoParseWithSQLCreatesUnexported(t *testing.T) {
	input := `package test

// ENUM(one, two, three)
type Number int
`
	g := NewGenerator(WithNoParse(), WithSQLDriver())
	f, err := parser.ParseFile(g.fileSet, "test.go", input, parser.ParseComments)
	require.NoError(t, err)

	output, err := g.Generate(f)
	require.NoError(t, err)
	require.NotNil(t, output)

	outputStr := string(output)

	// Should not contain public Parse method
	assert.NotContains(t, outputStr, "func ParseNumber(")

	// Should contain unexported parse method
	assert.Contains(t, outputStr, "func parseNumber(")

	// Should contain Scan method that uses the unexported parse
	assert.Contains(t, outputStr, "func (x *Number) Scan(")
	// The Scan method should reference parseNumber somewhere
	assert.Regexp(t, `parseNumber\([^)]+\)`, outputStr)
}

// TestNoParseWithFlagCreatesUnexported tests that NoParse with Flag creates unexported parse
func TestNoParseWithFlagCreatesUnexported(t *testing.T) {
	input := `package test

// ENUM(one, two, three)
type Number int
`
	g := NewGenerator(WithNoParse(), WithFlag())
	f, err := parser.ParseFile(g.fileSet, "test.go", input, parser.ParseComments)
	require.NoError(t, err)

	output, err := g.Generate(f)
	require.NoError(t, err)
	require.NotNil(t, output)

	outputStr := string(output)

	// Should not contain public Parse method
	assert.NotContains(t, outputStr, "func ParseNumber(")

	// Should contain unexported parse method
	assert.Contains(t, outputStr, "func parseNumber(")

	// Should contain Set method that uses the unexported parse
	assert.Contains(t, outputStr, "func (x *Number) Set(")
	assert.Contains(t, outputStr, "parseNumber(val)")
}

// TestNoParseWithMultipleFeaturesCreatesUnexported tests NoParse with multiple features
func TestNoParseWithMultipleFeaturesCreatesUnexported(t *testing.T) {
	input := `package test

// ENUM(one, two, three)
type Number int
`
	g := NewGenerator(WithNoParse(), WithMarshal(), WithSQLDriver(), WithFlag())
	f, err := parser.ParseFile(g.fileSet, "test.go", input, parser.ParseComments)
	require.NoError(t, err)

	output, err := g.Generate(f)
	require.NoError(t, err)
	require.NotNil(t, output)

	outputStr := string(output)

	// Should not contain public Parse method
	assert.NotContains(t, outputStr, "func ParseNumber(")

	// Should contain exactly one unexported parse method
	assert.Contains(t, outputStr, "func parseNumber(")

	// Should contain all the feature methods
	assert.Contains(t, outputStr, "func (x *Number) UnmarshalText(")
	assert.Contains(t, outputStr, "func (x *Number) Scan(")
	assert.Contains(t, outputStr, "func (x *Number) Set(")
}

// TestNoParseDoesNotAffectOtherMethods tests that NoParse doesn't break other functionality
func TestNoParseDoesNotAffectOtherMethods(t *testing.T) {
	input := `package test

// ENUM(one, two, three)
type Number int
`
	g := NewGenerator(
		WithNoParse(),
		WithNames(),
		WithValues(),
		WithPtr(),
	)
	f, err := parser.ParseFile(g.fileSet, "test.go", input, parser.ParseComments)
	require.NoError(t, err)

	output, err := g.Generate(f)
	require.NoError(t, err)
	require.NotNil(t, output)

	outputStr := string(output)

	// Should not contain Parse method
	assert.NotContains(t, outputStr, "func ParseNumber")
	assert.NotContains(t, outputStr, "func parseNumber")

	// Should contain other methods
	assert.Contains(t, outputStr, "func NumberNames()")
	assert.Contains(t, outputStr, "func NumberValues()")
	assert.Contains(t, outputStr, "func (x Number) Ptr()")
	assert.Contains(t, outputStr, "func (x Number) String()")
	assert.Contains(t, outputStr, "func (x Number) IsValid()")
}

// TestNoParseWithStringEnumAndMarshal tests NoParse with string enum and marshal
func TestNoParseWithStringEnumAndMarshal(t *testing.T) {
	input := `package test

// ENUM(alpha, beta, gamma)
type Greek string
`
	g := NewGenerator(WithNoParse(), WithMarshal())
	f, err := parser.ParseFile(g.fileSet, "test.go", input, parser.ParseComments)
	require.NoError(t, err)

	output, err := g.Generate(f)
	require.NoError(t, err)
	require.NotNil(t, output)

	outputStr := string(output)

	// Should not contain public Parse method
	assert.NotContains(t, outputStr, "func ParseGreek(")

	// Should contain unexported parse method
	assert.Contains(t, outputStr, "func parseGreek(")

	// Should contain marshal methods
	assert.Contains(t, outputStr, "func (x *Greek) UnmarshalText(")
	assert.Contains(t, outputStr, "parseGreek(string(text))")
}

// TestNoParseOmitsErrorVariable tests that NoParse without dependent features omits ErrInvalidXXX
func TestNoParseOmitsErrorVariable(t *testing.T) {
	input := `package test

// ENUM(one, two, three)
type Number int
`
	g := NewGenerator(WithNoParse())
	f, err := parser.ParseFile(g.fileSet, "test.go", input, parser.ParseComments)
	require.NoError(t, err)

	output, err := g.Generate(f)
	require.NoError(t, err)
	require.NotNil(t, output)

	outputStr := string(output)

	// Should NOT contain error variable since Parse is not generated
	assert.NotContains(t, outputStr, "var ErrInvalidNumber")
	assert.NotContains(t, outputStr, "ErrInvalidNumber")
}

// TestNoParseWithMarshalIncludesErrorVariable tests that error is generated when needed
func TestNoParseWithMarshalIncludesErrorVariable(t *testing.T) {
	input := `package test

// ENUM(one, two, three)
type Number int
`
	g := NewGenerator(WithNoParse(), WithMarshal())
	f, err := parser.ParseFile(g.fileSet, "test.go", input, parser.ParseComments)
	require.NoError(t, err)

	output, err := g.Generate(f)
	require.NoError(t, err)
	require.NotNil(t, output)

	outputStr := string(output)

	// Should contain error variable since parseNumber uses it
	assert.Contains(t, outputStr, "var ErrInvalidNumber")
}

// TestNoParseWithStringEnumOmitsErrorVariable tests error omission with string enums
func TestNoParseWithStringEnumOmitsErrorVariable(t *testing.T) {
	input := `package test

// ENUM(alpha, beta, gamma)
type Greek string
`
	g := NewGenerator(WithNoParse())
	f, err := parser.ParseFile(g.fileSet, "test.go", input, parser.ParseComments)
	require.NoError(t, err)

	output, err := g.Generate(f)
	require.NoError(t, err)
	require.NotNil(t, output)

	outputStr := string(output)

	// Should NOT contain error variable
	assert.NotContains(t, outputStr, "var ErrInvalidGreek")
	assert.NotContains(t, outputStr, "ErrInvalidGreek")
}

// TestStringEnumWithSQLIntIncludesErrorVariable tests that sqlint generates error even with noparse
func TestStringEnumWithSQLIntIncludesErrorVariable(t *testing.T) {
	input := `package test

// ENUM(alpha, beta, gamma)
type Greek string
`
	g := NewGenerator(WithNoParse(), WithSQLInt())
	f, err := parser.ParseFile(g.fileSet, "test.go", input, parser.ParseComments)
	require.NoError(t, err)

	output, err := g.Generate(f)
	require.NoError(t, err)
	require.NotNil(t, output)

	outputStr := string(output)

	// Should contain error variable because lookupSqlInt and Value use it
	assert.Contains(t, outputStr, "var ErrInvalidGreek")
	assert.Contains(t, outputStr, "lookupSqlIntGreek")
}

func TestInitialismParsing(t *testing.T) {
	tests := map[string]struct {
		input  []string
		result []string
		err    string
	}{
		"no initialisms": {
			result: nil,
		},
		"single entry": {
			input:  []string{"HTTP"},
			result: []string{"HTTP"},
		},
		"comma separated": {
			input:  []string{"HTTP,URL,ID"},
			result: []string{"HTTP", "URL", "ID"},
		},
		"multiple flags": {
			input:  []string{"HTTP", "URL,ID", "API"},
			result: []string{"HTTP", "URL", "ID", "API"},
		},
		"deduplication": {
			input:  []string{"HTTP,HTTP,URL"},
			result: []string{"HTTP", "URL"},
		},
		"invalid lowercase": {
			input: []string{"Http"},
			err:   `invalid initialism "Http": must be all uppercase ASCII letters`,
		},
		"invalid number": {
			input: []string{"H2"},
			err:   `invalid initialism "H2": must be all uppercase ASCII letters`,
		},
		"empty entries ignored": {
			input:  []string{"HTTP,,URL"},
			result: []string{"HTTP", "URL"},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result, err := ParseInitialisms(tc.input)
			if tc.err != "" {
				require.Error(t, err)
				require.EqualError(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.result, result)
			}
		})
	}
}

func TestWithInitialisms(t *testing.T) {
	config := &GeneratorConfig{}
	option := WithInitialisms("HTTP", "URL")
	option(config)
	assert.Equal(t, []string{"HTTP", "URL"}, config.Initialisms)

	// Test appending
	option2 := WithInitialisms("API")
	option2(config)
	assert.Equal(t, []string{"HTTP", "URL", "API"}, config.Initialisms)
}

func TestInitialismsInGeneration(t *testing.T) {
	input := `package test

// ENUM(
//   get_http_url,
//   post_api_request,
//   fetch_html_id,
// )
type Method int
`
	g := NewGenerator(WithInitialisms("HTTP", "URL", "API", "ID", "HTML"))
	f, err := parser.ParseFile(g.fileSet, "test.go", input, parser.ParseComments)
	require.NoError(t, err)

	output, err := g.Generate(f)
	require.NoError(t, err)

	outputStr := string(output)
	// Verify initialisms are fully uppercased in const names
	assert.Contains(t, outputStr, "MethodGetHTTPURL")
	assert.Contains(t, outputStr, "MethodPostAPIRequest")
	assert.Contains(t, outputStr, "MethodFetchHTMLID")
	// Verify string values are NOT affected (stored in _MethodName concatenation)
	assert.Contains(t, outputStr, "get_http_urlpost_api_requestfetch_html_id")
}

func TestInitialismsKfeaturesStyle(t *testing.T) {
	input := `package test

// ENUM(
//   bpf_lsm,
//   btf,
//   bpf_tracing,
//   ima,
// )
type Feature int
`
	g := NewGenerator(WithInitialisms("BPF", "LSM", "BTF", "IMA"))
	f, err := parser.ParseFile(g.fileSet, "test.go", input, parser.ParseComments)
	require.NoError(t, err)

	output, err := g.Generate(f)
	require.NoError(t, err)

	outputStr := string(output)
	assert.Contains(t, outputStr, "FeatureBPFLSM")
	assert.Contains(t, outputStr, "FeatureBTF")
	assert.Contains(t, outputStr, "FeatureBPFTracing")
	assert.Contains(t, outputStr, "FeatureIMA")
}

func TestInitialismsWithLeaveSnakeCase(t *testing.T) {
	input := `package test

// ENUM(get_http_url)
type Method int
`
	g := NewGenerator(WithoutSnakeToCamel(), WithInitialisms("HTTP", "URL"))
	f, err := parser.ParseFile(g.fileSet, "test.go", input, parser.ParseComments)
	require.NoError(t, err)

	output, err := g.Generate(f)
	require.NoError(t, err)

	outputStr := string(output)
	// With nocamel, snakeToCamelCase is skipped. cases.Title only uppercases
	// the first rune of the entire rawName, so underscore-separated parts after
	// the first remain lowercase. applyInitialisms finds no title-cased matches,
	// so initialisms have no effect in nocamel mode.
	assert.Contains(t, outputStr, "MethodGet_http_url")
}

func TestInitialismsWithNoPrefix(t *testing.T) {
	input := `package test

// ENUM(http_url)
type Method int
`
	g := NewGenerator(WithNoPrefix(), WithInitialisms("HTTP", "URL"))
	f, err := parser.ParseFile(g.fileSet, "test.go", input, parser.ParseComments)
	require.NoError(t, err)

	output, err := g.Generate(f)
	require.NoError(t, err)

	outputStr := string(output)
	assert.Contains(t, outputStr, "HTTPURL")
}

func TestInitialismsWithStringEnum(t *testing.T) {
	input := `package test

// ENUM(http_api, rest_url)
type Endpoint string
`
	g := NewGenerator(WithInitialisms("HTTP", "API", "URL", "REST"))
	f, err := parser.ParseFile(g.fileSet, "test.go", input, parser.ParseComments)
	require.NoError(t, err)

	output, err := g.Generate(f)
	require.NoError(t, err)

	outputStr := string(output)
	assert.Contains(t, outputStr, "EndpointHTTPAPI")
	assert.Contains(t, outputStr, "EndpointRESTURL")
	// String values unchanged
	assert.Contains(t, outputStr, `"http_api"`)
	assert.Contains(t, outputStr, `"rest_url"`)
}

func TestInitialismOrdering(t *testing.T) {
	input := `package test

// ENUM(ide, id_value)
type Thing int
`
	// ID and IDE overlap should resolve by exact token match.
	g := NewGenerator(WithInitialisms("ID", "IDE"))
	f, err := parser.ParseFile(g.fileSet, "test.go", input, parser.ParseComments)
	require.NoError(t, err)

	output, err := g.Generate(f)
	require.NoError(t, err)

	outputStr := string(output)
	assert.Contains(t, outputStr, "ThingIDE")
	assert.Contains(t, outputStr, "ThingIDValue")
}

func TestInitialismsDoNotReplaceSubstringsInsideTokens(t *testing.T) {
	input := `package test

// ENUM(apiary, ideology, id_value, api_id)
type Thing int
`
	g := NewGenerator(WithInitialisms("API", "IDE", "ID"))
	f, err := parser.ParseFile(g.fileSet, "test.go", input, parser.ParseComments)
	require.NoError(t, err)

	output, err := g.Generate(f)
	require.NoError(t, err)

	outputStr := string(output)
	assert.Contains(t, outputStr, "ThingApiary")
	assert.Contains(t, outputStr, "ThingIdeology")
	assert.Contains(t, outputStr, "ThingIDValue")
	assert.Contains(t, outputStr, "ThingAPIID")
	assert.NotContains(t, outputStr, "ThingAPIary")
	assert.NotContains(t, outputStr, "ThingIDEology")
}

func TestShouldSplitToken(t *testing.T) {
	tests := map[string]struct {
		value    string
		index    int
		expected bool
	}{
		"current rune underscore": {
			value:    "A_B",
			index:    1,
			expected: true,
		},
		"previous rune underscore": {
			value:    "A_B",
			index:    2,
			expected: true,
		},
		"digit to letter": {
			value:    "2A",
			index:    1,
			expected: true,
		},
		"letter to digit": {
			value:    "A2",
			index:    1,
			expected: true,
		},
		"lower to upper": {
			value:    "aB",
			index:    1,
			expected: true,
		},
		"upper run before trailing lower": {
			value:    "HTTPServer",
			index:    4, // P|S where next is lower e
			expected: true,
		},
		"upper run at end": {
			value:    "HTTP",
			index:    3,
			expected: false,
		},
		"upper to lower": {
			value:    "Ab",
			index:    1,
			expected: false,
		},
		"digit to digit": {
			value:    "22",
			index:    1,
			expected: false,
		},
		"lower to lower": {
			value:    "ab",
			index:    1,
			expected: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			runes := []rune(tc.value)
			require.GreaterOrEqual(t, tc.index, 1)
			require.Less(t, tc.index, len(runes))
			assert.Equal(t, tc.expected, shouldSplitToken(runes, tc.index))
		})
	}
}

func TestSplitIdentifierTokens(t *testing.T) {
	tests := map[string]struct {
		value    string
		expected []string
	}{
		"empty string": {
			value:    "",
			expected: nil,
		},
		"underscore boundaries": {
			value:    "API_ID",
			expected: []string{"API", "_", "ID"},
		},
		"digit boundaries": {
			value:    "V2API3ID",
			expected: []string{"V", "2", "API", "3", "ID"},
		},
		"camel boundaries": {
			value:    "MyHTTPServer",
			expected: []string{"My", "HTTP", "Server"},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expected, splitIdentifierTokens(tc.value))
		})
	}
}
