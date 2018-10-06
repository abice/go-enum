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

/* ENUM(
 test_lower
 Test_capital
 anotherLowerCaseStart
)
*/
type Cases int64

/* ENUM(
 test-Hyphen
 -hyphenStart
 _underscoreFirst
 0numberFirst
 123456789a
 123123-asdf
 ending-hyphen-
)
*/
type Sanitizing int64

/* ENUM(
 startWithNum=23
 nextNum
)
*/
type StartNotZero int64

// ENUM(
// Black, White, Red
// Green
// Blue=33 // Blue starts with 33
// grey=
// yellow
// )
type ColorWithComment int

/*ENUM(
Black, White, Red
Green
Blue=33 // Blue starts with 33
grey=
yellow
)*/
type ColorWithComment2 int

/* ENUM(
Black, White, Red
Green = 33 // Green starts with 33
*/
// Blue
// grey=
// yellow
// blue-green // blue-green comment
// red-orange
// )
type ColorWithComment3 int
