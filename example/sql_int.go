//go:generate ../bin/go-enum  --sqlnullint -b example

package example

// ENUM(jpeg, jpg, png, tiff, gif)
type ImageType int
