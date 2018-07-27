package errors

import (
	"sync/atomic"

	"github.com/mkideal/pkg/debug"
)

// Error aliases string
type Error string

func (e Error) Error() string { return string(e) }

var enableTrace int32 = 1

// SwitchTrace switchs trace on/off
func SwitchTrace(on bool) {
	if on {
		atomic.StoreInt32(&enableTrace, 1)
	} else {
		atomic.StoreInt32(&enableTrace, 0)
	}
}

// traceError wraps string and contains created stack information
type traceError struct {
	stack string
	text  string
}

func (e traceError) Error() string {
	return e.stack + e.text
}

// Throw throws an error which contains stack information
func Throw(text string) error {
	if atomic.LoadInt32(&enableTrace) != 0 {
		return traceError{stack: string(debug.Stack(2)) + "\n", text: text}
	}
	return traceError{text: text}
}

// wrappedError define an interface which wraps another error
type wrappedError interface {
	error
	Core() error
}

var _ = wrappedError(wrapStackError{})

// wrapStackError implements wrappedError
type wrapStackError struct {
	stack string
	err   error
}

func (e wrapStackError) String() string {
	return e.stack + "\n" + e.err.Error()
}

func (e wrapStackError) Error() string {
	if atomic.LoadInt32(&enableTrace) == 0 {
		return e.err.Error()
	}
	return e.String()
}

func (e wrapStackError) Core() error { return e.err }

// Wrap wraps another error
func Wrap(err error) error {
	return wrapStackError{stack: string(debug.Stack(2)), err: err}
}

// Core returns wrapped error
func Core(err error) error {
	const maxWrapLayer = 32
	c := 0
	for err != nil && c < maxWrapLayer {
		c++
		we, ok := err.(wrappedError)
		if !ok {
			break
		}
		err = we.Core()
	}
	return err
}
