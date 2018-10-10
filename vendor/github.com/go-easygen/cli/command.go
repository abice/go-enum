package cli

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/labstack/gommon/color"
	"github.com/mattn/go-colorable"
	"github.com/mkideal/pkg/debug"
)

var commandNameRegexp = regexp.MustCompile("^[a-zA-Z_0-9][a-zA-Z_\\-0-9]*$")

// IsValidCommandName validates name of command
func IsValidCommandName(commandName string) bool {
	return commandNameRegexp.MatchString(commandName)
}

type (
	// CommandFunc ...
	CommandFunc func(*Context) error

	// ArgvFunc ...
	ArgvFunc func() interface{}

	// NumCheckFunc represents function type which used to check num of args
	NumCheckFunc func(n int) bool

	// UsageFunc represents custom function of usage
	UsageFunc func() string
)

func ExactN(num int) NumCheckFunc  { return func(n int) bool { return n == num } }
func AtLeast(num int) NumCheckFunc { return func(n int) bool { return n >= num } }
func AtMost(num int) NumCheckFunc  { return func(n int) bool { return n <= num } }

type (
	// Command is the top-level instance in command-line app
	Command struct {
		Name    string   // Command name
		Aliases []string // Command aliases name
		Desc    string   // Command abstract
		Text    string   // Command detail description

		CanSubRoute bool
		NoHook      bool
		NoHTTP      bool
		Global      bool

		// functions
		Fn        CommandFunc // Command handler
		UsageFn   UsageFunc   // Custom usage function
		Argv      ArgvFunc    // Command argument factory function
		NumArg    NumCheckFunc
		NumOption NumCheckFunc

		HTTPRouters []string
		HTTPMethods []string

		// hooks for current command
		OnBefore func(*Context) error
		OnAfter  func(*Context) error

		// hooks for all commands if current command is root command
		OnRootPrepareError func(error) error
		OnRootBefore       func(*Context) error
		OnRootAfter        func(*Context) error

		routersMap map[string]string

		parent   *Command
		children []*Command

		isServer bool

		locker     sync.Mutex // protect following data
		usage      string
		usageStyle UsageStyle
	}

	// CommandTree represents a tree of commands
	CommandTree struct {
		command *Command
		forest  []*CommandTree
	}
)

// Register registers a child command
func (cmd *Command) Register(child *Command) *Command {
	if child == nil {
		debug.Panicf("command `%s` try register a nil command", cmd.Name)
	}
	if !IsValidCommandName(child.Name) {
		debug.Panicf("illegal command name `%s`", cmd.Name)
	}
	if cmd.children == nil {
		cmd.children = []*Command{}
	}
	if child.parent != nil {
		debug.Panicf("command `%s` has been child of `%s`", child.Name, child.parent.Name)
	}
	if cmd.findChild(child.Name) != nil {
		debug.Panicf("repeat register child `%s` for command `%s`", child.Name, cmd.Name)
	}
	if child.Aliases != nil {
		for _, alias := range child.Aliases {
			if cmd.findChild(alias) != nil {
				debug.Panicf("repeat register child `%s` for command `%s`", alias, cmd.Name)
			}
		}
	}
	cmd.children = append(cmd.children, child)
	child.parent = cmd

	return child
}

// RegisterFunc registers handler as child command
func (cmd *Command) RegisterFunc(name string, fn CommandFunc, argvFn ArgvFunc) *Command {
	return cmd.Register(&Command{Name: name, Fn: fn, Argv: argvFn})
}

// RegisterTree registers a command tree
func (cmd *Command) RegisterTree(forest ...*CommandTree) {
	for _, tree := range forest {
		cmd.Register(tree.command)
		if tree.forest != nil && len(tree.forest) > 0 {
			tree.command.RegisterTree(tree.forest...)
		}
	}
}

// Parent returns command's parent
func (cmd *Command) Parent() *Command {
	return cmd.parent
}

// IsServer returns command whether if run as server
func (cmd *Command) IsServer() bool {
	return cmd.isServer
}

// IsClient returns command whether if run as client
func (cmd *Command) IsClient() bool {
	return !cmd.IsServer()
}

// SetIsServer sets command running mode(server or not)
func (cmd *Command) SetIsServer(yes bool) {
	cmd.Root().isServer = yes
}

// Run runs the command with args
func (cmd *Command) Run(args []string) error {
	return cmd.RunWith(args, nil, nil)
}

