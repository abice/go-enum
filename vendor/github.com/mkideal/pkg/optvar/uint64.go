package optvar

import (
	"strconv"
)

type Uint64Var struct {
	name     string
	required bool
	value    *uint64
	lastSet  *Source
}

func Uint64(name string, ptr *uint64, dft uint64) *Uint64Var {
	i := &Uint64Var{
		name:  name,
		value: ptr,
	}
	*i.value = dft
	return i
}

func RequiredUint64(name string, ptr *uint64) *Uint64Var {
	return &Uint64Var{
		name:     name,
		value:    ptr,
		required: true,
	}
}

func (i *Uint64Var) Name() string { return i.name }
func (i *Uint64Var) IsNil() bool  { return i.value == nil }

func (i *Uint64Var) Get() uint64 {
	if i.value == nil {
		return 0
	}
	return *i.value
}

func (i *Uint64Var) Set(s string) error {
	if s == "" {
		return nil
	}
	i.lastSet = &Source{Str: s}
	v, err := strconv.ParseUint(s, 0, 64)
	if err != nil {
		i.lastSet.Err = err
		return err
	}
	*i.value = v
	return nil
}

func (i *Uint64Var) Source() string { return i.lastSet.Str }

func (i *Uint64Var) Error() error {
	return getError(i.required, i.lastSet, i.name)
}

func (i *Uint64Var) IsSet() bool { return i.lastSet != nil && i.lastSet.Err == nil }
