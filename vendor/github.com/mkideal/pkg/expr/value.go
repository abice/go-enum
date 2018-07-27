package expr

import (
	"errors"
	"strconv"
	"strings"
)

var (
	ErrFailedToParseInteger  = errors.New("failed to parse integer")
	ErrFailedToParseFloat    = errors.New("failed to parse float")
	ErrNotAnInteger          = errors.New("not an integer")
	ErrNotAFloat             = errors.New("not a float")
	ErrUnsupportedType       = errors.New("unsupported type")
	ErrTypeMismatchForOp     = errors.New("type mismatch for operater")
	ErrDivideZero            = errors.New("divide zero")
	ErrPowOfZero             = errors.New("power of zero")
	ErrComparedTypesMismatch = errors.New("compared types mismatch")
	ErrBadArgumentsSize      = errors.New("bad arguments size")
)

type Kind int

const (
	KindInvalid Kind = iota
	KindInt
	KindFloat
	KindString
)

var (
	nilValue = Value{kind: KindInvalid}

	varZero  = Value{kind: KindInt, rawValue: "0"}
	varTrue  = Value{kind: KindInt, intValue: 1, rawValue: "true"}
	varFalse = Value{kind: KindInt, intValue: 0, rawValue: "false"}
)

func Nil() Value   { return nilValue }
func Zero() Value  { return varZero }
func True() Value  { return varTrue }
func False() Value { return varFalse }
func Equal(v1, v2 Value) bool {
	eq, err := v1.Eq(v2)
	return err == nil && eq.Bool()
}

func Bool(ok bool) Value {
	if ok {
		return True()
	}
	return False()
}

func Int(i int64) Value     { return Value{kind: KindInt, intValue: i, rawValue: intRawString(i)} }
func Float(f float64) Value { return Value{kind: KindFloat, floatValue: f, rawValue: floatRawString(f)} }
func String(s string) Value { return Value{kind: KindString, rawValue: s} }

type Value struct {
	kind       Kind
	rawValue   string
	intValue   int64
	floatValue float64
}

func NewValue(kind Kind) Value {
	return Value{kind: kind}
}

func (v *Value) Set(s string) error {
	switch v.kind {
	case KindString:
		// donothing
	case KindInt:
		if i, err := strconv.ParseInt(s, 0, 64); err == nil {
			v.intValue = i
		} else {
			return ErrFailedToParseInteger
		}
	case KindFloat:
		if f, err := strconv.ParseFloat(s, 64); err == nil {
			v.floatValue = f
		} else {
			return ErrFailedToParseFloat
		}
	default:
		return ErrUnsupportedType
	}
	v.rawValue = s
	return nil
}

func (v Value) Kind() Kind     { return v.kind }
func (v Value) String() string { return v.rawValue }
func (v Value) Int() int64 {
	if v.kind == KindFloat {
		return int64(v.floatValue)
	}
	return v.intValue
}
func (v Value) Float() float64 {
	if v.kind == KindInt {
		return float64(v.intValue)
	}
	return v.floatValue
}
func (v Value) Bool() bool {
	switch v.kind {
	case KindString:
		return v.rawValue != ""
	case KindInt:
		return v.intValue != 0
	case KindFloat:
		return v.floatValue != 0
	}
	return false
}

func (v Value) Add(v2 Value) (Value, error) {
	if v.kind == KindString && v2.kind == KindString {
		return stringAdd(v, v2), nil
	}
	return binaryOp(v, v2, intAdd, floatAdd)
}

func (v Value) Sub(v2 Value) (Value, error) { return binaryOp(v, v2, intSub, floatSub) }
func (v Value) Mul(v2 Value) (Value, error) { return binaryOp(v, v2, intMul, floatMul) }
func (v Value) Quo(v2 Value) (Value, error) { return binaryOp(v, v2, intQuo, floatQuo) }
func (v Value) Rem(v2 Value) (Value, error) { return binaryOp(v, v2, intRem, floatRem) }
func (v Value) Pow(v2 Value) (Value, error) { return binaryOp(v, v2, intPow, floatPow) }
func (v Value) And(v2 Value) Value          { return Bool(v.Bool() && v2.Bool()) }
func (v Value) Or(v2 Value) Value           { return Bool(v.Bool() || v2.Bool()) }
func (v Value) Not() Value                  { return Bool(!v.Bool()) }
func (v Value) Eq(v2 Value) (Value, error)  { return compare(v, v2, stringEq, intEq, floatEq) }

func (v Value) Ne(v2 Value) (Value, error) {
	result, err := v.Eq(v2)
	if err == nil {
		result = result.Not()
	}
	return result, err
}

func (v Value) Gt(v2 Value) (Value, error) { return compare(v, v2, stringGt, intGt, floatGt) }
func (v Value) Ge(v2 Value) (Value, error) { return compare(v, v2, stringGe, intGe, floatGe) }
func (v Value) Lt(v2 Value) (Value, error) { return v2.Gt(v) }
func (v Value) Le(v2 Value) (Value, error) { return v2.Ge(v) }

func (v Value) Contains(v2 Value) Value {
	if v.kind == KindString && v2.kind == KindString {
		return Bool(strings.Contains(v.rawValue, v2.rawValue))
	}
	return False()
}

func ExpectNArg(got, want int) error {
	if got != want {
		return ErrBadArgumentsSize
	}
	return nil
}
