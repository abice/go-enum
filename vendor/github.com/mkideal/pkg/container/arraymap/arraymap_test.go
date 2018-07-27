package arraymap_test

import (
	"testing"

	"github.com/mkideal/pkg/container/arraymap"
	"github.com/stretchr/testify/assert"
)

func TestArrayMap(t *testing.T) {
	m := arraymap.New()
	k1 := m.Add(2)
	k2 := m.Add(3)
	v1, ok := m.Get(k1)
	assert.True(t, ok)
	assert.Equal(t, 2, v1)
	v2, ok := m.Get(k2)
	assert.True(t, ok)
	assert.Equal(t, 3, v2)
	m.Remove(k2)
	_, ok = m.Get(k2)
	assert.False(t, ok)
}
