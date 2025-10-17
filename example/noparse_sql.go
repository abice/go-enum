//go:build example
// +build example

//go:generate ../bin/go-enum --noparse -b example --sql

package example

// ENUM(
// A,
// B,
// C
// D
// E
// ).
type UnparsedSqlValues uint8

// ENUM(
// A,
// B,
// C
// D
// E
// ).
type UnparsedSqlString string
