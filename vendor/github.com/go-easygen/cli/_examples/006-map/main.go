package main

import (
	"os"

	"github.com/mkideal/cli"
)

type argT struct {
	Macros map[string]int `cli:"D" usage:"define macros"`
}

func main() {
	os.Exit(cli.Run(new(argT), func(ctx *cli.Context) error {
		ctx.JSONln(ctx.Argv())
		return nil
	}))
}
