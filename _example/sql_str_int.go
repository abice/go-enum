//go:generate ../bin/go-enum -f=$GOFILE --sqlint --sqlnullint --names

package example

// ENUM(_,zeus, apollo, athena=20, ares)
type GreekGod string

// ENUM(_,zeus, apollo, _=19, athena="20", ares)
type GreekGodCustom string
