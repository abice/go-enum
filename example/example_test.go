package example

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
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

	names := MakeNames()
	assert.Len(t, names, 11)
}

var makeTests = []struct {
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
		err:           errors.New("Porsche is not a valid Make, try [Toyota, Chevy, Ford, Tesla, Hyundai, Nissan, Jaguar, Audi, BMW, Mercedes-Benz, Volkswagon]"),
	},
}

func TestMakeUnmarshal(t *testing.T) {
	for _, test := range makeTests {
		t.Run(test.name, func(tt *testing.T) {
			x := &makeTest{}
			err := json.Unmarshal([]byte(test.input), x)
			if !test.errorExpected {
				require.NoError(tt, err, "failed unmarshalling the json.")
				assert.Equal(tt, test.output.M, x.M)

				// Marshal back
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

func TestFlagInterface(t *testing.T) {
	for _, test := range makeTests {
		t.Run(test.name, func(tt *testing.T) {
			if test.output != nil {
				var tmp Make
				m := &tmp
				require.Equal(tt, m.Type(), "Make")
				require.Equal(tt, Make(0), m.Get(), "Unset value should be default")
				require.NoError(tt, m.Set(test.output.M.String()), "failed setting flag value")
			}
		})
	}
}

func TestNoZeroValues(t *testing.T) {
	t.Run("base", func(tt *testing.T) {
		assert.Equal(tt, 20, int(NoZerosStart))
		assert.Equal(tt, 21, int(NoZerosMiddle))
		assert.Equal(tt, 22, int(NoZerosEnd))
		assert.Equal(tt, 23, int(NoZerosPs))
		assert.Equal(tt, 24, int(NoZerosPps))
		assert.Equal(tt, 25, int(NoZerosPpps))
		assert.Equal(tt, "ppps", NoZerosPpps.String())
		assert.Equal(tt, "NoZeros(4)", NoZeros(4).String())

		assert.Len(tt, NoZerosNames(), 6)

		_, err := ParseNoZeros("pppps")
		assert.EqualError(tt, err, "pppps is not a valid NoZeros, try [start, middle, end, ps, pps, ppps]")

		tmp, _ := ParseNoZeros("ppps")
		assert.Equal(tt, NoZerosPpps, tmp)

		val := map[string]*NoZeros{}

		err = json.Unmarshal([]byte(`{"nz":"pppps"}`), &val)
		assert.EqualError(tt, err, "pppps is not a valid NoZeros, try [start, middle, end, ps, pps, ppps]")

	})

	for _, name := range NoZerosNames() {
		t.Run(name, func(tt *testing.T) {
			val := map[string]*NoZeros{}

			rawJSON := fmt.Sprintf(`{"val":"%s"}`, name)

			require.NoError(tt, json.Unmarshal([]byte(rawJSON), &val), "Failed unmarshalling no zero")
			marshalled, err := json.Marshal(val)
			require.NoError(tt, err, "failed marshalling back to json")
			require.JSONEq(tt, rawJSON, string(marshalled), "marshalled json did not match")

			// Flag
			var tmp NoZeros
			nz := &tmp
			require.Equal(tt, nz.Type(), "NoZeros")
			require.Equal(tt, NoZeros(0), nz.Get(), "Unset value should be default")
			require.NoError(tt, nz.Set(name), "failed setting flag value")

		})

	}
}
