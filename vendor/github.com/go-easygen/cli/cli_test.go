package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"reflect"
	"testing"

	"github.com/labstack/gommon/color"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var donothing = func(*Context) error { return nil }

func TestTree(t *testing.T) {
	cmd := &Command{Name: "cmd", Fn: donothing}
	subcmd := &Command{Name: "subcmd", Fn: donothing}

	tree := Tree(cmd, Tree(subcmd))
	assert.Equal(t, tree.command.Name, "cmd")
	require.Len(t, tree.forest, 1)
	assert.Equal(t, tree.forest[0].command.Name, "subcmd")
	assert.Len(t, tree.forest[0].forest, 0)
}

func TestRoot(t *testing.T) {
	donothing := func(*Context) error { return nil }
	root := Root(&Command{Name: "cmd", Fn: donothing},
		Tree(&Command{Name: "subcmd", Fn: donothing}),
	)
	assert.Equal(t, root.Name, "cmd")
	require.Len(t, root.children, 1)
	assert.Equal(t, root.children[0].Name, "subcmd")
}

type argT struct {
	Short         bool   `cli:"s" usage:"short flag"`
	Short2        bool   `cli:"2" usage:"another short flag"`
	ShortAndLong  string `cli:"S,long" usage:"short and long flags"`
	ShortsAndLong int    `cli:"x,y,abcd,omitof" usage:"many short and long flags"`
	Long          uint   `cli:"long-flag" usage:"long flag"`
	Required      int8   `cli:"*required" usage:"required flag, note the *"`
	Default       uint8  `cli:"dft,default" usage:"default value" dft:"102"`
	UnName        uint16 `usage:"unname field"`

	Int8    int8    `cli:"i8" usage:"type int8"`
	Uint8   uint8   `cli:"u8" usage:"type uint8"`
	Int16   int16   `cli:"i16" usage:"type int16"`
	Uint16  uint16  `cli:"u16" usage:"type uint16"`
	Int32   int32   `cli:"i32" usage:"type int32"`
	Uint32  uint32  `cli:"u32" usage:"type uint32"`
	Int64   int64   `cli:"i64" usage:"type int64"`
	Uint64  uint64  `cli:"u64" usage:"type uint64"`
	Float32 float32 `cli:"f32" usage:"type float32"`
	Float64 float64 `cli:"f64" usage:"type float64"`
}

func toStr(i interface{}) string {
	return fmt.Sprintf("%v", i)
}

