package expr

import (
	"fmt"
	"math"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConstExpr(t *testing.T) {
	for _, x := range []struct {
		s   string
		val float64
	}{
		{"1", 1},
		{"1+2", 3},
		{"1*2+3", 5},
		{"1*(2+3)", 5},
	} {
		e, err := New(x.s, nil)
		if err != nil {
			t.Errorf("parse %q error: %v", x.s, err)
			continue
		}
		val, err := e.Eval(nil)
		if err != nil {
			t.Errorf("parse %q error: %v", x.s, err)
			continue
		}
		if math.Abs(val.Float()-x.val) > 1E-6 {
			t.Errorf("%q want %f, got %f", x.s, x.val, val)
		}
	}
}

func TestVarExpr(t *testing.T) {
	getter := Getter{
		"x": Float(1.5),
		"y": Float(2.5),
		"m": Float(5),
		"n": Float(2),
	}

	for _, x := range []struct {
		s     string
		val   float64
		isErr bool
	}{
		{"1", 1, false},
		{"x", 1.5, false},
		{"y", 2.5, false},
		{"x+y", 4, false},
		{"x-y", -1, false},
		{"x*y", 3.75, false},
		{"x/y", 0.6, false},
		{"m / n", 2.5, false},
		{"m / n // line comment", 2.5, false},
		{"m / /*multiline\n	comment*/ n", 2.5, false},
		{"m%n", 1, false},
		{"(1+x)*y/(m-n)^2", 4.34027776, false},
		{"(1+x)*y/((m-n)^2)", 0.69444444, false},

		{"min(x,y,m)", 1.5, false},
		{"max(x,y,m)", 5, false},
		{"max(1, x+y)", 4, false},
		{"min(max(1, x+y), 2, x)", 1.5, false},

		{"undefined_func(x,y,m)", 0, true},
		{"max()", 0, true},
		{"m!-!n", 0, true},
		{"m / undefined_var", 0, true},
	} {
		e, err := New(x.s, nil)
		if err != nil {
			if !x.isErr {
				t.Errorf("parse %q error: %v", x.s, err)
			}
			continue
		}
		val, err := e.Eval(getter)
		if err != nil {
			if !x.isErr {
				t.Errorf("parse %q error: %v", x.s, err)
			}
			continue
		}
		if math.Abs(val.Float()-x.val) > 1E-6 {
			t.Errorf("%q want %f, got %f", x.s, x.val, val.Float())
		}
	}
}

func TestCustomFactory(t *testing.T) {
	pool, _ := NewPool(map[string]Func{
		"constant": func(...Value) (Value, error) {
			return Int(123), nil
		},
		"sum": func(args ...Value) (Value, error) {
			var (
				sum = Zero()
				err error
			)
			for _, arg := range args {
				sum, err = sum.Add(arg)
				if err != nil {
					return sum, err
				}
			}
			return sum, nil
		},
		"average": func(args ...Value) (Value, error) {
			n := len(args)
			if n == 0 {
				return Zero(), fmt.Errorf("missing arguments for function `%s`", "average")
			}
			var (
				sum = Zero()
				err error
			)
			for _, arg := range args {
				sum, err = sum.Add(arg)
				if err != nil {
					return sum, err
				}
			}
			return sum.Quo(Float(float64(n)))
		},
	})
	getter := Getter{
		"x": Float(1.5),
		"y": Float(2.5),
	}

	for i, x := range []struct {
		s     string
		val   float64
		isErr bool
	}{
		{"constant()", 123, false},
		{"constant(x)", 123, false},
		{"sum(1,2,3)", 6, false},
		{"sum()", 0, false},
		{"sum(x, y, x)", 5.5, false},
		{"average(x)", 1.5, false},
		{"average(x,y)", 2, false},
		{"average()", 0, true},
	} {
		e, err := New(x.s, pool)
		if err != nil {
			if !x.isErr {
				t.Errorf("%dth: parse %q error: %v", i, x.s, err)
			}
			continue
		}
		val, err := e.Eval(getter)
		if err != nil {
			if !x.isErr {
				t.Errorf("%dth: parse %q error: %v", i, x.s, err)
			}
			continue
		}
		if math.Abs(val.Float()-x.val) > 1E-6 {
			t.Errorf("%dth: %q want %f, got %f", i, x.s, x.val, val.Float())
		}
	}
}

func TestOnVarMissing(t *testing.T) {
	defaults := map[string]Value{
		"a": Int(0),
		"b": Int(1),
	}
	pool, _ := NewPool()
	pool.SetOnVarMissing(func(varName string) (Value, error) {
		if dft, ok := defaults[varName]; ok {
			return dft, nil
		}
		return DefaultOnVarMissing(varName)
	})
	v, err := Eval("2 / b + a + x", map[string]Value{"x": Int(1)}, pool)
	assert.Nil(t, err)
	assert.Equal(t, float64(3), v.Float())

	v, err = Eval("2 / b + a + undefined", nil, pool)
	assert.NotNil(t, err)
}

func TestOp(t *testing.T) {
	for i, tc := range []struct {
		s      string
		result Value
		err    error
	}{
		{`1`, True(), nil},
		{`0`, False(), nil},
		{`0.0`, False(), nil},
		{`1.0`, True(), nil},
		{`"a"`, True(), nil},
		{`'a'`, True(), nil},
		{`'a'`, String("a"), nil},
		{`2 > 1`, True(), nil},
		{`2 >= 1`, True(), nil},
		{`1 >= 1`, True(), nil},
		{`1 <= 1`, True(), nil},
		{`1 == 1`, True(), nil},
		{`1 >= 2`, False(), nil},
		{`1 && 0`, False(), nil},
		{`1 && 2`, True(), nil},
		{`1 || 0`, True(), nil},
		{`1 < 2`, True(), nil},
		{`1 < 1`, False(), nil},
		{`1 <= 2`, True(), nil},
		{`1 <= 0`, False(), nil},
		{`1 != 0`, True(), nil},
		{`1 != 1`, False(), nil},
		{`"a" != "b"`, True(), nil},
		{`"a" != "a"`, False(), nil},
		{`"a" == "a"`, True(), nil},
		{`"ab" < "ac"`, True(), nil},
		{`"ab" <= "ac"`, True(), nil},
		{`"ab" <= "ab"`, True(), nil},
		{`"ab" > "ab"`, False(), nil},
		{`'a' == 'a'`, True(), nil},

		{`"a" + "b"`, String("ab"), nil},
		{`"ab" + "bc"`, String("abbc"), nil},
		{`2 - 1 > 0`, True(), nil},
		{`2 - 1 > 2`, False(), nil},
		{`2 + 1 > 2`, True(), nil},
		{`2 - 1`, Int(1), nil},

		{`"a" - "b"`, Nil(), ErrTypeMismatchForOp},
		{`"a" + 1`, Nil(), ErrTypeMismatchForOp},
		{`"a" + 1.2`, Nil(), ErrTypeMismatchForOp},
		{`"a" - 1`, Nil(), ErrTypeMismatchForOp},
		{`"a" == 1`, Nil(), ErrComparedTypesMismatch},
		{`"a" != 1`, Nil(), ErrComparedTypesMismatch},
		{`"a" > 1`, Nil(), ErrComparedTypesMismatch},
		{`"a" < 1`, Nil(), ErrComparedTypesMismatch},
		{`"a" >= 1`, Nil(), ErrComparedTypesMismatch},
		{`"a" <= 1`, Nil(), ErrComparedTypesMismatch},
		{`1/0`, Nil(), ErrDivideZero},
		{`1%0`, Nil(), ErrDivideZero},
		{`0^2`, Nil(), ErrPowOfZero},
		{`0.0^2`, Nil(), ErrPowOfZero},
	} {
		e, err := New(tc.s, nil)
		if err != nil {
			t.Errorf("%dth: invalid expression `%s'", i, tc.s)
			continue
		}
		got, err := e.Eval(nil)
		if err != nil {
			if err != tc.err {
				t.Errorf("%dth: want error `%v', got error `%v'", i, tc.err, err)
				continue
			}
		} else if tc.err != nil {
			t.Errorf("%dth: want error `%s', got nil", i, tc.err)
			continue
		}
		eq, _ := tc.result.Ne(got)
		if eq.Bool() {
			t.Errorf("%dth: result error, want `%s', got `%s'", i, tc.result.String(), got.String())
		}
	}
}

func TestOpWithGetter(t *testing.T) {
	pool := MustNewPool(map[string]Func{
		"contains": func(args ...Value) (Value, error) {
			if err := ExpectNArg(len(args), 2); err != nil {
				return Nil(), fmt.Errorf("expected number of arguments for function `contains' is %d, but got %d", 2, len(args))
			}
			if args[0].kind == KindString && args[1].kind == KindString {
				return Bool(strings.Contains(args[0].String(), args[1].String())), nil
			}
			return Nil(), ErrTypeMismatchForOp
		},
	})
	getter := Getter{
		"a": String("u"),
		"b": String("v"),
		"c": String("w"),
		"x": Int(1),
		"y": Int(2),
		"z": Int(0),
	}
	for i, tc := range []struct {
		s      string
		result Value
		err    error
	}{
		{`a + b`, String("uv"), nil},
		{`x + y`, Int(3), nil},

		{`a + 1`, nilValue, ErrTypeMismatchForOp},
		{`a > 1`, nilValue, ErrComparedTypesMismatch},
	} {
		e, err := New(tc.s, pool)
		if err != nil {
			t.Errorf("%dth: invalid expression `%s'", i, tc.s)
			continue
		}
		got, err := e.Eval(getter)
		if err != nil {
			if err != tc.err {
				t.Errorf("%dth: want error `%v', got error `%v'", i, tc.err, err)
				continue
			}
		} else if tc.err != nil {
			t.Errorf("%dth: want error `%s', got nil", i, tc.err)
			continue
		}
		eq, _ := tc.result.Ne(got)
		if eq.Bool() {
			t.Errorf("%dth: result error, want `%s', got `%s'", i, tc.result.String(), got.String())
		}
	}
}
