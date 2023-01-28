//go:generate ../bin/go-enum  --sql --sqlnullstr --sqlnullint --ptr --marshal --nocomments

package example

// ENUM(pending, inWork, completed, rejected)
type ProjectStatus int
