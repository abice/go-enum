package optvar

import (
	"strconv"
)

type Int16Var struct {
	name     string
	required bool
	value    *int16
	lastSet  *Source
}

func Int16(name string, ptr *int16, dft int16) *Int16Var {
	i := &Int16Var{
		name:  name,
		value: ptr,
	}
	*i.value = dft
	return i
}

func RequiredInt16(name string, ptr *int16) *Int16Var {
	return &Int16Var{
		name:     name,
		value:    ptr,
		required: true,
	}
}

func (i *Int16Var) Name() string { return i.name }
func (i *Int16Var) IsNil() bool  { return i.value == nil }

func (i *Int16Var) Get() int16 {
	if i.value == nil {
		return 0
	}
	return *i.value
}

func (i *Int16Var) Set(s string) error {
	if s == "" {
		return nil
	}
	i.lastSet = &Source{Str: s}
	v, err := strconv.ParseInt(s, 0, 16)
	if err != nil {
		i.lastSet.Err = err
		return err
	}
	*i.value = int16(v)
	return nil
}

func (i *Int16Var) Source() string { return i.lastSet.Str }

func (i *Int16Var) Error() error {
	return getError(i.required, i.lastSet, i.name)
}

func (i *Int16Var) IsSet() bool { return i.lastSet != nil && i.lastSet.Err == nil }
