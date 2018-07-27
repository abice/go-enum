package codec

import (
	"errors"
	"io"
)

const Unused = 0

var (
	ErrNegativeLength = errors.New("negative length")
	ErrTooBigLength   = errors.New("too big length")
)

type Dispatcher interface {
	Dispatch(id uint32, body Reader) (int, error)
}

type Message interface {
	MessageId() int
	MessageName() string
	Encode(Writer) (int, error)
	Decode(Reader) (int, error)
}

// sizedReader implements Reader interface
type sizedReader struct {
	reader Reader
	size   int
}

var EOF = errors.New("eof")

func (sr *sizedReader) Read(p []byte) (n int, err error) {
	size := len(p)
	if size > sr.size {
		n, err = sr.reader.Read(p[:sr.size])
		if err == nil {
			err = EOF
		}
	} else {
		n, err = sr.reader.Read(p)
	}
	sr.size -= n
	return n, err
}

func (sr *sizedReader) ReadByte() (byte, error) {
	if sr.size < 1 {
		return 0, EOF
	}
	b, err := sr.reader.ReadByte()
	if err == nil {
		sr.size -= 1
	}
	return b, err
}

// MessageReader read and dispatch Message
type MessageReader struct {
	sr         *sizedReader
	dispatcher Dispatcher
}

func NewMessageReader(reader io.Reader, dispatcher Dispatcher) *MessageReader {
	mr := &MessageReader{
		sr:         &sizedReader{reader: WrapReader(reader)},
		dispatcher: dispatcher,
	}
	return mr
}

func (mr *MessageReader) Read() (int, error) {
	total := 0
	// read message length
	mr.sr.size = 4
	size, n, err := Dec.DecodeInt32f(mr.sr)
	total += n
	if err != nil {
		return total, err
	}
	if size > MaxLength {
		return total, ErrTooBigLength
	}
	mr.sr.size = int(size)
	// read message id
	id, n, err := Dec.DecodeUint32v(mr.sr)
	total += n
	if err != nil {
		return total, err
	}
	// dispatch message body by id
	n, err = mr.dispatcher.Dispatch(id, mr.sr)
	total += n
	if err != nil && err != EOF {
		return total, err
	}
	return total, nil
}

const (
	MaxLength = 1 << 20 // 1M
)

var (
	errLengthTooMax = errors.New("length too max")

	Enc = NewEncoder()
	Dec = NewDecoder()
)

func NewEncoder() Encoder { return zencoder{} }
func NewDecoder() Decoder { return zdecoder{} }

type Writer interface {
	io.Writer
	io.ByteWriter
}

type extendByteWriter struct {
	w   io.Writer
	buf [1]byte
}

func (w extendByteWriter) Write(b []byte) (int, error) { return w.w.Write(b) }
func (w extendByteWriter) WriteByte(c byte) error {
	w.buf[0] = c
	_, err := w.w.Write(w.buf[:])
	return err
}

func WrapWriter(w io.Writer) Writer {
	if wc, ok := w.(Writer); ok {
		return wc
	}
	return extendByteWriter{w: w}
}

type Reader interface {
	io.Reader
	io.ByteReader
}

type extendByteReader struct {
	r   io.Reader
	buf [1]byte
}

func (r extendByteReader) Read(p []byte) (int, error) { return r.r.Read(p) }
func (r extendByteReader) ReadByte() (byte, error) {
	_, err := r.r.Read(r.buf[:])
	return r.buf[0], err
}

func WrapReader(r io.Reader) Reader {
	if rc, ok := r.(Reader); ok {
		return rc
	}
	return extendByteReader{r: r}
}

type Encoder interface {
	// Fixed-length integer
	EncodeInt8(w Writer, v int8) (n int, err error)
	EncodeInt16f(w Writer, v int16) (n int, err error)
	EncodeInt32f(w Writer, v int32) (n int, err error)
	EncodeInt64f(w Writer, v int64) (n int, err error)
	EncodeUint8(w Writer, v uint8) (n int, err error)
	EncodeUint16f(w Writer, v uint16) (n int, err error)
	EncodeUint32f(w Writer, v uint32) (n int, err error)
	EncodeUint64f(w Writer, v uint64) (n int, err error)

	// Variable-length integer
	EncodeInt16v(w Writer, v int16) (n int, err error)
	EncodeInt32v(w Writer, v int32) (n int, err error)
	EncodeInt64v(w Writer, v int64) (n int, err error)
	EncodeUint16v(w Writer, v uint16) (n int, err error)
	EncodeUint32v(w Writer, v uint32) (n int, err error)
	EncodeUint64v(w Writer, v uint64) (n int, err error)

	EncodeFloat32(w Writer, v float32) (n int, err error)
	EncodeFloat64(w Writer, v float64) (n int, err error)
	EncodeBool(w Writer, v bool) (n int, err error)
	EncodeString(w Writer, v string) (n int, err error)
	EncodeBytes(w Writer, v []byte) (n int, err error)
}

type Decoder interface {
	// Fixed-length integer
	DecodeInt8(r Reader) (v int8, n int, err error)
	DecodeInt16f(r Reader) (v int16, n int, err error)
	DecodeInt32f(r Reader) (v int32, n int, err error)
	DecodeInt64f(r Reader) (v int64, n int, err error)
	DecodeUint8(r Reader) (v uint8, n int, err error)
	DecodeUint16f(r Reader) (v uint16, n int, err error)
	DecodeUint32f(r Reader) (v uint32, n int, err error)
	DecodeUint64f(r Reader) (v uint64, n int, err error)

	// Variable-length integer
	DecodeInt16v(r Reader) (v int16, n int, err error)
	DecodeInt32v(r Reader) (v int32, n int, err error)
	DecodeInt64v(r Reader) (v int64, n int, err error)
	DecodeUint16v(r Reader) (v uint16, n int, err error)
	DecodeUint32v(r Reader) (v uint32, n int, err error)
	DecodeUint64v(r Reader) (v uint64, n int, err error)

	DecodeFloat32(r Reader) (v float32, n int, err error)
	DecodeFloat64(r Reader) (v float64, n int, err error)
	DecodeBool(r Reader) (v bool, n int, err error)
	DecodeString(r Reader) (v string, n int, err error)
	DecodeBytes(r Reader) (v []byte, n int, err error)
}

const (
	i06 = 1<<6 - 1
	i13 = 1<<13 - 1
	i20 = 1<<20 - 1
	i27 = 1<<27 - 1
	i34 = 1<<34 - 1
	i41 = 1<<41 - 1
	i48 = 1<<48 - 1
	i55 = 1<<55 - 1
	i62 = 1<<62 - 1

	m06 = 0x80
	m13 = 0xC000
	m20 = 0xE00000
	m27 = 0xF0000000
	m34 = 0xF800000000
	m41 = 0xFC0000000000
	m48 = 0xFE000000000000
	m55 = 0xFF00000000000000
	m62 = 0xFF80
	m69 = 0xFFC0

	u06 = 0x40
	u13 = 0x2000
	u20 = 0x100000
	u27 = 0x08000000
	u34 = 0x0400000000
	u41 = 0x020000000000
	u48 = 0x01000000000000
	u55 = 0x0080000000000000
	u62 = 0x0040
	u69 = 0x0020
)
