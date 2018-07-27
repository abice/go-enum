package optvar

import (
	"errors"
	"reflect"
)

var (
	ErrDataSourceIsNil = errors.New("data source is nil")
)

type MissingRequiredVarError string

func (e MissingRequiredVarError) Error() string {
	return "missing required `" + string(e) + "`"
}

func MissingRequiredVar(name string) error {
	return MissingRequiredVarError(name)
}

func getError(required bool, src *Source, name string) error {
	if src != nil {
		return src.Err
	} else if required {
		return MissingRequiredVar(name)
	}
	return nil
}

type UnsupportedTypeErrror struct {
	kind reflect.Kind
}

func (e UnsupportedTypeErrror) Error() string { return "unsupported type: " + e.kind.String() }

func UnsupportedType(kind reflect.Kind) error { return UnsupportedTypeErrror{kind} }