func TestParse(t *testing.T) {
	clr := color.Color{}
	for i, tab := range []struct {
		args        []string
		want        argT
		isErr       bool
		freedomArgs []string
	}{
		//Case: ENV
		{
			args: []string{"--required=0"},
			want: argT{Default: 102},
		},
		//Case: missing required
		{
			args:  []string{},
			isErr: true,
		},
		//Case: undefined flag
		{
			args:  []string{"--required=0", "-Q"},
			isErr: true,
		},
		//Case: undefined flag
		{
			args:  []string{"--required=0", "--KdjiiejdfwkHJH"},
			isErr: true,
		},
		//Case: short flag group
		{
			args: []string{"--required=0", "-s2"},
			want: argT{Default: 102, Short: true, Short2: true},
		},
		//Case: check default
		{
			args: []string{"--required=0"},
			want: argT{Default: 102},
		},
		//Case: modify default
		{
			args: []string{"--required=0", "--dft", "55"},
			want: argT{Default: 55},
		},
		//Case: modify default
		{
			args: []string{"--required=0", "--default", "55"},
			want: argT{Default: 55},
		},
		//Case: UnName
		{
			args: []string{"--required=0", "--UnName", "64"},
			want: argT{Default: 102, UnName: 64},
		},
		//Case: not a bool
		{
			args:  []string{"--required=0", "-s=not-a-bool"},
			isErr: true,
		},
		//Case: "" -> bool
		{
			args: []string{"--required=0", "-s"},
			want: argT{Default: 102, Short: true},
		},
		//Case: "true" -> bool
		{
			args: []string{"--required=0", "-s=true"},
			want: argT{Default: 102, Short: true},
		},
		//Case: non-zero integer -> bool
		{
			args: []string{"--required=0", "-s=1"},
			want: argT{Default: 102, Short: true},
		},
		//Case: zero -> bool
		{
			args: []string{"--required=0", "-s=0"},
			want: argT{Default: 102},
		},
		//Case: no -> bool
		{
			args: []string{"--required=0", "-s=no"},
			want: argT{Default: 102},
		},
		//Case: not -> bool
		{
			args: []string{"--required=0", "-s=not"},
			want: argT{Default: 102},
		},
		//Case: none -> bool
		{
			args: []string{"--required=0", "-s=none"},
			want: argT{Default: 102},
		},
		//Case: false -> bool
		{
			args: []string{"--required=0", "-s=false"},
			want: argT{Default: 102},
		},
		//Case: int64
		{
			args: []string{"--required=0", "--i64", toStr(12)},
			want: argT{Default: 102, Int64: 12},
		},
		//Case: int64 overflow
		{
			args:  []string{"--required=0", "--i64", toStr(uint64(math.MaxUint64))},
			isErr: true,
		},
		//Case: uint64
		{
			args: []string{"--required=0", "--u64", toStr(12)},
			want: argT{Default: 102, Uint64: 12},
		},
		//Case: max uint64
		{
			args: []string{"--required=0", "--u64", toStr(uint64(math.MaxUint64))},
			want: argT{Default: 102, Uint64: uint64(math.MaxUint64)},
		},
		//Case: negative -> uint64
		{
			args:  []string{"--required=0", "--u64", "-1"},
			isErr: true,
		},
		//Case: int32
		{
			args: []string{"--required=0", "--i32", toStr(12)},
			want: argT{Default: 102, Int32: 12},
		},
		//Case: int32 overflow
		{
			args:  []string{"--required=0", "--i32", toStr(uint32(math.MaxUint32))},
			isErr: true,
		},
		//Case: uint32
		{
			args: []string{"--required=0", "--u32", toStr(12)},
			want: argT{Default: 102, Uint32: 12},
		},
		//Case: max uint32
		{
			args: []string{"--required=0", "--u32", toStr(uint32(math.MaxUint32))},
			want: argT{Default: 102, Uint32: uint32(math.MaxUint32)},
		},
		//Case: negative -> uint32
		{
			args:  []string{"--required=0", "--u32", "-1"},
			isErr: true,
		},
		//Case: int16
		{
			args: []string{"--required=0", "--i16", toStr(12)},
			want: argT{Default: 102, Int16: 12},
		},
		//Case: int16 overflow
		{
			args:  []string{"--required=0", "--i16", toStr(uint16(math.MaxUint16))},
			isErr: true,
		},
		//Case: uint16
		{
			args: []string{"--required=0", "--u16", toStr(12)},
			want: argT{Default: 102, Uint16: 12},
		},
		//Case: max uint16
		{
			args: []string{"--required=0", "--u16", toStr(uint16(math.MaxUint16))},
			want: argT{Default: 102, Uint16: uint16(math.MaxUint16)},
		},
		//Case: negative -> uint16
		{
			args:  []string{"--required=0", "--u16", "-1"},
			isErr: true,
		},
		//Case: int8
		{
			args: []string{"--required=0", "--i8", toStr(12)},
			want: argT{Default: 102, Int8: 12},
		},
		//Case: int8 overflow
		{
			args:  []string{"--required=0", "--i8", toStr(uint8(math.MaxUint8))},
			isErr: true,
		},
		//Case: uint8
		{
			args: []string{"--required=0", "--u8", toStr(12)},
			want: argT{Default: 102, Uint8: 12},
		},
		//Case: max uint8
		{
			args: []string{"--required=0", "--u8", toStr(uint8(math.MaxUint8))},
			want: argT{Default: 102, Uint8: uint8(math.MaxUint8)},
		},
		//Case: negative -> uint8
		{
			args:  []string{"--required=0", "--u8", "-1"},
			isErr: true,
		},
		//Case: many invalid value
		{
			args:  []string{"--required=0", "--u8", "-1", "--u8", "256"},
			isErr: true,
		},
		//Case: too many value
		{
			args:  []string{"--required=0=1"},
			isErr: true,
		},
		//Case: float32
		{
			args: []string{"--required=0", "--f32", "12.34"},
			want: argT{Default: 102, Float32: 12.34},
		},
		//Case: not a float32
		{
			args:  []string{"--required=0", "--f32", "not-a-float32"},
			isErr: true,
		},
		//Case: float32 overflow
		{
			args:  []string{"--required=0", "--f32", "123456789123456789123456789123456789123456789"},
			isErr: true,
		},
		//Case: float32 overflow
		{
			args:  []string{"--required=0", "--f32", "-123456789123456789123456789123456789123456789"},
			isErr: true,
		},
		//Case: float64
		{
			args: []string{"--required=0", "--f64=-1234.5678"},
			want: argT{Default: 102, Float64: -1234.5678},
		},
		//Case: not a float64
		{
			args:  []string{"--required=0", "--f64=not-a-float64"},
			isErr: true,
		},
		//Case: fold flag not boolean
		{
			args:  []string{"--required=0", "-2y"},
			isErr: true,
		},
		//Case: test -F<value>
		{
			args: []string{"--required=0", "-Sshort-and-long"},
			want: argT{Default: 102, ShortAndLong: "short-and-long"},
		},
		//Case: test `--`
		{
			args:        []string{"--required=0", "--", "-Sshort-and-long"},
			want:        argT{Default: 102, ShortAndLong: ""},
			freedomArgs: []string{"-Sshort-and-long"},
		},
		//Case: test flags and args
		{
			args:        []string{"--required=0", "-Sshort-and-long", "abc"},
			want:        argT{Default: 102, ShortAndLong: "short-and-long"},
			freedomArgs: []string{"abc"},
		},
		{
			args:        []string{"--required=0", "-Sshort-and-long", "abc", "-s", "xyz"},
			want:        argT{Default: 102, Short: true, ShortAndLong: "short-and-long"},
			freedomArgs: []string{"abc", "xyz"},
		},
	} {
		if tab.args == nil {
			tab.args = []string{}
		}
		v := new(argT)
		flagSet := parseArgv(tab.args, v, clr)
		if tab.isErr {
			if flagSet.err == nil {
				t.Errorf("[%2d] want error, got nil", i)
			}
			continue
		}
		if flagSet.err != nil {
			t.Errorf("[%2d] parseArgv error: %v", i, flagSet.err)
			continue
		}
		if !reflect.DeepEqual(*v, tab.want) {
			t.Errorf("[%2d] want %v, got %v", i, tab.want, *v)
		}
		if !stringsEqual(tab.freedomArgs, flagSet.args) {
			t.Errorf("[%2d] want %v, got %v", i, tab.freedomArgs, flagSet.args)
		}
	}

	//Case parse non-pointer object
	if flagSet := parseArgv([]string{}, argT{}, clr); flagSet.err != errNotAPointer {
		t.Errorf("want %v, got %v", errNotAPointer, flagSet.err)
	}
	if usage([]interface{}{argT{}}, clr, NormalStyle) != "" {
		t.Errorf("want usage empty, but not")
	}

	//Case parse pointer, but not indirect a struct
	tmp := 0
	ptrInt := &tmp
	if flagSet := parseArgv([]string{}, ptrInt, clr); flagSet.err != errNotAPointerToStruct {
		t.Errorf("want %v, got %v", errNotAPointerToStruct, flagSet.err)
	}
	if usage([]interface{}{ptrInt}, clr, NormalStyle) != "" {
		t.Errorf("want usage empty, but not")
	}

	//Case repeat tag
	type tmpT struct {
		A bool `cli:"a"`
		B bool `cli:"a"`
	}
	if flagSet := parseArgv([]string{}, new(tmpT), clr); flagSet.err == nil {
		t.Errorf("want error, got nil")
	}
	if usage([]interface{}{new(tmpT)}, clr, NormalStyle) != "" {
		t.Errorf("want usage empty, but not")
	}

	type envT struct {
		DefaultEnv string `cli:"default-env" usage:"default value" dft:"$ENV_CLI_TEST"`
	}
	envV := new(envT)
	if flagSet := parseArgv([]string{}, envV, clr); flagSet.err != nil {
		t.Errorf(flagSet.err.Error())
	} else {
		if want := os.Getenv("ENV_CLI_TEST"); want != envV.DefaultEnv {
			t.Errorf("ENV_CLI_TEST want `%s`, got `%s`", want, envV.DefaultEnv)
		}
	}
}

