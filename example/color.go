//go:generate go-enum -f=$GOFILE --marshal --lower

package example

// Color is an enumeration of colors that are allowed.
/* ENUM(
Black, White, Red
Green = 33 // Green starts with 33
*/
// Blue
// grey=
// yellow
// blue-green
// red-orange
// )
type Color int
