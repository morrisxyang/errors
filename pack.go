package errors

import (
	stderrors "errors"
	"fmt"
	"math"
)

const (
	// UnknownCode is the error code for unspecified errors.
	UnknownCode = math.MinInt32
	// Success is the success prompt string.
	Success = "success"
)

// New creates an error with a stack trace using the provided message
func New(msg string) error {
	return &baseError{
		msg:   msg,
		stack: callers(),
	}
}

// Errorf formats according to a format specifier and returns the string
// as a value that satisfies error.
// Errorf also records the stack trace at the point it was called.
func Errorf(format string, args ...interface{}) error {
	return &baseError{
		msg:   fmt.Sprintf(format, args...),
		stack: callers(),
	}
}

// Newf creates a new error with the provided format specifier and arguments.
// It has the same functionality as New function
func Newf(format string, args ...interface{}) error {
	return &baseError{
		msg:   fmt.Sprintf(format, args...),
		stack: callers(),
	}
}

// NewWithCode creates a new error with a stack trace, using the provided code and message.
func NewWithCode(code int, msg string) error {
	return &baseError{
		msg:   msg,
		stack: callers(),
		code:  code,
	}
}

// NewWithCodef creates a new error with a stack trace, the provided code, format specifier and arguments.
// This function has the same functionality as the NewWithCode function.
func NewWithCodef(code int, format string, args ...interface{}) error {
	return &baseError{
		msg:   fmt.Sprintf(format, args...),
		stack: callers(),
		code:  code,
	}
}

// Wrap function wraps the incoming error with stack information and message.
// If the incoming err already has a stack, the stack will not be set again.
// If the incoming err is nil, Wrap will return nil.
func Wrap(e error, msg string) error {
	// check if e is nil
	if e == nil {
		return nil
	}
	wrapErr := &baseError{
		cause: e,
		msg:   msg,
	}
	var fd *baseError
	if !stderrors.As(e, &fd) {
		// If there is no error of the same type on the link, it means that it is the first time to package and add stack information
		wrapErr.stack = callers()
	}
	return wrapErr
}

// Wrapf function wraps the incoming error with stack information, format specifier and arguments.
// This function has the same functionality as the Wrap function.
func Wrapf(e error, format string, args ...interface{}) error {
	if e == nil {
		return nil
	}
	wrapErr := &baseError{
		cause: e,
		msg:   fmt.Sprintf(format, args...),
	}
	var fd *baseError
	if !stderrors.As(e, &fd) {
		// If there is no error of the same type on the link, it means that it is the first time to package and add stack information
		wrapErr.stack = callers()
	}
	return wrapErr
}

// WrapWithCode function wraps the incoming error with stack information, code and message.
// If the incoming err already has a stack, the stack will not be set again.
// If the incoming err is nil, WrapWithCode will return nil.
func WrapWithCode(e error, code int, msg string) error {
	if e == nil {
		return nil
	}
	wrapErr := &baseError{
		cause: e,
		msg:   msg,
		code:  code,
	}
	var fd *baseError
	if !stderrors.As(e, &fd) {
		// If there is no error of the same type on the link, it means that it is the first time to package and add stack information
		wrapErr.stack = callers()
	}
	return wrapErr
}

// WrapWithCodef function wraps the incoming error with stack information, code, format specifier and arguments.
// This function has the same functionality as the WrapWithCode function.
func WrapWithCodef(e error, code int, format string, args ...interface{}) error {
	if e == nil {
		return nil
	}
	wrapErr := &baseError{
		cause: e,
		msg:   fmt.Sprintf(format, args...),
		code:  code,
	}
	var fd *baseError
	if !stderrors.As(e, &fd) {
		// If there is no error of the same type on the link, it means that it is the first time to package and add stack information
		wrapErr.stack = callers()
	}
	return wrapErr
}

// Code function returns the error code associated with an error object if it is of type *baseError.
// If the error object is not of type *baseError, it returns the minimum value of int32.
func Code(e error) int {
	if e == nil {
		return 0
	}
	err, ok := e.(*baseError)
	if !ok {
		return UnknownCode
	}
	if err == (*baseError)(nil) {
		return 0
	}
	return err.Code()
}

// Msg function returns the error message associated with an error object if it is of type *baseError.
// If the error object is not of type *baseError, it returns it's Error().
func Msg(e error) string {
	if e == nil {
		return ""
	}
	err, ok := e.(*baseError)
	if !ok {
		return e.Error()
	}
	if err == (*baseError)(nil) {
		return Success
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
func Cause(e error) error {
	if e == nil {
		return e
	}
	type causer interface {
		Cause() error
	}
	for {
		cause, ok := e.(causer)
		if !ok {
			break
		}
		if cause.Cause() != nil {
			e = cause.Cause()
			continue
		}
		break
	}
	return e
}

// EffectiveCode returns the first valid error code from the error chain.
// If error object encountered is not of type *baseError, it will return UnknownCode.
func EffectiveCode(e error) int {
	if e == nil {
		return 0
	}
	type (
		causer interface {
			Cause() error
		}
	)
	for {
		err, ok := e.(*baseError)
		if !ok {
			break
		}
		if err == (*baseError)(nil) {
			return 0
		}
		// effective
		if err.Code() != 0 {
			return err.Code()
		}
		// through chain
		cause, ok := e.(causer)
		if !ok {
			break
		}
		if cause.Cause() != nil {
			e = cause.Cause()
			continue
		}
		break
	}
	// If no effective code was found, return UnknownCode.
	return UnknownCode
}
