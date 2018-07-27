package cli

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/labstack/gommon/color"
	"github.com/mkideal/pkg/expr"
)

type flag struct {
	field reflect.StructField
	value reflect.Value

	// isAssigned indicates whether the flag is set(contains default value)
	isAssigned bool

	// isSet indicates whether the flag is set
	isSet bool

	// tag properties
	tag tagProperty

	// actual flag name
	actualFlagName string

	isNeedDelaySet bool

	// last value for need delay set
	// flag maybe assigned too many times, like:
	//	-f xx -f yy -f zz
	// `zz` is the last value
	lastValue string
}

func newFlag(field reflect.StructField, value reflect.Value, tag *tagProperty, clr color.Color, dontSetValue bool) (fl *flag, err error) {
	fl = &flag{field: field, value: value}
	if !fl.value.CanSet() {
		return nil, fmt.Errorf("field %s can not set", clr.Bold(fl.field.Name))
	}
	fl.tag = *tag
	if fl.isPtr() && fl.value.IsNil() {
		fl.value.Set(reflect.New(fl.field.Type.Elem()))
	}
	isSliceDecoder := fl.value.Type().Implements(reflect.TypeOf((*SliceDecoder)(nil)).Elem())
	if !isSliceDecoder && fl.value.CanAddr() {
		isSliceDecoder = fl.value.Addr().Type().Implements(reflect.TypeOf((*SliceDecoder)(nil)).Elem())
	}
	fl.isNeedDelaySet = fl.tag.parserCreator != nil ||
		(fl.field.Type.Kind() != reflect.Slice && fl.field.Type.Kind() != reflect.Map && !isSliceDecoder)
	err = fl.init(clr, dontSetValue)
	return
}

func (fl *flag) init(clr color.Color, dontSetValue bool) error {
	var (
		isNumber  = fl.isInteger() || fl.isFloat()
		isDecoder = fl.value.Type().Implements(reflect.TypeOf((*Decoder)(nil)).Elem())
		dft       string
		err       error
	)
	if !isDecoder && fl.value.CanAddr() {
		isDecoder = fl.value.Addr().Type().Implements(reflect.TypeOf((*Decoder)(nil)).Elem())
	}
	dft, err = parseExpression(fl.tag.dft, isNumber)
	if err != nil {
		return err
	}
	if isNumber && !isDecoder {
		v, err := expr.Eval(dft, nil, nil)
		if err == nil {
			if fl.isInteger() {
				dft = fmt.Sprintf("%d", v.Int())
			} else if fl.isFloat() {
				dft = fmt.Sprintf("%f", v.Float())
			}
		}
	}
	if !dontSetValue && fl.tag.dft != "" && dft != "" {
		if fl.isPtr() || isDecoder || isEmpty(fl.value) {
			return fl.setDefault(dft, clr)
		}
	}
	return nil
}

func isEmpty(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}

func isWordByte(b byte) bool {
	return (b >= 'a' && b <= 'z') ||
		(b >= 'A' && b <= 'Z') ||
		(b >= '0' && b <= '9') ||
		b == '_'
}

func parseExpression(s string, isNumber bool) (string, error) {
	const escapeByte = '$'
	var (
		src       = []byte(s)
		escaping  = false
		exprBuf   bytes.Buffer
		envvarBuf bytes.Buffer
	)
	writeEnv := func(envName string) error {
		if envName == "" {
			return fmt.Errorf("unexpected end after %v", escapeByte)
		}
		env := os.Getenv(envName)
		if env == "" && isNumber {
			env = "0"
		}
		exprBuf.WriteString(env)
		return nil
	}
	for i, b := range src {
		if b == escapeByte {
			if escaping && envvarBuf.Len() == 0 {
				exprBuf.WriteByte(b)
				escaping = false
			} else {
				escaping = true
				if i+1 == len(src) {
					return "", fmt.Errorf("unexpected end after %v", escapeByte)
				}
				envvarBuf.Reset()
			}
			continue
		}
		if escaping {
			if isWordByte(b) {
				envvarBuf.WriteByte(b)
				if i+1 == len(src) {
					if err := writeEnv(envvarBuf.String()); err != nil {
						return "", err
					}
				}
			} else {
				if err := writeEnv(envvarBuf.String()); err != nil {
					return "", err
				}
				exprBuf.WriteByte(b)
				envvarBuf.Reset()
				escaping = false
			}
		} else {
			exprBuf.WriteByte(b)
		}
	}
	return exprBuf.String(), nil
}

