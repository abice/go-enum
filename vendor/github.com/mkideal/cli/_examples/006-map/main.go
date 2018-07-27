package main

import (
	"github.com/mkideal/cli"
)

type argT struct {
	Macros map[string]int `cli:"D" usage:"define macros"`
}

func main() {
	cli.Run(new(argT), func(ctx *cli.Context) error {
		ctx.JSONln(ctx.Argv())
		return nil
	})
}
