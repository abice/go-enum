package expr

import (
	"fmt"
	"math"
	"strconv"
)

func intRawString(i int64) string     { return strconv.FormatInt(i, 10) }
func floatRawString(f float64) string { return fmt.Sprintf("%f", f) }

func stringAdd(s1, s2 Value) Value {
	return Value{
		kind:     KindString,
		rawValue: s1.rawValue + s2.rawValue,
	}
}

func intAdd(v1, v2 Value) (Value, error) { return Int(v1.intValue + v2.intValue), nil }
func intSub(v1, v2 Value) (Value, error) { return Int(v1.intValue - v2.intValue), nil }
func intMul(v1, v2 Value) (Value, error) { return Int(v1.intValue * v2.intValue), nil }
func intQuo(v1, v2 Value) (Value, error) {
	if v2.intValue == 0 {
		return Zero(), ErrDivideZero
	}
	return Int(v1.intValue / v2.intValue), nil
}
func intRem(v1, v2 Value) (Value, error) {
	if v2.intValue == 0 {
		return Zero(), ErrDivideZero
	}
	return Int(v1.intValue % v2.intValue), nil
}
func intPow(v1, v2 Value) (Value, error) {
	if v1.intValue == 0 {
		return Zero(), ErrPowOfZero
	}
	return Int(int64(math.Pow(float64(v1.intValue), float64(v2.intValue)))), nil
}

func floatAdd(v1, v2 Value) (Value, error) { return Float(v1.floatValue + v2.floatValue), nil }
func floatSub(v1, v2 Value) (Value, error) { return Float(v1.floatValue - v2.floatValue), nil }
func floatMul(v1, v2 Value) (Value, error) { return Float(v1.floatValue * v2.floatValue), nil }
func floatQuo(v1, v2 Value) (Value, error) { return Float(v1.floatValue / v2.floatValue), nil }
func floatRem(v1, v2 Value) (Value, error) {
	return Float(math.Remainder(v1.floatValue, v2.floatValue)), nil
}
func floatPow(v1, v2 Value) (Value, error) {
	if v1.floatValue == 0 {
		return Zero(), ErrPowOfZero
	}
	return Float(math.Pow(v1.floatValue, v2.floatValue)), nil
}

type binaryOpFunc func(Value, Value) (Value, error)

func binaryOp(v1, v2 Value, iop, fop binaryOpFunc) (Value, error) {
	if v1.kind == KindString || v2.kind == KindString {
		return Zero(), ErrTypeMismatchForOp
	}
	if v1.kind == KindInvalid || v2.kind == KindInvalid {
		return Zero(), ErrUnsupportedType
	}
	if v1.kind == KindFloat {
		if v2.kind == KindInt {
			return fop(v1, Float(float64(v2.intValue)))
		}
		return fop(v1, v2)
	} else {
		if v2.kind == KindInt {
			return iop(v1, v2)
		}
		return fop(Float(float64(v1.intValue)), v2)
	}
}

type compareFunc func(Value, Value) Value

func stringEq(v1, v2 Value) Value { return Bool(v1.rawValue == v2.rawValue) }
func stringGt(v1, v2 Value) Value { return Bool(v1.rawValue > v2.rawValue) }
func stringGe(v1, v2 Value) Value { return Bool(v1.rawValue >= v2.rawValue) }

func intEq(v1, v2 Value) Value { return Bool(v1.intValue == v2.intValue) }
func intGt(v1, v2 Value) Value { return Bool(v1.intValue > v2.intValue) }
func intGe(v1, v2 Value) Value { return Bool(v1.intValue >= v2.intValue) }

func floatEq(v1, v2 Value) Value { return Bool(v1.floatValue == v2.floatValue) }
func floatGt(v1, v2 Value) Value { return Bool(v1.floatValue > v2.floatValue) }
func floatGe(v1, v2 Value) Value { return Bool(v1.floatValue >= v2.floatValue) }

func compare(v1, v2 Value, scmp, icmp, fcmp compareFunc) (Value, error) {
	switch v1.kind {
	case KindString:
		if v2.kind == KindString {
			return scmp(v1, v2), nil
		}
		return False(), ErrComparedTypesMismatch
	case KindInt:
		if v2.kind == KindInt {
			return icmp(v1, v2), nil
		} else if v2.kind == KindFloat {
			return fcmp(Float(float64(v1.intValue)), v2), nil
		}
		return False(), ErrComparedTypesMismatch
	case KindFloat:
		if v2.kind == KindInt {
			return fcmp(v1, Float(float64(v2.intValue))), nil
		} else if v2.kind == KindFloat {
			return fcmp(v1, v2), nil
		}
		return False(), ErrComparedTypesMismatch
	default:
		return False(), ErrUnsupportedType
	}
}