func TestUsage(t *testing.T) {
	clr := color.Color{}
	clr.Disable()
	got := usage([]interface{}{new(argT)}, clr, NormalStyle)
	want := fmt.Sprintf(
		`      -s                           short flag
      -2                           another short flag
      -S, --long                   short and long flags
  -x, -y, --abcd, --omitof         many short and long flags
          --long-flag              long flag
          --required              *required flag, note the *
          --dft, --default[=102]   default value
          --UnName                 unname field
          --i8                     type int8
          --u8                     type uint8
          --i16                    type int16
          --u16                    type uint16
          --i32                    type int32
          --u32                    type uint32
          --i64                    type int64
          --u64                    type uint64
          --f32                    type float32
          --f64                    type float64
`)
	assert.Equal(t, got, want)

	got = usage([]interface{}{new(argT)}, clr, ManualStyle)
	want = `  -s
      short flag

  -2
      another short flag

  -S, --long
      short and long flags

  -x, -y, --abcd, --omitof
      many short and long flags

  --long-flag
      long flag

  --required
      *required flag, note the *

  --dft, --default[=102]
      default value

  --UnName
      unname field

  --i8
      type int8

  --u8
      type uint8

  --i16
      type int16

  --u16
      type uint16

  --i32
      type int32

  --u32
      type uint32

  --i64
      type int64

  --u64
      type uint64

  --f32
      type float32

  --f64
      type float64
`
	assert.Equal(t, got, want)
}

