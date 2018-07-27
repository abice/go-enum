package cli

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateCommandName(t *testing.T) {
	for i, invalidName := range []string{
		"",
		"-",
		"-abc",
		"~@#$%^&*+/\\<>.,:;'\"",
		"无效的",
		"無效的",
		"недействительный",
		"無効",
		"잘못된",
	} {
		assert.False(t, IsValidCommandName(invalidName), "%d: %s", i, invalidName)
	}
	for _, validName := range []string{
		"a",
		"abc",
		"abc-def",
		"abc_def",
		"_abc",
		"A",
		"Abc",
		"aBC",
	} {
		assert.True(t, IsValidCommandName(validName))
	}
}

func TestCommandTree(t *testing.T) {
	app := &Command{}

	type argT struct {
		Help    bool   `cli:"h,help" usage:"show help"`
		Version string `cli:"v,version" usage:"show version" dft:"v0.0.0"`
	}

	sub1 := app.Register(&Command{
		Name: "sub1",
		Fn: func(ctx *Context) error {
			assert.Equal(t, ctx.Path(), "sub1")
			argv := ctx.Argv().(*argT)
			assert.True(t, argv.Help)
			assert.Equal(t, argv.Version, "v0.0.0")
			assert.Equal(t, ctx.Command().Name, "sub1")
			return nil
		},
		Desc: "sub1 command describe",
		Argv: func() interface{} { return new(argT) },
	})

	sub1.Register(&Command{
		Name: "sub11",
		Fn: func(ctx *Context) error {
			assert.Equal(t, ctx.Path(), "sub1 sub11")
			argv := ctx.Argv().(*argT)
			assert.False(t, argv.Help)
			assert.Equal(t, argv.Version, "v1.0.0")
			return nil
		},
		Desc: "sub11 desc",
		Text: "sub11 text",
		Argv: func() interface{} { return new(argT) },
	})

	assert.Nil(t, app.Run([]string{"sub1", "-h"}))
	assert.Nil(t, app.Run([]string{"sub1", "sub11", "--version=v1.0.0"}))
	assert.Equal(t, sub1.ChildrenDescriptions("", " "), "sub11 sub11 desc\n")
}

func TestRegisterNilCommand(t *testing.T) {
	cmd := &Command{Name: "root"}
	assert.Panics(t, func() { cmd.Register(nil) })
}

func TestRegisterCommandWithInvalidName(t *testing.T) {
	cmd := &Command{Name: "root"}
	assert.Panics(t, func() { cmd.Register(&Command{Name: "-invalid-", Fn: donothing}) })
}

func TestRegisterCommandWhichHasHadParent(t *testing.T) {
	cmd := &Command{Name: "root"}
	child := &Command{
		Name:   "sub",
		Fn:     donothing,
		parent: &Command{Name: "parent"},
	}
	assert.Panics(t, func() { cmd.Register(child) })
}

func TestResgisterRepeatedCommand(t *testing.T) {
	cmd := &Command{Name: "root"}
	cmd.Register(&Command{Name: "sub", Fn: donothing})
	assert.Panics(t, func() { cmd.Register(&Command{Name: "sub", Fn: donothing}) })
	assert.Panics(t, func() { cmd.Register(&Command{Name: "hello", Aliases: []string{"sub"}, Fn: donothing}) })
}

func TestRegisterTree(t *testing.T) {
	cmd := &Command{Name: "root"}
	tree := Tree(&Command{Name: "sub", Fn: donothing}, Tree(&Command{Name: "sub2", Fn: donothing}))
	cmd.RegisterTree(tree)
	require.Len(t, cmd.children, 1)
	assert.Equal(t, cmd.children[0].Name, "sub")
	assert.Equal(t, cmd.children[0].Parent().Name, "root")
	require.Len(t, cmd.children[0].children, 1)
	assert.Equal(t, cmd.children[0].children[0].Name, "sub2")
	assert.Equal(t, cmd.children[0].children[0].Parent().Name, "sub")
}

func TestRegisterFunc(t *testing.T) {
	root := &Command{Name: "root"}
	out := ""
	sayHello := func(*Context) error {
		out = "hello"
		return nil
	}
	argvFn := func() interface{} { return "argv" }
	root.RegisterFunc("cmd", sayHello, argvFn)
	require.Len(t, root.children, 1)
	assert.Equal(t, root.children[0].Name, "cmd")

	assert.Equal(t, out, "")
	root.children[0].Fn(nil)
	assert.Equal(t, out, "hello")

	assert.Equal(t, root.children[0].Argv(), "argv")
}

func TestCommandNotFound(t *testing.T) {
	root := &Command{Name: "root"}
	sub := &Command{Name: "sub", Fn: donothing}
	root.Register(sub)
	err := root.RunWith([]string{"not-found"}, nil, nil)
	if e, ok := err.(wrapError); ok {
		err = e.err
	}
	assert.IsType(t, commandNotFoundError{}, err)
}

type testValidator struct {
	Value int `cli:"v" usage:"must be range int [1,10)"`
}

func (argv *testValidator) Validate(ctx *Context) error {
	if argv.Value >= 1 && argv.Value < 10 {
		return nil
	}
	return fmt.Errorf("out of range")
}

func TestValidator(t *testing.T) {
	getCmd := func() *Command {
		return &Command{
			Name: "root",
			Argv: func() interface{} { return new(testValidator) },
			Fn: func(ctx *Context) error {
				argv := ctx.Argv().(*testValidator)
				assert.True(t, argv.Value >= 1 && argv.Value < 10)
				return nil
			},
		}
	}
	assert.Error(t, getCmd().RunWith([]string{"-v=20"}, nil, nil))
	assert.Nil(t, getCmd().RunWith([]string{"-v=2"}, nil, nil))
}

//TODO: TestCommandHooks

func TestCommandMisc(t *testing.T) {
	root := &Command{Name: "root"}
	root.SetIsServer(true)
	sub := &Command{Name: "sub", Fn: donothing}
	sub2 := &Command{Name: "sub2", Fn: donothing}
	root.Register(sub)
	sub.Register(sub2)

	assert.True(t, root.IsServer())
	assert.False(t, root.IsClient())

	assert.Equal(t, sub.Parent().Name, "root")
	assert.Equal(t, sub.Path(), "sub")
	assert.Equal(t, sub2.Path(), "sub sub2")
	assert.Equal(t, root.findChild("sub"), sub)
	assert.Equal(t, root.Root(), root)
	assert.Equal(t, sub.Root(), root)
	assert.Equal(t, sub2.Root(), root)

	cmd, deep := root.SubRoute([]string{"not", "found"})
	assert.Equal(t, cmd, root)
	assert.Equal(t, deep, 0)
	cmd, deep = root.SubRoute([]string{"sub", "no"})
	assert.Equal(t, cmd, sub)
	assert.Equal(t, deep, 1)
	cmd, deep = root.SubRoute([]string{"sub", "sub2"})
	assert.Equal(t, cmd, sub2)
	assert.Equal(t, deep, 2)

	assert.Nil(t, root.Route([]string{"not", "found"}))
	assert.Equal(t, root.Route([]string{}), root)
	assert.Equal(t, root.Route([]string{"sub"}), sub)
	assert.Nil(t, root.Route([]string{"sub", "not", "found"}))
	assert.Equal(t, root.Route([]string{"sub", "sub2"}), sub2)

	assert.Equal(t, root.Suggestions("su"), []string{"sub"})
}
