package optvar

import (
	"strconv"
)

type Uint32Var struct {
	name     string
	required bool
	value    *uint32
	lastSet  *Source
}

func Uint32(name string, ptr *uint32, dft uint32) *Uint32Var {
	i := &Uint32Var{
		name:  name,
		value: ptr,
	}
	*i.value = dft
	return i
}

func (i *Uint32Var) Name() string { return i.name }
func (i *Uint32Var) IsNil() bool  { return i.value == nil }

func (i *Uint32Var) Get() uint32 {
	if i.value == nil {
		return 0
	}
	return *i.value
}

func RequiredUint32(name string, ptr *uint32) *Uint32Var {
	return &Uint32Var{
		name:     name,
		value:    ptr,
		required: true,
	}
}

func (i *Uint32Var) Set(s string) error {
	if s == "" {
		return nil
	}
	i.lastSet = &Source{Str: s}
	v, err := strconv.ParseUint(s, 0, 32)
	if err != nil {
		i.lastSet.Err = err
		return err
	}
	*i.value = uint32(v)
	return nil
}

func (i *Uint32Var) Source() string { return i.lastSet.Str }

func (i *Uint32Var) Error() error {
	return getError(i.required, i.lastSet, i.name)
}

func (i *Uint32Var) IsSet() bool { return i.lastSet != nil && i.lastSet.Err == nil }
