//go:generate ../bin/go-enum -a "+:Plus,#:Sharp"

package example

// Animal x ENUM(
// Cat,
// Dog,
// Fish
// Fish++
// Fish#
// ).
type Animal int32
