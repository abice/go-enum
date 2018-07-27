package container

import (
	"sort"
)

func Len(c Container) int                      { return c.Len() }
func Contains(c Container, v interface{}) bool { return c.Contains(v) }

type CompareFunc func(left, right interface{}) bool

func Sort(begin, end RandomAccessIterator, less CompareFunc, native []int) {
	sort.Sort(&iteratorSorter{less, begin, end, native})
}
