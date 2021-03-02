//go:generate ../bin/go-enum -f=$GOFILE --sql --sqlnullstr

package example

// ENUM(pending, processing, completed, failed)
type JobState int
