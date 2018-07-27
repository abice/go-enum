package optvar

type StringVar struct {
	name     string
	required bool
	value    *string
	isSet    bool
}

func String(name string, ptr *string, dft string) *StringVar {
	i := &StringVar{
		name:  name,
		value: ptr,
	}
	*i.value = dft
	return i
}

func RequiredString(name string, ptr *string) *StringVar {
	return &StringVar{
		name:     name,
		value:    ptr,
		required: true,
	}
}

func (i *StringVar) Name() string { return i.name }
func (i *StringVar) IsNil() bool  { return i.value == nil }

func (i *StringVar) Get() string {
	if i.value == nil {
		return ""
	}
	return *i.value
}

func (i *StringVar) Set(s string) error {
	if s == "" {
		return nil
	}
	i.isSet = true
	*i.value = s
	return nil
}

func (i *StringVar) Source() string {
	if i.value != nil {
		return *i.value
	}
	return ""
}

func (i *StringVar) Error() error {
	if i.required && i.value == nil {
		return MissingRequiredVar(i.name)
	}
	return nil
}

func (i *StringVar) IsSet() bool { return i.isSet }
