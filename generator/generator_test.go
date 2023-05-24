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
	testExample = `example_test.go`
)

// TestNoStructInputFile
func TestNoStructFile(t *testing.T) {
	input := `package test
	// Behavior
	type SomeInterface interface{

	}
	`
	g := NewGenerator().
		WithoutSnakeToCamel()
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
	g := NewGenerator().
		WithoutSnakeToCamel()
	// Parse the file given in arguments
	_, err := g.GenerateFromFile("")
	assert.NotNil(t, err, "Error generating formatted code")
}

// TestExampleFile
func TestExampleFile(t *testing.T) {
	g := NewGenerator().
		WithMarshal().
		WithSQLDriver().
		WithCaseInsensitiveParse().
		WithNames().
		WithoutSnakeToCamel()
	// Parse the file given in arguments
	imported, err := g.GenerateFromFile(testExample)
	require.Nil(t, err, "Error generating formatted code")

	outputLines := strings.Split(string(imported), "\n")
	cupaloy.SnapshotT(t, outputLines)

	if false {
		fmt.Println(string(imported))
	}
}

// TestExampleFileMoreOptions
func TestExampleFileMoreOptions(t *testing.T) {
	g := NewGenerator().
		WithMarshal().
		WithSQLDriver().
		WithCaseInsensitiveParse().
		WithNames().
		WithoutSnakeToCamel().
		WithMustParse().
		WithForceLower().
		WithTemplates(`../example/user_template.tmpl`)
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
	g := NewGenerator().
		WithMarshal().
		WithLowercaseVariant().
		WithNoPrefix().
		WithFlag().
		WithoutSnakeToCamel()
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
func TestReplacePrefixExampleFile(t *testing.T) {
	g := NewGenerator().
		WithMarshal().
		WithLowercaseVariant().
		WithNoPrefix().
		WithPrefix("MyPrefix_").
		WithFlag().
		WithoutSnakeToCamel()
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
	g := NewGenerator().
		WithMarshal().
		WithLowercaseVariant().
		WithNoPrefix().
		WithFlag()

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
	g := NewGenerator().
		WithMarshal().
		WithLowercaseVariant().
		WithNoPrefix().
		WithFlag().
		WithoutSnakeToCamel().
		WithPtr().
		WithSQLNullInt().
		WithSQLNullStr().
		WithPrefix("Custom_prefix_")
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
			replacementNames = map[string]string{}
			err := ParseAliases(tc.input)
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
	g := NewGenerator().
		WithoutSnakeToCamel()
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
	g := NewGenerator().
		WithoutSnakeToCamel()
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
	g := NewGenerator().
		WithoutSnakeToCamel()
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
	g := NewGenerator().
		WithoutSnakeToCamel()
	_ = ParseAliases([]string{"CDEF:C"})
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
