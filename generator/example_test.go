package generator

// X is doc'ed
type X struct {
}

// Color is an enumeration of colors that are allowed.
// ENUM(
// Black, White, Red
// Green
// Blue=33
// grey=
// yellow
// )
type Color int

// Animal x ENUM(
// Cat,
// Dog,
// Fish
// )
type Animal int32

// Model x ENUM(Toyota,_,Chevy,_,Ford)
type Model int32

/* ENUM(
 Coke
 Pepsi
 MtnDew
)
*/
type Soda int64
