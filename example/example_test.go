package example

import (
	"encoding/json"
	"errors"
	"flag"
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type makeTest struct {
	M Make `json:"make"`
}

func TestMake(t *testing.T) {
	assert.Equal(t, "Ford", MakeFord.String())
	assert.Equal(t, "Make(99)", Make(99).String())
	ford := MakeFord
	assert.Implements(t, (*flag.Value)(nil), &ford)
	assert.Implements(t, (*flag.Getter)(nil), &ford)
	assert.Implements(t, (*pflag.Value)(nil), &ford)
}

func TestMakeUnmarshal(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		output        *makeTest
		errorExpected bool
		err           error
	}{
		{
			name:          "toyota",
			input:         `{"make":"Toyota"}`,
			output:        &makeTest{M: MakeToyota},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "chevy",
			input:         `{"make":"Chevy"}`,
			output:        &makeTest{M: MakeChevy},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "ford",
			input:         `{"make":"Ford"}`,
			output:        &makeTest{M: MakeFord},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "tesla",
			input:         `{"make":"Tesla"}`,
			output:        &makeTest{M: MakeTesla},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "hyundai",
			input:         `{"make":"Hyundai"}`,
			output:        &makeTest{M: MakeHyundai},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "nissan",
			input:         `{"make":"Nissan"}`,
			output:        &makeTest{M: MakeNissan},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "jaguar",
			input:         `{"make":"Jaguar"}`,
			output:        &makeTest{M: MakeJaguar},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "audi",
			input:         `{"make":"Audi"}`,
			output:        &makeTest{M: MakeAudi},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "bmw",
			input:         `{"make":"BMW"}`,
			output:        &makeTest{M: MakeBMW},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "mercedes",
			input:         `{"make":"Mercedes-Benz"}`,
			output:        &makeTest{M: MakeMercedesBenz},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "volkswagon",
			input:         `{"make":"Volkswagon"}`,
			output:        &makeTest{M: MakeVolkswagon},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "porsche",
			input:         `{"make":"Porsche"}`,
			output:        &makeTest{M: MakeVolkswagon},
			errorExpected: true,
			err:           errors.New("Porsche is not a valid Make"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			x := &makeTest{}
			err := json.Unmarshal([]byte(test.input), x)
			if !test.errorExpected {
				require.NoError(tt, err, "failed unmarshalling the json.")
				assert.Equal(tt, test.output.M, x.M)
			} else {
				require.Error(tt, err)
				assert.EqualError(tt, err, test.err.Error())
			}
		})
	}
}

func TestMakeMarshal(t *testing.T) {
	tests := []struct {
		name          string
		input         *makeTest
		output        string
		errorExpected bool
		err           error
	}{
		{
			name:          "toyota",
			output:        `{"make":"Toyota"}`,
			input:         &makeTest{M: MakeToyota},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "chevy",
			output:        `{"make":"Chevy"}`,
			input:         &makeTest{M: MakeChevy},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "ford",
			output:        `{"make":"Ford"}`,
			input:         &makeTest{M: MakeFord},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "tesla",
			output:        `{"make":"Tesla"}`,
			input:         &makeTest{M: MakeTesla},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "hyundai",
			output:        `{"make":"Hyundai"}`,
			input:         &makeTest{M: MakeHyundai},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "nissan",
			output:        `{"make":"Nissan"}`,
			input:         &makeTest{M: MakeNissan},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "jaguar",
			output:        `{"make":"Jaguar"}`,
			input:         &makeTest{M: MakeJaguar},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "audi",
			output:        `{"make":"Audi"}`,
			input:         &makeTest{M: MakeAudi},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "bmw",
			output:        `{"make":"BMW"}`,
			input:         &makeTest{M: MakeBMW},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "mercedes",
			output:        `{"make":"Mercedes-Benz"}`,
			input:         &makeTest{M: MakeMercedesBenz},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "volkswagon",
			output:        `{"make":"Volkswagon"}`,
			input:         &makeTest{M: MakeVolkswagon},
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

func TestNoZeroValues(t *testing.T) {
	assert.Equal(t, 20, int(NoZerosStart))
	assert.Equal(t, 21, int(NoZerosMiddle))
	assert.Equal(t, 22, int(NoZerosEnd))
	assert.Equal(t, 23, int(NoZerosPs))
	assert.Equal(t, 24, int(NoZerosPps))
	assert.Equal(t, 25, int(NoZerosPpps))
}
