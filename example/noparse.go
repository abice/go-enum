//go:build example
// +build example

//go:generate ../bin/go-enum --noparse -b example

package example

// ENUM(
// A,
// B,
// C
// D
// E
// ).
type UnparsedValues uint8

// ENUM(
// A,
// B,
// C
// D
// E
// ).
type UnparsedString string
