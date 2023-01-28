//go:generate ../bin/go-enum  --ptr --marshal --flag --nocase --mustparse --sqlnullstr --sql --names --values --nocomments

package example

// ENUM(pending, running, completed, failed)
type StrState string
