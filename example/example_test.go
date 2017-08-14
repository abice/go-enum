package example

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModel(t *testing.T) {
	assert.Equal(t, "Ford", ModelFord.String())
}
