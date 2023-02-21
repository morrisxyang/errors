package errors

import (
	stderrors "errors"
	"fmt"
	"math"
)

// New 使用传入的信息创建一个携带堆栈的错误
func New(msg string) error {
	return &fundamentalError{
		msg:   msg,
		stack: callers(),
	}
}

// Errorf formats according to a format specifier and returns the string
// as a value that satisfies error.
// Errorf also records the stack trace at the point it was called.
func Errorf(format string, args ...interface{}) error {
	return &fundamentalError{
		msg:   fmt.Sprintf(format, args...),
		stack: callers(),
	}
}

// Newf 使用 format 格式的信息创建一个携带堆栈的错误
func Newf(format string, args ...interface{}) error {
	return &fundamentalError{
		msg:   fmt.Sprintf(format, args...),
		stack: callers(),
	}
}

// NewWithCode 使用传入的信息创建一个携带堆栈的错误
func NewWithCode(code int, msg string) error {
	return &fundamentalError{
		msg:   msg,
		stack: callers(),
		code:  code,
	}
}

// NewWithCodef 使用 format 格式的信息创建一个携带堆栈的错误
func NewWithCodef(code int, format string, args ...interface{}) error {
	return &fundamentalError{
		msg:   fmt.Sprintf(format, args...),
		stack: callers(),
		code:  code,
	}
}

// Wrap 使用传入的信息包装错误, 携带堆栈信息
// 如果传入的 err 已经有错误码, 直接延用
// 如果传入的 err 已经有堆栈, 不再设置堆栈
// 如果传入的 err 为 nil, Wrap 将返回 nil
func Wrap(err error, msg string) error {
	if err == nil {
		return nil
	}
	wrapErr := &fundamentalError{
		cause: err,
		msg:   msg,
	}
	var fd *fundamentalError
	if stderrors.As(err, &fd) {
		// 链路上有同类型错误的时候，延用 code
		wrapErr.code = fd.code
	} else {
		// 链路上没有同类型错误的时候，证明是首次包装, 添加上堆栈信息
		wrapErr.stack = callers()
	}
	return wrapErr
}

// Wrapf 使用 format 格式的信息包装错误, 携带堆栈信息
// 如果传入的 err 已经有错误码, 直接延用
// 如果传入的 err 已经有堆栈, 不再设置堆栈
// 如果传入的 err 为 nil, Wrapf 将返回 nil
func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	wrapErr := &fundamentalError{
		cause: err,
		msg:   fmt.Sprintf(format, args...),
	}
	var fd *fundamentalError
	if stderrors.As(err, &fd) {
		// 链路上有同类型错误的时候，延用 code
		wrapErr.code = fd.code
	} else {
		// 链路上没有同类型错误的时候，证明是首次包装, 添加上堆栈信息
		wrapErr.stack = callers()
	}
	return wrapErr
}

// WrapWithCode 使用传入的信息包装错误, 携带堆栈信息
// 即使传入的 err 已经有错误码, 不会延用, 而是使用传入的 code
// 如果传入的 err 已经有堆栈, 不再设置堆栈
// 如果传入的 err 为 nil, WrapWithCode 将返回 nil
func WrapWithCode(err error, code int, msg string) error {
	if err == nil {
		return nil
	}
	wrapErr := &fundamentalError{
		cause: err,
		msg:   msg,
		code:  code,
	}
	var fd *fundamentalError
	if !stderrors.As(err, &fd) {
		// 链路上没有同类型错误的时候，证明是首次包装, 添加上堆栈信息
		wrapErr.stack = callers()
	}
	return wrapErr
}

// WrapWithCodef 使用 format 格式的信息包装错误, 携带堆栈信息
// 即使传入的 err 已经有错误码, 不会延用, 而是使用传入的 code
// 如果传入的 err 已经有堆栈, 不再设置堆栈
// 如果传入的 err 为 nil, WrapWithCodef 将返回 nil
func WrapWithCodef(err error, code int, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	wrapErr := &fundamentalError{
		cause: err,
		msg:   fmt.Sprintf(format, args...),
		code:  code,
	}
	var fd *fundamentalError
	if !stderrors.As(err, &fd) {
		// 链路上没有同类型错误的时候，证明是首次包装, 添加上堆栈信息
		wrapErr.stack = callers()
	}
	return wrapErr
}

// Code 获取 code
func Code(e error) int {
	if e == nil {
		return 0
	}
	// int32最小值
	const unknownCode = math.MinInt32
	err, ok := e.(*fundamentalError)
	if !ok {
		return unknownCode
	}
	return err.Code()
}

// Msg 获取 msg
func Msg(e error) string {
	if e == nil {
		return ""
	}
	const unknownMsg = "unknown error"
	err, ok := e.(*fundamentalError)
	if !ok {
		return unknownMsg
	}
	return err.Msg()
}

// Cause returns the underlying cause of the error, if possible.
// An error value has a cause if it implements the following
// interface:
//
//	type causer interface {
//	       Cause() error
//	}
//
// If the error does not implement Cause, the original error will
// be returned. If the error is nil, nil will be returned without further
// investigation.
func Cause(err error) error {
	type causer interface {
		Cause() error
	}

	if err == nil {
		return err
	}

	for {
		cause, ok := err.(causer)
		if !ok {
			break
		}
		if cause.Cause() != nil {
			err = cause.Cause()
			continue
		}
		break
	}
	return err
}
