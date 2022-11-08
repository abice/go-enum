//go:generate ../bin/go-enum  --sqlnullint

package example

// ENUM(jpeg, jpg, png, tiff, gif)
type ImageType int
