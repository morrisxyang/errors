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
