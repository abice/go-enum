package optvar

import (
	"strconv"
)

type BoolVar struct {
	name     string
	required bool
	value    *bool
	lastSet  *Source
}

func Bool(name string, ptr *bool, dft bool) *BoolVar {
	i := &BoolVar{
		name:  name,
		value: ptr,
	}
	*i.value = dft
	return i
}

func RequiredBool(name string, ptr *bool) *BoolVar {
	return &BoolVar{
		name:     name,
		value:    ptr,
		required: true,
	}
}

func (i *BoolVar) Name() string { return i.name }
func (i *BoolVar) IsNil() bool  { return i.value == nil }

func (i *BoolVar) Get() bool {
	if i.value == nil {
		return false
	}
	return *i.value
}

func (i *BoolVar) Set(s string) error {
	if s == "" {
		return nil
	}
	i.lastSet = &Source{Str: s}
	v, err := strconv.ParseBool(s)
	if err != nil {
		i.lastSet.Err = err
		return err
	}
	*i.value = v
	return nil
}

func (i *BoolVar) Source() string { return i.lastSet.Str }

func (i *BoolVar) Error() error {
	return getError(i.required, i.lastSet, i.name)
}

func (i *BoolVar) IsSet() bool { return i.lastSet != nil && i.lastSet.Err == nil }
