package main

import (
	"fmt"
	"github.com/mkideal/cli"
	"os"
)

func main() {
	root := cli.Root(
		&cli.Command{Name: "tree"},
		cli.Tree(cmd1,
			cli.Tree(cmd11),
			cli.Tree(cmd12),
		),
		cli.Tree(cmd2,
			cli.Tree(cmd21),
			cli.Tree(cmd22,
				cli.Tree(cmd221),
				cli.Tree(cmd222),
				cli.Tree(cmd223),
			),
		),
	)

	if err := root.Run(os.Args[1:]); err != nil {
		fmt.Println(err)
	}
}

func log(ctx *cli.Context) error {
	ctx.String("path: `%s`\n", ctx.Path())
	return nil
}

var (
	cmd1  = &cli.Command{Name: "cmd1", Fn: log}
	cmd11 = &cli.Command{Name: "cmd11", Fn: log}
	cmd12 = &cli.Command{Name: "cmd12", Fn: log}

	cmd2   = &cli.Command{Name: "cmd2", Fn: log}
	cmd21  = &cli.Command{Name: "cmd21", Fn: log}
	cmd22  = &cli.Command{Name: "cmd22", Fn: log}
	cmd221 = &cli.Command{Name: "cmd221", Fn: log}
	cmd222 = &cli.Command{Name: "cmd222", Fn: log}
	cmd223 = &cli.Command{Name: "cmd223", Fn: log}
)
