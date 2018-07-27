package main

import (
	"github.com/mkideal/cli"
)

type argT struct {
	cli.Helper
	Port int  `cli:"p,port" usage:"short and long format flags both are supported"`
	X    bool `cli:"x" usage:"boolean type"`
	Y    bool `cli:"y" usage:"boolean type, too"`
}

func main() {
	cli.Run(new(argT), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)
		ctx.String("port=%d, x=%v, y=%v\n", argv.Port, argv.X, argv.Y)
		return nil
	})
}
