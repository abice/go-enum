//go:build example
// +build example

package example

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnum32Bit(t *testing.T) {
	tests := map[string]struct {
		input  string
		output Enum32bit
	}{
		"E2P15": {
			input:  `E2P15`,
			output: Enum32bitE2P15,
		},
		"E2P30": {
			input:  `E2P30`,
			output: Enum32bitE2P30,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			output, err := ParseEnum32bit(tc.input)
			assert.NoError(t, err)
			assert.Equal(t, tc.output, output)

			assert.Equal(t, tc.input, output.String())
		})
	}

	t.Run("basics", func(t *testing.T) {
		assert.Equal(t, "E2P23", Enum32bitE2P23.String())
		assert.Equal(t, "Enum32bit(99)", Enum32bit(99).String())
		_, err := ParseEnum32bit("-1")
		assert.Error(t, err)

		names := Enum32bitNames()
		assert.Len(t, names, 12)
	})
}
