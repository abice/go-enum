//go:generate ../bin/go-enum -f=$GOFILE --marshal --prefix=AcmeInt_ --noprefix --nocamel --names

package example

// Shops ENUM(
// SOME_PLACE_AWESOME,
// SomewhereElse,
// LocationUnknown
// )
type IntShop int
