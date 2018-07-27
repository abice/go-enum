package main

import (
	"github.com/mkideal/cli"
)

var _ = app.Register(&cli.Command{
	Name: "clean",
	Desc: "Clean build data",
	Argv: func() interface{} { return new(cleanT) },
	Fn:   clean,
})

type cleanT struct {
	cli.Helper
	Recursion bool `cli:"R,recursion" usage:"clean recursion or not" dft:"true"`
}

func clean(ctx *cli.Context) error {
	argv := ctx.Argv().(*cleanT)
	ctx.String("%s: %v", ctx.Path(), jsonIndent(argv))
	return nil
}
