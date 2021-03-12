//go:generate ../bin/go-enum -f=$GOFILE --sqlnullint

package example

// ENUM(jpeg, jpg, png, tiff, gif)
type ImageType int
