package optvar

import (
	"strconv"
)

type UintVar struct {
	name     string
	required bool
	value    *uint
	lastSet  *Source
}

func Uint(name string, ptr *uint, dft uint) *UintVar {
	i := &UintVar{
		name:  name,
		value: ptr,
	}
	*i.value = dft
	return i
}

func RequiredUint(name string, ptr *uint) *UintVar {
	return &UintVar{
		name:     name,
		value:    ptr,
		required: true,
	}
}

func (i *UintVar) Name() string { return i.name }
func (i *UintVar) IsNil() bool  { return i.value == nil }

func (i *UintVar) Get() uint {
	if i.value == nil {
		return 0
	}
	return *i.value
}

func (i *UintVar) Set(s string) error {
	if s == "" {
		return nil
	}
	i.lastSet = &Source{Str: s}
	v, err := strconv.ParseUint(s, 0, 64)
	if err != nil {
		i.lastSet.Err = err
		return err
	}
	*i.value = uint(v)
	return nil
}

func (i *UintVar) Source() string { return i.lastSet.Str }

func (i *UintVar) Error() error {
	return getError(i.required, i.lastSet, i.name)
}

func (i *UintVar) IsSet() bool { return i.lastSet != nil && i.lastSet.Err == nil }
