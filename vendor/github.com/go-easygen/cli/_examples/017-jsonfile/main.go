package main

import (
	"os"

	"github.com/mkideal/cli"
)

type config struct {
	A string
	B int
	C bool
}

type argT struct {
	JSON config `cli:"c,config" usage:"parse json from file" parser:"jsonfile"`
}

func main() {
	os.Exit(cli.Run(new(argT), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)
		ctx.JSONIndentln(argv.JSON, "", "    ")
		return nil
	}))
}
