//go:generate ../bin/go-enum  --sql --sqlnullstr --sqlnullint --ptr --marshal --nocomments -b example

package example

// ENUM(pending, inWork, completed, rejected)
type ProjectStatus int
