//go:generate ../bin/go-enum --flag --marshal --names --ptr --sql --sqlnullstr --values

package example

// ENUM(
// First,
// Second,
// Third,
// )
type privateInt int

// ENUM(
// First="a",
// Second="b",
// Third="c",
// )
type privateStr string
