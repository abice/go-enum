package main

import (
	"fmt"
	"os"
	"time"

	"github.com/mkideal/cli"
)

type argT struct {
	cli.Helper
	Wait  uint `cli:"wait" usage:"seconds for waiting" dft:"10"`
	Error bool `cli:"e" usage:"create an error"`
}

const successResponsePrefix = "start ok"

func main() {
	if err := cli.Root(root,
		cli.Tree(daemon),
	).Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var root = &cli.Command{
	Argv: func() interface{} { return new(argT) },
	Fn: func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)
		if argv.Error {
			err := fmt.Errorf("occurs error")
			cli.DaemonResponse(err.Error())
			return err
		}
		cli.DaemonResponse(successResponsePrefix)
		<-time.After(time.Duration(argv.Wait) * time.Second)
		return nil
	},
}

var daemon = &cli.Command{
	Name: "daemon",
	Argv: func() interface{} { return new(argT) },
	Fn: func(ctx *cli.Context) error {
		return cli.Daemon(ctx, successResponsePrefix)
	},
}
