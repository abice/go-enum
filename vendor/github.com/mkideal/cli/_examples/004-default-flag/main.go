package main

import (
	"github.com/mkideal/cli"
)

type argT struct {
	cli.Helper
	Basic  int    `cli:"basic" usage:"basic usage of default" dft:"2"`
	Env    string `cli:"env" usage:"env variable as default" dft:"$HOME"`
	Expr   int    `cli:"expr" usage:"expression as default" dft:"$BASE_PORT+1000"`
	DevDir string `cli:"devdir" usage:"directory of developer" dft:"$HOME/dev"`
}

func main() {
	cli.Run(new(argT), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)
		ctx.String("%d, %s, %d, %s\n", argv.Basic, argv.Env, argv.Expr, argv.DevDir)
		return nil
	})
}
