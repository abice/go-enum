//go:generate ../bin/go-enum -f=$GOFILE --ptr --marshal --flag --nocase --mustparse --sqlnullint --names

package example

// ENUM(mp4=1, mp3, ogg, flac)
type MediaTypeInt string
