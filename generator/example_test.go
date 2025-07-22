package generator

// X is doc'ed
type X struct{}

// Color is an enumeration of colors that are allowed.
// ENUM(
// Black, White, Red
// Green
// Blue=33
// grey=
// yellow
// ).
type Color int

// Animal x ENUM(
// Cat,
// Dog,
// Fish
// ) Some other line of info
type Animal int32

// Model x ENUM(Toyota,_,Chevy,_,Ford).
type Model int32

/*
	ENUM(
	Coke
	Pepsi
	MtnDew

).
*/
type Soda int64

/*
	ENUM(
	test_lower
	Test_capital
	anotherLowerCaseStart

)
*/
type Cases int64

/*
	ENUM(
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

/*
	ENUM(
	startWithNum=23
	nextNum

)
*/
type StartNotZero int64

// ENUM(
// Black, White, Red
// Green
// Blue=33 // Blue starts with 33.
// grey=
// yellow
// )
type ColorWithComment int

/*
ENUM(
Black, White, Red
Green
Blue=33 // Blue starts with 33
grey=
yellow
)
*/
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
// red-orange-blue
// )
type ColorWithComment3 int

/* ENUM(
	_, // Placeholder
Black, White, Red
Green = 33 // Green starts with 33
*/
// Blue
// grey=
// yellow // Where did all the (somewhat) bad fish go? (something else that goes in parentheses at the end of the line)
// blue-green // blue-green comment
// red-orange // has a , in it!?!
// )
type ColorWithComment4 int

/*
	ENUM(

Unknown= 0
E2P15					= 32768
E2P16					= 65536
E2P17					= 131072
E2P18					= 262144
E2P19					= 524288
E2P20					= 1048576
E2P21					= 2097152
E2P22					= 33554432
E2P23					= 67108864
E2P28					= 536870912
E2P30					= 1073741824
E2P31					= 2147483648
E2P32					= 4294967296
E2P33					= 8454967296
)
*/
type Enum64bit uint64

// NonASCII
// ENUM(
// Продам = 1114
// 車庫 = 300
// էժան = 1
// )
type NonASCII int

// StringEnum.
// ENUM(
// random = 1114
// values = 300
// here  = 1
// )
type StringEnum string
