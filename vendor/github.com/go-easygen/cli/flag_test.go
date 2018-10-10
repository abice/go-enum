package cli

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseExpression(t *testing.T) {
	os.Setenv("CLI_EXPR_TEST1", "1")
	for i, tt := range []struct {
		s        string
		isNumber bool
		expr     string
		hasErr   bool
	}{
		{"$", false, "", true},
		{"$$", false, "$", false},
		{"$$$", false, "", true},
		{"$$$$", false, "$$", false},
		{"$HOME", false, os.Getenv("HOME"), false},
		{"$HOME$", false, "", true},
		{"$HOME-abc", false, os.Getenv("HOME") + "-abc", false},
		{"abc$HOME", false, "abc" + os.Getenv("HOME"), false},
		{"2+$CLI_EXPR_TEST1", true, "2+1", false},
		{"2+$CLI_EXPR_UNKNOWN", true, "2+0", false},
	} {
		got, err := parseExpression(tt.s, tt.isNumber)
		if tt.hasErr {
			if err == nil {
				t.Errorf("%dth: error wanted, but not", i)
			}
			continue
		}
		assert.Equal(t, tt.expr, got)
	}
}
