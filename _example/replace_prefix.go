//go:generate ../bin/go-enum  --marshal --prefix=AcmeInc_ --noprefix --nocamel --names

package example

// Shops ENUM(
// SOME_PLACE_AWESOME,
// SomewhereElse,
// LocationUnknown
// )
type Shop string
