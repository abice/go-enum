package main

import (
	"fmt"
	"os"

	"github.com/mkideal/cli"
)

func main() {
	if err := cli.Root(root,
		cli.Tree(child1),
		cli.Tree(child2),
	).Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type argT struct {
	Error bool `cli:"e" usage:"return error"`
}

var root = &cli.Command{
	Name: "app",
	Argv: func() interface{} { return new(argT) },
	OnRootBefore: func(ctx *cli.Context) error {
		ctx.String("OnRootBefore invoked\n")
		return nil
	},
	OnRootAfter: func(ctx *cli.Context) error {
		ctx.String("OnRootAfter invoked\n")
		return nil
	},
	Fn: func(ctx *cli.Context) error {
		ctx.String("exec root command\n")
		argv := ctx.Argv().(*argT)
		if argv.Error {
			return fmt.Errorf("root command returns error")
		}
		return nil
	},
}

var child1 = &cli.Command{
	Name: "child1",
	Argv: func() interface{} { return new(argT) },
	OnBefore: func(ctx *cli.Context) error {
		ctx.String("child1 OnBefore invoked\n")
		return nil
	},
	OnAfter: func(ctx *cli.Context) error {
		ctx.String("child1 OnAfter invoked\n")
		return nil
	},
	Fn: func(ctx *cli.Context) error {
		ctx.String("exec child1 command\n")
		argv := ctx.Argv().(*argT)
		if argv.Error {
			return fmt.Errorf("child1 command returns error")
		}
		return nil
	},
}

var child2 = &cli.Command{
	Name:   "child2",
	NoHook: true,
	Fn: func(ctx *cli.Context) error {
		ctx.String("exec child2 command\n")
		return nil
	},
}
