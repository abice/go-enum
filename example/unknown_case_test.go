//go:build example
// +build example

package example

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStrGenderUnmarshal(t *testing.T) {
	type testData struct {
		Gender StrGender `json:"gender"`
	}
	tests := []struct {
		name          string
		input         string
		output        *testData
		errorExpected bool
	}{
		{
			name:          "male",
			input:         `{"gender":"male"}`,
			output:        &testData{Gender: StrGenderMale},
			errorExpected: false,
		},
		{
			name:          "female",
			input:         `{"gender":"female"}`,
			output:        &testData{Gender: StrGenderFemale},
			errorExpected: false,
		},
		{
			name:          "malformed json",
			input:         `{"gender":"male}`,
			output:        nil,
			errorExpected: true,
		},
		{
			name:          "unknown value",
			input:         `{"gender":"unknown"}`,
			output:        &testData{Gender: StrGenderUnknown},
			errorExpected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			x := &testData{}
			err := json.Unmarshal([]byte(test.input), x)
			if !test.errorExpected {
				require.NoError(tt, err, "failed unmarshalling the json.")
				assert.Equal(tt, test.output.Gender, x.Gender)
			} else {
				require.Error(tt, err)
			}
		})
	}
}
