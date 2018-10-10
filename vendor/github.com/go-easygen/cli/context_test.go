package cli

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContextGetArgvList(t *testing.T) {
	type rootT struct {
		A string `cli:"a"`
	}
	type helloT struct {
		B int `cli:"b"`
	}
	type worldT struct {
		C bool `cli:"c"`
	}
	root := &Command{
		Name:   "root",
		Argv:   func() interface{} { return new(rootT) },
		Global: true,
		Fn: func(ctx *Context) error {
			rootArgv := ctx.RootArgv().(*rootT)
			assert.Equal(t, rootArgv.A, "Root-A")

			assert.Nil(t, ctx.GetArgvList(nil))

			var pA = &rootT{}
			assert.Nil(t, ctx.GetArgvList(pA))
			assert.Equal(t, pA.A, "Root-A")
			return nil
		},
	}
	hello := &Command{
		Name: "hello",
		Argv: func() interface{} { return new(helloT) },
		Fn: func(ctx *Context) error {
			argv := ctx.Argv().(*helloT)
			assert.Equal(t, argv.B, 100)

			pA := new(rootT)
			pB := new(helloT)
			assert.Nil(t, ctx.GetArgvList(pB, pA))
			assert.Equal(t, pA.A, "Root-A")
			assert.Equal(t, pB.B, argv.B)
			return nil
		},
	}
	world := &Command{
		Name: "world",
		Argv: func() interface{} { return new(worldT) },
		Fn: func(ctx *Context) error {
			argv := ctx.Argv().(*worldT)
			assert.Equal(t, argv.C, true)

			pA := new(rootT)
			pC := new(worldT)
			assert.Nil(t, ctx.GetArgvList(pC, nil, pA))
			assert.Equal(t, pA.A, "Root-A")
			assert.Equal(t, pC.C, true)
			assert.Error(t, ctx.GetArgvList(pC, new(helloT), pA))
			return nil
		},
	}
	root.Register(hello)
	hello.Register(world)
	assert.Nil(t, root.RunWith([]string{"-a=Root-A"}, nil, nil))
	assert.Nil(t, root.RunWith([]string{"hello", "-a=Root-A", "-b", "100"}, nil, nil))
	assert.Nil(t, root.RunWith([]string{"hello", "world", "-a=Root-A", "-c"}, nil, nil))
}

func TestContextMisc(t *testing.T) {
	type argT struct {
		Hello string `cli:"hello"`
		Age   int    `cli:"age" dft:"10"`
	}
	root := &Command{Name: "root"}
	parent := &Command{Name: "parent", Fn: donothing}
	cmd := &Command{
		Name: "cmd",
		Argv: func() interface{} { return new(argT) },
		Fn: func(ctx *Context) error {
			argv := ctx.Argv().(*argT)
			assert.Equal(t, ctx.Args(), []string{"a", "b", "c"})
			assert.Equal(t, argv.Hello, "world")
			assert.Equal(t, argv.Age, 10)
			assert.Equal(t, ctx.NativeArgs(), []string{"--hello=world", "a", "b", "c"})
			assert.Equal(t, ctx.Command().Name, "cmd")
			assert.Equal(t, ctx.IsSet("--hello"), true)
			assert.Equal(t, ctx.IsSet("--age"), false)
			assert.Equal(t, ctx.NArg(), 3)
			assert.Equal(t, ctx.Path(), "parent cmd")
			assert.Equal(t, ctx.Router(), []string{"parent", "cmd"})

			ctx.writer = nil
			assert.NotNil(t, ctx.Writer())
			return nil
		},
	}
	root.Register(parent)
	parent.Register(cmd)
	assert.Nil(t, root.RunWith([]string{"parent", "cmd", "--hello=world", "a", "b", "c"}, nil, nil))
}

func TestContextWriter(t *testing.T) {
	w := bytes.NewBufferString("")
	assert.Nil(t, (&Command{
		Name:        "root",
		CanSubRoute: true,
		Fn: func(ctx *Context) error {
			ctx.String("String")
			ctx.JSON(struct{ A int }{10})
			ctx.JSONln(struct{ B int }{10})
			ctx.JSONIndent(struct{ C string }{"11"}, "", "  ")
			ctx.JSONIndentln(struct{ D bool }{true}, "", "  ")
			assert.Equal(t, w, ctx.Writer())

			n, err := ctx.Write([]byte("end"))
			assert.Equal(t, 3, n)
			return err
		},
	}).RunWith([]string{"a", "b"}, w, nil))
	assert.Equal(t, w.String(), `String{"A":10}{"B":10}
{
  "C": "11"
}{
  "D": true
}
end`)
}
