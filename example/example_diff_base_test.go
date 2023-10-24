//go:build example
// +build example

package example

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDiffBase(t *testing.T) {

	tests := []struct {
		name     string
		actual   int
		expected DiffBase
	}{
		{
			name:     "DiffBaseB3",
			actual:   3,
			expected: DiffBaseB3,
		},
		{
			name:     "DiffBaseB4",
			actual:   4,
			expected: DiffBaseB4,
		},
		{
			name:     "DiffBaseB5",
			actual:   5,
			expected: DiffBaseB5,
		}, {
			name:     "DiffBaseB6",
			actual:   6,
			expected: DiffBaseB6,
		}, {
			name:     "DiffBaseB7",
			actual:   7,
			expected: DiffBaseB7,
		}, {
			name:     "DiffBaseB8",
			actual:   8,
			expected: DiffBaseB8,
		}, {
			name:     "DiffBaseB9",
			actual:   9,
			expected: DiffBaseB9,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			assert.Equal(tt, int(test.expected), test.actual)
		})
	}
}
