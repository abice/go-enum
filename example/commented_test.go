package example

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type commentedData struct {
	CommentedX Commented `json:"Commented"`
}

func TestCommentedEnumString(t *testing.T) {
	x := Commented(109)
	assert.Equal(t, "Commented(109)", x.String())
	x = Commented(1)
	assert.Equal(t, "value2", x.String())

	y, err := ParseCommented("value3")
	require.NoError(t, err, "Failed parsing cat")
	assert.Equal(t, CommentedValue3, y)

	z, err := ParseCommented("value4")
	require.Error(t, err, "Shouldn't parse a snake")
	assert.Equal(t, Commented(0), z)
}

func TestCommentedUnmarshal(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		output        *commentedData
		errorExpected bool
		err           error
	}{
		{
			name:          "value1",
			input:         `{"Commented":"value1"}`,
			output:        &commentedData{CommentedX: CommentedValue1},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "value2",
			input:         `{"Commented":"value2"}`,
			output:        &commentedData{CommentedX: CommentedValue2},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "value3",
			input:         `{"Commented":"value3"}`,
			output:        &commentedData{CommentedX: CommentedValue3},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "notanCommented",
			input:         `{"Commented":"value4"}`,
			output:        nil,
			errorExpected: true,
			err:           errors.New("value4 is not a valid Commented"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			x := &commentedData{}
			err := json.Unmarshal([]byte(test.input), x)
			if !test.errorExpected {
				require.NoError(tt, err, "failed unmarshalling the json.")
				assert.Equal(tt, test.output.CommentedX, x.CommentedX)
				raw, err := json.Marshal(test.output)
				require.NoError(tt, err, "failed marshalling back to json")
				require.JSONEq(tt, test.input, string(raw), "json didn't match")
			} else {
				require.Error(tt, err)
				assert.EqualError(tt, err, test.err.Error())
			}
		})
	}
}
