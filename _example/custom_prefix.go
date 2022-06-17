//go:generate ../bin/go-enum -f=$GOFILE --prefix=AcmeInc

package example

// Products of AcmeInc ENUM(
// Anvil,
// Dynamite,
// Glue
// )
type Product int32
