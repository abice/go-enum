package generator

// X is doc'ed
type X struct {
}

// Color is an enumeration of colors that are allowed.
// ENUM(
// Black, White, Red
// Green
// Blue
// grey
// yellow
// )
type Color int

// Animal x ENUM(
// Cat,
// Dog,
// Fish
// )
type Animal int32

// Model x ENUM(Toyota,Chevy,Ford)
type Model int32
