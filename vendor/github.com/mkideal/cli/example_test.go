package cli_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"

	"github.com/mkideal/cli"
)

func ExampleRun() {
	type argT struct {
		Flag string `cli:"f"`
	}
	cli.RunWithArgs(new(argT), []string{"app", "-f=xxx"}, func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)
		ctx.String("flag: %s\n", argv.Flag)
		return nil
	})
	// Output:
	// flag: xxx
}

// This example demonstrates how to use short and long format flag
func ExampleParse_shortAndLongFlagName() {
	// argument object
	type argT struct {
		Port int `cli:"p,port" usage:"listening port"`
	}

	for _, args := range [][]string{
		[]string{"app", "-p", "8080"},
		[]string{"app", "-p8081"},
		[]string{"app", "-p=8082"},
		[]string{"app", "--port", "8083"},
		[]string{"app", "--port=8084"},
	} {
		cli.RunWithArgs(&argT{}, args, func(ctx *cli.Context) error {
			argv := ctx.Argv().(*argT)
			ctx.String("port=%d\n", argv.Port)
			return nil
		})
	}
	// Output:
	// port=8080
	// port=8081
	// port=8082
	// port=8083
	// port=8084
}

// This example demonstrates how to use default value
func ExampleParse_defaultValue() {
	type argT1 struct {
		Port int `cli:"p,port" usage:"listening port" dft:"8080"`
	}
	type argT2 struct {
		Port int `cli:"p,port" usage:"listening port" dft:"$CLI_TEST_HTTP_PORT"`
	}
	type argT3 struct {
		Port int `cli:"p,port" usage:"listening port" dft:"$CLI_TEST_HTTP_PORT+800"`
	}
	type argT4 struct {
		DevDir string `cli:"dir" usage:"develope directory" dft:"$CLI_TEST_DEV_PARENT_DIR/dev"`
	}

	os.Setenv("CLI_TEST_DEV_PARENT_DIR", "/home")
	os.Setenv("CLI_TEST_HTTP_PORT", "8000")

	for _, tt := range []struct {
		argv interface{}
		args []string
	}{
		{new(argT1), []string{"app"}},
		{new(argT2), []string{"app"}},
		{new(argT3), []string{"app"}},
		{new(argT4), []string{"app"}},
		{new(argT4), []string{"app", "--dir=/dev"}},
	} {
		cli.RunWithArgs(tt.argv, tt.args, func(ctx *cli.Context) error {
			ctx.String("argv=%v\n", ctx.Argv())
			return nil
		})
	}
	// Output:
	// argv=&{8080}
	// argv=&{8000}
	// argv=&{8800}
	// argv=&{/home/dev}
	// argv=&{/dev}
}

// This example demonstrates to use Slice and Map
func ExampleParse_sliceAndMap() {
	type argT1 struct {
		Slice []uint32 `cli:"U,u32-slice" usage:"uint32 slice"`
	}
	type argT2 struct {
		Slice []string `cli:"S,str-slice" usage:"string slice"`
	}
	type argT3 struct {
		Slice []bool `cli:"B,bool-slice" usage:"boolean slice"`
	}
	type argT4 struct {
		MapA map[string]int  `cli:"A" usage:"string => int"`
		MapB map[int]int     `cli:"B" usage:"int => int"`
		MapC map[int]string  `cli:"C" usage:"int => string"`
		MapD map[string]bool `cli:"D" usage:"string => bool"`
	}

	for _, tt := range []struct {
		argv interface{}
		args []string
	}{
		{new(argT1), []string{"app", "-U1", "-U2"}},
		{new(argT1), []string{"app", "-U", "1", "-U", "2"}},
		{new(argT1), []string{"app", "--u32-slice", "1", "--u32-slice", "2"}},
		{new(argT2), []string{"app", "-Shello", "-Sworld"}},
		{new(argT2), []string{"app", "-S", "hello", "-S", "world"}},
		{new(argT2), []string{"app", "--str-slice", "hello", "--str-slice", "world"}},
		{new(argT3), []string{"app", "-Btrue", "-Bfalse"}},
		{new(argT3), []string{"app", "-B", "true", "-B", "false"}},
		{new(argT3), []string{"app", "--bool-slice", "true", "--bool-slice", "false"}},

		{new(argT4), []string{"app",
			"-Ax=1",
			"-B", "1=2",
			"-C1=a",
			"-Dx",
		}},
	} {
		cli.RunWithArgs(tt.argv, tt.args, func(ctx *cli.Context) error {
			ctx.String("argv=%v\n", ctx.Argv())
			return nil
		})
	}
	// Output:
	// argv=&{[1 2]}
	// argv=&{[1 2]}
	// argv=&{[1 2]}
	// argv=&{[hello world]}
	// argv=&{[hello world]}
	// argv=&{[hello world]}
	// argv=&{[true false]}
	// argv=&{[true false]}
	// argv=&{[true false]}
	// argv=&{map[x:1] map[1:2] map[1:a] map[x:true]}
}

