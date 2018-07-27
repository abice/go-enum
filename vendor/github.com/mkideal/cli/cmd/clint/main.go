package main

import (
	"github.com/mkideal/cli"
)

type argT struct {
	cli.Helper
	Version bool `cli:"!v,version" usage:"display version info"`
}

func main() {
	cli.Run(new(argT), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)
		if argv.Version {
			ctx.String("%v\n", appVersion)
			return nil
		}
		return nil
	})
}
