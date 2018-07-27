package codec

import (
	"bytes"
	"math"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEncoderFixedLengthInteger(t *testing.T) {
	w := new(bytes.Buffer)

	en := NewEncoder()

	en.EncodeBool(w, true)
	assert.Equal(t, []byte{1}, w.Bytes())
	w.Reset()

	en.EncodeBool(w, false)
	assert.Equal(t, []byte{0}, w.Bytes())
	w.Reset()

	en.EncodeInt8(w, math.MaxInt8)
	assert.Equal(t, []byte{0x7F}, w.Bytes())
	w.Reset()

	en.EncodeInt8(w, math.MinInt8)
	assert.Equal(t, []byte{0x80}, w.Bytes())
	w.Reset()

	en.EncodeInt16f(w, math.MaxInt16)
	assert.Equal(t, []byte{0x7F, 0xFF}, w.Bytes())
	w.Reset()

	en.EncodeInt16f(w, math.MinInt16)
	assert.Equal(t, []byte{0x80, 0x00}, w.Bytes())
	w.Reset()

	en.EncodeInt32f(w, math.MaxInt32)
	assert.Equal(t, []byte{0x7F, 0xFF, 0xFF, 0xFF}, w.Bytes())
	w.Reset()

	en.EncodeInt32f(w, math.MinInt32)
	assert.Equal(t, []byte{0x80, 0x00, 0x00, 0x00}, w.Bytes())
	w.Reset()

	en.EncodeInt64f(w, math.MaxInt64)
	assert.Equal(t, []byte{0x7F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, w.Bytes())
	w.Reset()

	en.EncodeInt64f(w, math.MinInt64)
	assert.Equal(t, []byte{0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, w.Bytes())
	w.Reset()

	en.EncodeUint8(w, math.MaxUint8)
	assert.Equal(t, []byte{0xFF}, w.Bytes())
	w.Reset()

	en.EncodeUint16f(w, math.MaxUint16)
	assert.Equal(t, []byte{0xFF, 0xFF}, w.Bytes())
	w.Reset()

	en.EncodeUint32f(w, math.MaxUint32)
	assert.Equal(t, []byte{0xFF, 0xFF, 0xFF, 0xFF}, w.Bytes())
	w.Reset()

	en.EncodeUint64f(w, math.MaxUint64)
	assert.Equal(t, []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, w.Bytes())
	w.Reset()
}

func TestEncoderVariableLengthInteger(t *testing.T) {
	w := new(bytes.Buffer)

	en := NewEncoder()

	//-------
	// int16/uint16
	//-------
	en.EncodeInt16v(w, i06)
	assert.Equal(t, []byte{0xBF}, w.Bytes())
	w.Reset()

	en.EncodeInt16v(w, -i06)
	assert.Equal(t, []byte{0x7F}, w.Bytes())
	w.Reset()

	en.EncodeInt16v(w, i13)
	assert.Equal(t, []byte{0xDF, 0xFF}, w.Bytes())
	w.Reset()

	en.EncodeInt16v(w, -i13)
	assert.Equal(t, []byte{0x3F, 0xFF}, w.Bytes())
	w.Reset()

	en.EncodeInt16v(w, math.MaxInt16)
	assert.Equal(t, []byte{0xE0, 0x7F, 0xFF}, w.Bytes())
	w.Reset()

	en.EncodeInt16v(w, math.MinInt16)
	assert.Equal(t, []byte{0x10, 0x80, 0x00}, w.Bytes())
	w.Reset()

	en.EncodeUint16v(w, i06)
	assert.Equal(t, []byte{0xBF}, w.Bytes())
	w.Reset()

	en.EncodeUint16v(w, i13)
	assert.Equal(t, []byte{0xDF, 0xFF}, w.Bytes())
	w.Reset()

	en.EncodeUint16v(w, math.MaxUint16)
	assert.Equal(t, []byte{0xE0, 0xFF, 0xFF}, w.Bytes())
	w.Reset()

	//-------
	// int32
	//-------
	en.EncodeInt32v(w, i06)
	assert.Equal(t, []byte{0xBF}, w.Bytes())
	w.Reset()

	en.EncodeInt32v(w, -i06)
	assert.Equal(t, []byte{0x7F}, w.Bytes())
	w.Reset()

	en.EncodeInt32v(w, i13)
	assert.Equal(t, []byte{0xDF, 0xFF}, w.Bytes())
	w.Reset()

	en.EncodeInt32v(w, -i13)
	assert.Equal(t, []byte{0x3F, 0xFF}, w.Bytes())
	w.Reset()

	en.EncodeInt32v(w, i20)
	assert.Equal(t, []byte{0xEF, 0xFF, 0xFF}, w.Bytes())
	w.Reset()

	en.EncodeInt32v(w, -i20)
	assert.Equal(t, []byte{0x1f, 0xFF, 0xFF}, w.Bytes())
	w.Reset()

	en.EncodeInt32v(w, i27)
	assert.Equal(t, []byte{0xF7, 0xFF, 0xFF, 0xFF}, w.Bytes())
	w.Reset()

	en.EncodeInt32v(w, -i27)
	assert.Equal(t, []byte{0x0F, 0xFF, 0xFF, 0xFF}, w.Bytes())
	w.Reset()

	en.EncodeInt32v(w, math.MaxInt32)
	assert.Equal(t, []byte{0xF8, 0x7F, 0xFF, 0xFF, 0xFF}, w.Bytes())
	w.Reset()

	en.EncodeInt32v(w, math.MinInt32)
	assert.Equal(t, []byte{0x04, 0x80, 0x00, 0x00, 0x00}, w.Bytes())
	w.Reset()

	en.EncodeUint32v(w, i06)
	assert.Equal(t, []byte{0xBF}, w.Bytes())
	w.Reset()

	en.EncodeUint32v(w, i13)
	assert.Equal(t, []byte{0xDF, 0xFF}, w.Bytes())
	w.Reset()

	en.EncodeUint32v(w, i20)
	assert.Equal(t, []byte{0xEF, 0xFF, 0xFF}, w.Bytes())
	w.Reset()

	en.EncodeUint32v(w, i27)
	assert.Equal(t, []byte{0xF7, 0xFF, 0xFF, 0xFF}, w.Bytes())
	w.Reset()

	en.EncodeUint32v(w, math.MaxUint32)
	assert.Equal(t, []byte{0xF8, 0xFF, 0xFF, 0xFF, 0xFF}, w.Bytes())
	w.Reset()

	//-------
	// int64
	//-------
	en.EncodeInt64v(w, i06)
	assert.Equal(t, []byte{0xBF}, w.Bytes())
	w.Reset()

	en.EncodeInt64v(w, -i06)
	assert.Equal(t, []byte{0x7F}, w.Bytes())
	w.Reset()

	en.EncodeInt64v(w, i13)
	assert.Equal(t, []byte{0xDF, 0xFF}, w.Bytes())
	w.Reset()

	en.EncodeInt64v(w, -i13)
	assert.Equal(t, []byte{0x3F, 0xFF}, w.Bytes())
	w.Reset()

	en.EncodeInt64v(w, i20)
	assert.Equal(t, []byte{0xEF, 0xFF, 0xFF}, w.Bytes())
	w.Reset()

	en.EncodeInt64v(w, -i20)
	assert.Equal(t, []byte{0x1f, 0xFF, 0xFF}, w.Bytes())
	w.Reset()

	en.EncodeInt64v(w, i27)
	assert.Equal(t, []byte{0xF7, 0xFF, 0xFF, 0xFF}, w.Bytes())
	w.Reset()

	en.EncodeInt64v(w, -i27)
	assert.Equal(t, []byte{0x0F, 0xFF, 0xFF, 0xFF}, w.Bytes())
	w.Reset()

	en.EncodeInt64v(w, i34)
	assert.Equal(t, []byte{0xFB, 0xFF, 0xFF, 0xFF, 0xFF}, w.Bytes())
	w.Reset()

	en.EncodeInt64v(w, -i34)
	assert.Equal(t, []byte{0x07, 0xFF, 0xFF, 0xFF, 0xFF}, w.Bytes())
	w.Reset()

	en.EncodeInt64v(w, i41)
	assert.Equal(t, []byte{0xFD, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, w.Bytes())
	w.Reset()

	en.EncodeInt64v(w, -i41)
	assert.Equal(t, []byte{0x03, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, w.Bytes())
	w.Reset()

	en.EncodeInt64v(w, i48)
	assert.Equal(t, []byte{0xFE, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, w.Bytes())
	w.Reset()

	en.EncodeInt64v(w, -i48)
	assert.Equal(t, []byte{0x01, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, w.Bytes())
	w.Reset()

	en.EncodeInt64v(w, i55)
	assert.Equal(t, []byte{0xFF, 0x7F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, w.Bytes())
	w.Reset()

	en.EncodeInt64v(w, -i55)
	assert.Equal(t, []byte{0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, w.Bytes())
	w.Reset()

	en.EncodeInt64v(w, i62)
	assert.Equal(t, []byte{0xFF, 0xBF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, w.Bytes())
	w.Reset()

	en.EncodeInt64v(w, -i62)
	assert.Equal(t, []byte{0x00, 0x7F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, w.Bytes())
	w.Reset()

	en.EncodeInt64v(w, math.MaxInt64)
	assert.Equal(t, []byte{0xFF, 0xC0, 0x7F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, w.Bytes())
	w.Reset()

	en.EncodeInt64v(w, math.MinInt64)
	assert.Equal(t, []byte{0x00, 0x20, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, w.Bytes())
	w.Reset()

	en.EncodeUint64v(w, i06)
	assert.Equal(t, []byte{0xBF}, w.Bytes())
	w.Reset()

	en.EncodeUint64v(w, i13)
	assert.Equal(t, []byte{0xDF, 0xFF}, w.Bytes())
	w.Reset()

	en.EncodeUint64v(w, i20)
	assert.Equal(t, []byte{0xEF, 0xFF, 0xFF}, w.Bytes())
	w.Reset()

	en.EncodeUint64v(w, i27)
	assert.Equal(t, []byte{0xF7, 0xFF, 0xFF, 0xFF}, w.Bytes())
	w.Reset()

	en.EncodeUint64v(w, i34)
	assert.Equal(t, []byte{0xFB, 0xFF, 0xFF, 0xFF, 0xFF}, w.Bytes())
	w.Reset()

	en.EncodeUint64v(w, i41)
	assert.Equal(t, []byte{0xFD, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, w.Bytes())
	w.Reset()

	en.EncodeUint64v(w, i48)
	assert.Equal(t, []byte{0xFE, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, w.Bytes())
	w.Reset()

	en.EncodeUint64v(w, i55)
	assert.Equal(t, []byte{0xFF, 0x7F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, w.Bytes())
	w.Reset()

	en.EncodeUint64v(w, i62)
	assert.Equal(t, []byte{0xFF, 0xBF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, w.Bytes())
	w.Reset()

	en.EncodeUint64v(w, math.MaxUint64)
	assert.Equal(t, []byte{0xFF, 0xC0, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, w.Bytes())
	w.Reset()
}

func TestDecoderFixedLengthInteger(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	wr := new(bytes.Buffer)
	en := NewEncoder()
	de := NewDecoder()

	const SIZE = 1000

	// uint-fixed
	for i := 0; i < SIZE; i++ {
		value := uint8(rand.Intn(math.MaxInt8))
		n, err := en.EncodeUint8(wr, value)
		assert.Nil(t, err)
		v2, n2, err2 := de.DecodeUint8(wr)
		assert.Nil(t, err2)
		assert.Equal(t, n, n2)
		assert.Equal(t, value, v2)
	}

	for i := 0; i < SIZE; i++ {
		value := uint16(rand.Intn(math.MaxInt16))
		n, err := en.EncodeUint16f(wr, value)
		assert.Nil(t, err)
		v2, n2, err2 := de.DecodeUint16f(wr)
		assert.Nil(t, err2)
		assert.Equal(t, n, n2)
		assert.Equal(t, value, v2)
	}

	for i := 0; i < SIZE; i++ {
		value := uint32(rand.Intn(math.MaxInt32))
		n, err := en.EncodeUint32f(wr, value)
		assert.Nil(t, err)
		v2, n2, err2 := de.DecodeUint32f(wr)
		assert.Nil(t, err2)
		assert.Equal(t, n, n2)
		assert.Equal(t, value, v2)
	}

	for i := 0; i < SIZE; i++ {
		value := uint64(rand.Intn(math.MaxInt64))
		n, err := en.EncodeUint64f(wr, value)
		assert.Nil(t, err)
		v2, n2, err2 := de.DecodeUint64f(wr)
		assert.Nil(t, err2)
		assert.Equal(t, n, n2)
		assert.Equal(t, value, v2)
	}

	// int-fixed
	for i := 0; i < SIZE; i++ {
		value := int8(rand.Intn(math.MaxInt8))
		n, err := en.EncodeInt8(wr, value)
		assert.Nil(t, err)
		v2, n2, err2 := de.DecodeInt8(wr)
		assert.Nil(t, err2)
		assert.Equal(t, n, n2)
		assert.Equal(t, value, v2)
	}

	for i := 0; i < SIZE; i++ {
		value := int16(rand.Intn(math.MaxInt16))
		n, err := en.EncodeInt16f(wr, value)
		assert.Nil(t, err)
		v2, n2, err2 := de.DecodeInt16f(wr)
		assert.Nil(t, err2)
		assert.Equal(t, n, n2)
		assert.Equal(t, value, v2)
	}

	for i := 0; i < SIZE; i++ {
		value := int32(rand.Intn(math.MaxInt32))
		n, err := en.EncodeInt32f(wr, value)
		assert.Nil(t, err)
		v2, n2, err2 := de.DecodeInt32f(wr)
		assert.Nil(t, err2)
		assert.Equal(t, n, n2)
		assert.Equal(t, value, v2)
	}

	for i := 0; i < SIZE; i++ {
		value := int64(rand.Intn(math.MaxInt64))
		n, err := en.EncodeInt64f(wr, value)
		assert.Nil(t, err)
		v2, n2, err2 := de.DecodeInt64f(wr)
		assert.Nil(t, err2)
		assert.Equal(t, n, n2)
		assert.Equal(t, value, v2)
	}

	// uint-variable
	for i := 0; i < SIZE; i++ {
		value := uint16(rand.Intn(math.MaxInt16))
		n, err := en.EncodeUint16v(wr, value)
		assert.Nil(t, err)
		v2, n2, err2 := de.DecodeUint16v(wr)
		assert.Nil(t, err2)
		assert.Equal(t, n, n2)
		assert.Equal(t, value, v2)
	}

	for i := 0; i < SIZE; i++ {
		value := uint32(rand.Intn(math.MaxInt32))
		n, err := en.EncodeUint32v(wr, value)
		assert.Nil(t, err)
		v2, n2, err2 := de.DecodeUint32v(wr)
		assert.Nil(t, err2)
		assert.Equal(t, n, n2)
		assert.Equal(t, value, v2)
	}

	for i := 0; i < SIZE; i++ {
		value := uint64(rand.Intn(math.MaxInt64))
		n, err := en.EncodeUint64v(wr, value)
		assert.Nil(t, err)
		v2, n2, err2 := de.DecodeUint64v(wr)
		assert.Nil(t, err2)
		assert.Equal(t, n, n2)
		assert.Equal(t, value, v2)
	}

	// int-variable
	for i := 0; i < SIZE; i++ {
		v := int16(rand.Intn(math.MaxInt16))
		for _, value := range []int16{v, -v} {
			n, err := en.EncodeInt16v(wr, value)
			assert.Nil(t, err)
			v2, n2, err2 := de.DecodeInt16v(wr)
			assert.Nil(t, err2)
			assert.Equal(t, n, n2)
			assert.Equal(t, value, v2)
		}
	}

	for i := 0; i < SIZE; i++ {
		v := int32(rand.Intn(math.MaxInt32))
		for _, value := range []int32{v, -v} {
			n, err := en.EncodeInt32v(wr, value)
			assert.Nil(t, err)
			v2, n2, err2 := de.DecodeInt32v(wr)
			assert.Nil(t, err2)
			assert.Equal(t, n, n2)
			assert.Equal(t, value, v2)
		}
	}

	for i := 0; i < SIZE; i++ {
		v := int64(rand.Intn(math.MaxInt64))
		for _, value := range []int64{v, -v} {
			n, err := en.EncodeInt64v(wr, value)
			assert.Nil(t, err)
			v2, n2, err2 := de.DecodeInt64v(wr)
			assert.Nil(t, err2)
			assert.Equal(t, n, n2)
			assert.Equal(t, value, v2)
		}
	}

	// bool
	for _, value := range []bool{true, false} {
		n, err := en.EncodeBool(wr, value)
		assert.Nil(t, err)
		v2, n2, err2 := de.DecodeBool(wr)
		assert.Nil(t, err2)
		assert.Equal(t, n, n2)
		assert.Equal(t, value, v2)
	}
}

func TestEncodeDecodeStringBytes(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	wr := new(bytes.Buffer)

	tests := []struct {
		value  string
		length int
	}{
		{"hello", 6},
		{"zcode", 6},
	}

	for _, c := range tests {
		n, err := Enc.EncodeString(wr, c.value)
		assert.Nil(t, err)
		assert.Equal(t, c.length, n)
		v, n, err := Dec.DecodeString(wr)
		assert.Nil(t, err)
		assert.Equal(t, c.length, n)
		assert.Equal(t, c.value, v)

		n, err = Enc.EncodeBytes(wr, []byte(c.value))
		assert.Nil(t, err)
		assert.Equal(t, c.length, n)
		b, n, err := Dec.DecodeBytes(wr)
		assert.Nil(t, err)
		assert.Equal(t, c.length, n)
		assert.Equal(t, c.value, string(b))
	}

	randString := func() string {
		buf := new(bytes.Buffer)
		n := rand.Intn(10000) + 1000
		for i := 0; i < n; i++ {
			buf.WriteByte(byte(rand.Intn(128)))
		}
		return buf.String()
	}

	for i := 0; i < 100; i++ {
		s := randString()
		n, err := Enc.EncodeString(wr, s)
		assert.Nil(t, err)
		v, n2, err := Dec.DecodeString(wr)
		assert.Nil(t, err)
		assert.Equal(t, n, n2)
		assert.Equal(t, s, v)

		n, err = Enc.EncodeBytes(wr, []byte(s))
		assert.Nil(t, err)
		b, n2, err := Dec.DecodeBytes(wr)
		assert.Nil(t, err)
		assert.Equal(t, n, n2)
		assert.Equal(t, s, string(b))
	}
}

func TestEncodeDecodeFloat(t *testing.T) {
	wr := new(bytes.Buffer)
	for _, v := range []float32{
		0.0,
		1.2,
		1234.5678,
		-1234.5678,
		-1.2,
		math.MaxFloat32,
		-math.MaxFloat32,
		math.SmallestNonzeroFloat32,
		-math.SmallestNonzeroFloat32,
	} {
		n, err := Enc.EncodeFloat32(wr, v)
		assert.Nil(t, err)
		v2, n2, err2 := Dec.DecodeFloat32(wr)
		assert.Nil(t, err2)
		assert.Equal(t, n, n2)
		assert.Equal(t, v, v2)
	}
	for _, v := range []float64{
		0.0,
		1.2,
		1234.5678,
		-1234.5678,
		-1.2,
		math.MaxFloat64,
		-math.MaxFloat64,
		math.SmallestNonzeroFloat64,
		-math.SmallestNonzeroFloat64,
	} {
		n, err := Enc.EncodeFloat64(wr, v)
		assert.Nil(t, err)
		v2, n2, err2 := Dec.DecodeFloat64(wr)
		assert.Nil(t, err2)
		assert.Equal(t, n, n2)
		assert.Equal(t, v, v2)
	}
}
