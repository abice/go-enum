//go:generate ../bin/go-enum  --prefix=AcmeInc -b example

package example

// Products of AcmeInc ENUM(
// Anvil,
// Dynamite,
// Glue
// )
type Product int32
