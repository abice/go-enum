package codec

import (
	"math"
)

// Implements Decoder
type zdecoder struct{}

func (de zdecoder) DecodeInt8(r Reader) (v int8, n int, err error) {
	var b byte
	b, err = r.ReadByte()
	if err == nil {
		n = 1
	}
	v = int8(b)
	return
}

func (de zdecoder) DecodeInt16f(r Reader) (v int16, n int, err error) {
	var uv uint16
	uv, n, err = de.DecodeUint16f(r)
	v = int16(uv)
	return
}

func (de zdecoder) DecodeInt32f(r Reader) (v int32, n int, err error) {
	var uv uint32
	uv, n, err = de.DecodeUint32f(r)
	v = int32(uv)
	return
}

func (de zdecoder) DecodeInt64f(r Reader) (v int64, n int, err error) {
	var uv uint64
	uv, n, err = de.DecodeUint64f(r)
	v = int64(uv)
	return
}

func (de zdecoder) DecodeUint8(r Reader) (v uint8, n int, err error) {
	var b byte
	b, err = r.ReadByte()
	if err == nil {
		n = 1
	}
	v = uint8(b)
	return
}

func (de zdecoder) DecodeUint16f(r Reader) (v uint16, n int, err error) {
	var uv uint64
	uv, n, err = de.readUint64(r, 2)
	v = uint16(uv)
	return
}

func (de zdecoder) DecodeUint32f(r Reader) (v uint32, n int, err error) {
	var uv uint64
	uv, n, err = de.readUint64(r, 4)
	v = uint32(uv)
	return
}

func (de zdecoder) DecodeUint64f(r Reader) (v uint64, n int, err error) {
	return de.readUint64(r, 8)
}

func (de zdecoder) DecodeInt16v(r Reader) (v int16, n int, err error) {
	var uv uint64
	uv, n, err = de.zdecodeUint64v(r)
	v = int16(int64(uv))
	return
}

func (de zdecoder) DecodeInt32v(r Reader) (v int32, n int, err error) {
	var uv uint64
	uv, n, err = de.zdecodeUint64v(r)
	v = int32(int64(uv))
	return
}

func (de zdecoder) DecodeInt64v(r Reader) (v int64, n int, err error) {
	var uv uint64
	uv, n, err = de.zdecodeUint64v(r)
	v = int64(uv)
	return
}

func (de zdecoder) DecodeUint16v(r Reader) (v uint16, n int, err error) {
	var uv uint64
	uv, n, err = de.zdecodeUint64v(r)
	v = uint16(uv)
	return
}

func (de zdecoder) DecodeUint32v(r Reader) (v uint32, n int, err error) {
	var uv uint64
	uv, n, err = de.zdecodeUint64v(r)
	v = uint32(uv)
	return
}

func (de zdecoder) DecodeUint64v(r Reader) (v uint64, n int, err error) {
	return de.zdecodeUint64v(r)
}

func (de zdecoder) DecodeFloat32(r Reader) (v float32, n int, err error) {
	b, n, err := de.DecodeUint32f(r)
	v = math.Float32frombits(b)
	return v, n, err
}

func (de zdecoder) DecodeFloat64(r Reader) (v float64, n int, err error) {
	b, n, err := de.DecodeUint64f(r)
	v = math.Float64frombits(b)
	return v, n, err
}

func (de zdecoder) DecodeBool(r Reader) (v bool, n int, err error) {
	var uv uint8
	uv, n, err = de.DecodeUint8(r)
	v = uv != 0
	return
}

func (de zdecoder) DecodeString(r Reader) (v string, n int, err error) {
	var l uint32
	var m int
	l, n, err = de.DecodeUint32v(r)
	if err != nil {
		return
	}
	if l > MaxLength {
		err = errLengthTooMax
	}
	p := make([]byte, l)
	m, err = r.Read(p)
	if err != nil {
		return
	}
	n += m
	v = string(p)
	return
}

func (de zdecoder) DecodeBytes(r Reader) (v []byte, n int, err error) {
	var l uint32
	var m int
	l, n, err = de.DecodeUint32v(r)
	if err != nil {
		return
	}
	if l > MaxLength {
		err = errLengthTooMax
	}
	v = make([]byte, l)
	m, err = r.Read(v)
	if err != nil {
		return
	}
	n += m
	return
}

func (de zdecoder) readUint64(r Reader, byteSize uint64) (v uint64, n int, err error) {
	var b byte
	for i := uint64(0); i < byteSize; i++ {
		b, err = r.ReadByte()
		if err != nil {
			return
		}
		n++
		v |= uint64(b) << ((byteSize - i - 1) << 3)
	}
	return
}

func (de zdecoder) zdecodeUint64v(r Reader) (v uint64, n int, err error) {
	var (
		positive bool
		b        byte
		byteSize = 1
	)
	// read first byte
	b, err = r.ReadByte()
	if err != nil {
		return
	}
	n++
	positive = b&0x80 != 0
	// read byte size
	index := 1
	for {
		fromIndex := index
		index = readBit(b, uint(fromIndex), positive)
		byteSize += index - fromIndex
		if index < 8 {
			break
		}
		b, err = r.ReadByte()
		if err != nil {
			return
		}
		n++
		// reset index
		index = 0
	}
	remainByteSize := uint64(byteSize - n)
	if index <= 6 {
		off := uint64(index) + 1
		v |= uint64(b<<off>>off) << (remainByteSize << 3)
	}
	// read value
	for i := uint64(0); i < remainByteSize; i++ {
		b, err = r.ReadByte()
		if err != nil {
			return
		}
		n++
		v |= uint64(b) << ((remainByteSize - i - 1) << 3)
	}
	if !positive {
		v = uint64(-v)
	}
	return
}

func readBit(b byte, fromIndex uint, positive bool) int {
	for i := fromIndex; i < 8; i++ {
		v := b & (0x80 >> i)
		if positive == (v == 0) {
			return int(i)
		}
	}
	return 8
}
