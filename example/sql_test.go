package example

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSQLExtras(t *testing.T) {

	assert.Equal(t, "ProjectStatus(22)", ProjectStatus(22).String(), "String value is not correct")

	_, err := ParseProjectStatus(`NotAStatus`)
	assert.Error(t, err, "Should have had an error parsing a non status")
}
