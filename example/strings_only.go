//go:generate ../bin/go-enum  --ptr --marshal --flag --nocase --mustparse --sqlnullstr --sql --names --values --nocomments -b example

package example

// ENUM(pending, running, completed, failed=error)
type StrState string
