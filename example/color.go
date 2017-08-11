//go:generate go-enum -f=color.go --marshal --lower

package example

// Color is an enumeration of colors that are allowed.
// ENUM(
// Black, White, Red
// Green
// Blue
// grey
// yellow
// )
type Color int