func ExampleCommand() {
	root := &cli.Command{
		Name: "app",
	}

	type childT struct {
		S string `cli:"s" usage:"string flag"`
		B bool   `cli:"b" usage:"boolean flag"`
	}
	root.Register(&cli.Command{
		Name:        "child",
		Aliases:     []string{"sub"},
		Desc:        "child command",
		Text:        "detailed description for command",
		Argv:        func() interface{} { return new(childT) },
		CanSubRoute: true,
		NoHook:      true,
		NoHTTP:      true,
		NumArg:      cli.ExactN(1),
		HTTPRouters: []string{"/v1/child", "/v2/child"},
		HTTPMethods: []string{"GET", "POST"},

		OnRootPrepareError: func(err error) error {
			return err
		},
		OnBefore: func(ctx *cli.Context) error {
			ctx.String("OnBefore\n")
			return nil
		},
		OnAfter: func(ctx *cli.Context) error {
			ctx.String("OnAfter\n")
			return nil
		},
		OnRootBefore: func(ctx *cli.Context) error {
			ctx.String("OnRootBefore\n")
			return nil
		},
		OnRootAfter: func(ctx *cli.Context) error {
			ctx.String("OnRootAfter\n")
			return nil
		},

		Fn: func(ctx *cli.Context) error {
			return nil
		},
	})
}

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
				if valField.Kind() == reflect.Int && valField.CanSet() {
					valField.SetInt(2)
				}
				if valField.Kind() == reflect.String && valField.CanSet() {
					valField.SetString("B")
				}
			}
		}
	}
	return nil
}

type config3 struct {
	A int
	B string
}

// This example demonstrates how to use custom parser
func ExampleRegisterFlagParser() {
	// register parser factory function
	cli.RegisterFlagParser("myparser", newMyParser)

	type argT struct {
		Cfg3 config3 `cli:"cfg3" parser:"myparser"`
	}

	args := []string{"app",
		`--cfg3`, `hello`,
	}

	cli.RunWithArgs(new(argT), args, func(ctx *cli.Context) error {
		ctx.JSON(ctx.Argv())
		return nil
	})
	// Output:
	// {"Cfg3":{"A":2,"B":"B"}}
}

type helloT struct {
	cli.Helper
	Name string `cli:"name" usage:"tell me your name" dft:"world"`
	Age  uint8  `cli:"a,age" usage:"tell me your age" dft:"100"`
}

// This is a HelloWorld example
func Example_hello() {
	args := []string{"app", "--name=Cliper"}
	cli.RunWithArgs(new(helloT), args, func(ctx *cli.Context) error {
		argv := ctx.Argv().(*helloT)
		ctx.String("Hello, %s! Your age is %d?\n", argv.Name, argv.Age)
		return nil
	})
	// Output:
	// Hello, Cliper! Your age is 100?
}

func ExampleNumArgFunc_exactN() {
	type argT struct {
		cli.Helper
		I int `cli:"i" usage:"int"`
	}

	app := func() *cli.Command {
		return &cli.Command{
			Name:    "hahaha",
			Argv:    func() interface{} { return new(argT) },
			NumArg:  cli.ExactN(1),
			UsageFn: func() string { return "usage function" },
			Fn: func(ctx *cli.Context) error {
				return nil
			},
		}
	}
	fmt.Println(app().Run([]string{}))
	fmt.Println(app().Run([]string{"-i", "1", "b"}))
	fmt.Println(app().Run([]string{"-i", "1", "b", "c"}))

	// Output:
	// usage function<nil>
	// <nil>
	// usage function<nil>
}

func ExampleNumArgFunc_atLeast() {
	type argT struct {
		cli.Helper
		I int `cli:"i" usage:"int"`
	}

	app := func() *cli.Command {
		return &cli.Command{
			Name:    "hahaha",
			Argv:    func() interface{} { return new(argT) },
			NumArg:  cli.AtLeast(1),
			UsageFn: func() string { return "usage function" },
			Fn: func(ctx *cli.Context) error {
				return nil
			},
		}
	}
	fmt.Println(app().Run([]string{}))
	fmt.Println(app().Run([]string{"-i", "1", "b"}))
	fmt.Println(app().Run([]string{"-i", "1", "b", "c"}))

	// Output:
	// usage function<nil>
	// <nil>
	// <nil>
}

