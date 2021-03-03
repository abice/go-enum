//go:generate ../bin/go-enum -f=$GOFILE --sql --sqlnullstr --sqlnullint --ptr --marshal

package example

// ENUM(pending, inWork, completed, rejected)
type ProjectStatus int
