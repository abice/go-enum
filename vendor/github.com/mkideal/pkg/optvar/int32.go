package optvar

import (
	"strconv"
)

type Int32Var struct {
	name     string
	required bool
	value    *int32
	lastSet  *Source
}

func Int32(name string, ptr *int32, dft int32) *Int32Var {
	i := &Int32Var{
		name:  name,
		value: ptr,
	}
	*i.value = dft
	return i
}

func RequiredInt32(name string, ptr *int32) *Int32Var {
	return &Int32Var{
		name:     name,
		value:    ptr,
		required: true,
	}
}

func (i *Int32Var) Name() string { return i.name }
func (i *Int32Var) IsNil() bool  { return i.value == nil }

func (i *Int32Var) Get() int32 {
	if i.value == nil {
		return 0
	}
	return *i.value
}

func (i *Int32Var) Set(s string) error {
	if s == "" {
		return nil
	}
	i.lastSet = &Source{Str: s}
	v, err := strconv.ParseInt(s, 0, 32)
	if err != nil {
		i.lastSet.Err = err
		return err
	}
	*i.value = int32(v)
	return nil
}

func (i *Int32Var) Source() string { return i.lastSet.Str }

func (i *Int32Var) Error() error {
	return getError(i.required, i.lastSet, i.name)
}

func (i *Int32Var) IsSet() bool { return i.lastSet != nil && i.lastSet.Err == nil }
