package main

import (
	"fmt"

	"github.com/mkideal/cli"
)

type argT struct {
	cli.Helper
	Age    int    `cli:"age" usage:"your age"`
	Gender string `cli:"g,gender" usage:"your gender" dft:"male"`
}

// Validate implements cli.Validator interface
func (argv *argT) Validate(ctx *cli.Context) error {
	if argv.Age < 0 || argv.Age > 300 {
		return fmt.Errorf("age %d out of range", argv.Age)
	}
	if argv.Gender != "male" && argv.Gender != "female" {
		return fmt.Errorf("invalid gender %s", ctx.Color().Yellow(argv.Gender))
	}
	return nil
}

func main() {
	cli.Run(new(argT), func(ctx *cli.Context) error {
		ctx.JSONln(ctx.Argv())
		return nil
	})
}
