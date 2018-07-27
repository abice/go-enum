package codec

import (
	"math"
	"testing"
)

type mockWriter int

func (w mockWriter) Write(p []byte) (int, error) { return len(p), nil }
func (w mockWriter) WriteByte(b byte) error      { return nil }

func Benchmark_EncodeInt32v_Big(b *testing.B) {
	en := NewEncoder()
	w := mockWriter(0)

	for i := 0; i < b.N; i++ {
		en.EncodeInt32v(w, math.MaxInt32)
	}
}

func Benchmark_EncodeInt64v_Big(b *testing.B) {
	en := NewEncoder()
	w := mockWriter(0)

	for i := 0; i < b.N; i++ {
		en.EncodeInt64v(w, math.MaxInt64)
	}
}

func Benchmark_EncodeInt32f_Big(b *testing.B) {
	en := NewEncoder()
	w := mockWriter(0)

	for i := 0; i < b.N; i++ {
		en.EncodeInt32f(w, math.MaxInt32)
	}
}

func Benchmark_EncodeInt64f_Big(b *testing.B) {
	en := NewEncoder()
	w := mockWriter(0)

	for i := 0; i < b.N; i++ {
		en.EncodeInt64f(w, math.MaxInt64)
	}
}

func Benchmark_EncodeInt32v_Small(b *testing.B) {
	en := NewEncoder()
	w := mockWriter(0)

	for i := 0; i < b.N; i++ {
		en.EncodeInt32v(w, 1)
	}
}

func Benchmark_EncodeInt64v_Small(b *testing.B) {
	en := NewEncoder()
	w := mockWriter(0)

	for i := 0; i < b.N; i++ {
		en.EncodeInt64v(w, 1)
	}
}

func Benchmark_EncodeInt32f_Small(b *testing.B) {
	en := NewEncoder()
	w := mockWriter(0)

	for i := 0; i < b.N; i++ {
		en.EncodeInt32f(w, 1)
	}
}

func Benchmark_EncodeInt64f_Small(b *testing.B) {
	en := NewEncoder()
	w := mockWriter(0)

	for i := 0; i < b.N; i++ {
		en.EncodeInt64f(w, 1)
	}
}
