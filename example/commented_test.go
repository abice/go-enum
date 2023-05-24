//go:build example
// +build example

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

type complexCommentedData struct {
	ComplexCommentedX ComplexCommented `json:"ComplexCommented,omitempty"`
}

func TestComplexCommentedEnumString(t *testing.T) {
	x := ComplexCommented(109)
	assert.Equal(t, "ComplexCommented(109)", x.String())
	x = ComplexCommented(1)
	assert.Equal(t, "value1", x.String())

	y, err := ParseComplexCommented("value3")
	require.NoError(t, err, "Failed parsing value3")
	assert.Equal(t, ComplexCommentedValue3, y)

	z, err := ParseComplexCommented("value4")
	require.Error(t, err, "Shouldn't parse a value4")
	assert.Equal(t, ComplexCommented(0), z)
}

func TestComplexCommentedUnmarshal(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		output        *complexCommentedData
		errorExpected bool
		err           error
	}{
		{
			name:          "value1",
			input:         `{"ComplexCommented":"value1"}`,
			output:        &complexCommentedData{ComplexCommentedX: ComplexCommentedValue1},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "value2",
			input:         `{"ComplexCommented":"value2"}`,
			output:        &complexCommentedData{ComplexCommentedX: ComplexCommentedValue2},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "value3",
			input:         `{"ComplexCommented":"value3"}`,
			output:        &complexCommentedData{ComplexCommentedX: ComplexCommentedValue3},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "notanCommented",
			input:         `{"ComplexCommented":"value4"}`,
			output:        nil,
			errorExpected: true,
			err:           errors.New("value4 is not a valid ComplexCommented"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			x := &complexCommentedData{}
			err := json.Unmarshal([]byte(test.input), x)
			if !test.errorExpected {
				require.NoError(tt, err, "failed unmarshalling the json.")
				assert.Equal(tt, test.output.ComplexCommentedX, x.ComplexCommentedX)
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
