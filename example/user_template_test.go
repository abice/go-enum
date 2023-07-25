//go:build example
// +build example

package example

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserTemplateColor(t *testing.T) {
	assert.Equal(t, OceanColor(0), OceanColorCerulean)
	assert.Equal(t, "OceanColor is an enumeration of ocean colors that are allowed.", ParseOceanColorDescription())
	assert.Equal(t, true, ParseOceanColorExample())
	assert.Equal(t, true, ParseOceanColorGlobbedExample())
	assert.Equal(t, true, ParseOceanColorGlobbedExample2())

	val, err := ParseOceanColor("Cerulean")
	assert.NoError(t, err)
	assert.Equal(t, "Cerulean", val.String())

	assert.Equal(t, "OceanColor(99)", OceanColor(99).String())
	_, err = ParseOceanColor("-1")
	assert.Error(t, err)
}
