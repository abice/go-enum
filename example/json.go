//go:generate go-enum -f=json.go --json

package example

// Planet x ENUM(
// Mercury,
// Venus,
// Earth,
// Mars,
// Jupiter,
// Saturn,
// Uranus,
// Neptune
// )
type Planet int32
