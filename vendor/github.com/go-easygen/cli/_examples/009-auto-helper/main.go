package main

import (
	"os"

	"github.com/mkideal/cli"
)

type argT struct {
	Help bool `cli:"h,help" usage:"show help"`
}

// AutoHelp implements cli.AutoHelper interface
// NOTE: cli.Helper is a predefined type which implements cli.AutoHelper
func (argv *argT) AutoHelp() bool {
	return argv.Help
}

func main() {
	os.Exit(cli.Run(new(argT), func(ctx *cli.Context) error {
		return nil
	}))
}
