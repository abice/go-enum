package cli

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/labstack/gommon/color"
)

// Run runs a single command app
func Run(argv interface{}, fn CommandFunc, descs ...string) {
	RunWithArgs(argv, os.Args, fn, descs...)
}

// RunWithArgs is similar to Run, but with args instead of os.Args
func RunWithArgs(argv interface{}, args []string, fn CommandFunc, descs ...string) {
	desc := ""
	if len(descs) > 0 {
		desc = strings.Join(descs, "\n")
	}
	err := (&Command{
		Name:        args[0],
		Desc:        desc,
		Argv:        func() interface{} { return argv },
		CanSubRoute: true,
		Fn:          fn,
	}).Run(args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

// Root registers forest for root and returns root
func Root(root *Command, forest ...*CommandTree) *Command {
	root.RegisterTree(forest...)
	return root
}

// Tree creates a CommandTree
func Tree(cmd *Command, forest ...*CommandTree) *CommandTree {
	return &CommandTree{
		command: cmd,
		forest:  forest,
	}
}

// Parse parses args to object argv
func Parse(args []string, argv interface{}) error {
	clr := color.Color{}
	fset := parseArgv(args, argv, clr)
	return fset.err
}

func parseArgv(args []string, argv interface{}, clr color.Color) *flagSet {
	return parseArgvList(args, []interface{}{argv}, clr)
}

func parseArgvList(args []string, argvList []interface{}, clr color.Color) *flagSet {
	flagSet := newFlagSet()
	for _, argv := range argvList {
		if argv == nil {
			continue
		}
		var (
			typ = reflect.TypeOf(argv)
			val = reflect.ValueOf(argv)
		)
		switch typ.Kind() {
		case reflect.Ptr:
			if reflect.Indirect(val).Type().Kind() != reflect.Struct {
				flagSet.err = errNotAPointerToStruct
				return flagSet
			}
			initFlagSet(typ, val, flagSet, clr, false)
			if flagSet.err != nil {
				return flagSet
			}
		default:
			flagSet.err = errNotAPointer
			return flagSet
		}
	}
	parseArgsToFlagSet(args, flagSet, clr)
	return flagSet
}

func usage(argvList []interface{}, clr color.Color, style UsageStyle) string {
	flagSet := newFlagSet()
	buf := bytes.NewBufferString("")
	for i := len(argvList) - 1; i >= 0; i-- {
		v := argvList[i]
		if v == nil {
			continue
		}
		var (
			typ = reflect.TypeOf(v)
			val = reflect.ValueOf(v)
		)
		if typ.Kind() == reflect.Ptr &&
			reflect.Indirect(val).Type().Kind() == reflect.Struct {
			// initialize flagSet
			initFlagSet(typ, val, flagSet, clr, true)
			if flagSet.err != nil {
				return ""
			}
		}
	}
	buf.WriteString(flagSlice(flagSet.flagSlice).StringWithStyle(clr, style))
	return buf.String()
}

func initFlagSet(typ reflect.Type, val reflect.Value, flagSet *flagSet, clr color.Color, dontSetValue bool) {
	var (
		typElem  = typ.Elem()
		valElem  = val.Elem()
		numField = valElem.NumField()
	)
	for i := 0; i < numField; i++ {
		var (
			typField          = typElem.Field(i)
			valField          = valElem.Field(i)
			tag, isEmpty, err = parseTag(typField.Name, typField.Tag)
		)
		if err != nil {
			flagSet.err = err
			return
		}
		if tag == nil {
			continue
		}

		// if `cli` tag is empty and the field is a struct
		if isEmpty && valField.Kind() == reflect.Struct {
			var (
				subObj   = valField.Addr().Interface()
				subType  = reflect.TypeOf(subObj)
				subValue = reflect.ValueOf(subObj)
			)
			initFlagSet(subType, subValue, flagSet, clr, dontSetValue)
			if flagSet.err != nil {
				return
			}
			continue
		}
		fl, err := newFlag(typField, valField, tag, clr, dontSetValue)
		if flagSet.err = err; err != nil {
			return
		}
		// ignored flag
		if fl == nil {
			continue
		}
		flagSet.flagSlice = append(flagSet.flagSlice, fl)

		// encode flag value
		value := ""
		if fl.isAssigned {
			if !valField.CanInterface() {
				flagSet.err = fmt.Errorf("field %s cannot interface", typField.Name)
				return
			}
			intf := valField.Interface()
			if encoder, ok := intf.(Encoder); ok {
				value = encoder.Encode()
			} else {
				value = fmt.Sprintf("%v", intf)
			}
		}

		names := append(fl.tag.shortNames, fl.tag.longNames...)
		for i, name := range names {
			if _, ok := flagSet.flagMap[name]; ok {
				flagSet.err = fmt.Errorf("option %s repeated", clr.Bold(name))
				return
			}
			flagSet.flagMap[name] = fl
			if dontSetValue {
				continue
			}
			if fl.isAssigned && i == 0 {
				flagSet.values[name] = []string{value}
			}
		}
	}
}

func parseArgsToFlagSet(args []string, flagSet *flagSet, clr color.Color) {
	size := len(args)
	for i := 0; i < size; i++ {
		arg := args[i]
		if !strings.HasPrefix(arg, dashOne) {
			// append a free argument
			flagSet.args = append(flagSet.args, arg)
			continue
		}

		var (
			next   = ""
			offset = 0
		)
		if i+1 < size {
			if !strings.HasPrefix(args[i+1], dashOne) {
				next = args[i+1]
				offset = 1
			}
		}

		if arg == dashOne {
			flagSet.err = fmt.Errorf("unexpected single dash")
			return
		}

		// terminate the flag parse while occur `--`
		if arg == dashTwo {
			flagSet.args = append(flagSet.args, args[i+1:]...)
			break
		}

		// split arg by "="(key=value)
		strs := []string{arg}
		index := strings.Index(arg, "=")
		if index >= 0 {
			strs = []string{arg[:index], arg[index+1:]}
		}

		arg = strs[0]
		fl, ok := flagSet.flagMap[arg]

		// found in flagMap
		if ok {
			retOffset := parseToFoundFlag(flagSet, fl, strs, arg, next, offset, clr)
			if flagSet.err != nil {
				return
			}
			i += retOffset
			continue
		}

		// not found in flagMap
		// it's an invalid flag if arg has prefix `--`
		if strings.HasPrefix(arg, dashTwo) {
			flagSet.err = fmt.Errorf("undefined option %s", clr.Bold(arg))
			return
		}

		// try parse `-F<value>`
		if _, ok := parseSiameseFlag(flagSet, arg[0:2], args[i][2:], clr); ok {
			continue
		} else if flagSet.err != nil {
			return
		}

		// other cases, find flag char by char
		arg = strings.TrimPrefix(arg, dashOne)
		parseFlagCharByChar(flagSet, arg, clr)
		if flagSet.err != nil {
			return
		}
		continue
	}

	// read delay flags
	for _, fl := range flagSet.flagSlice {
		if fl.isNeedDelaySet && fl.isAssigned {
			err := setWithProperType(fl, fl.field.Type, fl.value, fl.lastValue, clr, false)
			if flagSet.err == nil && err != nil {
				flagSet.err = err
			}
		}
		if fl.tag.isForce && fl.getBool() {
			flagSet.hasForce = true
		}
	}

	// read prompt flags
	if !flagSet.hasForce {
		if flagSet.err != nil {
			return
		}
		flagSet.readPrompt(os.Stdout, clr)
		if flagSet.err != nil {
			return
		}
		flagSet.readEditor(clr)
		if flagSet.err != nil {
			return
		}
	} else {
		flagSet.err = nil
	}

	buff := bytes.NewBufferString("")
	for _, fl := range flagSet.flagSlice {
		if !fl.isAssigned && fl.tag.isRequired {
			if buff.Len() > 0 {
				buff.WriteByte('\n')
			}
			fmt.Fprintf(buff, "required parameter %s missing", clr.Bold(fl.name()))
		}
	}
	if buff.Len() > 0 && !flagSet.hasForce {
		flagSet.err = fmt.Errorf(buff.String())
	}
}

func parseToFoundFlag(flagSet *flagSet, fl *flag, strs []string, arg, next string, offset int, clr color.Color) int {
	retOffset := 0
	l := len(strs)
	if l == 1 {
		if fl.isBoolean() {
			flagSet.err = fl.set(arg, "true", clr)
		} else if fl.isCounter() {
			fl.counterIncr("", clr)
		} else if offset > 0 {
			flagSet.err = fl.set(arg, next, clr)
			retOffset = offset
		} else {
			//flagSet.err = fmt.Errorf("missing argument")
			flagSet.err = fl.set(arg, "", clr)
		}
	} else if l == 2 {
		flagSet.err = fl.set(arg, strs[1], clr)
	} else {
		flagSet.err = fmt.Errorf("too many(%d) arguments", l)
	}
	if flagSet.err != nil {
		flagSet.err = fmt.Errorf("parameter %s invalid: %v", clr.Bold(arg), flagSet.err)
		return retOffset
	}
	flagSet.values[arg] = []string{fmt.Sprintf("%v", fl.value.Interface())}
	return retOffset
}

func parseFlagCharByChar(flagSet *flagSet, arg string, clr color.Color) {
	// NOTE: each fold flag should be boolean
	chars := []byte(arg)
	for _, c := range chars {
		tmp := dashOne + string([]byte{c})
		fl, ok := flagSet.flagMap[tmp]
		if !ok {
			flagSet.err = fmt.Errorf("undefined option %s", clr.Bold(tmp))
			return
		}

		if fl.isBoolean() {
			fl.set(tmp, "true", clr)
			flagSet.values[tmp] = []string{"true"}
		} else if fl.isCounter() {
			fl.counterIncr("", clr)
		} else {
			flagSet.err = fmt.Errorf("each fold option should be boolean, but %s not", clr.Bold(tmp))
			return
		}
	}
}

func parseSiameseFlag(flagSet *flagSet, firstHalf, latterHalf string, clr color.Color) (*flag, bool) {
	// NOTE: fl must be not a boolean
	key, val := firstHalf, latterHalf
	if fl, ok := flagSet.flagMap[key]; ok && !fl.isBoolean() {
		if fl.isCounter() {
			return nil, false
		}
		if flagSet.err = fl.set(key, val, clr); flagSet.err != nil {
			return fl, false
		}
		return fl, true
	}
	return nil, false
}
