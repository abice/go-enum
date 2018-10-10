package cli

import (
	"testing"

	"github.com/labstack/gommon/color"
)

func TestJSONParser(t *testing.T) {
	type T struct {
		A string
		B int
	}
	type argT struct {
		Value T `cli:"t" parser:"json"`
	}

	v := new(argT)
	clr := color.Color{}
	flagSet := parseArgv([]string{`-t`, `{"a": "string", "b": 2}`}, v, clr)
	if flagSet.err != nil {
		t.Errorf("error: %v", flagSet.err)
		return
	}
	want := T{A: "string", B: 2}
	if v.Value != want {
		t.Errorf("want %v, got %v", want, v.Value)
	}
}
