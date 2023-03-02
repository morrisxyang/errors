// Package errors provides wrapped errors with stack trace.
package errors

import (
	"bytes"
	"fmt"
	"io"
)

// baseError defines an error that includes a stack trace.
type baseError struct {
	cause error       // cause is the nested error, building an error chain
	code  int         // code is the error code
	msg   string      // msg is the error description
	stack *StackTrace // stack is the error stack, if the error chain already has a stack, it will not be set again
}

// Error implements the Error interface to print the error chain information.
func (b *baseError) Error() string {
	var buffer bytes.Buffer
	if b.code != 0 {
		buffer.WriteString(fmt.Sprintf("%d", b.code))
	}
	if b.msg != "" {
		if buffer.Len() > 0 {
			buffer.WriteString(", ")
		}
		buffer.WriteString(b.msg)
	}
	if b.cause != nil {
		if buffer.Len() > 0 {
			buffer.WriteString(GetCfg().ErrorConnectionFlag)
		}
		buffer.WriteString(b.Cause().Error())
	}
	return buffer.String()
}

// Format implements the Format interface for printing.
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

// StackTrace returns the error chain stack trace.
// The deepest error created will carry the stack information and shallow errors will not repeat the record.
func (b *baseError) StackTrace() StackTrace {
	e := b
	for e != nil {
		if e.stack != nil {
			break
		}
		e, _ = e.Cause().(*baseError)
	}
	return *e.stack
}

// Code returns the code.
func (b *baseError) Code() int {
	return b.code
}

// Msg returns the message.
func (b *baseError) Msg() string {
	return b.msg
}

// Cause returns the nested error.
func (b *baseError) Cause() error { return b.cause }

// Unwrap supports Go 1.13+ error chains.
func (b *baseError) Unwrap() error { return b.cause }
