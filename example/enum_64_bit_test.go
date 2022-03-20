package example

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnum64Bit(t *testing.T) {

	tests := map[string]struct {
		input  string
		output Enum64bit
	}{
		"E2P15": {
			input:  `E2P15`,
			output: Enum64bitE2P15,
		},
		"E2P63": {
			input:  `E2P63`,
			output: Enum64bitE2P63,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			output, err := ParseEnum64bit(tc.input)
			assert.NoError(t, err)
			assert.Equal(t, tc.output, output)

			assert.Equal(t, tc.input, output.String())
		})
	}

	t.Run("basics", func(t *testing.T) {
		assert.Equal(t, "E2P23", Enum64bitE2P23.String())
		assert.Equal(t, "Enum64bit(99)", Enum64bit(99).String())
		_, err := ParseEnum64bit("-1")
		assert.Error(t, err)

		names := Enum64bitNames()
		assert.Len(t, names, 16)
	})

}
