package main

import (
	"github.com/mkideal/cli"
)

type argT struct {
	Version  bool `cli:"!v" usage:"force flag, note the !"`
	Required int  `cli:"*r" usage:"required flag"`
}

func main() {
	cli.Run(new(argT), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)
		if argv.Version {
			ctx.String("v0.0.1\n")
		}
		return nil
	})
}
