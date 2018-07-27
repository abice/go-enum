package main

import (
	"fmt"
	"os"

	"github.com/mkideal/cli"
)

func main() {
	if err := cli.Root(root,
		cli.Tree(help),
		cli.Tree(daemon),
		cli.Tree(ping),
		cli.Tree(api,
			cli.Tree(build),
			cli.Tree(install),
		),
	).Run(os.Args[1:]); err != nil {
		fmt.Println(err)
	}
}

//------
// root
//------
var root = &cli.Command{
	Fn: func(ctx *cli.Context) error {
		ctx.WriteUsage()
		return nil
	},
}

//------
// help
//------
var help = &cli.Command{
	Name:        "help",
	Desc:        "display help",
	CanSubRoute: true,
	HTTPRouters: []string{"/help", "/v1/help"},
	HTTPMethods: []string{"GET"},

	Fn: cli.HelpCommandFn,
}

//--------
// daemon
//--------
type daemonT struct {
	cli.Helper
	Port uint16 `cli:"p,port" usage:"http port" dft:"8080"`
}

func (t *daemonT) Validate(ctx *cli.Context) error {
	if t.Port == 0 {
		return fmt.Errorf("please don't use 0 as http port")
	}
	return nil
}

var daemon = &cli.Command{
	Name: "daemon",
	Desc: "startup app as daemon",
	Argv: func() interface{} { return new(daemonT) },
	Fn: func(ctx *cli.Context) error {
		argv := ctx.Argv().(*daemonT)
		addr := fmt.Sprintf(":%d", argv.Port)
		ctx.String("http addr: %s\n", addr)

		r := ctx.Command().Root()
		if err := r.RegisterHTTP(ctx); err != nil {
			return err
		}
		return r.ListenAndServeHTTP(addr)
	},
}

//------
// ping
//------
var ping = &cli.Command{
	Name: "ping",
	Desc: "ping server",
	Fn: func(ctx *cli.Context) error {
		ctx.String("pong\n")
		return nil
	},
}

//-----
// api
//-----
var api = &cli.Command{
	Name: "api",
	Desc: "display all api",
	Fn: func(ctx *cli.Context) error {
		ctx.String("Commands:\n")
		ctx.String("    build\n")
		ctx.String("    install\n")
		return nil
	},
}

//-------
// build
//-------
type buildT struct {
	cli.Helper
	Dir string `cli:"dir" usage:"dest path" dft:"./"`
}

var build = &cli.Command{
	Name: "build",
	Desc: "build application",
	Argv: func() interface{} { return new(buildT) },
	Fn: func(ctx *cli.Context) error {
		argv := ctx.Argv().(*buildT)
		ctx.JSONIndentln(argv, "", "    ")
		return nil
	},
}

//---------
// install
//---------
var install = &cli.Command{
	Name: "install",
	Desc: "install application",
	Argv: func() interface{} { return new(buildT) },
	Fn: func(ctx *cli.Context) error {
		argv := ctx.Argv().(*buildT)
		ctx.JSONIndentln(argv, "", "    ")
		return nil
	},
}
