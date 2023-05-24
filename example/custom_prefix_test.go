//go:build example
// +build example

package example

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type productData struct {
	ProductX Product `json:"product"`
}

func TestProductString(t *testing.T) {
	x := Product(109)
	assert.Equal(t, "Product(109)", x.String())
	x = Product(1)
	assert.Equal(t, "Dynamite", x.String())

	y, err := ParseProduct("Anvil")
	require.NoError(t, err, "Failed parsing anvil")
	assert.Equal(t, AcmeIncProductAnvil, y)

	z, err := ParseProduct("Snake")
	require.Error(t, err, "Shouldn't parse a snake")
	assert.Equal(t, Product(0), z)
}

func TestProductUnmarshal(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		output        *productData
		errorExpected bool
		err           error
	}{
		{
			name:          "anvil",
			input:         `{"product":0}`,
			output:        &productData{ProductX: AcmeIncProductAnvil},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "dynamite",
			input:         `{"product":1}`,
			output:        &productData{ProductX: AcmeIncProductDynamite},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "glue",
			input:         `{"product":2}`,
			output:        &productData{ProductX: AcmeIncProductGlue},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "notanproduct",
			input:         `{"product":22}`,
			output:        &productData{ProductX: Product(22)},
			errorExpected: false,
			err:           nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			x := &productData{}
			err := json.Unmarshal([]byte(test.input), x)
			if !test.errorExpected {
				require.NoError(tt, err, "failed unmarshalling the json.")
				assert.Equal(tt, test.output.ProductX, x.ProductX)
			} else {
				require.Error(tt, err)
				assert.EqualError(tt, err, test.err.Error())
			}
		})
	}
}

func TestProductMarshal(t *testing.T) {
	tests := []struct {
		name          string
		input         *productData
		output        string
		errorExpected bool
		err           error
	}{
		{
			name:          "anvil",
			output:        `{"product":0}`,
			input:         &productData{ProductX: AcmeIncProductAnvil},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "dynamite",
			output:        `{"product":1}`,
			input:         &productData{ProductX: AcmeIncProductDynamite},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "glue",
			output:        `{"product":2}`,
			input:         &productData{ProductX: AcmeIncProductGlue},
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