func (fl *flag) name() string {
	if fl.actualFlagName != "" {
		return fl.actualFlagName
	}
	if len(fl.tag.longNames) > 0 {
		return fl.tag.longNames[0]
	}
	if len(fl.tag.shortNames) > 0 {
		return fl.tag.shortNames[0]
	}
	return ""
}

func (fl *flag) isBoolean() bool {
	return fl.field.Type.Kind() == reflect.Bool
}

func (fl *flag) isInteger() bool {
	switch fl.field.Type.Kind() {
	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64:
		return true
	}
	return false
}

func (fl *flag) isSlice() bool {
	return fl.field.Type.Kind() == reflect.Slice
}

func (fl *flag) isMap() bool {
	return fl.field.Type.Kind() == reflect.Map
}

func (fl *flag) isFloat() bool {
	kind := fl.field.Type.Kind()
	return kind == reflect.Float32 || kind == reflect.Float64
}

func (fl *flag) isString() bool {
	return fl.field.Type.Kind() == reflect.String
}

func (fl *flag) isPtr() bool {
	return fl.field.Type.Kind() == reflect.Ptr
}

func (fl *flag) getBool() bool {
	if !fl.isBoolean() {
		return false
	}
	return fl.value.Bool()
}

func (fl *flag) setDefault(s string, clr color.Color) error {
	fl.isAssigned = true
	if fl.isNeedDelaySet {
		fl.lastValue = s
		return nil
	}
	return setWithProperType(fl, fl.field.Type, fl.value, s, clr, false)
}

func (fl *flag) set(actualFlagName, s string, clr color.Color) error {
	fl.isSet = true
	fl.isAssigned = true
	fl.actualFlagName = actualFlagName
	if fl.isNeedDelaySet {
		fl.lastValue = s
		return nil
	}
	return setWithProperType(fl, fl.field.Type, fl.value, s, clr, false)
}

func (fl *flag) counterIncr(s string, clr color.Color) error {
	return setWithProperType(fl, fl.field.Type, fl.value, s, clr, false)
}

func (fl *flag) isCounter() bool {
	if decoder := tryGetDecoder(fl.value.Type().Kind(), fl.value); decoder != nil {
		if _, ok := decoder.(CounterDecoder); ok {
			return true
		}
	}
	return false
}

func (fl *flag) setWithNoDelay(actualFlagName, s string, clr color.Color) error {
	fl.isSet = true
	fl.isAssigned = true
	fl.actualFlagName = actualFlagName
	return setWithProperType(fl, fl.field.Type, fl.value, s, clr, false)
}

func tryGetDecoder(kind reflect.Kind, val reflect.Value) Decoder {
	if val.CanInterface() {
		var addrVal = val
		if kind != reflect.Ptr && val.CanAddr() {
			addrVal = val.Addr()
		}
		// try Decoder
		if addrVal.CanInterface() {
			if i := addrVal.Interface(); i != nil {
				if decoder, ok := i.(Decoder); ok {
					return decoder
				}
			}
		}
	}
	return nil
}

