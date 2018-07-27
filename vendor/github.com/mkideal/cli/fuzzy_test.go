package cli

import (
	"testing"
)

func TestEditDistance(t *testing.T) {
	for _, arg := range []struct {
		s, t string
		dist float32
	}{
		{"", "", 0},
		{"a", "a", 0},
		{"", "a", 1},
		{"", "abc", 3},
		{"a", "b", 1},
		{"aa", "b", 2},
		{"aa", "bb", 2},
		{"abc", "ac", 1},
		{"abc", "acc", 1},
		{"cli", "clli", 1},
		{"publish", "pub", 4},
		{"publish", "pbish", 2},
	} {
		dist := editDistance([]byte(arg.s), []byte(arg.t))
		if dist != arg.dist {
			t.Errorf("dist of between %s and %s: want %f, got %f", arg.s, arg.t, arg.dist, dist)
		}
	}
}
