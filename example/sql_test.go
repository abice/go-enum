package example

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSQLExtras(t *testing.T) {

	assert.Equal(t, "ProjectStatus(22)", ProjectStatus(22).String(), "String value is not correct")

	_, err := ParseProjectStatus(`NotAStatus`)
	assert.Error(t, err, "Should have had an error parsing a non status")

	var (
		intVal int    = 3
		strVal string = "completed"
	)

	tests := map[string]struct {
		input  interface{}
		result NullProjectStatus
	}{
		"nil": {},
		"val": {
			input: ProjectStatusRejected,
			result: NullProjectStatus{
				ProjectStatus: ProjectStatusRejected,
				Valid:         true,
			},
		},
		"ptr": {
			input: ProjectStatusCompleted.Ptr(),
			result: NullProjectStatus{
				ProjectStatus: ProjectStatusCompleted,
				Valid:         true,
			},
		},
		"string": {
			input: strVal,
			result: NullProjectStatus{
				ProjectStatus: ProjectStatusCompleted,
				Valid:         true,
			},
		},
		"*string": {
			input: &strVal,
			result: NullProjectStatus{
				ProjectStatus: ProjectStatusCompleted,
				Valid:         true,
			},
		},
		"invalid string": {
			input: "random value",
		},
		"[]byte": {
			input: []byte(ProjectStatusInWork.String()),
			result: NullProjectStatus{
				ProjectStatus: ProjectStatusInWork,
				Valid:         true,
			},
		},
		"int": {
			input: intVal,
			result: NullProjectStatus{
				ProjectStatus: ProjectStatusRejected,
				Valid:         true,
			},
		},
		"*int": {
			input: &intVal,
			result: NullProjectStatus{
				ProjectStatus: ProjectStatusRejected,
				Valid:         true,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			status := NewNullProjectStatus(tc.input)
			assert.Equal(t, status, tc.result)

		})
	}

}
