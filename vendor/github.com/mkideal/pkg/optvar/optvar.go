package optvar

type Value interface {
	Name() string
	Set(string) error
	IsNil() bool
	IsSet() bool
	Source() string
	Error() error
}

type Source struct {
	Str string
	Err error
}

type Options interface {
	Name() string
	Groups() []Options
	Values() []Value
}

// group implements Options interface
type group struct {
	name   string
	values []Value
}

func (g group) Name() string      { return g.name }
func (g group) Groups() []Options { return nil }
func (g group) Values() []Value   { return g.values }

// Values creates a `Options`: group
func Values(name string, values ...Value) Options {
	return group{
		name:   name,
		values: values,
	}
}

// options implements Options interface
type options struct {
	name   string
	groups []Options
}

func (t options) Name() string      { return t.name }
func (t options) Groups() []Options { return t.groups }
func (t options) Values() []Value   { return nil }

// Groups creates a `Options`: options
func Groups(name string, groups ...Options) Options {
	return options{
		name:   name,
		groups: groups,
	}
}
