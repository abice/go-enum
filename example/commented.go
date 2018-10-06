//go:generate go-enum -f=$GOFILE --marshal --lower

package example

// Commented is an enumeration of commented values
/*
ENUM(
value1 // Commented value 1
value2
value3 // Commented value 3
)
*/
type Commented int
