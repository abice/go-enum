package debug

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"

	"github.com/labstack/gommon/color"
	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
)

var enableDebug = false

func Switch(on bool) {
	enableDebug = on
}

var gopaths = func() []string {
	paths := strings.Split(os.Getenv("GOPATH"), ";")
	for i, path := range paths {
		paths[i] = filepath.Join(path, "src") + "/"
	}
	if goroot := runtime.GOROOT(); goroot != "" {
		paths = append(paths, filepath.Join(goroot, "src")+"/")
	}
	return paths
}()

var debugOut, debugColor = func() (io.Writer, color.Color) {
	clr := color.Color{}
	out := colorable.NewColorableStdout()
	ColorSwitch(&clr, out, os.Stdout.Fd())
	return out, clr
}()

func Debugf(format string, args ...interface{}) {
	if !enableDebug {
		return
	}
	_, file, line, _ := runtime.Caller(1)
	fileline := makeFileLine(file, line)
	fmt.Fprintf(debugOut, "[DEBUG]"+debugColor.Bold(fileline)+" "+format+"\n", args...)
}

func Panicf(format string, args ...interface{}) {
	buff := bytes.NewBufferString("")
	buff.WriteString(fmt.Sprintf(format, args...))
	buff.WriteString("\n\n[stack]\n")
	skip := 1
	for {
		_, file, line, ok := runtime.Caller(skip)
		if !ok {
			break
		}
		skip++
		buff.WriteString(makeFileLine(file, line))
		buff.WriteString("\n")
	}
	panic(buff.String())
}

func TypeName(i interface{}) string {
	t := reflect.TypeOf(i)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.Name()
}

func JSON(v interface{}) string {
	data, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return ""
	}
	return string(data)
}

func makeFileLine(file string, line int) string {
	for _, path := range gopaths {
		if strings.HasPrefix(file, path) {
			file = strings.TrimPrefix(strings.TrimPrefix(file, path), "/")
			break
		}
	}
	return fmt.Sprintf("%s:%d", file, line)
}

func ColorSwitch(clr *color.Color, w io.Writer, fds ...uintptr) {
	clr.Disable()
	if len(fds) > 0 {
		if isatty.IsTerminal(fds[0]) {
			clr.Enable()
		}
	} else if w, ok := w.(*os.File); ok && isatty.IsTerminal(w.Fd()) {
		clr.Enable()
	}
}
