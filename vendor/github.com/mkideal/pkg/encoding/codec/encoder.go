package codec

import (
	"io"
	"math"
)

// Implements Encoder
type zencoder struct{}

func (en zencoder) EncodeInt8(w Writer, v int8) (n int, err error) {
	err = w.WriteByte(byte(v))
	if err == nil {
		n = 1
	}
	return
}

func (en zencoder) EncodeInt16f(w Writer, v int16) (n int, err error) {
	return en.EncodeUint16f(w, uint16(v))
}

func (en zencoder) EncodeInt32f(w Writer, v int32) (n int, err error) {
	return en.writeUint64(w, 4, uint64(v))
}

func (en zencoder) EncodeInt64f(w Writer, v int64) (n int, err error) {
	return en.writeUint64(w, 8, uint64(v))
}

func (en zencoder) EncodeUint8(w Writer, v uint8) (n int, err error) {
	err = w.WriteByte(byte(v))
	if err == nil {
		n = 1
	}
	return
}

func (en zencoder) EncodeUint16f(w Writer, v uint16) (n int, err error) {
	err = w.WriteByte(byte(v >> 8))
	if err != nil {
		return
	}
	n += 1
	err = w.WriteByte(byte(v & 0x00FF))
	if err != nil {
		return
	}
	n += 1
	return
}

func (en zencoder) EncodeUint32f(w Writer, v uint32) (n int, err error) {
	return en.writeUint64(w, 4, uint64(v))
}

func (en zencoder) EncodeUint64f(w Writer, v uint64) (n int, err error) {
	return en.writeUint64(w, 8, uint64(v))
}

func (en zencoder) EncodeInt16v(w Writer, v int16) (n int, err error) {
	if v < 0 {
		return en.zencodeUint64v(w, false, uint64(-int32(v)))
	}
	return en.zencodeUint64v(w, true, uint64(v))
}

func (en zencoder) EncodeInt32v(w Writer, v int32) (n int, err error) {
	if v < 0 {
		return en.zencodeUint64v(w, false, uint64(-int64(v)))
	}
	return en.zencodeUint64v(w, true, uint64(v))
}

func (en zencoder) EncodeInt64v(w Writer, v int64) (n int, err error) {
	if v < 0 {
		return en.zencodeUint64v(w, false, uint64(-v))
	}
	return en.zencodeUint64v(w, true, uint64(v))
}

func (en zencoder) EncodeUint16v(w Writer, v uint16) (n int, err error) {
	return en.zencodeUint64v(w, true, uint64(v))
}

func (en zencoder) EncodeUint32v(w Writer, v uint32) (n int, err error) {
	return en.zencodeUint64v(w, true, uint64(v))
}

func (en zencoder) EncodeUint64v(w Writer, v uint64) (n int, err error) {
	return en.zencodeUint64v(w, true, uint64(v))
}

func (en zencoder) EncodeFloat32(w Writer, v float32) (n int, err error) {
	return en.EncodeUint32f(w, math.Float32bits(v))
}

func (en zencoder) EncodeFloat64(w Writer, v float64) (n int, err error) {
	return en.EncodeUint64f(w, math.Float64bits(v))
}

func (en zencoder) EncodeBool(w Writer, v bool) (n int, err error) {
	if v {
		return en.EncodeUint8(w, 1)
	}
	return en.EncodeUint8(w, 0)
}

func (en zencoder) EncodeString(w Writer, v string) (n int, err error) {
	l := uint32(len(v))
	n, err = en.EncodeUint32v(w, l)
	if err != nil {
		return
	}
	n2, err2 := io.WriteString(w, v)
	n += n2
	err = err2
	return
}

func (en zencoder) EncodeBytes(w Writer, v []byte) (n int, err error) {
	l := uint32(len(v))
	n, err = en.EncodeUint32v(w, l)
	if err != nil {
		return
	}
	n2, err2 := w.Write(v)
	n += n2
	err = err2
	return
}

func (en zencoder) zencodeUint64v(w Writer, positive bool, v uint64) (n int, err error) {
	var byteSize uint64
	switch {
	case v <= i06:
		if positive {
			v |= m06
		} else {
			v |= u06
		}
		byteSize = 1
	case v <= i13:
		if positive {
			v |= m13
		} else {
			v |= u13
		}
		byteSize = 2
	case v <= i20:
		if positive {
			v |= m20
		} else {
			v |= u20
		}
		byteSize = 3
	case v <= i27:
		if positive {
			v |= m27
		} else {
			v |= u27
		}
		byteSize = 4
	case v <= i34:
		if positive {
			v |= m34
		} else {
			v |= u34
		}
		byteSize = 5
	case v <= i41:
		if positive {
			v |= m41
		} else {
			v |= u41
		}
		byteSize = 6
	case v <= i48:
		if positive {
			v |= m48
		} else {
			v |= u48
		}
		byteSize = 7
	case v <= i55:
		if positive {
			v |= m55
		} else {
			v |= u55
		}
		byteSize = 8
	case v <= i62:
		if positive {
			// m62=0xFF80
			err = w.WriteByte(0xFF)
			v |= 0x80 << 56
		} else {
			// u62=0x0040
			err = w.WriteByte(0x00)
			v |= 0x40 << 56
		}
		if err != nil {
			return
		}
		n += 1
		byteSize = 8
		n2, err2 := en.writeUint64(w, byteSize, v)
		n += n2
		if err2 != nil {
			err = err2
		}
		return
	default:
		if positive {
			n, err = en.writeUint64(w, 2, m69)
		} else {
			n, err = en.writeUint64(w, 2, u69)
		}
		if err != nil {
			return
		}
		byteSize = 8
		n2, err2 := en.writeUint64(w, byteSize, v)
		n += n2
		if err2 != nil {
			err = err2
		}
		return
	}
	return en.writeUint64(w, byteSize, v)
}

func (en zencoder) writeUint64(w Writer, byteSize uint64, v uint64) (n int, err error) {
	const c = 0xFF
	for i := uint64(0); i < byteSize; i++ {
		off := (byteSize - i - 1) << 3
		b := (v & (c << off)) >> off
		err = w.WriteByte(byte(b))
		if err != nil {
			return
		}
		n++
	}
	return
}
