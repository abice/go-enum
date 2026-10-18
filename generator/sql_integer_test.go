package generator

import (
	"fmt"
	"go/parser"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// Exercise the generated Scanner, including conversions that wrap to a named value.
func TestSQLIntegerOverflow(t *testing.T) {
	input := "package bounds\n\ntype CustomSigned int8\ntype CustomUnsigned uint8\n"
	types := []string{"int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "uintptr", "byte", "rune", "CustomSigned", "CustomUnsigned"}
	for i, typ := range types {
		input += fmt.Sprintf("// ENUM(Zero, One)\ntype Number%d %s\n", i, typ)
	}
	for _, mode := range []struct {
		name           string
		option         Option
		numericStrings bool
	}{
		{"sql", WithSQLDriver(), false},
		{"sqlnullstr", WithSQLNullStr(), false},
		{"sqlnullint", WithSQLNullInt(), true},
	} {
		t.Run(mode.name, func(t *testing.T) {
			g := NewGenerator(mode.option)
			f, err := parser.ParseFile(g.fileSet, "bounds.go", input, parser.ParseComments)
			require.NoError(t, err)
			output, err := g.Generate(f)
			require.NoError(t, err)
			dir := t.TempDir()
			for name, contents := range map[string]string{
				"go.mod":         "module bounds\n\ngo 1.25.0\n",
				"bounds.go":      input,
				"bounds_enum.go": string(output),
				"bounds_test.go": strings.ReplaceAll(sqlIntegerBoundsTest, "NUMERIC_STRINGS", fmt.Sprint(mode.numericStrings)),
			} {
				require.NoError(t, os.WriteFile(filepath.Join(dir, name), []byte(contents), 0600))
			}
			cmd := exec.Command("go", "test", "-count=1", "./...")
			cmd.Dir = dir
			result, err := cmd.CombinedOutput()
			require.NoError(t, err, "%s", result)
		})
	}
}

const sqlIntegerBoundsTest = `package bounds
import (
 "fmt"
 "math"
 "reflect"
 "strconv"
 "testing"
)
type scanner interface { Scan(interface{}) error }
func TestBounds(t *testing.T) {
 cases := []struct {
  name string
  create func() scanner
  invalid []interface{}
  maximum interface{}
 }{
{"int", func() scanner { x := Number0One; return &x }, []interface{}{uint64(math.MaxUint64)}, int64(1 << (strconv.IntSize-2))},
{"int8", func() scanner { x := Number1One; return &x }, []interface{}{uint64(math.MaxUint64), int64(math.MaxInt8)+1, int64(math.MinInt8)-1}, int64(math.MaxInt8)},
{"int16", func() scanner { x := Number2One; return &x }, []interface{}{uint64(math.MaxUint64), int64(math.MaxInt16)+1, int64(math.MinInt16)-1}, int64(math.MaxInt16)},
{"int32", func() scanner { x := Number3One; return &x }, []interface{}{uint64(math.MaxUint64), int64(math.MaxInt32)+1, int64(math.MinInt32)-1}, int64(math.MaxInt32)},
{"int64", func() scanner { x := Number4One; return &x }, []interface{}{uint64(math.MaxUint64)}, int64(math.MaxInt64)},
{"uint", func() scanner { x := Number5One; return &x }, []interface{}{int64(-1)}, uint64(1 << (strconv.IntSize-1))},
{"uint8", func() scanner { x := Number6One; return &x }, []interface{}{int64(-1), uint64(1)<<8}, uint64(math.MaxUint8)},
{"uint16", func() scanner { x := Number7One; return &x }, []interface{}{int64(-1), uint64(1)<<16}, uint64(math.MaxUint16)},
{"uint32", func() scanner { x := Number8One; return &x }, []interface{}{int64(-1), uint64(1)<<32}, uint64(math.MaxUint32)},
{"uint64", func() scanner { x := Number9One; return &x }, []interface{}{int64(-1)}, uint64(math.MaxUint64)},
{"CustomSigned", func() scanner { x := Number13One; return &x }, []interface{}{int64(-129), int64(128), uint64(256)}, int64(127)},
{"CustomUnsigned", func() scanner { x := Number14One; return &x }, []interface{}{int64(-1), int64(256), uint64(256)}, uint64(255)},
{"byte", func() scanner { x := Number11One; return &x }, []interface{}{int64(-1), int64(256), uint64(256)}, uint64(math.MaxUint8)},
{"rune", func() scanner { x := Number12One; return &x }, []interface{}{int64(math.MinInt32)-1, int64(math.MaxInt32)+1, uint64(math.MaxUint64)}, int64(math.MaxInt32)},
{"uintptr", func() scanner { x := Number10One; return &x }, []interface{}{int64(-1)}, uint64(1 << (strconv.IntSize-1))},
 }
 for _, tc := range cases {
  t.Run(tc.name, func(t *testing.T) {
   for _, invalid := range append(tc.invalid, math.NaN(), math.Inf(1), math.Inf(-1), math.Exp2(64), -math.Exp2(64)) {
    inputs := []interface{}{invalid}
    ptr := reflect.New(reflect.TypeOf(invalid)); ptr.Elem().Set(reflect.ValueOf(invalid))
    inputs = append(inputs, ptr.Interface())
    if NUMERIC_STRINGS && reflect.TypeOf(invalid).Kind() != reflect.Float64 {
     text := fmt.Sprint(invalid)
     inputs = append(inputs, text, []byte(text), &text)
    }
    for _, value := range inputs {
     x := tc.create()
     if err := x.Scan(value); err == nil { t.Errorf("Scan(%T %v) accepted overflow", value, value) }
     if number(x) != "1" { t.Errorf("Scan(%T %v) changed receiver: %s", value, value, number(x)) }
    }
   }
   inputs := []interface{}{int(1), int64(1), uint(1), uint64(1), float64(1.75), "One", []byte("One"), tc.maximum}
   for _, value := range append([]interface{}{}, inputs...) {
    if reflect.TypeOf(value).Kind() != reflect.Slice {
     ptr := reflect.New(reflect.TypeOf(value)); ptr.Elem().Set(reflect.ValueOf(value))
     inputs = append(inputs, ptr.Interface())
    }
   }
   for _, value := range inputs {
    x := tc.create()
    if err := x.Scan(value); err != nil { t.Errorf("Scan(%T %v): %v", value, value, err) }
    expected := "1"
    plain := value
    if reflect.TypeOf(value).Kind() == reflect.Ptr { plain = reflect.ValueOf(value).Elem().Interface() }
    if reflect.TypeOf(plain) == reflect.TypeOf(tc.maximum) && fmt.Sprint(plain) == fmt.Sprint(tc.maximum) { expected = fmt.Sprint(tc.maximum) }
    if got := number(x); got != expected { t.Errorf("Scan(%T %v) = %s, want %s", value, value, got, expected) }
   }
   x := tc.create()
   if err := x.Scan(nil); err != nil || number(x) != "0" { t.Errorf("nil Scan = %s, %v", number(x), err) }
  })
 }
}
func number(x scanner) string {
 v := reflect.ValueOf(x).Elem()
 if v.Kind() >= reflect.Uint && v.Kind() <= reflect.Uintptr { return strconv.FormatUint(v.Uint(),10) }
 return strconv.FormatInt(v.Int(),10)
}
`
