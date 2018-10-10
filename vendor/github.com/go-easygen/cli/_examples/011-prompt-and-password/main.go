package main

import (
	"os"

	"github.com/mkideal/cli"
)

type argT struct {
	cli.Helper
	Username string `cli:"u,username" usage:"github account" prompt:"type github account"`
	Password string `pw:"p,password" usage:"password of github account" prompt:"type the password"`
}

func main() {
	os.Exit(cli.Run(new(argT), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)
		ctx.String("username=%s, password=%s\n", argv.Username, argv.Password)
		return nil
	}))
}
