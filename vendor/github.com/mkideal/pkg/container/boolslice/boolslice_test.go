package boolslice_test

import (
	"math/rand"
	"testing"

	"github.com/mkideal/pkg/container/boolslice"
	"github.com/stretchr/testify/assert"
)

var valuesSet = func() [][]bool {
	n := 100
	ret := make([][]bool, 0, n)
	for i := 0; i < n; i++ {
		size := 1024 + i
		list := make([]bool, 0, size)
		for j := 0; j < size; j++ {
			list = append(list, rand.Intn(2) == 1)
		}
		ret = append(ret, list)
	}
	return ret
}()

func TestPushAndPop(t *testing.T) {
	for _, values := range valuesSet {
		s := boolslice.New()
		for _, v := range values {
			s.Push(v)
		}
		assert.Equal(t, len(values), s.Len())
		for i, n := 0, s.Len(); i < n; i++ {
			assert.Equal(t, values[i], s.Get(i))
		}
		assert.Equal(t, values[s.Len()-1], s.Pop())
	}
	s1 := boolslice.NewWithSlice([]bool{true, true, true, true, true, true, true, true})
	s2 := boolslice.NewWithSlice([]bool{true, true, true, true, true, true, true})
	s1.Pop()
	assert.True(t, s1.Equal(s2))
}

func TestSet(t *testing.T) {
	for _, values := range valuesSet {
		s := boolslice.NewWithSize(len(values), len(values))
		assert.Equal(t, len(values), s.Len())
		for i, v := range values {
			s.Set(i, v)
		}
		for i, v := range values {
			assert.Equal(t, v, s.Get(i))
		}
	}
}

func TestInsert(t *testing.T) {
	s := boolslice.NewWithSlice([]bool{true, false, true, false})
	s.Insert(0, false)
	s.Insert(s.Len(), true)
	s.Insert(1, false)
	assert.Equal(t, 7, s.Len())
	assert.Equal(t, false, s.Get(0))
	assert.Equal(t, false, s.Get(1))
	assert.Equal(t, true, s.Get(2))
	assert.Equal(t, false, s.Get(3))
	assert.Equal(t, true, s.Get(4))
	assert.Equal(t, false, s.Get(5))
	assert.Equal(t, true, s.Get(6))
}

func TestClone(t *testing.T) {
	for _, values := range valuesSet {
		s1 := boolslice.NewWithSlice(values)
		s2 := s1.Clone()
		assert.True(t, s1.Equal(s2))
	}
}

func TestTruncate(t *testing.T) {
	for _, values := range valuesSet {
		s1 := boolslice.NewWithSlice(values)
		mid := s1.Len() / 2

		s2 := s1.Clone()
		s2.Truncate(0, s2.Len())
		assert.True(t, s1.Equal(s2))

		s2 = s1.Clone()
		s2.Truncate(0, mid)
		assert.Equal(t, mid, s2.Len())
		assert.True(t, s2.EqualToSlice(values[:mid]))

		s2 = s1.Clone()
		s2.Truncate(mid, s2.Len())
		assert.Equal(t, s1.Len()-mid, s2.Len())
		assert.True(t, s2.EqualToSlice(values[mid:]))

		s2 = s1.Clone()
		end := (mid + s2.Len()) / 2
		s2.Truncate(mid, end)
		assert.Equal(t, end-mid, s2.Len())
		assert.True(t, s2.EqualToSlice(values[mid:end]))

		s2 = s1.Clone()
		s2.Truncate(mid, mid)
		assert.Equal(t, 0, s2.Len())
	}
}
