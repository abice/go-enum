package example

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testData struct {
	ColorX Color `json:"color"`
}

func TestColorString(t *testing.T) {
	x := Color(109)
	assert.Equal(t, "Color(109)", x.String())
}

func TestColorUnmarshal(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		output        *testData
		errorExpected bool
		err           error
	}{
		{
			name:          "black",
			input:         `{"color":"Black"}`,
			output:        &testData{ColorX: ColorBlack},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "blacklower",
			input:         `{"color":"black"}`,
			output:        &testData{ColorX: ColorBlack},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "white",
			input:         `{"color":"White"}`,
			output:        &testData{ColorX: ColorWhite},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "whitelower",
			input:         `{"color":"white"}`,
			output:        &testData{ColorX: ColorWhite},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "red",
			input:         `{"color":"Red"}`,
			output:        &testData{ColorX: ColorRed},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "redlower",
			input:         `{"color":"red"}`,
			output:        &testData{ColorX: ColorRed},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "green",
			input:         `{"color":"Green"}`,
			output:        &testData{ColorX: ColorGreen},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "greenlower",
			input:         `{"color":"green"}`,
			output:        &testData{ColorX: ColorGreen},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "blue",
			input:         `{"color":"Blue"}`,
			output:        &testData{ColorX: ColorBlue},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "bluelower",
			input:         `{"color":"blue"}`,
			output:        &testData{ColorX: ColorBlue},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "grey",
			input:         `{"color":"Grey"}`,
			output:        &testData{ColorX: ColorGrey},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "greylower",
			input:         `{"color":"grey"}`,
			output:        &testData{ColorX: ColorGrey},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "yellow",
			input:         `{"color":"Yellow"}`,
			output:        &testData{ColorX: ColorYellow},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "yellowlower",
			input:         `{"color":"yellow"}`,
			output:        &testData{ColorX: ColorYellow},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "magenta",
			input:         `{"color":"Magenta"}`,
			output:        &testData{ColorX: ColorYellow},
			errorExpected: true,
			err:           errors.New("Magenta is not a valid Color"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			x := &testData{}
			err := json.Unmarshal([]byte(test.input), x)
			if !test.errorExpected {
				require.NoError(tt, err, "failed unmarshalling the json.")
				assert.Equal(tt, test.output.ColorX, x.ColorX)
			} else {
				require.Error(tt, err)
				assert.EqualError(tt, err, test.err.Error())
			}
		})
	}
}

func TestColorMarshal(t *testing.T) {
	tests := []struct {
		name          string
		input         *testData
		output        string
		errorExpected bool
		err           error
	}{
		{
			name:          "black",
			output:        `{"color":"Black"}`,
			input:         &testData{ColorX: ColorBlack},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "white",
			output:        `{"color":"White"}`,
			input:         &testData{ColorX: ColorWhite},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "red",
			output:        `{"color":"Red"}`,
			input:         &testData{ColorX: ColorRed},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "green",
			output:        `{"color":"Green"}`,
			input:         &testData{ColorX: ColorGreen},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "blue",
			output:        `{"color":"Blue"}`,
			input:         &testData{ColorX: ColorBlue},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "grey",
			output:        `{"color":"Grey"}`,
			input:         &testData{ColorX: ColorGrey},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "yellow",
			output:        `{"color":"Yellow"}`,
			input:         &testData{ColorX: ColorYellow},
			errorExpected: false,
			err:           nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			raw, err := json.Marshal(test.input)
			require.NoError(tt, err, "failed marshalling to json")
			assert.JSONEq(tt, test.output, string(raw))
		})
	}
}
