package main

import (
	"github.com/mkideal/cli"
)

type config struct {
	A string
	B int
	C bool
}

type argT struct {
	JSON config `cli:"c,config" usage:"parse json string" parser:"json"`
}

func main() {
	cli.Run(new(argT), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)
		ctx.JSONIndentln(argv.JSON, "", "    ")
		return nil
	})
}
