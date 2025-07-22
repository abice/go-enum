//go:build example
// +build example

//go:generate ../bin/go-enum  -b example --no-iota
package example

// ENUM(
// A = 0
// B = 2
// C = 1
// )
type Buggy uint
