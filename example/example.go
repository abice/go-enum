//go:generate go-enum -f=example.go --marshal --lower

package example

// X is doc'ed
type X struct {
}

// Animal x ENUM(
// Cat,
// Dog,
// Fish
// )
type Animal int32

// Model x ENUM(Toyota,Chevy,_,Ford)
type Model int32
