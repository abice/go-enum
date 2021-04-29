package example

import (
	"testing"

	"github.com/stretchr/testify/assert"
)
func TestUserTemplateColor(t *testing.T) {
	assert.Equal(t, OceanColor(0), OceanColorCerulean)
	assert.Equal(t, true, ParseOceanColorExample())
}