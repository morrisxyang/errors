// Package errors provides wrapped errors with stack trace.
package errors

import (
	"fmt"
	"io"
)

// baseError defines an error that includes a stack trace.
type baseError struct {
	cause error       // cause 内部嵌套错误, 构造错误链
	code  int         // code 错误码
	msg   string      // msg 错误描述
	stack *StackTrace // stack 错误堆栈, 如果错误链已有堆栈, 则不再设置
}

// Error 实现 Error 接口, 打印错误链路信息
func (b *baseError) Error() string {
	if b.msg != "" && b.cause != nil {
		if b.code != 0 {
			return fmt.Sprintf("%d, %s"+GetCfg().ErrorConnectionFlag+"%s",
				b.code, b.msg, b.cause.Error())
		}
		return fmt.Sprintf("%s"+GetCfg().ErrorConnectionFlag+"%s", b.msg, b.cause.Error())
	}
	if b.msg != "" {
		if b.code != 0 {
			return fmt.Sprintf("%d, %s", b.code, b.msg)
		}
		return fmt.Sprintf("%s", b.msg)
	}
	if b.cause != nil {
		return fmt.Sprintf("%s", b.cause.Error())
	}
	return ""
}

// Format 实现 Format 接口
func (b *baseError) Format(s fmt.State, verb rune) {
	var stackTrace *StackTrace
	defer func() {
		if stackTrace != nil {
			stackTrace.Format(s, verb)
		}
	}()
	switch verb {
	case 'v':
		if s.Flag('+') {
			if b.msg != "" {
				_, _ = io.WriteString(s, b.msg)
			}
			if b.Cause() != nil {
				_, _ = fmt.Fprintf(s, GetCfg().ErrorConnectionFlag+"%+v", b.Cause())
			}
			if b.stack != nil {
				stackTrace = b.stack
			}
			return
		}
		fallthrough
	case 's':
		_, _ = io.WriteString(s, b.Error())
	case 'q':
		_, _ = fmt.Fprintf(s, "%q", b.Error())
	default:
		_, _ = fmt.Fprintf(s, "unsupported format: %%!%c, use %%s: %s", verb, b.Error())
	}
}

// StackTrace ...
func (b *baseError) StackTrace() StackTrace {
	f := b
	for f != nil {
		if f.stack != nil {
			break
		}
		f, _ = f.Cause().(*baseError)
	}
	return *f.stack
}

// Code 返回 code
func (b *baseError) Code() int {
	return b.code
}

// Msg 返回 msg
func (b *baseError) Msg() string {
	return b.msg
}

// Cause 返回内部的错误
func (b *baseError) Cause() error { return b.cause }

// Unwrap 支持Go 1.13+ error chains.
func (b *baseError) Unwrap() error { return b.cause }
