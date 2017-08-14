package example

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type modelTest struct {
	M Model `json:"model"`
}

func TestModel(t *testing.T) {
	assert.Equal(t, "Ford", ModelFord.String())
	assert.Equal(t, "Model(99)", Model(99).String())
}

func TestModelUnmarshal(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		output        *modelTest
		errorExpected bool
		err           error
	}{
		{
			name:          "toyota",
			input:         `{"model":"Toyota"}`,
			output:        &modelTest{M: ModelToyota},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "chevy",
			input:         `{"model":"Chevy"}`,
			output:        &modelTest{M: ModelChevy},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "ford",
			input:         `{"model":"Ford"}`,
			output:        &modelTest{M: ModelFord},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "tesla",
			input:         `{"model":"Tesla"}`,
			output:        &modelTest{M: ModelTesla},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "hyundai",
			input:         `{"model":"Hyundai"}`,
			output:        &modelTest{M: ModelHyundai},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "nissan",
			input:         `{"model":"Nissan"}`,
			output:        &modelTest{M: ModelNissan},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "jaguar",
			input:         `{"model":"Jaguar"}`,
			output:        &modelTest{M: ModelJaguar},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "audi",
			input:         `{"model":"Audi"}`,
			output:        &modelTest{M: ModelAudi},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "bmw",
			input:         `{"model":"BMW"}`,
			output:        &modelTest{M: ModelBMW},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "mercedes",
			input:         `{"model":"Mercedes"}`,
			output:        &modelTest{M: ModelMercedes},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "volkswagon",
			input:         `{"model":"Volkswagon"}`,
			output:        &modelTest{M: ModelVolkswagon},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "porsche",
			input:         `{"model":"Porsche"}`,
			output:        &modelTest{M: ModelVolkswagon},
			errorExpected: true,
			err:           errors.New("Porsche is not a valid Model"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			x := &modelTest{}
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

func TestModelMarshal(t *testing.T) {
	tests := []struct {
		name          string
		input         *modelTest
		output        string
		errorExpected bool
		err           error
	}{
		{
			name:          "toyota",
			output:        `{"model":"Toyota"}`,
			input:         &modelTest{M: ModelToyota},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "chevy",
			output:        `{"model":"Chevy"}`,
			input:         &modelTest{M: ModelChevy},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "ford",
			output:        `{"model":"Ford"}`,
			input:         &modelTest{M: ModelFord},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "tesla",
			output:        `{"model":"Tesla"}`,
			input:         &modelTest{M: ModelTesla},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "hyundai",
			output:        `{"model":"Hyundai"}`,
			input:         &modelTest{M: ModelHyundai},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "nissan",
			output:        `{"model":"Nissan"}`,
			input:         &modelTest{M: ModelNissan},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "jaguar",
			output:        `{"model":"Jaguar"}`,
			input:         &modelTest{M: ModelJaguar},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "audi",
			output:        `{"model":"Audi"}`,
			input:         &modelTest{M: ModelAudi},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "bmw",
			output:        `{"model":"BMW"}`,
			input:         &modelTest{M: ModelBMW},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "mercedes",
			output:        `{"model":"Mercedes"}`,
			input:         &modelTest{M: ModelMercedes},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "volkswagon",
			output:        `{"model":"Volkswagon"}`,
			input:         &modelTest{M: ModelVolkswagon},
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