func setWithProperType(fl *flag, typ reflect.Type, val reflect.Value, s string, clr color.Color, isSubField bool) error {
	kind := typ.Kind()

	// try parser first of all
	if fl.tag.parserCreator != nil && val.CanInterface() {
		if kind != reflect.Ptr && val.CanAddr() {
			val = val.Addr()
		}
		return fl.tag.parserCreator(val.Interface()).Parse(s)
	}

	if decoder := tryGetDecoder(kind, val); decoder != nil {
		return decoder.Decode(s)
	}

	switch kind {
	case reflect.Bool:
		if v, err := getBool(s, clr); err == nil {
			val.SetBool(v)
		} else {
			return err
		}

	case reflect.String:
		val.SetString(s)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if v, err := getInt(s, clr); err == nil {
			if minmaxIntCheck(kind, v) {
				val.SetInt(v)
			} else {
				return errors.New(clr.Red("value overflow"))
			}
		} else {
			return err
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if v, err := getUint(s, clr); err == nil {
			if minmaxUintCheck(kind, v) {
				val.SetUint(uint64(v))
			} else {
				return errors.New(clr.Red("value overflow"))
			}
		} else {
			return err
		}

	case reflect.Float32, reflect.Float64:
		if v, err := getFloat(s, clr); err == nil {
			if minmaxFloatCheck(kind, v) {
				val.SetFloat(float64(v))
			} else {
				return errors.New(clr.Red("value overflow"))
			}
		} else {
			return err
		}

	case reflect.Slice:
		if isSubField {
			return fmt.Errorf("unsupported type %s as a sub field", kind.String())
		}
		sliceOf := typ.Elem()
		if val.IsNil() {
			slice := reflect.MakeSlice(typ, 0, 4)
			val.Set(slice)
		}
		index := val.Len()
		sliceCap := val.Cap()
		if index+1 <= sliceCap {
			val.SetLen(index + 1)
		} else {
			slice := reflect.MakeSlice(typ, index+1, index+sliceCap/2+1)
			for k := 0; k < index; k++ {
				slice.Index(k).Set(val.Index(k))
			}
			val.Set(slice)
		}
		return setWithProperType(fl, sliceOf, val.Index(index), s, clr, true)

	case reflect.Map:
		if isSubField {
			return fmt.Errorf("unsupported type %s as a sub field", kind.String())
		}
		keyString, valString, err := splitKeyVal(s, fl.tag.sep)
		if err != nil {
			return err
		}
		keyType := typ.Key()
		valType := typ.Elem()
		if val.IsNil() {
			val.Set(reflect.MakeMap(typ))
		}
		k, v := reflect.New(keyType), reflect.New(valType)
		if err := setWithProperType(fl, keyType, k.Elem(), keyString, clr, true); err != nil {
			return err
		}
		if err := setWithProperType(fl, valType, v.Elem(), valString, clr, true); err != nil {
			return err
		}
		val.SetMapIndex(k.Elem(), v.Elem())

	default:
		return fmt.Errorf("unsupported type: %s", kind.String())
	}
	return nil
}

func splitKeyVal(s, sep string) (key, val string, err error) {
	if s == "" {
		err = fmt.Errorf("empty key,val pair")
		return
	}
	index := strings.Index(s, sep)
	if index == -1 {
		return s, "", nil
	}
	return s[:index], s[index+1:], nil
}

func minmaxIntCheck(kind reflect.Kind, v int64) bool {
	switch kind {
	case reflect.Int, reflect.Int64:
		return v >= int64(math.MinInt64) && v <= int64(math.MaxInt64)
	case reflect.Int8:
		return v >= int64(math.MinInt8) && v <= int64(math.MaxInt8)
	case reflect.Int16:
		return v >= int64(math.MinInt16) && v <= int64(math.MaxInt16)
	case reflect.Int32:
		return v >= int64(math.MinInt32) && v <= int64(math.MaxInt32)
	}
	return true
}

func minmaxUintCheck(kind reflect.Kind, v uint64) bool {
	switch kind {
	case reflect.Uint, reflect.Uint64:
		return v <= math.MaxUint64
	case reflect.Uint8:
		return v <= math.MaxUint8
	case reflect.Uint16:
		return v <= math.MaxUint16
	case reflect.Uint32:
		return v <= math.MaxUint32
	}
	return true
}

func minmaxFloatCheck(kind reflect.Kind, v float64) bool {
	switch kind {
	case reflect.Float32:
		return v >= -float64(math.MaxFloat32) && v <= float64(math.MaxFloat32)
	case reflect.Float64:
		return v >= -float64(math.MaxFloat64) && v <= float64(math.MaxFloat64)
	}
	return true
}

func getBool(s string, clr color.Color) (bool, error) {
	s = strings.ToLower(s)
	if s == "true" || s == "yes" || s == "y" || s == "" {
		return true, nil
	}
	if s == "false" || s == "none" || s == "no" || s == "not" || s == "n" {
		return false, nil
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return false, fmt.Errorf("`%s' couldn't converted to a %s", s, clr.Bold("bool"))
	}
	return i != 0, nil
}

func getInt(s string, clr color.Color) (int64, error) {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("`%s' couldn't converted to an %s", s, clr.Bold("int"))
	}
	return i, nil
}

func getUint(s string, clr color.Color) (uint64, error) {
	i, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("`%s' couldn't converted to an %s", s, clr.Bold("uint"))
	}
	return i, nil
}

func getFloat(s string, clr color.Color) (float64, error) {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("`%s' couldn't converted to a %s", s, clr.Bold("float"))
	}
	return f, nil
}
