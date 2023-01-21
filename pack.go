package errors

import (
	"fmt"
)

// New 使用传入的信息创建一个携带堆栈的错误
func New(msg string) error {
	return &fundamentalError{
		msg:    msg,
		detail: fmt.Sprintf("%v, %v", callerFuncInfo(), msg),
		stack:  callers(),
	}
}

// Newf 使用 format 格式的信息创建一个携带堆栈的错误
func Newf(format string, args ...interface{}) error {
	return &fundamentalError{
		msg:    fmt.Sprintf(format, args...),
		detail: fmt.Sprintf("%v, %v", callerFuncInfo(), fmt.Sprintf(format, args...)),
		stack:  callers(),
	}
}

// NewWithCode 使用传入的信息创建一个携带堆栈的错误
func NewWithCode(code int, msg string) error {
	return &fundamentalError{
		msg:    msg,
		detail: fmt.Sprintf("%v, %v", callerFuncInfo(), msg),
		stack:  callers(),
		code:   code,
	}
}

// NewWithCodef 使用 format 格式的信息创建一个携带堆栈的错误
func NewWithCodef(code int, format string, args ...interface{}) error {
	return &fundamentalError{
		msg:    fmt.Sprintf(format, args...),
		detail: fmt.Sprintf("%v, %v", callerFuncInfo(), fmt.Sprintf(format, args...)),
		stack:  callers(),
		code:   code,
	}
}
