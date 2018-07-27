package main

import (
	"github.com/mkideal/cli"
)

type argT struct {
	cli.Helper
	Msg string `edit:"m" usage:"message"`
}

func main() {
	cli.Run(new(argT), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)
		ctx.String("msg: %s", argv.Msg)
		return nil
	})
}