func TestStructField(t *testing.T) {
	type BaseT struct {
		Help bool `cli:"!h,help" usage:"display help"`
	}
	type subT struct {
		BaseT
		Version string `cli:"v" usage:"display version" dft:"v0.0.1"`
	}
	args := []string{
		"-h",
		"-v=v1.1.1",
	}
	argv := new(subT)
	clr := color.Color{}
	flagSet := parseArgv(args, argv, clr)
	if flagSet.err != nil {
		t.Errorf("parseArgv error: %v", flagSet.err)
	} else {
		if argv.BaseT.Help != true {
			t.Errorf("Help want true, but got false")
		}
		if argv.Version != "v1.1.1" {
			t.Errorf("Version want v1.1.1, but got %s", argv.Version)
		}
	}
}

func stringsEqual(ss1, ss2 []string) bool {
	if len(ss1) != len(ss2) {
		return false
	}
	for i := 0; i < len(ss1); i++ {
		if ss1[i] != ss2[i] {
			return false
		}
	}
	return true
}

func intsEqual(ss1, ss2 []uint32) bool {
	if ss1 == nil && ss2 == nil {
		return true
	}
	if len(ss1) != len(ss2) {
		return false
	}
	for i := 0; i < len(ss1); i++ {
		if ss1[i] != ss2[i] {
			return false
		}
	}
	return true
}

func mapEqual(m1, m2 map[uint32]string) bool {
	if m1 == nil && m2 == nil {
		return true
	}
	for k, v := range m1 {
		if v2, ok := m2[k]; !ok {
			return false
		} else if v2 != v {
			return false
		}
	}
	for k, v := range m2 {
		if v1, ok := m1[k]; !ok {
			return false
		} else if v1 != v {
			return false
		}
	}
	return true
}

func TestSliceAndMap(t *testing.T) {
	clr := color.Color{}
	type T struct {
		Slice []uint32          `cli:"s,slice"`
		Map   map[uint32]string `cli:"m,map"`
	}
	for _, tab := range []struct {
		args  []string
		want  T
		isErr bool
	}{
		{
			args: []string{"--slice=12", "--slice", "23"},
			want: T{Slice: []uint32{12, 23}},
		},
		{
			args: []string{"-s12", "-s23"},
			want: T{Slice: []uint32{12, 23}},
		},
		{
			args: []string{"-s=12", "-s=23"},
			want: T{Slice: []uint32{12, 23}},
		},
		{
			args: []string{"-s", "12", "-s", "23"},
			want: T{Slice: []uint32{12, 23}},
		},
		{
			args:  []string{"-s", "12", "-s", "not-a-number"},
			isErr: true,
		},
		{
			args: []string{"-m2=s", "-m3=y"},
			want: T{Map: map[uint32]string{
				2: "s",
				3: "y",
			}},
		},
		{
			args: []string{"-m2=s", "-m", "3=y"},
			want: T{Map: map[uint32]string{
				2: "s",
				3: "y",
			}},
		},
		{
			args:  []string{"-m"},
			isErr: true,
		},
		{
			args:  []string{"-ms=2", "-my=3"},
			isErr: true,
		},
	} {
		v := new(T)
		flagSet := parseArgv(tab.args, v, clr)
		if tab.isErr {
			if flagSet.err == nil {
				t.Errorf("want error, but not. tab=%v", tab)
			}
			continue
		}
		if flagSet.err != nil {
			t.Errorf("TestSliceAndMap: error: %v", flagSet.err)
			continue
		}
		if !intsEqual(tab.want.Slice, v.Slice) {
			t.Errorf("want %v, got %v", tab.want.Slice, v.Slice)
		}
		if !mapEqual(tab.want.Map, v.Map) {
			t.Errorf("want %v, got %v", tab.want.Map, v.Map)
		}
	}
}

