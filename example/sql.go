//go:generate go-enum -f=sql.go --sql

package example

// ENUM(pending, inWork, completed, rejected)
type ProjectStatus int
