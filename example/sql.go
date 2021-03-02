//go:generate ../bin/go-enum -f=$GOFILE --sql --sqlnullstr --sqlnullint -ptr

package example

// ENUM(pending, inWork, completed, rejected)
type ProjectStatus int
