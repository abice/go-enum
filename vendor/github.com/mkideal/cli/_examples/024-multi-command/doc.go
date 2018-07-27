package main

import (
	"github.com/mkideal/cli"
)

var _ = app.Register(&cli.Command{
	Name: "doc",
	Desc: "Generate documents",
	Argv: func() interface{} { return new(docT) },
	Fn:   doc,
})

type docT struct {
	cli.Helper
	Suffix string `cli:"suffix" usage:"source file suffix" dft:".go,.c,.s"`
	Out    string `cli:"o,out" usage:"output filename"`
}

func doc(ctx *cli.Context) error {
	argv := ctx.Argv().(*docT)
	ctx.String("%s: %v", ctx.Path(), jsonIndent(argv))
	return nil
}