func ExampleNumArgFunc_atMost() {
	type argT struct {
		cli.Helper
		I int `cli:"i" usage:"int"`
	}

	app := func() *cli.Command {
		return &cli.Command{
			Name:    "hahaha",
			Argv:    func() interface{} { return new(argT) },
			NumArg:  cli.AtMost(2),
			UsageFn: func() string { return "usage function" },
			Fn: func(ctx *cli.Context) error {
				return nil
			},
		}
	}
	fmt.Println(app().Run([]string{}))
	fmt.Println(app().Run([]string{"-i", "1", "b"}))
	fmt.Println(app().Run([]string{"-i", "1", "b", "c"}))
	fmt.Println(app().Run([]string{"-i", "1", "b", "c", "d"}))

	// Output:
	// <nil>
	// <nil>
	// <nil>
	// usage function<nil>
}

func ExampleNumOptionFunc_exactN() {
	type argT struct {
		cli.Helper
		I int `cli:"i" usage:"int"`
	}

	app := func() *cli.Command {
		return &cli.Command{
			Name:      "hahaha",
			Argv:      func() interface{} { return new(argT) },
			NumOption: cli.ExactN(1),
			UsageFn:   func() string { return "usage function" },
			Fn: func(ctx *cli.Context) error {
				return nil
			},
		}
	}
	fmt.Println(app().Run([]string{}))
	fmt.Println(app().Run([]string{"-i", "1"}))

	// Output:
	// usage function<nil>
	// <nil>
}

func ExampleNumOptionFunc_atLeast() {
	type argT struct {
		cli.Helper
		I int `cli:"i" usage:"int"`
	}

	app := func() *cli.Command {
		return &cli.Command{
			Name:      "hahaha",
			Argv:      func() interface{} { return new(argT) },
			NumOption: cli.AtLeast(1),
			UsageFn:   func() string { return "usage function" },
			Fn: func(ctx *cli.Context) error {
				return nil
			},
		}
	}
	fmt.Println(app().Run([]string{}))
	fmt.Println(app().Run([]string{"-i", "1"}))

	// Output:
	// usage function<nil>
	// <nil>
}

func ExampleNumOptionFunc_atMost() {
	type argT struct {
		cli.Helper
		I int `cli:"i" usage:"int"`
		J int `cli:"j" usage:"int"`
		K int `cli:"k" usage:"int"`
	}

	app := func() *cli.Command {
		return &cli.Command{
			Name:      "hahaha",
			Argv:      func() interface{} { return new(argT) },
			NumOption: cli.AtMost(2),
			UsageFn:   func() string { return "usage function" },
			Fn: func(ctx *cli.Context) error {
				return nil
			},
		}
	}
	fmt.Println(app().Run([]string{}))
	fmt.Println(app().Run([]string{"-i", "1"}))
	fmt.Println(app().Run([]string{"-i", "1", "-j=2"}))
	fmt.Println(app().Run([]string{"-i", "1", "-j=2", "-k=3"}))

	// Output:
	// <nil>
	// <nil>
	// <nil>
	// usage function<nil>
}

type config1 struct {
	A string
	B int
}

type config2 struct {
	C string
	D bool
}

// This example demonstrates how to use builtin praser(json,jsonfile)
func ExampleFlagParser() {
	type argT struct {
		Cfg1 config1 `cli:"cfg1" parser:"json"`
		Cfg2 config2 `cli:"cfg2" parser:"jsonfile"`
	}
	jsonfile := "1.json"
	args := []string{"app",
		`--cfg1`, `{"A": "hello", "B": 2}`,
		`--cfg2`, jsonfile,
	}
	ioutil.WriteFile(jsonfile, []byte(`{"C": "world", "D": true}`), 0644)
	defer os.Remove(jsonfile)

	cli.RunWithArgs(new(argT), args, func(ctx *cli.Context) error {
		ctx.JSON(ctx.Argv())
		return nil
	})
	// Output:
	// {"Cfg1":{"A":"hello","B":2},"Cfg2":{"C":"world","D":true}}
}

func ExampleHelper() {
	type argT struct {
		cli.Helper
	}
	cli.RunWithArgs(new(argT), []string{"app", "-h"}, func(ctx *cli.Context) error {
		return nil
	})
	// Output:
	// Options:
	//
	//   -h, --help   display help information
}
