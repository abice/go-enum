package main

import (
	"os"
	"reflect"

	"github.com/mkideal/cli"
)

type myParser struct {
	ptr interface{}
}

func newMyParser(ptr interface{}) cli.FlagParser {
	return &myParser{ptr}
}

// Parse implements FlagParser.Parse interface
func (parser *myParser) Parse(s string) error {
	typ := reflect.TypeOf(parser.ptr)
	val := reflect.ValueOf(parser.ptr)
	if typ.Kind() == reflect.Ptr {
		kind := reflect.Indirect(val).Type().Kind()
		if kind == reflect.Struct {
			typElem, valElem := typ.Elem(), val.Elem()
			numField := valElem.NumField()
			for i := 0; i < numField; i++ {
				_, valField := typElem.Field(i), valElem.Field(i)
				if valField.Kind() == reflect.Int &&
					valField.CanSet() {
					valField.SetInt(2)
				}
				if valField.Kind() == reflect.String &&
					valField.CanSet() {
					valField.SetString("B")
				}
			}
		}
	}
	return nil
}

type config struct {
	A int
	B string
}

type argT struct {
	Cfg config `cli:"cfg" parser:"myparser"`
}

func main() {
	// register parser factory function
	cli.RegisterFlagParser("myparser", newMyParser)

	os.Exit(cli.Run(new(argT), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)
		ctx.String("%v\n", argv.Cfg)
		return nil
	}))
}
