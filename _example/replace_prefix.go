//go:generate ../bin/go-enum -f=$GOFILE --marshal --prefix=AcmeInc_ --noprefix --nocamel --names

package example

// Shops ENUM(
// SOME_PLACE_AWESOME,
// SomewhereElse,
// LocationUnknown
// )
type Shop string
