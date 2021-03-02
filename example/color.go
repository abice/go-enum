//go:generate ../bin/go-enum -f=$GOFILE --marshal --lower --ptr

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
// yellow_green
// red-orange-blue
// )
type Color int
