package optvar

import (
	"strconv"
)

type Int64Var struct {
	name     string
	required bool
	value    *int64
	lastSet  *Source
}

func Int64(name string, ptr *int64, dft int64) *Int64Var {
	i := &Int64Var{
		name:  name,
		value: ptr,
	}
	*i.value = dft
	return i
}

func RequiredInt64(name string, ptr *int64) *Int64Var {
	return &Int64Var{
		name:     name,
		value:    ptr,
		required: true,
	}
}

func (i *Int64Var) Name() string { return i.name }
func (i *Int64Var) IsNil() bool  { return i.value == nil }

func (i *Int64Var) Get() int64 {
	if i.value == nil {
		return 0
	}
	return *i.value
}

func (i *Int64Var) Set(s string) error {
	if s == "" {
		return nil
	}
	i.lastSet = &Source{Str: s}
	v, err := strconv.ParseInt(s, 0, 64)
	if err != nil {
		i.lastSet.Err = err
		return err
	}
	*i.value = v
	return nil
}

func (i *Int64Var) Source() string { return i.lastSet.Str }

func (i *Int64Var) Error() error {
	return getError(i.required, i.lastSet, i.name)
}

func (i *Int64Var) IsSet() bool { return i.lastSet != nil && i.lastSet.Err == nil }
