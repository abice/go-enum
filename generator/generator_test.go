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
		WithLowercaseVariant().
		WithNames().
		WithoutSnakeToCamel()
	// Parse the file given in arguments
	imported, err := g.GenerateFromFile(testExample)
	require.Nil(t, err, "Error generating formatted code")

	outputLines := strings.Split(string(imported), "\n")
	err = cupaloy.Snapshot(outputLines)
	assert.NoError(t, err, "Output must match snapshot")

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
	err = cupaloy.Snapshot(outputLines)
	assert.NoError(t, err, "Output must match snapshot")

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
	err = cupaloy.Snapshot(outputLines)
	assert.NoError(t, err, "Output must match snapshot")

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
		WithPrefix("Custom_prefix_")
	// Parse the file given in arguments
	imported, err := g.GenerateFromFile(testExample)
	require.Nil(t, err, "Error generating formatted code")

	outputLines := strings.Split(string(imported), "\n")
	err = cupaloy.Snapshot(outputLines)
	assert.NoError(t, err, "Output must match snapshot")

	if false {
		fmt.Println(string(imported))
	}
}
