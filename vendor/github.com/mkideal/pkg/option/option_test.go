package option_test

import (
	"fmt"
	"testing"

	"github.com/mkideal/pkg/option"
	"github.com/stretchr/testify/assert"
)

func TestOption(t *testing.T) {
	assert.Equal(t, true, option.Bool(true))
	assert.Equal(t, false, option.Bool(false))
	assert.Equal(t, true, option.Bool(true, true))
	assert.Equal(t, false, option.Bool(true, false))
	assert.Equal(t, true, option.Bool(false, true))
	assert.Equal(t, false, option.Bool(false, false))
}

func ExampleOption() {
	fn := func(name string, ageOpt ...int) {
		age := option.Int(20, ageOpt...)
		fmt.Printf("name: %s, age: %d\n", name, age)
	}
	fn("A")
	fn("B", 10)
	fn("C", 15, 25) // `25` ignored
	// Output:
	// name: A, age: 20
	// name: B, age: 10
	// name: C, age: 15
}
