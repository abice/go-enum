//go:build example
// +build example

package example

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatusString(t *testing.T) {
	tests := map[string]struct {
		input  string
		output Status
	}{
		"bad": {
			input:  `Bad`,
			output: StatusBad,
		},
		"unknown": {
			input:  `Unknown`,
			output: StatusUnknown,
		},
		"good": {
			input:  `Good`,
			output: StatusGood,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			output, err := ParseStatus(tc.input)
			assert.NoError(t, err)
			assert.Equal(t, tc.output, output)

			assert.Equal(t, tc.input, output.String())
		})
	}

	t.Run("failures", func(t *testing.T) {
		assert.Equal(t, "Status(99)", Status(99).String())
		failedStatus, err := ParseStatus("")
		assert.Error(t, err)

		assert.Equal(t, Status(0), failedStatus)
		t.Run("cased", func(t *testing.T) {
			actual, err := ParseStatus("BAD")
			assert.NoError(t, err)
			assert.Equal(t, StatusBad, actual)
		})
	})
}

func TestNegativeString(t *testing.T) {
	tests := map[string]struct {
		input  string
		output AllNegative
	}{
		"unknown": {
			input:  `Unknown`,
			output: AllNegativeUnknown,
		},
		"good": {
			input:  `Good`,
			output: AllNegativeGood,
		},
		"bad": {
			input:  `Bad`,
			output: AllNegativeBad,
		},
		"ugly": {
			input:  `Ugly`,
			output: AllNegativeUgly,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			output, err := ParseAllNegative(tc.input)
			assert.NoError(t, err)
			assert.Equal(t, tc.output, output)

			assert.Equal(t, tc.input, output.String())
		})
	}

	t.Run("failures", func(t *testing.T) {
		assert.Equal(t, "AllNegative(99)", AllNegative(99).String())
		allN, err := ParseAllNegative("")
		assert.Error(t, err)

		assert.Equal(t, AllNegative(0), allN)
	})
	t.Run("cased", func(t *testing.T) {
		actual, err := ParseAllNegative("UGLY")
		assert.NoError(t, err)
		assert.Equal(t, AllNegativeUgly, actual)
	})
}
