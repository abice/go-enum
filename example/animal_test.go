package example

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type animalData struct {
	AnimalX Animal `json:"animal"`
}

func TestAnimalString(t *testing.T) {
	x := Animal(109)
	assert.Equal(t, "Animal(109)", x.String())
	x = Animal(1)
	assert.Equal(t, "Dog", x.String())

	y, err := ParseAnimal("Cat")
	require.NoError(t, err, "Failed parsing cat")
	assert.Equal(t, AnimalCat, y)

	z, err := ParseAnimal("Snake")
	require.Error(t, err, "Shouldn't parse a snake")
	assert.Equal(t, Animal(0), z)
}

func TestAnimalUnmarshal(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		output        *animalData
		errorExpected bool
		err           error
	}{
		{
			name:          "cat",
			input:         `{"animal":0}`,
			output:        &animalData{AnimalX: AnimalCat},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "dog",
			input:         `{"animal":1}`,
			output:        &animalData{AnimalX: AnimalDog},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "fish",
			input:         `{"animal":2}`,
			output:        &animalData{AnimalX: AnimalFish},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "notananimal",
			input:         `{"animal":22}`,
			output:        &animalData{AnimalX: Animal(22)},
			errorExpected: false,
			err:           nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			x := &animalData{}
			err := json.Unmarshal([]byte(test.input), x)
			if !test.errorExpected {
				require.NoError(tt, err, "failed unmarshalling the json.")
				assert.Equal(tt, test.output.AnimalX, x.AnimalX)
			} else {
				require.Error(tt, err)
				assert.EqualError(tt, err, test.err.Error())
			}
		})
	}
}

func TestAnimalMarshal(t *testing.T) {
	tests := []struct {
		name          string
		input         *animalData
		output        string
		errorExpected bool
		err           error
	}{
		{
			name:          "cat",
			output:        `{"animal":0}`,
			input:         &animalData{AnimalX: AnimalCat},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "dog",
			output:        `{"animal":1}`,
			input:         &animalData{AnimalX: AnimalDog},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "fish",
			output:        `{"animal":2}`,
			input:         &animalData{AnimalX: AnimalFish},
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
