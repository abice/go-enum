//go:generate ../bin/go-enum  --sql --sqlnullstr -b example

package example

// ENUM(pending, processing, completed, failed)
type JobState int
