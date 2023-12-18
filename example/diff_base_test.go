//go:build example
// +build example

package example

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDiffBase(t *testing.T) {
	tests := map[string]struct {
		actual   int
		expected DiffBase
	}{
		"DiffBaseB3": {
			actual:   3,
			expected: DiffBaseB3,
		},
		"DiffBaseB4": {
			actual:   4,
			expected: DiffBaseB4,
		},
		"DiffBaseB5": {
			actual:   5,
			expected: DiffBaseB5,
		},
		"DiffBaseB6": {
			actual:   6,
			expected: DiffBaseB6,
		},
		"DiffBaseB7": {
			actual:   7,
			expected: DiffBaseB7,
		},
		"DiffBaseB8": {
			actual:   8,
			expected: DiffBaseB8,
		},
		"DiffBaseB9": {
			actual:   9,
			expected: DiffBaseB9,
		},
		"DiffBaseB10": {
			actual:   11,
			expected: DiffBaseB10,
		},
		"DiffBaseB11": {
			actual:   43,
			expected: DiffBaseB11,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, int(tc.expected), tc.actual)
		})
	}
}