// RunWith runs the command with args and writer,httpMethods
func (cmd *Command) RunWith(args []string, writer io.Writer, resp http.ResponseWriter, httpMethods ...string) error {
	fds := []uintptr{}
	if writer == nil {
		writer = colorable.NewColorableStdout()
		fds = append(fds, os.Stdout.Fd())
	}
	clr := color.Color{}
	colorSwitch(&clr, writer, fds...)

	var ctx *Context
	var suggestion string
	ctx, suggestion, err := cmd.prepare(clr, args, writer, resp, httpMethods...)
	if err == ExitError {
		return nil
	}

	if err != nil {
		if cmd.OnRootPrepareError != nil {
			err = cmd.OnRootPrepareError(err)
		}
		if err != nil {
			return wrapErr(err, suggestion, clr)
		}
		return nil
	}

	if ctx.command.NoHook {
		return ctx.command.Fn(ctx)
	}

	funcs := []func(*Context) error{
		ctx.command.OnBefore,
		cmd.OnRootBefore,
		ctx.command.Fn,
		cmd.OnRootAfter,
		ctx.command.OnAfter,
	}
	for _, f := range funcs {
		if f != nil {
			if err := f(ctx); err != nil {
				if err == ExitError {
					return nil
				}
				return err
			}
		}
	}
	return nil
}

func isEmptyArgvList(argvList []interface{}) bool {
	if argvList == nil {
		return true
	}
	for _, argv := range argvList {
		if argv != nil {
			return false
		}
	}
	return true
}

func (cmd *Command) argvList() []interface{} {
	argvList := make([]interface{}, 0, 1)
	if cmd.Argv != nil {
		argvList = append(argvList, cmd.Argv())
	} else {
		argvList = append(argvList, nil)
	}
	next := cmd.parent
	for next != nil {
		if next.Argv != nil && next.Global {
			argvList = append(argvList, next.Argv())
		} else {
			argvList = append(argvList, nil)
		}
		next = next.parent
	}
	return argvList
}

func (cmd *Command) prepare(clr color.Color, args []string, writer io.Writer, resp http.ResponseWriter, httpMethods ...string) (ctx *Context, suggestion string, err error) {
	// split args
	router := []string{}
	for _, arg := range args {
		if strings.HasPrefix(arg, dashOne) {
			break
		}
		router = append(router, arg)
	}
	path := strings.Join(router, " ")
	child, end := cmd.SubRoute(router)

	// if route fail
	if !child.CanSubRoute && end != len(router) {
		suggestions := cmd.Suggestions(path)
		buff := bytes.NewBufferString("")
		if suggestions != nil && len(suggestions) > 0 {
			if len(suggestions) == 1 {
				fmt.Fprintf(buff, "\nDid you mean %s?", clr.Bold(suggestions[0]))
			} else {
				fmt.Fprintf(buff, "\n\nDid you mean one of these?\n")
				for _, sug := range suggestions {
					fmt.Fprintf(buff, "    %s\n", sug)
				}
			}
		}
		suggestion = buff.String()
		err = throwCommandNotFound(clr.Yellow(path))
		return
	}

	methodAllowed := false
	if len(httpMethods) == 0 ||
		child.HTTPMethods == nil ||
		len(child.HTTPMethods) == 0 {
		methodAllowed = true
	} else {
		method := httpMethods[0]
		for _, m := range child.HTTPMethods {
			if method == m {
				methodAllowed = true
				break
			}
		}
	}
	if !methodAllowed {
		err = throwMethodNotAllowed(clr.Yellow(httpMethods[0]))
		return
	}

	// create argvList
	argvList := child.argvList()

	// create Context
	path = child.Path()
	ctx, err = newContext(path, router[:end], args[end:], argvList, clr)
	ctx.command = child
	ctx.writer = writer
	if !ctx.flagSet.hasForce {
		if !child.checkNumOption(ctx.NOpt()) || !ctx.command.checkNumArg(ctx.NArg()) {
			ctx.WriteUsage()
			err = ExitError
			return
		}
	}
	if err != nil {
		return
	}
	ctx.HTTPResponse = resp

	// auto help
	for _, argv := range argvList {
		if argv != nil {
			if helper, ok := argv.(AutoHelper); ok && helper.AutoHelp() {
				ctx.WriteUsage()
				err = ExitError
				return
			}
		}
	}

	if len(router) == 0 && cmd.Fn == nil {
		err = throwCommandNotFound(clr.Yellow(cmd.Name))
		return
	}

	if !ctx.flagSet.hasForce {
		for _, argv := range argvList {
			// validate argv if argv implements interface Validator
			if argv != nil {
				if validator, ok := argv.(Validator); ok {
					err = validator.Validate(ctx)
					if err != nil {
						return
					}
				}
			}
		}
	}

	return
}

func (cmd *Command) checkNumArg(num int) bool {
	return cmd.NumArg == nil || cmd.NumArg(num)
}

func (cmd *Command) checkNumOption(num int) bool {
	return cmd.NumOption == nil || cmd.NumOption(num)
}

// Usage returns the usage string of command
func (cmd *Command) Usage(ctx *Context) string {
	if cmd.UsageFn != nil {
		return cmd.UsageFn()
	}
	return cmd.defaultUsageFn(ctx)
}

