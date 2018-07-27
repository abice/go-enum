package optvar

import (
	"reflect"
)

func Reflect(name string, v interface{}, tagName string) (Options, error) {
	typ := reflect.TypeOf(v)
	val := reflect.ValueOf(v)
	switch typ.Kind() {
	case reflect.Ptr:
		elem := typ.Elem()
		switch elem.Kind() {
		case reflect.Ptr, reflect.Map:
			return Reflect(name, val.Elem(), tagName)
		case reflect.Struct:
			return reflectStruct(name, val.Elem(), tagName)
		default:
			return nil, UnsupportedType(elem.Kind())
		}
	case reflect.Map:
		return reflectMap(name, val)
	default:
		return nil, UnsupportedType(typ.Kind())
	}
}

func reflectMap(name string, val reflect.Value) (Options, error) {
	var values []Value
	for _, mkey := range val.MapKeys() {
		if mkey.Type().Kind() != reflect.String {
			return nil, UnsupportedType(mkey.Type().Kind())
		}
		mval := val.MapIndex(mkey)
		key := mkey.Interface().(string)
		value, err := reflectValue(key, mval)
		if err != nil {
			return nil, err
		}
		values = append(values, value)
	}
	return Values(name, values...), nil
}

func reflectStruct(name string, val reflect.Value, tagName string) (Options, error) {
	numField := val.NumField()
	for i := 0; i < numField; i++ {
		field := deepElem(val.Field(i))
		fieldType := field.Type()
		switch fieldType.Kind() {
		case reflect.Ptr:
		case reflect.Struct:
		}
	}
	return nil, nil
}

func reflectValue(key string, val reflect.Value) (Value, error) {
	// TODO
	return nil, nil
}

func deepElem(val reflect.Value) reflect.Value {
	for val.Type().Kind() != reflect.Ptr {
		if val.IsNil() {
			val = reflect.New(val.Type())
		}
		val = val.Elem()
	}
	return val
}

//	switch typ.Kind() {
//	case reflect.Bool:
//	case reflect.Int:
//	case reflect.Int8:
//	case reflect.Int16:
//	case reflect.Int32:
//	case reflect.Int64:
//	case reflect.Uint:
//	case reflect.Uint8:
//	case reflect.Uint16:
//	case reflect.Uint32:
//	case reflect.Uint64:
//	case reflect.Uintptr:
//	case reflect.Float32:
//	case reflect.Float64:
//	case reflect.Complex64:
//	case reflect.Complex128:
//	case reflect.Array:
//	case reflect.Chan:
//	case reflect.Func:
//	case reflect.Interface:
//	case reflect.Map:
//	case reflect.Ptr:
//	case reflect.Slice:
//	case reflect.String:
//	case reflect.Struct:
//	case reflect.UnsafePointer:
//	default:
//	}
