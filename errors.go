// Package errors 封装了携带堆栈的统一错误.
package errors

import (
	"fmt"
	"io"
)

// fundamentalError 定义了包含 stack 的 error
type fundamentalError struct {
	cause  error      // cause 内部嵌套错误, 构造错误链
	code   int        // code 传入的错误码, 选填
	msg    string     // msg 传入的错误描述, 可对外暴露
	detail string     // detail 在错误描述 msg 的基础上, 增加文件名、行数、调用函数名等信息, 不对外暴露, 服务内部使用
	stack  StackTrace // stack 错误堆栈, 如果内部嵌套错误 cause 已有堆栈, 则不再设置
}

// Error 实现 Error 接口, 打印链路 detail 信息, 包含文件名、行数等
func (fd *fundamentalError) Error() string {
	if fd.detail != "" && fd.cause != nil {
		return fmt.Sprintf("%v"+defaultCfg.ErrorConnectionFlag+"%v", fd.detail, fd.cause.Error())
	}
	if fd.detail != "" {
		return fmt.Sprintf("%v", fd.detail)
	}
	if fd.cause != nil {
		return fmt.Sprintf("%v", fd.cause.Error())
	}
	return ""
}

// Cause 返回内部的错误
func (fd *fundamentalError) Cause() error { return fd.cause }

// Unwrap 支持Go 1.13+ error chains.
func (fd *fundamentalError) Unwrap() error { return fd.cause }

// Code 返回 code
func (fd *fundamentalError) Code() int {
	return fd.code
}

// Msg 返回 msg
func (fd *fundamentalError) Msg() string {
	return fd.msg
}

// Format 实现 Format 接口
func (fd *fundamentalError) Format(s fmt.State, verb rune) {
	var stackTrace StackTrace
	defer func() {
		if stackTrace != nil {
			stackTrace.Format(s, verb)
		}
	}()
	switch verb {
	case 'v':
		if s.Flag('+') {
			if fd.detail != "" {
				_, _ = io.WriteString(s, fd.detail)
			}
			if fd.Cause() != nil {
				_, _ = fmt.Fprintf(s, defaultCfg.ErrorConnectionFlag+"%+v", fd.Cause())
			}
			if fd.stack != nil {
				stackTrace = fd.stack
			}
			return
		}
		fallthrough
	case 's':
		_, _ = io.WriteString(s, fd.Error())
	case 'q':
		_, _ = fmt.Fprintf(s, "%q", fd.Error())
	default:
		_, _ = fmt.Fprintf(s, "unsupported format: %%!%c, use %%s: %s)", verb, fd.Error())
	}
}
