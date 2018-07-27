package container

import (
	"reflect"
)

type Container interface {
	Len() int
	Iter() Iterator
	Contains(interface{}) bool
}

type Array interface {
	Container
}

type Map interface {
	Container
}

// Iterator ...
type Iterator interface {
	Next() (k, v interface{})
}

// emptyIterator
type emptyIterator struct{}

func (i emptyIterator) Next() (interface{}, interface{}) { return nil, nil }

var EmptyIterator = emptyIterator{}

type RandomAccessIterator interface {
	Iterator
	Add(n int) RandomAccessIterator
	Diff(iter RandomAccessIterator) int
	Get() interface{}
	Set(v interface{})

	Value(i int) reflect.Value
}

type iteratorSorter struct {
	less       CompareFunc
	begin, end RandomAccessIterator
	native     []int
}

func (sorter *iteratorSorter) Len() int { return sorter.end.Diff(sorter.begin) }
func (sorter *iteratorSorter) Less(i, j int) bool {
	return sorter.less(sorter.native[i], sorter.native[j])
	//return sorter.less(sorter.begin.Value(i).Interface(), sorter.begin.Value(j).Interface())
}
func (sorter *iteratorSorter) Swap(i, j int) {
	sorter.native[i], sorter.native[j] = sorter.native[j], sorter.native[i]
	//vi, vj := sorter.begin.Value(i), sorter.begin.Value(j)
	//tmp := vi.Interface()
	//vi.Set(reflect.ValueOf(vj.Interface()))
	//vj.Set(reflect.ValueOf(tmp))
}

type sliceIterator struct {
	slice        reflect.Value
	currentIndex int
}

// $zcheck: slice IS SLICE
func IterFromSlice(slice interface{}) RandomAccessIterator {
	value := reflect.ValueOf(slice)
	return &sliceIterator{
		slice:        value,
		currentIndex: 0,
	}
}

func (iter *sliceIterator) Next() (k, v interface{}) {
	if iter.currentIndex >= iter.slice.Len() {
		return nil, nil
	}
	v = iter.slice.Index(iter.currentIndex).Interface()
	k = iter.currentIndex
	iter.currentIndex++
	return
}

func (iter sliceIterator) Add(n int) RandomAccessIterator {
	if iter.currentIndex+n < 0 {
		panic("index < 0")
	}
	return &sliceIterator{
		slice:        iter.slice,
		currentIndex: iter.currentIndex + n,
	}
}

func (iter *sliceIterator) AddWith(n int) RandomAccessIterator {
	if iter.currentIndex+n < 0 {
		panic("index < 0")
	}
	iter.currentIndex += n
	return iter
}

func (iter sliceIterator) Diff(iter2 RandomAccessIterator) int {
	sliceIter2, ok := iter2.(*sliceIterator)
	if !ok {
		panic("diff a different type iterator")
	}
	return iter.currentIndex - sliceIter2.currentIndex
}

func (iter sliceIterator) Get() interface{} {
	if iter.currentIndex >= iter.slice.Len() {
		return nil
	}
	return iter.slice.Index(iter.currentIndex).Interface()
}

func (iter sliceIterator) Value(i int) reflect.Value {
	index := iter.currentIndex + i
	return iter.slice.Index(index)
}

func (iter *sliceIterator) Set(v interface{}) {
	if iter.currentIndex >= iter.slice.Len() {
		panic("out of range")
	}
	iter.slice.Index(iter.currentIndex).Set(reflect.ValueOf(v))
}
