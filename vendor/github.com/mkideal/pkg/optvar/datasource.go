package optvar

import (
	"flag"
	"net/url"
)

// DataSource represents a data source which can apply to Options
type DataSource interface {
	Apply(Options) error
}

type KeyValueSource interface {
	Get(string) string
}

type KeyValueSource2 interface {
	Get(string) (string, bool)
}

// FromKV converts KeyValueSource to DataSource
func FromKV(kv KeyValueSource) DataSource {
	return kvSource{kv}
}

// FromKV2 converts KeyValueSource2 to DataSource
func FromKV2(kv KeyValueSource2) DataSource {
	return kvSource2{kv}
}

type kvSource struct {
	kv KeyValueSource
}

func (ds kvSource) Apply(t Options) error {
	if ds.kv == nil {
		return ErrDataSourceIsNil
	}
	for _, value := range t.Values() {
		if err := value.Set(ds.kv.Get(value.Name())); err != nil {
			return err
		}
	}
	return nil
}

type kvSource2 struct {
	kv KeyValueSource2
}

func (ds kvSource2) Apply(t Options) error {
	if ds.kv == nil {
		return ErrDataSourceIsNil
	}
	for _, value := range t.Values() {
		s, ok := ds.kv.Get(value.Name())
		if ok {
			if err := value.Set(s); err != nil {
				return err
			}
		}
	}
	return nil
}

// Form creates DataSource from url.Values
func Form(values url.Values) DataSource {
	return FromKV(values)
}

// FormString creates DataSource from raw form string
func FormString(s string) (DataSource, error) {
	values, err := url.ParseQuery(s)
	if err != nil {
		return nil, err
	}
	return Form(values), nil
}

type mapSource map[string]string

func (m mapSource) Get(key string) (string, bool) {
	value, ok := m[key]
	return value, ok
}

// Map creates DataSource from string->string map
func Map(m map[string]string) DataSource {
	return FromKV2(mapSource(m))
}

// FlagSet creates DataSource from flag.FlagSet
func FlagSet(flagSet *flag.FlagSet) DataSource {
	m := map[string]string{}
	flagSet.Visit(func(fl *flag.Flag) {
		if fl.Value != nil {
			m[fl.Name] = fl.Value.String()
		}
	})
	return Map(m)
}

// CommandLine creates DataSource from flag.CommandLine
func CommandLine() DataSource {
	if !flag.Parsed() {
		flag.Parse()
	}
	return FlagSet(flag.CommandLine)
}
