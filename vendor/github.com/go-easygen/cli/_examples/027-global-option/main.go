package main

import (
	"fmt"
	"os"

	"github.com/mkideal/cli"
)

func main() {
	if err := cli.Root(root, cli.Tree(sub)).Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// root command
type rootT struct {
	cli.Helper
	Self *rootT `json:"-" cli:"c,config" usage:"config" parser:"jsonfile" dft:"1.json"`
	Host string `cli:"H,host" usage:"host addr" dft:"$HOST"`
	Port int    `cli:"p,port" usage:"listening port"`
}

var root = &cli.Command{
	Name:   "app",
	Desc:   "application",
	Global: true,
	Argv: func() interface{} {
		t := new(rootT)
		t.Self = t
		return t
	},
	Fn: func(ctx *cli.Context) error {
		ctx.JSON(ctx.RootArgv())
		ctx.JSON(ctx.Argv())
		return nil
	},
}

// sub command
type subT struct {
	World string `cli:"w" usage:"world is a sub flag"`
}

var sub = &cli.Command{
	Name: "sub",
	Desc: "subcommand",
	Argv: func() interface{} { return new(subT) },
	Fn: func(ctx *cli.Context) error {
		ctx.JSON(ctx.RootArgv())
		ctx.JSON(ctx.Argv())
		ctx.String("\n")

		var argv = &subT{}
		var parentArgv = &rootT{}
		if err := ctx.GetArgvList(argv, parentArgv); err != nil {
			return err
		}
		ctx.JSON(parentArgv)
		ctx.JSON(argv)
		return nil
	},
}
