package main

import (
	"github.com/mkideal/cli"
)

type argT struct {
	// []bool, []int, []float32, ... supported too.
	Friends []string `cli:"F" usage:"my friends"`
}

func main() {
	cli.Run(new(argT), func(ctx *cli.Context) error {
		ctx.JSONln(ctx.Argv())
		return nil
	})
}
