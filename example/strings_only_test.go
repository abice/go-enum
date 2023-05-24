//go:build example
// +build example

package example

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStrState(t *testing.T) {
	x := StrState("")
	assert.Equal(t, "", x.String())

	assert.Equal(t, StrState("pending"), StrStatePending)
	assert.Equal(t, StrState("running"), StrStateRunning)
	assert.Equal(t, &x, StrState("").Ptr())
}

func TestStrStateMustParse(t *testing.T) {
	x := `avocado`

	assert.PanicsWithError(t, x+" is not a valid StrState, try [pending, running, completed, error]", func() { MustParseStrState(x) })
	assert.NotPanics(t, func() { MustParseStrState(StrStateFailed.String()) })
}

func TestStrStateUnmarshal(t *testing.T) {
	type testData struct {
		StrStateX StrState `json:"state"`
	}
	tests := []struct {
		name          string
		input         string
		output        *testData
		errorExpected bool
		err           error
	}{
		{
			name:          "pending",
			input:         `{"state":"Pending"}`,
			output:        &testData{StrStateX: StrStatePending},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "pendinglower",
			input:         `{"state":"pending"}`,
			output:        &testData{StrStateX: StrStatePending},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "running",
			input:         `{"state":"RUNNING"}`,
			output:        &testData{StrStateX: StrStateRunning},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "running",
			input:         `{"state":"running"}`,
			output:        &testData{StrStateX: StrStateRunning},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "completed",
			input:         `{"state":"Completed"}`,
			output:        &testData{StrStateX: StrStateCompleted},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "completedlower",
			input:         `{"state":"completed"}`,
			output:        &testData{StrStateX: StrStateCompleted},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "failed",
			input:         `{"state":"Error"}`,
			output:        &testData{StrStateX: StrStateFailed},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "failedlower",
			input:         `{"state":"error"}`,
			output:        &testData{StrStateX: StrStateFailed},
			errorExpected: false,
			err:           nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			x := &testData{}
			err := json.Unmarshal([]byte(test.input), x)
			if !test.errorExpected {
				require.NoError(tt, err, "failed unmarshalling the json.")
				assert.Equal(tt, test.output.StrStateX, x.StrStateX)
			} else {
				require.Error(tt, err)
				assert.EqualError(tt, err, test.err.Error())
			}
		})
	}
}

func TestStrStateMarshal(t *testing.T) {
	type testData struct {
		StrStateX StrState `json:"state"`
	}
	tests := []struct {
		name          string
		input         *testData
		output        string
		errorExpected bool
		err           error
	}{
		{
			name:          "black",
			output:        `{"state":"pending"}`,
			input:         &testData{StrStateX: StrStatePending},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "white",
			output:        `{"state":"running"}`,
			input:         &testData{StrStateX: StrStateRunning},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "red",
			output:        `{"state":"completed"}`,
			input:         &testData{StrStateX: StrStateCompleted},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "green",
			output:        `{"state":"error"}`,
			input:         &testData{StrStateX: StrStateFailed},
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

func TestStrStateSQLExtras(t *testing.T) {
	_, err := ParseStrState(`NotAState`)
	assert.Error(t, err, "Should have had an error parsing a non status")

	var (
		intVal       int      = 3
		strVal       string   = "completed"
		enumVal      StrState = StrStateCompleted
		nullInt      *int
		nullInt64    *int64
		nullUint     *uint
		nullUint64   *uint64
		nullString   *string
		nullStrState *StrState
	)

	tests := map[string]struct {
		input  interface{}
		result NullStrState
	}{
		"nil": {},
		"val": {
			input: StrStatePending,
			result: NullStrState{
				StrState: StrStatePending,
				Valid:    true,
			},
		},
		"ptr": {
			input: &enumVal,
			result: NullStrState{
				StrState: StrStateCompleted,
				Valid:    true,
			},
		},
		"string": {
			input: strVal,
			result: NullStrState{
				StrState: StrStateCompleted,
				Valid:    true,
			},
		},
		"*string": {
			input: &strVal,
			result: NullStrState{
				StrState: StrStateCompleted,
				Valid:    true,
			},
		},
		"invalid string": {
			input: "random value",
		},
		"[]byte": {
			input: []byte(StrStateRunning.String()),
			result: NullStrState{
				StrState: StrStateRunning,
				Valid:    true,
			},
		},
		"int": {
			input: intVal,
		},
		"*int": {
			input: &intVal,
		},
		"nullInt": {
			input: nullInt,
		},
		"nullInt64": {
			input: nullInt64,
		},
		"nullUint": {
			input: nullUint,
		},
		"nullUint64": {
			input: nullUint64,
		},
		"nullString": {
			input: nullString,
		},
		"nullImageType": {
			input: nullStrState,
		},
		"int as []byte": { // must have --sqlnullint flag to get this feature.
			input: []byte("3"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			status := NewNullStrState(tc.input)
			assert.Equal(t, status, tc.result)
		})
	}
}
