package optvar

import (
	"strconv"
)

type IntVar struct {
	name     string
	required bool
	value    *int
	lastSet  *Source
}

func Int(name string, ptr *int, dft int) *IntVar {
	i := &IntVar{
		name:  name,
		value: ptr,
	}
	*i.value = dft
	return i
}

func RequiredInt(name string, ptr *int) *IntVar {
	return &IntVar{
		name:     name,
		value:    ptr,
		required: true,
	}
}

func (i *IntVar) Name() string { return i.name }
func (i *IntVar) IsNil() bool  { return i.value == nil }

func (i *IntVar) Get() int {
	if i.value == nil {
		return 0
	}
	return *i.value
}

func (i *IntVar) Set(s string) error {
	if s == "" {
		return nil
	}
	i.lastSet = &Source{Str: s}
	v, err := strconv.ParseInt(s, 0, 64)
	if err != nil {
		i.lastSet.Err = err
		return err
	}
	*i.value = int(v)
	return nil
}

func (i *IntVar) Source() string { return i.lastSet.Str }

func (i *IntVar) Error() error {
	return getError(i.required, i.lastSet, i.name)
}

func (i *IntVar) IsSet() bool { return i.lastSet != nil && i.lastSet.Err == nil }
