package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/labstack/gommon/color"
	"github.com/mkideal/cli"
)

type argT struct {
	cli.Helper
	File  string `cli:"F,file" usage:"create source file for new command, default <commandName>.go" name:"NAME"`
	Force bool   `cli:"f,force" usage:"force create file if exists" dft:"false"`

	Package      string `cli:"p,package" usage:"dest package name, default <basedir FILE>"`
	Desc         string `cli:"s,desc" usage:"command description"`
	CanSubRoute  bool   `cli:"csr,can-sub-route" usage:"set CanSubRoute attribute for new command" dft:"false"`
	ArgvTypeName string `cli:"argv-type-name" usage:"argv type, default <commandName>T, e.g. command name is hello, then defaut argv type is helloT"`

	Name string `cli:"-"`
}

func (argv *argT) Validate(ctx *cli.Context) error {
	yellow := ctx.Color().Yellow
	args := ctx.Args()
	if len(args) == 0 {
		return fmt.Errorf("command name missing")
	}
	if len(args) > 1 {
		return fmt.Errorf("too many(%d) commands, can only create one command", len(args))
	}
	argv.Name = ctx.Args()[0]
	if cli.IsValidCommandName(argv.Name) {
		return fmt.Errorf("invalid command name: %s", yellow(argv.Name))
	}
	if argv.File == "" {
		argv.File = argv.Name + ".go"
	}
	if argv.ArgvTypeName == "" {
		argv.ArgvTypeName = argv.Name + "T"
	}
	if argv.Package == "" {
		path, _ := filepath.Split(argv.File)
		pwd, err := os.Getwd()
		if err != nil {
			return err
		}
		if path != "" {
			if err := os.MkdirAll(path, 0755); err != nil {
				return err
			}
			if err := os.Chdir(path); err != nil {
				return err
			}
		}
		cmd := exec.Command("go", "list", "-f", "{{.Name}}")
		output, err := cmd.Output()
		if err != nil {
			argv.Package = filepath.Base(path)
		} else {
			argv.Package = string(output)
		}
		if err := os.Chdir(pwd); err != nil {
			return err
		}
		argv.Package = strings.TrimSpace(argv.Package)
	}
	return nil
}

func run(ctx *cli.Context, argv *argT) error {
	flag := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	if !argv.Force {
		flag |= os.O_EXCL
	}
	file, err := os.OpenFile(argv.File, flag, 0666)
	if err != nil {
		return err
	}
	t, err := template.New("clier").Parse(commandTpl)
	if err != nil {
		return err
	}
	return t.Execute(file, argv)
}

var commandTpl = `package {{.Package}}

import (
	"github.com/mkideal/cli"
)

type {{.ArgvTypeName}} struct {
	cli.Helper
}

var {{.Name}} = &cli.Command{
	Name: "{{.Name}}",
	Desc: "{{.Desc}}",
	Argv: func() interface{} { return new({{.ArgvTypeName}}) },
	CanSubRoute: {{.CanSubRoute}},

	Fn: func(ctx *cli.Context) error {
		argv := ctx.Argv().(*{{.ArgvTypeName}})
		if argv.Help {
			ctx.WriteUsage()
			return nil
		}
		//TODO: write your code here...
		return nil
	},
}
`

func main() {
	os.Exit(cli.Run(new(argT), func(ctx *cli.Context) error {
		cli.SetUsageStyle(cli.ManualStyle)
		argv := ctx.Argv().(*argT)
		if argv.Help {
			ctx.WriteUsage()
			return nil
		}
		return run(ctx, argv)
	}, fmt.Sprintf(`%s used to create a new command for github.com/mkideal/cli

%s: clier [OPTIONS] COMMAND-NAME

%s:
	clier hello
	clier -f -s "balabalabala" hello
	clier -p balabala hello`, color.Bold("clier"), color.Bold("Usage"), color.Bold("Examples"))))
}
