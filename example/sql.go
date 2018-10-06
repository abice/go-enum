//go:generate go-enum -f=$GOFILE --sql

package example

// ENUM(pending, inWork, completed, rejected)
type ProjectStatus int
