package main

import (
	"github.com/mkideal/cli"
)

var _ = app.Register(&cli.Command{
	Name: "build",
	Desc: "Build golang application",
	Argv: func() interface{} { return new(buildT) },
	Fn:   build,
})

type buildT struct {
	cli.Helper
	Dir    string `cli:"dir" usage:"source code root dir" dft:"./"`
	Suffix string `cli:"suffix" usage:"source file suffix" dft:".go,.c,.s"`
	Out    string `cli:"o,out" usage:"output filename"`
}

func build(ctx *cli.Context) error {
	argv := ctx.Argv().(*buildT)
	ctx.String("%s: %v", ctx.Path(), jsonIndent(argv))
	return nil
}
