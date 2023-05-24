//go:build example
// +build example

package example

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestForceLowerString(t *testing.T) {
	tests := map[string]struct {
		input  string
		output ForceLowerType
	}{
		"dataswap": {
			input:  `dataswap`,
			output: ForceLowerTypeDataSwap,
		},
		"bootnode": {
			input:  `bootnode`,
			output: ForceLowerTypeBootNode,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			output, err := ParseForceLowerType(tc.input)
			assert.NoError(t, err)
			assert.Equal(t, tc.output, output)

			assert.Equal(t, tc.input, output.String())
		})
	}

	t.Run("failures", func(t *testing.T) {
		assert.Equal(t, "ForceLowerType(99)", ForceLowerType(99).String())
		_, err := ParseForceLowerType("-1")
		assert.Error(t, err)
	})
}
