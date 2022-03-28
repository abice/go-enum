//go:generate ../bin/go-enum -f=$GOFILE --sql --sqlnullstr --nocase

package example

// ENUM(pending, processing, completed, failed)
type JobState int

// ENUM(pending, processing, completed, failed)
type JobStateStr string
