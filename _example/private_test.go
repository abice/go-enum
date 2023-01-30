package example

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrivateInt(t *testing.T) {
	assert.Equal(t, privateIntFirst, privateInt(0))

	assert.Equal(t, privateIntNames(), _privateIntNames)

	assert.Contains(t, privateIntValues(), privateIntSecond)

	_, err := parsePrivateInt("invalid")
	assert.ErrorIs(t, err, errInvalidPrivateInt)
}

func TestPrivateStr(t *testing.T) {
	assert.Equal(t, privateStrFirst, privateStr("a"))

	assert.Equal(t, privateStrNames(), _privateStrNames)

	assert.Contains(t, privateStrValues(), privateStrSecond)

	_, err := parsePrivateStr("invalid")
	assert.ErrorIs(t, err, errInvalidPrivateStr)
}
