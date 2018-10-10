package main

import (
	"os"

	"github.com/mkideal/cli"
)

type argT struct {
	cli.Helper
	Msg string `edit:"m" usage:"message"`
}

func main() {
	os.Exit(cli.Run(new(argT), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)
		ctx.String("msg: %s", argv.Msg)
		return nil
	}))
}
