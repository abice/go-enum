//go:build example
// +build example

package example

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnum32bitfield(t *testing.T) {
	t.Parallel()

	allFields := Enum32bitfieldValues()

	// Generate all pairs (i, j) where i < j to avoid duplicates
	testCases := make([]struct {
		name   string
		field1 Enum32bitfield
		field2 Enum32bitfield
	}, 0, len(allFields)*(len(allFields)-1)/2)

	for i := 0; i < len(allFields); i++ {
		for j := i + 1; j < len(allFields); j++ {
			testCases = append(testCases, struct {
				name   string
				field1 Enum32bitfield
				field2 Enum32bitfield
			}{
				name:   allFields[i].String() + "_" + allFields[j].String(),
				field1: allFields[i],
				field2: allFields[j],
			})
		}
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Collapse the two views
			collapsed := tc.field1 | tc.field2

			assert.Equal(t, tc.field1, collapsed^tc.field2)
			assert.Equal(t, tc.field2, collapsed^tc.field1)
		})
	}
}
