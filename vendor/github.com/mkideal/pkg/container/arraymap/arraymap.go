package arraymap

import (
	"github.com/mkideal/pkg/container/boolslice"
)

type T interface{}

type ArrayMap struct {
	data       []T
	validFlags *boolslice.BoolSlice // validFlags.Len() = len(data)
	holes      []int                // len(holes) <= len(data)
}

func New() *ArrayMap {
	return &ArrayMap{
		data:       []T{},
		validFlags: boolslice.New(),
		holes:      []int{},
	}
}

func NewWithSize(size, cap int) *ArrayMap {
	return &ArrayMap{
		data:       make([]T, size, cap),
		validFlags: boolslice.NewWithSize(size, cap),
		holes:      []int{},
	}
}

func (m *ArrayMap) isValidKey(key int) bool {
	return key < len(m.data) && m.validFlags.Get(key)
}

func (m *ArrayMap) Len() int { return len(m.data) - len(m.holes) }
func (m *ArrayMap) Cap() int { return cap(m.data) }

// Get gets the value by key
func (m *ArrayMap) Get(key int) (value T, ok bool) {
	if m.isValidKey(key) {
		return m.data[key], true
	}
	return nil, false
}

// Update updates <key,value>
// Returns true if key is valid
func (m *ArrayMap) Update(key int, value T) bool {
	if m.isValidKey(key) {
		m.data[key] = value
		return true
	}
	return false
}

// Add adds a new element and returns the allocated key
func (m *ArrayMap) Add(value T) (key int) {
	if n := len(m.holes); n > 0 {
		key = m.holes[n-1]
		m.holes = m.holes[:n-1]
		m.data[key] = value
		m.validFlags.Set(key, true)
		return
	}
	key = len(m.data)
	m.data = append(m.data, value)
	m.validFlags.Push(true)
	return
}

// Remove removes the value by key
func (m *ArrayMap) Remove(key int) (value T, ok bool) {
	if m.isValidKey(key) {
		value = m.data[key]
		ok = true
		m.validFlags.Set(key, false)
		m.holes = append(m.holes, key)
	}
	return
}

// For traversal the map
func (m *ArrayMap) For(visitor func(key int, value T) (broken bool)) {
	for key, value := range m.data {
		if m.validFlags.Get(key) {
			if visitor(key, value) {
				break
			}
		}
	}
}
