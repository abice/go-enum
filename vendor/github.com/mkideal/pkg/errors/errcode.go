package errors

import (
	"fmt"
)

type Status int

const StatusInvalid Status = -65535

func (status Status) Print(msg string) ErrorCode {
	return ErrorCode{
		code: status,
		msg:  msg,
	}
}

func (status Status) Printf(format string, args ...interface{}) ErrorCode {
	return status.Print(fmt.Sprintf(format, args...))
}

type ErrorCode struct {
	code Status
	msg  string
}

func (err ErrorCode) Error() string { return err.msg }
func (err ErrorCode) Code() Status  { return err.code }

func (err ErrorCode) Print(msg string) ErrorCode {
	return ErrorCode{
		code: err.code,
		msg:  err.msg + ": " + msg,
	}
}

func (err ErrorCode) Printf(format string, args ...interface{}) ErrorCode {
	return err.Print(fmt.Sprintf(format, args...))
}

func Code(err error) Status {
	err = Core(err)
	if err != nil {
		if e, ok := err.(ErrorCode); ok {
			return e.Code()
		}
	}
	return StatusInvalid
}

var (
	e1 = Status(1).Print("e1")
)
