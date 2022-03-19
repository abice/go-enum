//go:build go1.18
// +build go1.18

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

var (
	testExampleFiles = map[string]string{
		"og":   `example_test.go`,
		"1.18": `example_1.18_test.go`,
	}
)

// TestNoStructInputFile
func Test118NoStructFile(t *testing.T) {
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
func Test118NoFile(t *testing.T) {
	g := NewGenerator().
		WithoutSnakeToCamel()
	// Parse the file given in arguments
	_, err := g.GenerateFromFile("")
	assert.NotNil(t, err, "Error generating formatted code")
}

// TestExampleFile
func Test118ExampleFile(t *testing.T) {
	g := NewGenerator().
		WithMarshal().
		WithSQLDriver().
		WithCaseInsensitiveParse().
		WithNames().
		WithoutSnakeToCamel()

	for name, testExample := range testExampleFiles {
		t.Run(name, func(t *testing.T) {
			// Parse the file given in arguments
			imported, err := g.GenerateFromFile(testExample)
			require.Nil(t, err, "Error generating formatted code")

			outputLines := strings.Split(string(imported), "\n")
			cupaloy.SnapshotT(t, outputLines)

			if false {
				fmt.Println(string(imported))
			}
		})
	}
}

// TestExampleFile
func Test118NoPrefixExampleFile(t *testing.T) {
	g := NewGenerator().
		WithMarshal().
		WithLowercaseVariant().
		WithNoPrefix().
		WithFlag().
		WithoutSnakeToCamel()
	for name, testExample := range testExampleFiles {
		t.Run(name, func(t *testing.T) {
			// Parse the file given in arguments
			imported, err := g.GenerateFromFile(testExample)
			require.Nil(t, err, "Error generating formatted code")

			outputLines := strings.Split(string(imported), "\n")
			cupaloy.SnapshotT(t, outputLines)

			if false {
				fmt.Println(string(imported))
			}
		})
	}
}

// TestExampleFile
func Test118NoPrefixExampleFileWithSnakeToCamel(t *testing.T) {
	g := NewGenerator().
		WithMarshal().
		WithLowercaseVariant().
		WithNoPrefix().
		WithFlag()

	for name, testExample := range testExampleFiles {
		t.Run(name, func(t *testing.T) {
			// Parse the file given in arguments
			imported, err := g.GenerateFromFile(testExample)
			require.Nil(t, err, "Error generating formatted code")

			outputLines := strings.Split(string(imported), "\n")
			cupaloy.SnapshotT(t, outputLines)

			if false {
				fmt.Println(string(imported))
			}
		})
	}
}

// TestCustomPrefixExampleFile
func Test118CustomPrefixExampleFile(t *testing.T) {
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
	for name, testExample := range testExampleFiles {
		t.Run(name, func(t *testing.T) {
			// Parse the file given in arguments
			imported, err := g.GenerateFromFile(testExample)
			require.Nil(t, err, "Error generating formatted code")

			outputLines := strings.Split(string(imported), "\n")
			cupaloy.SnapshotT(t, outputLines)

			if false {
				fmt.Println(string(imported))
			}
		})
	}
}

func Test118AliasParsing(t *testing.T) {
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
