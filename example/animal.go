//go:generate ../bin/go-enum -f=$GOFILE

package example

// Animal x ENUM(
// Cat,
// Dog,
// Fish
// )
type Animal int32
