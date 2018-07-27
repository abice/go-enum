package tuple

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTuple(t *testing.T) {
	assert.Equal(t, 1, FirstOfTwo(1, 2).(int))
	assert.Equal(t, 2, SecondOfTwo(1, 2).(int))
	assert.Equal(t, 1, FirstOfThree(1, 2, 3).(int))
	assert.Equal(t, 2, SecondOfThree(1, 2, 3).(int))
	assert.Equal(t, 3, ThirdOfThree(1, 2, 3).(int))

	assert.Nil(t, FirstOfTwo(nil, 2))

	assert.Nil(t, FirstError(nil, 2))
	assert.Nil(t, SecondError(1, nil))
	err := errors.New("error")
	assert.Equal(t, err, FirstError(err, 2))
	assert.Equal(t, err, SecondError(1, err))
	assert.Equal(t, true, FirstBool(true, 2))
	assert.Equal(t, false, FirstBool(false, 2))
	assert.Equal(t, true, SecondBool(1, true))
	assert.Equal(t, false, SecondBool(1, false))
}
