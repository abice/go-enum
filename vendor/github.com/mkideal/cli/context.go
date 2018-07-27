package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/labstack/gommon/color"
	"github.com/mattn/go-colorable"
	"github.com/mkideal/pkg/debug"
)

type (
	// Context provides running context
	Context struct {
		router     []string
		path       string
		argvList   []interface{}
		nativeArgs []string
		flagSet    *flagSet
		command    *Command
		writer     io.Writer
		color      color.Color

		HTTPRequest  *http.Request
		HTTPResponse http.ResponseWriter
	}

	// Validator validates flag before running command
	Validator interface {
		Validate(*Context) error
	}

	// AutoHelper represents interface for showing help information automatically
	AutoHelper interface {
		AutoHelp() bool
	}
)

func newContext(path string, router, args []string, argvList []interface{}, clr color.Color) (*Context, error) {
	ctx := &Context{
		path:       path,
		router:     router,
		argvList:   argvList,
		nativeArgs: args,
		color:      clr,
		flagSet:    newFlagSet(),
	}
	if !isEmptyArgvList(argvList) {
		ctx.flagSet = parseArgvList(args, argvList, ctx.color)
		if ctx.flagSet.err != nil {
			return ctx, ctx.flagSet.err
		}
	}
	return ctx, nil
}

// Path returns full command name
// `./app hello world -a --xyz=1` will returns "hello world"
func (ctx *Context) Path() string {
	return ctx.path
}

// Router returns full command name with string array
// `./app hello world -a --xyz=1` will returns ["hello" "world"]
func (ctx *Context) Router() []string {
	return ctx.router
}

// NativeArgs returns native args
// `./app hello world -a --xyz=1` will return ["-a" "--xyz=1"]
func (ctx *Context) NativeArgs() []string {
	return ctx.nativeArgs
}

// Args returns free args
// `./app hello world -a=1 abc xyz` will return ["abc" "xyz"]
func (ctx *Context) Args() []string {
	return ctx.flagSet.args
}

// NArg returns length of Args
func (ctx *Context) NArg() int {
	return len(ctx.flagSet.args)
}

// NOpt returns num of options
func (ctx *Context) NOpt() int {
	if ctx.flagSet == nil || ctx.flagSet.flagSlice == nil {
		return 0
	}
	n := 0
	for _, fl := range ctx.flagSet.flagSlice {
		if fl.isSet {
			n++
		}
	}
	return n
}

// Argv returns parsed args object
func (ctx *Context) Argv() interface{} {
	if ctx.argvList == nil || len(ctx.argvList) == 0 {
		return nil
	}
	return ctx.argvList[0]
}

func (ctx *Context) RootArgv() interface{} {
	if isEmptyArgvList(ctx.argvList) {
		return nil
	}
	index := len(ctx.argvList) - 1
	return ctx.argvList[index]
}

func (ctx *Context) GetArgvList(curr interface{}, parents ...interface{}) error {
	if isEmptyArgvList(ctx.argvList) {
		return argvError{isEmpty: true}
	}
	for i, argv := range append([]interface{}{curr}, parents...) {
		if argv == nil {
			continue
		}
		if i >= len(ctx.argvList) {
			return argvError{isOutOfRange: true}
		}
		if ctx.argvList[i] == nil {
			return argvError{ith: i, msg: "source is nil"}
		}

		buf := bytes.NewBufferString("")
		if err := json.NewEncoder(buf).Encode(ctx.argvList[i]); err != nil {
			return err
		}
		if err := json.NewDecoder(buf).Decode(argv); err != nil {
			return err
		}
	}
	return nil
}

// IsSet determins whether `flag` is set
func (ctx *Context) IsSet(flag string, aliasFlags ...string) bool {
	fl, ok := ctx.flagSet.flagMap[flag]
	if ok {
		return fl.isSet
	}
	for _, alias := range aliasFlags {
		if fl, ok := ctx.flagSet.flagMap[alias]; ok {
			return fl.isSet
		}
	}
	return false
}

// FormValues returns parsed args as url.Values
func (ctx *Context) FormValues() url.Values {
	if ctx.flagSet == nil {
		debug.Panicf("ctx.flagSet == nil")
	}
	return ctx.flagSet.values
}

// Command returns current command instance
func (ctx *Context) Command() *Command {
	return ctx.command
}

// Usage returns current command's usage with current context
func (ctx *Context) Usage() string {
	return ctx.command.Usage(ctx)
}

// WriteUsage writes usage to writer
func (ctx *Context) WriteUsage() {
	ctx.String(ctx.Usage())
}

// Writer returns writer
func (ctx *Context) Writer() io.Writer {
	if ctx.writer == nil {
		ctx.writer = colorable.NewColorableStdout()
	}
	return ctx.writer
}

// Write implements io.Writer
func (ctx *Context) Write(data []byte) (n int, err error) {
	return ctx.Writer().Write(data)
}

// Color returns color instance
func (ctx *Context) Color() *color.Color {
	return &ctx.color
}

// String writes formatted string to writer
func (ctx *Context) String(format string, args ...interface{}) *Context {
	fmt.Fprintf(ctx.Writer(), format, args...)
	return ctx
}

// JSON writes json string of obj to writer
func (ctx *Context) JSON(obj interface{}) *Context {
	data, err := json.Marshal(obj)
	if err == nil {
		fmt.Fprint(ctx.Writer(), string(data))
	}
	return ctx
}

// JSONln writes json string of obj end with "\n" to writer
func (ctx *Context) JSONln(obj interface{}) *Context {
	return ctx.JSON(obj).String("\n")
}

// JSONIndent writes pretty json string of obj to writer
func (ctx *Context) JSONIndent(obj interface{}, prefix, indent string) *Context {
	data, err := json.MarshalIndent(obj, prefix, indent)
	if err == nil {
		fmt.Fprint(ctx.Writer(), string(data))
	}
	return ctx
}

// JSONIndentln writes pretty json string of obj end with "\n" to writer
func (ctx *Context) JSONIndentln(obj interface{}, prefix, indent string) *Context {
	return ctx.JSONIndent(obj, prefix, indent).String("\n")
}
