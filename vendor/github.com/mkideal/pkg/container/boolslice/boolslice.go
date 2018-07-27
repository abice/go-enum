package boolslice

import (
	"github.com/mkideal/pkg/container"
)

type BoolSlice struct {
	data   []byte
	length int
}

func New() *BoolSlice {
	return &BoolSlice{data: []byte{}}
}

func NewWithSize(size, cap int) *BoolSlice {
	return &BoolSlice{data: make([]byte, (size+7)>>3, (cap+7)>>3), length: size}
}

func NewWithSlice(slice []bool) *BoolSlice {
	size := len(slice)
	s := NewWithSize(size, size)
	for i := 0; i < size; i++ {
		s.Set(i, slice[i])
	}
	return s
}

func (s *BoolSlice) ij(index int) (i int, j byte) { return index >> 3, byte(index & 0x7) }

func (s *BoolSlice) Get(index int) bool {
	i, j := s.ij(index)
	return s.data[i]&(1<<j) != 0
}

func (s *BoolSlice) Set(index int, value bool) {
	i, j := s.ij(index)
	if value {
		s.data[i] |= 1 << j
	} else {
		s.data[i] &= ^(1 << j)
	}
}

func (s *BoolSlice) Push(value bool) {
	i, j := s.ij(s.length)
	if i >= len(s.data) {
		s.data = append(s.data, 0)
	}
	if value {
		s.data[i] |= 1 << j
	} else {
		s.data[i] &= ^(1 << j)
	}
	s.length++
}

func (s *BoolSlice) Pop() (value bool) {
	if s.length == 0 {
		panic("length == 0")
	}
	s.length--
	value = s.Get(s.length)
	if n := (s.length >> 3) + 1; n < len(s.data) {
		s.data = s.data[:n]
	}
	return
}

func (s *BoolSlice) Insert(index int, value bool) {
	s.Push(value)
	for i := s.length - 1; i > index; i-- {
		s.Set(i, s.Get(i-1))
	}
	s.Set(index, value)
}

func (s *BoolSlice) Truncate(from, to int) {
	if to < s.length {
		s.length = to
		if n := (s.length >> 3) + 1; n < len(s.data) {
			s.data = s.data[:n]
		}
	}
	if from > 0 {
		l := to - from
		for i := 0; i < l; i++ {
			s.Set(i, s.Get(i+from))
		}
		s.length = l
		if n := (s.length >> 3) + 1; n < len(s.data) {
			s.data = s.data[:n]
		}
	}
}

func (s *BoolSlice) Clone() *BoolSlice {
	s2 := &BoolSlice{data: make([]byte, s.length), length: s.length}
	copy(s2.data, s.data)
	return s2
}

func (s *BoolSlice) Equal(s2 *BoolSlice) bool {
	if s.length != s2.length {
		return false
	}
	n := len(s.data)
	if n == 0 {
		return true
	}
	for i := 0; i+1 < n; i++ {
		if s.data[i] != s2.data[i] {
			return false
		}
	}
	j := byte(s.length & 0x7)
	return s.data[n-1]<<(8-j) == s2.data[n-1]<<(8-j)
}

func (s *BoolSlice) EqualToSlice(s2 []bool) bool {
	if s.Len() != len(s2) {
		return false
	}
	for i, n := 0, s.Len(); i < n; i++ {
		if s.Get(i) != s2[i] {
			return false
		}
	}
	return true
}

// container.Container adaptor
type iterator struct {
	s      *BoolSlice
	cursor int
}

func (iter *iterator) Next() (k, v interface{}) {
	if iter.cursor >= iter.s.length {
		return iter.cursor, nil
	}
	index := iter.cursor
	iter.cursor++
	return index, iter.s.Get(index)
}

func (s *BoolSlice) Len() int                 { return s.length }
func (s *BoolSlice) Iter() container.Iterator { return &iterator{s: s} }
func (s *BoolSlice) Contains(value interface{}) bool {
	v := value.(bool)
	for i := 0; i < s.length; i++ {
		if s.Get(i) == v {
			return true
		}
	}
	return false
}
