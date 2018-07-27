Command line interface
======================

[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/mkideal/cli/master/LICENSE) [![Travis branch](https://img.shields.io/travis/mkideal/cli/master.svg)](https://travis-ci.org/mkideal/cli) [![Coverage Status](https://coveralls.io/repos/github/mkideal/cli/badge.svg?branch=master)](https://coveralls.io/github/mkideal/cli?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/mkideal/cli)](https://goreportcard.com/report/github.com/mkideal/cli) [![GoDoc](https://godoc.org/github.com/mkideal/cli?status.svg)](https://godoc.org/github.com/mkideal/cli)

Screenshot
----------

![screenshot2](http://www.mkideal.com/images/screenshot2.png)

Key features
------------

-	Lightweight and easy to use.
-	Defines flag by tag, e.g. flag name(short or/and long), description, default value, password, prompt and so on.
-	Type safety.
-	Output looks very nice.
-	Supports custom Validator.
-	Supports slice and map as a flag.
-	Supports any type as a flag field which implements cli.Decoder interface.
-	Supports any type as a flag field which uses FlagParser.
-	Suggestions for command.(e.g. `hl` => `help`, "veron" => "version").
-	Supports default value for flag, even expression about env variable(e.g. `dft:"$HOME/dev"`).
-	Supports editor like `git commit` command.(See example [21](http://www.mkideal.com/golang/cli-examples.html#example-21-editor) and [22](http://www.mkideal.com/golang/cli-examples.html#example-22-custom-editor)\)

API documentation
-----------------

See [**godoc**](https://godoc.org/github.com/mkideal/cli)

Examples
--------

-	[Example 1: Hello world](#example-1-hello)
-	[Example 2: How to use **flag**](#example-2-flag)
-	[Example 3: How to use **required** flag](#example-3-required-flag)
-	[Example 4: How to use **default** flag](#example-4-default-flag)
-	[Example 5: How to use **slice**](#example-5-slice)
-	[Example 6: How to use **map**](#example-6-map)
-	[Example 7: Usage of **force** flag](#example-7-force-flag)
-	[Example 8: Usage of **child command**](#example-8-child-command)
-	[Example 9: **Auto help**](#example-9-auto-help)
-	[Example 10: Usage of **Validator**](#example-10-usage-of-validator)
-	[Example 11: **Prompt** and **Password**](#example-11-prompt-and-password)
-	[Example 12: How to use **Decoder**](#example-12-decoder)
-	[Example 13: Builtin Decoder: **PidFile**](#example-13-pid-file)
-	[Example 14: Builtin Decoder: **Time** and **Duration**](#example-14-time-and-duration)
-	[Example 15: Builtin Decoder: **File**](#example-15-file)
-	[Example 16: **Parser**](#example-16-parser)
-	[Example 17: Builtin Parser: **JSONFileParser**](#example-17-json-file)
-	[Example 18: How to use **custom parser**](#example-18-custom-parser)
-	[Example 19: How to use **Hooks**](#example-19-hooks)
-	[Example 20: How to use **Daemon**](#example-20-daemon)
-	[Example 21: How to use **Editor**](#example-21-editor)
-	[Example 22: Custom **Editor**](#example-22-custom-editor)

### Example 1: Hello

[back to **examples**](#examples)

```go
// main.go
// This is a HelloWorld-like example

package main

import (
	"github.com/mkideal/cli"
)

type argT struct {
	Name string `cli:"name" usage:"tell me your name"`
}

func main() {
	cli.Run(new(argT), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)
		ctx.String("Hello, %s!\n", argv.Name)
		return nil
	})
}
```

```sh
$ go build -o hello
$ ./hello --name Clipher
Hello, Clipher!
```

### Example 2: Flag

[back to **examples**](#examples)

```go
// main.go
// This example show basic usage of flag

package main

import (
	"github.com/mkideal/cli"
)

type argT struct {
	cli.Helper
	Port int  `cli:"p,port" usage:"short and long format flags both are supported"`
	X    bool `cli:"x" usage:"boolean type"`
	Y    bool `cli:"y" usage:"boolean type, too"`
}

func main() {
	cli.Run(new(argT), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)
		ctx.String("port=%d, x=%v, y=%v\n", argv.Port, argv.X, argv.Y)
		return nil
	})
}
```

```sh
$ go build -o app
$ ./app -h
Options:

  -h, --help     display help information
  -p, --port     short and long format flags both are supported
  -x             boolean type
  -y             boolean type, too
$ ./app -p=8080 -x
port=8080, x=true, y=false
$ ./app -p 8080 -x=true
port=8080, x=true, y=false
$ ./app -p8080 -y true
port=8080, x=false, y=true
$ ./app --port=8080 -xy
port=8080, x=true, y=true
$ ./app --port 8080 -yx
port=8080, x=true, y=true
```

### Example 3: Required flag

[back to **examples**](#examples)

```go
// main.go
// This example show how to use required flag

package main

import (
	"github.com/mkideal/cli"
)

type argT struct {
	cli.Helper
	Id uint8 `cli:"*id" usage:"this is a required flag, note the *"`
}

func main() {
	cli.Run(new(argT), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)
		ctx.String("%d\n", argv.Id)
		return nil
	})
}
```

```sh
$ go build -o app
$ ./app
ERR! required argument --id missing
$ ./app --id=2
2
```

### Example 4: Default flag

[back to **examples**](#examples)

```go
// main.go
// This example show how to use default flag

package main

import (
	"github.com/mkideal/cli"
)

type argT struct {
	cli.Helper
	Basic  int    `cli:"basic" usage:"basic usage of default" dft:"2"`
	Env    string `cli:"env" usage:"env variable as default" dft:"$HOME"`
	Expr   int    `cli:"expr" usage:"expression as default" dft:"$BASE_PORT+1000"`
	DevDir string `cli:"devdir" usage:"directory of developer" dft:"$HOME/dev"`
}

func main() {
	cli.Run(new(argT), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)
		ctx.String("%d, %s, %d, %s\n", argv.Basic, argv.Env, argv.Expr, argv.DevDir)
		return nil
	})
}
```

```sh
$ go build -o app
$ ./app -h
Options:

  -h, --help                       display help information
      --basic[=2]                  basic usage of default
      --env[=$HOME]                env variable as default
      --expr[=$BASE_PORT+1000]     expression as default
      --devdir[=$HOME/dev]         directory of developer
$ ./app
2, /Users/wang, 1000, /Users/wang/dev
$ BASE_PORT=8000 ./app --basic=3
3, /Users/wang, 9000, /Users/wang/dev
```

### Example 5: Slice

[back to **examples**](#examples)

```go
// main.go
// This example show how to use slice as a flag

package main

import (
	"github.com/mkideal/cli"
)

type argT struct {
	// []bool, []int, []float32, ... supported too.
	Friends []string `cli:"F" usage:"my friends"`
}

func main() {
	cli.Run(new(argT), func(ctx *cli.Context) error {
		ctx.JSONln(ctx.Argv())
		return nil
	})
}
```

```sh
$ go build -o app
$ ./app
{"Friends":null}
$ ./app -FAlice -FBob -F Charlie
{"Friends":["Alice","Bob","Charlie"]}
```

### Example 6: Map

[back to **examples**](#examples)

```go
// main.go
// This example show how to use map as a flag

package main

import (
	"github.com/mkideal/cli"
)

type argT struct {
	Macros map[string]int `cli:"D" usage:"define macros"`
}

func main() {
	cli.Run(new(argT), func(ctx *cli.Context) error {
		ctx.JSONln(ctx.Argv())
		return nil
	})
}
```

```sh
$ go build -o app
$ ./app
{"Macros":null}
$ ./app -Dx=not-a-number
ERR! `not-a-number` couldn't converted to an int value
$ ./app -Dx=1 -D y=2
{"Macros":{"x":1,"y":2}}
```

### Example 7: Force flag

[back to **examples**](#examples)

```go
// main.go
// This example show usage of force flag
// Force flag has prefix !, and must be a boolean.
// Will prevent validating flags if some force flag assigned true

package main

import (
	"github.com/mkideal/cli"
)

type argT struct {
	Version  bool `cli:"!v" usage:"force flag, note the !"`
	Required int  `cli:"*r" usage:"required flag"`
}

func main() {
	cli.Run(new(argT), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)
		if argv.Version {
			ctx.String("v0.0.1\n")
		}
		return nil
	})
}
```

```sh
$ go build -o app
$ ./app
ERR! required argument -r missing

# -v is a force flag, and assigned true, so `ERR` disappear.
$ ./app -v
v0.0.1
```

### Example 8: Child command

[back to **examples**](#examples)

```go
// main.go
// This example demonstrates usage of child command

package main

import (
	"fmt"
	"os"

	"github.com/mkideal/cli"
)

func main() {
	if err := cli.Root(root,
		cli.Tree(help),
		cli.Tree(child),
	).Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var help = cli.HelpCommand("display help information")

// root command
type rootT struct {
	cli.Helper
	Name string `cli:"name" usage:"your name"`
}

var root = &cli.Command{
	Desc: "this is root command",
	// Argv is a factory function of argument object
	// ctx.Argv() is if Command.Argv == nil or Command.Argv() is nil
	Argv: func() interface{} { return new(rootT) },
	Fn: func(ctx *cli.Context) error {
		argv := ctx.Argv().(*rootT)
		ctx.String("Hello, root command, I am %s\n", argv.Name)
		return nil
	},
}

// child command
type childT struct {
	cli.Helper
	Name string `cli:"name" usage:"your name"`
}

var child = &cli.Command{
	Name: "child",
	Desc: "this is a child command",
	Argv: func() interface{} { return new(childT) },
	Fn: func(ctx *cli.Context) error {
		argv := ctx.Argv().(*childT)
		ctx.String("Hello, child command, I am %s\n", argv.Name)
		return nil
	},
}
```

```sh
$ go build -o app

# help for root
# equivalent to "./app -h"
$ ./app help
this is root command

Options:

  -h, --help     display help information
      --name     your name

Commands:
  help    display help information
  child   this is a child command

# help for specific command
# equivalent to "./app child -h"
$ ./app help child
this is a child command

Options:

  -h, --help     display help information
      --name     your name

# execute root command
$ ./app --name 123
Hello, root command, I am 123

# execute child command
$ ./app child --name=123
Hello, child command, I am 123

# something wrong, but got a suggestion.
$ ./app chd
ERR! command chd not found
Did you mean child?
```

### Example 9: Auto help

[back to **examples**](#examples)

```go
// main.go
// This example demonstrates cli.AutoHelper

package main

import (
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
	cli.Run(new(argT), func(ctx *cli.Context) error {
		return nil
	})
}
```

```sh
$ go build -o app
$ ./app -h
Options:

  -h, --help     show help
```

Try comment `AutoHelp` method and rerun it.

### Example 10: Usage of Validator

[back to **examples**](#examples)

```go
// main.go
// This example demonstrates how to utilize Validator

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
```

```sh
$ go build -o app
$ ./app --age=-1
ERR! age -1 out of range
$ ./app --age=1000
ERR! age 1000 out of range
$ ./app -g balabala
ERR! invalid gender balabala
$ ./app --age 88 --gender female
{"Help":false,"Age":88,"Gender":"female"}
```

### Example 11: Prompt and Password

[back to **examples**](#examples)

```go
// main.go
// This example introduce prompt and pw tag

package main

import (
	"github.com/mkideal/cli"
)

type argT struct {
	cli.Helper
	Username string `cli:"u,username" usage:"github account" prompt:"type github account"`
	Password string `pw:"p,password" usage:"password of github account" prompt:"type the password"`
}

func main() {
	cli.Run(new(argT), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)
		ctx.String("username=%s, password=%s\n", argv.Username, argv.Password)
		return nil
	})
}
```

```sh
$ go build -o app
$ ./app
type github account: hahaha # visible
type the password: # invisible because of `pw` tag
username=hahaha, password=123456
```

### Example 12: Decoder

[back to **examples**](#examples)

```go
// main.go
// This example show how to use decoder

package main

import (
	"strings"

	"github.com/mkideal/cli"
)

type exampleDecoder struct {
	list []string
}

// Decode implements cli.Decoder interface
func (d *exampleDecoder) Decode(s string) error {
	d.list = strings.Split(s, ",")
	return nil
}

type argT struct {
	Example exampleDecoder `cli:"d" usage:"example decoder"`
}

func main() {
	cli.Run(new(argT), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)
		ctx.JSONln(argv.Example.list)
		return nil
	})
}
```

```sh
$ go build -o app
$ ./app -d a,b,c
["a","b","c"]
```

### Example 13: Pid file

[back to **examples**](#examples)

```go
// main.go
// This example show how to use builtin Decoder: PidFile

package main

import (
	"github.com/mkideal/cli"
	clix "github.com/mkideal/cli/ext"
)

type argT struct {
	cli.Helper
	PidFile clix.PidFile `cli:"pid" usage:"pid file" dft:"013-pidfile.pid"`
}

func main() {
	cli.Run(new(argT), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)

		if err := argv.PidFile.New(); err != nil {
			return err
		}
		defer argv.PidFile.Remove()

		return nil
	})
}
```

### Example 14: Time and Duration

[back to **examples**](#examples)

```go
// main.go
// This example show how to use builtin Decoder: Time and Duration

package main

import (
	"github.com/mkideal/cli"
	clix "github.com/mkideal/cli/ext"
)

type argT struct {
	Time     clix.Time     `cli:"t" usage:"time"`
	Duration clix.Duration `cli:"d" usage:"duration"`
}

func main() {
	cli.Run(new(argT), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)
		ctx.String("time=%v, duration=%v\n", argv.Time, argv.Duration)
		return nil
	})
}
```

```sh
$ go build -o app
$ ./app -t '2016-1-2 3:5' -d=10ms
time=2016-01-02 03:05:00 +0800 CST, duration=10ms
```

### Example 15: File

[back to **examples**](#examples)

```go
// main.go
// This example show how to use builtin Decoder: File

package main

import (
	"github.com/mkideal/cli"
	clix "github.com/mkideal/cli/ext"
)

type argT struct {
	Content clix.File `cli:"f,file" usage:"read content from file or stdin"`
}

func main() {
	cli.Run(new(argT), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)
		ctx.String(argv.Content.String())
		return nil
	})
}
```

```sh
$ go build -o app
# read from stdin
$ echo hello | ./app -f
hello
# read from file
$ echo hello > test.txt && ./app -f test.txt
hello
$ rm test.txt
```

### Example 16: Parser

[back to **examples**](#examples)

```go
// main.go
// This example introduce Parser
// `Parser` is another way to use custom type of data.
// Unlike `Decoder`, `Parser` used to parse string according to specific rule,
// like json,yaml and so on.
//
// Builtin parsers:
// * json
// * jsonfile

package main

import (
	"github.com/mkideal/cli"
)

type config struct {
	A string
	B int
	C bool
}

type argT struct {
	JSON config `cli:"c,config" usage:"parse json string" parser:"json"`
}

func main() {
	cli.Run(new(argT), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)
		ctx.JSONIndentln(argv.JSON, "", "    ")
		return nil
	})
}
```

```sh
$ go build -o app
$ ./app
{
    "A": "",
    "B": 0,
    "C": false
}
$ ./app -c '{"A": "hello", "b": 22, "C": true}'
{
    "A": "hello",
    "B": 22,
    "C": true
}
```

### Example 17: JSON file

[back to **examples**](#examples)

```go
// main.go
// This example show how to use builtin parser: jsonfile
// It's similar to json, but read string from file.

package main

import (
	"github.com/mkideal/cli"
)

type config struct {
	A string
	B int
	C bool
}

type argT struct {
	JSON config `cli:"c,config" usage:"parse json from file" parser:"jsonfile"`
}

func main() {
	cli.Run(new(argT), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)
		ctx.JSONIndentln(argv.JSON, "", "    ")
		return nil
	})
}
```

```sh
$ go build -o app
$ echo '{"A": "hello", "b": 22, "C": true}' > test.json
$ ./app -c test.json
{
    "A": "hello",
    "B": 22,
    "C": true
}
$ rm test.json
```

### Example 18: Custom parser

[back to **examples**](#examples)

```go
// main.go
// This example demonstrates how to use custom parser

package main

import (
	"reflect"

	"github.com/mkideal/cli"
)

type myParser struct {
	ptr interface{}
}

func newMyParser(ptr interface{}) cli.FlagParser {
	return &myParser{ptr}
}

// Parse implements FlagParser.Parse interface
func (parser *myParser) Parse(s string) error {
	typ := reflect.TypeOf(parser.ptr)
	val := reflect.ValueOf(parser.ptr)
	if typ.Kind() == reflect.Ptr {
		kind := reflect.Indirect(val).Type().Kind()
		if kind == reflect.Struct {
			typElem, valElem := typ.Elem(), val.Elem()
			numField := valElem.NumField()
			for i := 0; i < numField; i++ {
				_, valField := typElem.Field(i), valElem.Field(i)
				if valField.Kind() == reflect.Int &&
					valField.CanSet() {
					valField.SetInt(2)
				}
				if valField.Kind() == reflect.String &&
					valField.CanSet() {
					valField.SetString("B")
				}
			}
		}
	}
	return nil
}

type config struct {
	A int
	B string
}

type argT struct {
	Cfg config `cli:"cfg" parser:"myparser"`
}

func main() {
	// register parser factory function
	cli.RegisterFlagParser("myparser", newMyParser)

	cli.Run(new(argT), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)
		ctx.String("%v\n", argv.Cfg)
		return nil
	})
}
```

```sh
$ go build -o app
$ ./app
{0 }
$ ./app --cfg xxx
{2 B}
```

### Example 19: Hooks

[back to **examples**](#examples)

```go
// main.go
// This example demonstrates how to use hooks

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
		ctx.String("child1's OnBefore invoked\n")
		return nil
	},
	OnAfter: func(ctx *cli.Context) error {
		ctx.String("child1's OnAfter invoked\n")
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
```

```sh
$ go build -o app
# OnRootBefore => Fn => OnRootAfter
$ ./app
OnRootBefore invoked
exec root command
OnRootAfter invoked
# OnBefore => OnRootBefore => Fn => OnRootAfter => OnAfter
$ ./app child1
child1 OnBefore invoked
OnRootBefore invoked
exec child1 command
OnRootAfter invoked
child1 OnAfter invoked
# No hooks
$ ./app child2
exec child2 command
# OnRootBefore => Fn --> Error
$ ./app -e
OnRootBefore invoked
exec root command
root command returns error
# OnBefore => OnRootBefore => Fn --> Error
$ ./app child1 -e
child1 OnBefore invoked
OnRootBefore invoked
exec child1 command
child1 command returns error
```

### Example 20: Daemon

[back to **examples**](#examples)

```go
// main.go
// This example demonstrates how to use `Daemon`

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
```

```sh
$ go build -o daemon-app
$ ./daemone-app daemon
start ok
# Within 10 seconds, you will see process "./daemon-app"
$ ps | grep daemon-app
11913 ttys002    0:00.01 ./daemon-app
11915 ttys002    0:00.00 grep daemon-app
# After 10 seconds
$ ps | grep daemon-app
11936 ttys002    0:00.00 grep daemon-app
# try again with an error
$ ./daemon-app daemon -e
occurs error
$ ps | grep daemon-app
11936 ttys002    0:00.00 grep daemon-app
```

### Example 21: Editor

[back to **examples**](#examples)

```go
// main.go
// This example demonstrates how to use `editor`. This similar to git commit

package main

import (
	"github.com/mkideal/cli"
)

type argT struct {
	cli.Helper
	Msg string `edit:"m" usage:"message"`
}

func main() {
	cli.Run(new(argT), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)
		ctx.String("msg: %s", argv.Msg)
		return nil
	})
}
```

```sh
$ go build -o app
$ ./app -m "hello, editor"
msg: hello, editor
$ ./app # Then, launch a editor(default is vim) and type `hello, editor`, quit the editor
msg: hello, editor
```

### Example 22: Custom Editor

[back to **examples**](#examples)

```go
// main.go
// This example demonstrates specific editor.

package main

import (
	"os"

	"github.com/mkideal/cli"
)

type argT struct {
	cli.Helper
	Msg string `edit:"m" usage:"message"`
}

func main() {
	cli.GetEditor = func() (string, error) {
		if editor := os.Getenv("EDITOR"); editor != "" {
			return editor, nil
		}
		return cli.DefaultEditor, nil
	}
	cli.Run(new(argT), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)
		ctx.String("msg: %s", argv.Msg)
		return nil
	})
}
```

```sh
$ go build -o app
$ ./app -m "hello, editor"
msg: hello, editor
$ EDITOR=nano ./app # Then, launch nano and type `hello, editor`, quit the editor
msg: hello, editor
```
