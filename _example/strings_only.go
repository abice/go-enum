//go:generate ../bin/go-enum -f=$GOFILE --ptr --marshal --flag --nocase --mustparse --sqlnullstr --sql --names

package example

// ENUM(pending, running, completed, failed)
type StrState string
