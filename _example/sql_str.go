//go:generate ../bin/go-enum  --sql --sqlnullstr

package example

// ENUM(pending, processing, completed, failed)
type JobState int
