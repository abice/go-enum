//go:generate ../bin/go-enum -f=$GOFILE --ptr --marshal --flag --nocase --mustparse --sqlnullint --sqluint --names

package example

// ENUM(mp4=1, mp3, ogg, flac)
type MediaType string
