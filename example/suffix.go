//go:generate ../bin/go-enum -b example --output-suffix .enum.gen

//go:build example
// +build example

package example

// Suffix ENUM(gen)
type Suffix string
