package main

import (
	"github.com/mkideal/cli"
	clix "github.com/mkideal/cli/ext"
)

type argT struct {
	Time     clix.Time     `cli:"t" usage:"time"`
	Duration clix.Duration `cli:"d" usage:"duration"`
}

func main() {
	cli.Run(new(argT), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)
		ctx.String("time=%v, duration=%v\n", argv.Time, argv.Duration)
		return nil
	})
}
