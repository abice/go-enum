//go:generate go-enum -f=animal.go

package example

// Animal x ENUM(
// Cat,
// Dog,
// Fish
// )
type Animal int32
