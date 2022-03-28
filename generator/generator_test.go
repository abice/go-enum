//go:build !go1.18
// +build !go1.18

package generator

import (
	"fmt"
	"go/parser"
	"strings"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
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

// TestStringEnum
func TestStringEnum(t *testing.T) {
	input := `package test
	/* ENUM(
		Cherry
		Apple
		Grape
	 )
	 */
	 type StringValues string
	`
	g := NewGenerator().
		WithoutSnakeToCamel()
	f, err := parser.ParseFile(g.fileSet, "TestRequiredErrors", input, parser.ParseComments)
	assert.Nil(t, err, "Error parsing no struct input")

	output, err := g.Generate(f)
	assert.Nil(t, err, "Error generating formatted code")
	if true { // Debugging statement
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
			defer func() {
				replacementNames = map[string]string{}
			}()
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
