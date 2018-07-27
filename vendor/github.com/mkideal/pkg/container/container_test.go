package container

import (
	"math/rand"
	"sort"
	"testing"
	"time"
)

const N = 10000
const N100 = N * 100

var s1 = func() []int {
	rand.Seed(time.Now().UnixNano())
	s := make([]int, N)
	for i := 0; i < N; i++ {
		s[i] = rand.Intn(N100)
	}
	return s
}()

var s2 = func() []int {
	s := make([]int, len(s1))
	copy(s, s1)
	return s
}()

func TestIterSort(t *testing.T) {
	s := []int{4, 1, 3, 5, 6, 83, 2}
	iter := IterFromSlice(s)
	Sort(iter, iter.Add(3), func(l, r interface{}) bool {
		return l.(int) < r.(int)
	}, s)
	t.Logf("sorted slice: %v", s)
}

func Benchmark_StdSort(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sort.Ints(s1)
	}
}

func Benchmark_IterSort(b *testing.B) {
	less := func(l, r interface{}) bool { return l.(int) < r.(int) }
	for i := 0; i < b.N; i++ {
		iter := IterFromSlice(s2)
		Sort(iter, iter.Add(len(s2)), less, s2)
	}
}