func (cmd *Command) defaultUsageFn(ctx *Context) string {
	var (
		style = GetUsageStyle()
		clr   = *(ctx.Color())
	)

	// get usage form cache
	cmd.locker.Lock()
	tmpUsage := cmd.usage
	usageStyle := cmd.usageStyle
	cmd.locker.Unlock()
	if tmpUsage != "" && usageStyle == style {
		debug.Debugf("get usage of command %s from cache", clr.Bold(cmd.Name))
		return tmpUsage
	}

	buff := bytes.NewBufferString("")
	if cmd.Desc != "" {
		fmt.Fprintf(buff, "%s\n\n", cmd.Desc)
	}
	if cmd.Text != "" {
		fmt.Fprintf(buff, "%s\n\n", cmd.Text)
	}
	argvList := cmd.argvList()
	isEmpty := isEmptyArgvList(argvList)
	if !isEmpty {
		fmt.Fprintf(buff, "%s:\n\n%s", clr.Bold("Options"), usage(argvList, clr, style))
	}
	if cmd.children != nil && len(cmd.children) > 0 {
		if !isEmpty {
			buff.WriteByte('\n')
		}
		fmt.Fprintf(buff, "%s:\n\n%v", clr.Bold("Commands"), cmd.ChildrenDescriptions("  ", "   "))
	}
	tmpUsage = buff.String()
	cmd.locker.Lock()
	cmd.usage = tmpUsage
	cmd.usageStyle = style
	cmd.locker.Unlock()
	return tmpUsage
}

// Path returns space-separated command full name
func (cmd *Command) Path() string {
	return cmd.pathWithSep(" ")
}

func (cmd *Command) pathWithSep(sep string) string {
	var (
		path = ""
		cur  = cmd
	)
	for cur.parent != nil {
		if cur.Name != "" {
			if path == "" {
				path = cur.Name
			} else {
				path = cur.Name + sep + path
			}
		}
		cur = cur.parent
	}
	return path
}

// Root returns command's ancestor
func (cmd *Command) Root() *Command {
	ancestor := cmd
	for ancestor.parent != nil {
		ancestor = ancestor.parent
	}
	return ancestor
}

// Route finds command full matching router
func (cmd *Command) Route(router []string) *Command {
	child, end := cmd.SubRoute(router)
	if end != len(router) {
		return nil
	}
	return child
}

// SubRoute finds command partial matching router
func (cmd *Command) SubRoute(router []string) (*Command, int) {
	cur := cmd
	for i, name := range router {
		child := cur.findChild(name)
		if child == nil {
			return cur, i
		}
		cur = child
	}
	return cur, len(router)
}

// findChild finds child command by name
func (cmd *Command) findChild(name string) *Command {
	if cmd.nochild() {
		return nil
	}
	for _, child := range cmd.children {
		if child.Name == name {
			return child
		}
		if child.Aliases != nil {
			for _, alias := range child.Aliases {
				if alias == name {
					return child
				}
			}
		}
	}
	return nil
}

// ListChildren returns all names of command children
func (cmd *Command) ListChildren() []string {
	if cmd.nochild() {
		return []string{}
	}

	ret := make([]string, 0, len(cmd.children))
	for _, child := range cmd.children {
		ret = append(ret, child.Name)
	}
	return ret
}

// ChildrenDescriptions returns all children's brief infos by one string
func (cmd *Command) ChildrenDescriptions(prefix, indent string) string {
	if cmd.nochild() {
		return ""
	}
	buff := bytes.NewBufferString("")
	length := 0
	for _, child := range cmd.children {
		if len(child.Name) > length {
			length = len(child.Name)
		}
	}
	format := fmt.Sprintf("%s%%-%ds%s%%s%%s\n", prefix, length, indent)
	for _, child := range cmd.children {
		aliases := ""
		if child.Aliases != nil && len(child.Aliases) > 0 {
			aliasesBuff := bytes.NewBufferString("(aliases ")
			aliasesBuff.WriteString(strings.Join(child.Aliases, ","))
			aliasesBuff.WriteString(")")
			aliases = aliasesBuff.String()
		}
		fmt.Fprintf(buff, format, child.Name, child.Desc, aliases)
	}
	return buff.String()
}

func (cmd *Command) nochild() bool {
	return cmd.children == nil || len(cmd.children) == 0
}

// Suggestions returns all similar commands
func (cmd *Command) Suggestions(path string) []string {
	if cmd.parent != nil {
		return cmd.Root().Suggestions(path)
	}

	var (
		cmds    = []*Command{cmd}
		targets = []string{}
	)
	for len(cmds) > 0 {
		if cmds[0].nochild() {
			cmds = cmds[1:]
		} else {
			for _, child := range cmds[0].children {
				targets = append(targets, child.Path())
			}
			cmds = append(cmds[0].children, cmds[1:]...)
		}
	}

	dists := []editDistanceRank{}
	for i, size := 0, len(targets); i < size; i++ {
		if d, ok := match(path, targets[i]); ok {
			dists = append(dists, editDistanceRank{s: targets[i], d: d})
		}
	}
	sort.Sort(editDistanceRankSlice(dists))
	for i := 0; i < len(dists); i++ {
		targets[i] = dists[i].s
	}
	return targets[:len(dists)]
}