func TestErrorSliceType(t *testing.T) {
	type A struct {
		Value string
	}
	type T struct {
		Slice []A `cli:"s"`
	}
	clr := color.Color{}
	v := new(T)
	flagSet := parseArgv([]string{"-s4"}, v, clr)
	if flagSet.err == nil {
		t.Errorf("want error, but not")
	}
}

func TestErrorMapType(t *testing.T) {
	type A struct {
		Value string
	}
	type T struct {
		Map map[string]A `cli:"m"`
	}
	clr := color.Color{}
	v := new(T)
	flagSet := parseArgv([]string{"-mkey=val"}, v, clr)
	if flagSet.err == nil {
		t.Errorf("want error, but not")
	}
}

func TestIsSet(t *testing.T) {
	type argT struct {
		A int `cli:"a,aa" dft:"1"`
		B int `cli:"b"`
	}
	for i, tt := range []struct {
		args   []string
		isSetA bool
		isSetB bool
	}{
		{[]string{"app", "-a=1", "-b=1"}, true, true},
		{[]string{"app", "--aa=1", "-b=1"}, true, true},
		{[]string{"app", "-a=1"}, true, false},
		{[]string{"app", "-b=1"}, false, true},
		{[]string{"app"}, false, false},
	} {
		RunWithArgs(new(argT), tt.args, func(ctx *Context) error {
			assert.Equal(t, ctx.IsSet("-a"), tt.isSetA, "case %d", i)
			assert.Equal(t, ctx.IsSet("-a", "--aa"), tt.isSetA, "case %d", i)
			assert.Equal(t, ctx.IsSet("-b"), tt.isSetB, "case %d", i)
			return nil
		})
	}
}

func TestHelpCommand(t *testing.T) {
	w := bytes.NewBufferString("")
	root := &Command{Name: "root"}
	help := HelpCommand("help command")
	root.Register(help)
	assert.Nil(t, root.RunWith([]string{"help"}, w, nil))
	assert.Equal(t, w.String(), "Commands:\n\n  help   help command\n")
	assert.Error(t, root.RunWith([]string{"help", "not-found"}, nil, nil))
}

func TestError(t *testing.T) {
	assert.Equal(t, ExitError.Error(), "exit")
	assert.Equal(t, throwCommandNotFound("cmd").Error(), "command cmd not found")
	assert.Equal(t, throwMethodNotAllowed("POST").Error(), "method POST not allowed")
	assert.Equal(t, throwRouterRepeat("R").Error(), "router R repeat")
	clr := color.Color{}
	clr.Disable()
	assert.Equal(t, wrapErr(throwCommandNotFound("cmd"), "_end", clr).Error(), `ERR! command cmd not found_end`)

	assert.Equal(t, argvError{isEmpty: true}.Error(), "argv list is empty")
	assert.Equal(t, argvError{isOutOfRange: true}.Error(), "argv list out of range")
	assert.Equal(t, argvError{ith: 1, msg: "ERROR MSG"}.Error(), "1th argv: ERROR MSG")
}

type customT struct {
	K1 string
	K2 int
}

func (t *customT) Decode(s string) error {
	return json.Unmarshal([]byte(s), t)
}

func TestDecoder(t *testing.T) {
	type argT struct {
		D customT `cli:"d"`
	}
	v := new(argT)
	clr := color.Color{}
	flagSet := parseArgv([]string{`-d`, `{"k1": "string", "k2": 2}`}, v, clr)
	assert.Nil(t, flagSet.err)
	assert.Equal(t, v.D, customT{K1: "string", K2: 2})
}
