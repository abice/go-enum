package flag

import (
	"strconv"
	"time"
)

type Value interface {
	String() string
	Set(string) error
}

// intValue wrap int as Value
type intValue int

func NewInt(val int, ptr *int) Value {
	*ptr = val
	return (*intValue)(ptr)
}

func (i *intValue) Set(s string) error {
	v, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	*i = intValue(v)
	return nil
}

func (i intValue) String() string {
	return strconv.Itoa(int(i))
}

// int8Value wrap int8 as Value
type int8Value int8

func NewInt8(val int8, ptr *int8) Value {
	*ptr = val
	return (*int8Value)(ptr)
}

func (i *int8Value) Set(s string) error {
	v, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	*i = int8Value(v)
	return nil
}

func (i int8Value) String() string {
	return strconv.Itoa(int(i))
}

// int16Value wrap int16 as Value
type int16Value int16

func NewInt16(val int16, ptr *int16) Value {
	*ptr = val
	return (*int16Value)(ptr)
}

func (i *int16Value) Set(s string) error {
	v, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	*i = int16Value(v)
	return nil
}

func (i int16Value) String() string {
	return strconv.Itoa(int(i))
}

// int32Value wrap int32 as Value
type int32Value int32

func NewInt32(val int32, ptr *int32) Value {
	*ptr = val
	return (*int32Value)(ptr)
}

func (i *int32Value) Set(s string) error {
	v, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	*i = int32Value(v)
	return nil
}

func (i int32Value) String() string {
	return strconv.Itoa(int(i))
}

// int64Value wrap int64 as Value
type int64Value int64

func NewInt64(val int64, ptr *int64) Value {
	*ptr = val
	return (*int64Value)(ptr)
}

func (i *int64Value) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 64)
	if err != nil {
		return err
	}
	*i = int64Value(v)
	return nil
}

func (i int64Value) String() string {
	return strconv.FormatInt(int64(i), 10)
}

// uintValue wrap uint as Value
type uintValue uint

func NewUint(val uint, ptr *uint) Value {
	*ptr = val
	return (*uintValue)(ptr)
}

func (i *uintValue) Set(s string) error {
	v, err := strconv.ParseUint(s, 0, 64)
	if err != nil {
		return err
	}
	*i = uintValue(v)
	return nil
}

func (i uintValue) String() string {
	return strconv.FormatUint(uint64(i), 10)
}

// uint8Value wrap uint8 as Value
type uint8Value uint8

func NewUint8(val uint8, ptr *uint8) Value {
	*ptr = val
	return (*uint8Value)(ptr)
}

func (i *uint8Value) Set(s string) error {
	v, err := strconv.ParseUint(s, 0, 8)
	if err != nil {
		return err
	}
	*i = uint8Value(v)
	return nil
}

func (i uint8Value) String() string {
	return strconv.FormatUint(uint64(i), 10)
}

// uint16Value wrap uint16 as Value
type uint16Value uint16

func NewUint16(val uint16, ptr *uint16) Value {
	*ptr = val
	return (*uint16Value)(ptr)
}

func (i *uint16Value) Set(s string) error {
	v, err := strconv.ParseUint(s, 0, 16)
	if err != nil {
		return err
	}
	*i = uint16Value(v)
	return nil
}

func (i uint16Value) String() string {
	return strconv.FormatUint(uint64(i), 10)
}

// uint32Value wrap uint32 as Value
type uint32Value uint32

func NewUint32(val uint32, ptr *uint32) Value {
	*ptr = val
	return (*uint32Value)(ptr)
}

func (i *uint32Value) Set(s string) error {
	v, err := strconv.ParseUint(s, 0, 32)
	if err != nil {
		return err
	}
	*i = uint32Value(v)
	return nil
}

func (i uint32Value) String() string {
	return strconv.FormatUint(uint64(i), 10)
}

// uint64Value wrap uint64 as Value
type uint64Value uint64

func NewUint64(val uint64, ptr *uint64) Value {
	*ptr = val
	return (*uint64Value)(ptr)
}

func (i *uint64Value) Set(s string) error {
	v, err := strconv.ParseUint(s, 0, 64)
	if err != nil {
		return err
	}
	*i = uint64Value(v)
	return nil
}

func (i uint64Value) String() string {
	return strconv.FormatUint(uint64(i), 10)
}

// boolValue wrap bool as Value
type boolValue bool

func NewBool(val bool, ptr *bool) Value {
	*ptr = val
	return (*boolValue)(ptr)
}

func (i *boolValue) Set(s string) error {
	v, err := strconv.ParseBool(s)
	if err != nil {
		return err
	}
	*i = boolValue(v)
	return nil
}

func (i boolValue) String() string {
	return strconv.FormatBool(bool(i))
}

// byteValue wrap byte as Value
type byteValue byte

func NewByte(val byte, ptr *byte) Value {
	*ptr = val
	return (*byteValue)(ptr)
}

func (i *byteValue) Set(s string) error {
	v, err := strconv.ParseUint(s, 0, 8)
	if err != nil {
		return err
	}
	*i = byteValue(v)
	return nil
}

func (i byteValue) String() string {
	return strconv.FormatUint(uint64(i), 10)
}

// stringValue wrap string as Value
type stringValue string

func NewString(val string, ptr *string) Value {
	*ptr = val
	return (*stringValue)(ptr)
}

func (i *stringValue) Set(s string) error {
	*i = stringValue(s)
	return nil
}

func (i stringValue) String() string { return string(i) }

// durationValue wrap time.Duration as Value
type durationValue time.Duration

func NewDuration(val time.Duration, ptr *time.Duration) Value {
	*ptr = val
	return (*durationValue)(ptr)
}

func (i *durationValue) Set(s string) error {
	d, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	*i = durationValue(d)
	return nil
}

func (i durationValue) String() string {
	return time.Duration(i).String()
}
