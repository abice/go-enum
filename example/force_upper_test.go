//go:build example
// +build example

package example

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestForceUpperString(t *testing.T) {
	tests := map[string]struct {
		input  string
		output ForceUpperType
	}{
		"dataswap": {
			input:  `DATASWAP`,
			output: ForceUpperTypeDataSwap,
		},
		"bootnode": {
			input:  `BOOTNODE`,
			output: ForceUpperTypeBootNode,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			output, err := ParseForceUpperType(tc.input)
			assert.NoError(t, err)
			assert.Equal(t, tc.output, output)

			assert.Equal(t, tc.input, output.String())
		})
	}

	t.Run("failures", func(t *testing.T) {
		assert.Equal(t, "ForceUpperType(99)", ForceUpperType(99).String())
		_, err := ParseForceUpperType("-1")
		assert.Error(t, err)
	})
}
