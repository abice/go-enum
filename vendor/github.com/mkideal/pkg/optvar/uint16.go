package optvar

import (
	"strconv"
)

type Uint16Var struct {
	name     string
	required bool
	value    *uint16
	lastSet  *Source
}

func Uint16(name string, ptr *uint16, dft uint16) *Uint16Var {
	i := &Uint16Var{
		name:  name,
		value: ptr,
	}
	*i.value = dft
	return i
}

func RequiredUint16(name string, ptr *uint16) *Uint16Var {
	return &Uint16Var{
		name:     name,
		value:    ptr,
		required: true,
	}
}

func (i *Uint16Var) Name() string { return i.name }
func (i *Uint16Var) IsNil() bool  { return i.value == nil }

func (i *Uint16Var) Get() uint16 {
	if i.value == nil {
		return 0
	}
	return *i.value
}

func (i *Uint16Var) Set(s string) error {
	if s == "" {
		return nil
	}
	i.lastSet = &Source{Str: s}
	v, err := strconv.ParseUint(s, 0, 16)
	if err != nil {
		i.lastSet.Err = err
		return err
	}
	*i.value = uint16(v)
	return nil
}

func (i *Uint16Var) Source() string { return i.lastSet.Str }

func (i *Uint16Var) Error() error {
	return getError(i.required, i.lastSet, i.name)
}

func (i *Uint16Var) IsSet() bool { return i.lastSet != nil && i.lastSet.Err == nil }
