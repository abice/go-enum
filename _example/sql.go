//go:generate ../bin/go-enum  --sql --sqlnullstr --sqlnullint --ptr --marshal

package example

// ENUM(pending, inWork, completed, rejected)
type ProjectStatus int
