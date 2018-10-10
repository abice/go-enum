package cli

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/labstack/gommon/color"
)

var (
	errNotAPointerToStruct = errors.New("not a pointer to struct")
	errNotAPointer         = errors.New("argv is not a pointer")
	errCliTagTooMany       = errors.New("cli tag too many")
)

type (
	exitError struct{}

	commandNotFoundError struct {
		command string
	}

	methodNotAllowedError struct {
		method string
	}

	routerRepeatError struct {
		router string
	}

	wrapError struct {
		err error
		msg string
	}

	argvError struct {
		isEmpty      bool
		isOutOfRange bool

		ith int
		msg string
	}
)

func (e exitError) Error() string { return "exit" }

// ExitError is a special error, should be ignored but return
var ExitError = exitError{}

func throwCommandNotFound(command string) commandNotFoundError {
	return commandNotFoundError{command: command}
}

func throwMethodNotAllowed(method string) methodNotAllowedError {
	return methodNotAllowedError{method: method}
}

func throwRouterRepeat(router string) routerRepeatError {
	return routerRepeatError{router: router}
}

func (e commandNotFoundError) Error() string {
	return fmt.Sprintf("command %s not found", e.command)
}

func (e methodNotAllowedError) Error() string {
	return fmt.Sprintf("method %s not allowed", e.method)
}

func (e routerRepeatError) Error() string {
	return fmt.Sprintf("router %s repeat", e.router)
}

func (e wrapError) Error() string {
	return e.msg
}

func wrapErr(err error, appendString string, clr color.Color) error {
	if err == nil {
		return err
	}
	errs := strings.Split(err.Error(), "\n")
	buff := bytes.NewBufferString("")
	errPrefix := clr.Red("ERR!") + " "
	for i, e := range errs {
		if i != 0 {
			buff.WriteByte('\n')
		}
		buff.WriteString(errPrefix)
		buff.WriteString(e)
	}
	buff.WriteString(appendString)
	return wrapError{err: err, msg: buff.String()}
}

func (e argvError) Error() string {
	if e.isEmpty {
		return "argv list is empty"
	}
	if e.isOutOfRange {
		return "argv list out of range"
	}
	return fmt.Sprintf("%dth argv: %s", e.ith, e.msg)
}
