package main

import (
	"os"
	"strings"

	"github.com/mkideal/cli"
)

type exampleDecoder struct {
	list []string
}

// Decode implements cli.Decoder interface
func (d *exampleDecoder) Decode(s string) error {
	d.list = strings.Split(s, ",")
	return nil
}

type argT struct {
	Example exampleDecoder `cli:"d" usage:"example decoder"`
}

func main() {
	os.Exit(cli.Run(new(argT), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)
		ctx.JSONln(argv.Example.list)
		return nil
	}))
}
