//go:build example
// +build example

//go:generate ../bin/go-enum -a "+:Plus,#:Sharp" -b example

package example

// Animal x ENUM(
// Cat,
// Dog,
// Fish
// Fish++
// Fish#
// ).
type Animal int32
