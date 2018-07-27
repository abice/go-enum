package random

import (
	"bytes"
	crand "crypto/rand"
	"encoding/binary"
	"math/rand"
	"reflect"
	"sort"
	"time"
)

var (
	digits         = []byte("0123456789")
	lowercaseChars = []byte("abcdefghijklmnopqrstuvwxyz")
	uppercaseChars = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	specialChars   = []byte("~!@#$%^&*")
)

const (
	O_DIGIT = 1 << iota
	O_LOWER_CHAR
	O_UPPER_CHAR
	O_SPECIAL_CHAR
)

type Source interface {
	Int63() int64
}

// defaultSource represents a unsafe random source
type defaultSource struct{}

func (source defaultSource) Int63() int64 { return rand.Int63() }

func newDefaultSource(seed int64) defaultSource {
	rand.Seed(seed)
	return defaultSource{}
}

// defaultSource represents a safe random source
type cryptoSource struct{}

func (source cryptoSource) Int63() int64 {
	var b [8]byte
	_, err := crand.Read(b[:])
	if err != nil {
		return DefaultSource.Int63()
	}
	b[0] &= 0x7F
	return int64(binary.BigEndian.Uint64(b[:]))
}

var (
	DefaultSource Source = newDefaultSource(time.Now().UnixNano())
	CryptoSource  Source = cryptoSource{}
)

func Int63(source Source) int64 {
	if source == nil {
		source = DefaultSource
	}
	return source.Int63()
}

func Intn(n int, source Source) int {
	return int(Int63(source) % int64(n))
}

func Bool(source Source) bool {
	return Intn(2, source) == 1
}

func String(length int, source Source, modes ...int) string {
	if length <= 0 {
		return ""
	}
	var mode int
	for _, m := range modes {
		mode |= m
	}
	if mode&(O_DIGIT|O_LOWER_CHAR|O_UPPER_CHAR|O_SPECIAL_CHAR) == 0 {
		mode = O_LOWER_CHAR | O_UPPER_CHAR
	}
	size := 0
	if mode&O_DIGIT != 0 {
		size += len(digits)
	}
	if mode&O_LOWER_CHAR != 0 {
		size += len(lowercaseChars)
	}
	if mode&O_UPPER_CHAR != 0 {
		size += len(uppercaseChars)
	}
	if mode&O_SPECIAL_CHAR != 0 {
		size += len(specialChars)
	}
	var buf bytes.Buffer
	for i := 0; i < length; i++ {
		index := Intn(size, source)
		tmpSize := 0
		if mode&O_DIGIT != 0 {
			tmpSize = len(digits)
			if index < tmpSize {
				buf.WriteByte(digits[index])
				continue
			}
			index -= tmpSize
		}
		if mode&O_LOWER_CHAR != 0 {
			tmpSize = len(lowercaseChars)
			if index < tmpSize {
				buf.WriteByte(lowercaseChars[index])
				continue
			}
			index -= tmpSize
		}
		if mode&O_UPPER_CHAR != 0 {
			tmpSize = len(uppercaseChars)
			if index < tmpSize {
				buf.WriteByte(uppercaseChars[index])
				continue
			}
			index -= tmpSize
		}
		if mode&O_SPECIAL_CHAR != 0 {
			tmpSize = len(specialChars)
			if index < tmpSize {
				buf.WriteByte(specialChars[index])
				continue
			}
		}
		buf.WriteByte('-')
	}
	return buf.String()
}

func Perm(n int, source Source) []int {
	s := make([]int, n)
	for i := 0; i < n; i++ {
		j := Intn(i+1, source)
		s[i] = s[j]
		s[j] = i
	}
	return s
}

type SwapableSlice interface {
	Len() int
	Swap(i, j int)
}

func Shuffle(orders SwapableSlice, source Source) {
	for i := orders.Len() - 1; i >= 0; i-- {
		if source == nil {
			orders.Swap(i, rand.Intn(i+1))
		} else {
			orders.Swap(i, int(source.Int63())%(i+1))
		}
	}
}

func ShuffleInts(orders []int, source Source)       { Shuffle(sort.IntSlice(orders), source) }
func ShuffleFloats(orders []float64, source Source) { Shuffle(sort.Float64Slice(orders), source) }
func ShuffleStrings(orders []string, source Source) { Shuffle(sort.StringSlice(orders), source) }

type swapableSlice struct {
	swapper func(int, int)
	length  int
}

func (s swapableSlice) Len() int      { return s.length }
func (s swapableSlice) Swap(i, j int) { s.swapper(i, j) }

func ShuffleSlice(slice interface{}, source Source) {
	Shuffle(swapableSlice{
		swapper: reflect.Swapper(slice),
		length:  reflect.ValueOf(slice).Len(),
	}, source)
}

// FiniteDistribution represents probability distribution
type FiniteDistribution interface {
	Len() int
	Probability(i int) int
}

type SummedFiniteDistribution interface {
	FiniteDistribution
	SumProbability() int
}

type IntsFiniteDistribution []int

func (d IntsFiniteDistribution) Len() int              { return len(d) }
func (d IntsFiniteDistribution) Probability(i int) int { return d[i] }

func sumFiniteDistribution(d FiniteDistribution) int {
	if summed, ok := d.(SummedFiniteDistribution); ok {
		return summed.SumProbability()
	}
	sum := 0
	for i, n := 0, d.Len(); i < n; i++ {
		sum += d.Probability(i)
	}
	return sum
}

func Index(d FiniteDistribution, source Source) int {
	sum := sumFiniteDistribution(d)
	if sum <= 0 {
		panic("sum < 0")
	}
	var (
		value = Intn(sum, source)
		acc   = 0
		end   = d.Len() - 1
	)
	for i := 0; i < end; i++ {
		acc += d.Probability(i)
		if acc > value {
			return i
		}
	}
	return end
}

func IndexInts(d []int, source Source) int {
	return Index(IntsFiniteDistribution(d), source)
}

type probabilityDistribution struct {
	probability func(int) int
	length      int
}

func (d probabilityDistribution) Len() int              { return d.length }
func (d probabilityDistribution) Probability(i int) int { return d.probability(i) }

func IndexSlice(d interface{}, probability func(int) int, source Source) int {
	return Index(probabilityDistribution{probability, reflect.ValueOf(d).Len()}, source)
}
