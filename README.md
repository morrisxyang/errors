# errors
[![Go Reference](https://pkg.go.dev/badge/github.com/morrisxyang/errors.svg)](https://pkg.go.dev/github.com/morrisxyang/errors)
![Static Badge](https://img.shields.io/badge/License-BSD2-Green)
[![Coverage Status](https://coveralls.io/repos/github/morrisxyang/errors/badge.svg?branch=master)](https://coveralls.io/github/morrisxyang/errors?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/morrisxyang/errors)](https://goreportcard.com/report/github.com/morrisxyang/errors)
![Static Badge](https://img.shields.io/badge/go%20verion-%3E%3D1.15-blue)

A simple error library that supports **error stacks**, **error codes**, and **error chains**:

- Supports carrying **stacks** and constructing nested error **chains**

- Supports carrying **error codes**

- Supports customizing the **depth of stack** printing and error chain **printing format**

- Uses CallersFrames instead of FuncForPC to generate stacks, avoiding issues such as "line number errors" in special cases, see [runtime: strongly encourage using CallersFrames over FuncForPC with Callers result](https://github.com/golang/go/issues/19426)

- Simplifies stack information when using multiple `Wrap` operations by only keeping **the deepest stack** in a chain and printing it only once.

[中文README](https://github.com/morrisxyang/errors/blob/master/README_CN.md)

## Installation and Docs

Install using `go get github.com/morrisxyang/errors`.

Full documentation is available at https://pkg.go.dev/github.com/morrisxyang/errors

## Quick Start

Construct error chain

```go
func a() error {
	err := b()
	err = Wrap(err, "a failed reason")
	return err
}

func b() error {
	err := c()
	err = Wrap(err, "b failed reason")
	return err
}

func c() error {
	_, err := os.Open("test")
	if err != nil {
		return WrapWithCode(err, 123, "c failed reason")
	}
	return nil
}
```

Print error message. `%+v` will print the error stack trace and `%v` only prints the error message.

```go
a failed reason
Caused by: b failed reason
Caused by: 123, c failed reason
Caused by: open test: no such file or directory
github.com/morrisxyang/errors.c
	/Users/morrisyang/Nutstore Files/go-proj/githuberrors/errors_test.go:94
github.com/morrisxyang/errors.b
	/Users/morrisyang/Nutstore Files/go-proj/githuberrors/errors_test.go:86
github.com/morrisxyang/errors.a
	/Users/morrisyang/Nutstore Files/go-proj/githuberrors/errors_test.go:80
....
```

## Core Methods

### Error Chain

- [func New(msg string) error](https://pkg.go.dev/github.com/morrisxyang/errors#New)
- [func Newf(format string, args ...interface{}) error](https://pkg.go.dev/github.com/morrisxyang/errors#Newf)
- [func NewWithCode(code int, msg string) error](https://pkg.go.dev/github.com/morrisxyang/errors#NewWithCode)
- [func NewWithCodef(code int, format string, args ...interface{}) error](https://pkg.go.dev/github.com/morrisxyang/errors#NewWithCodef)
- [func Wrap(e error, msg string) error](https://pkg.go.dev/github.com/morrisxyang/errors#Wrap)
- [func Wrapf(e error, format string, args ...interface{}) error](https://pkg.go.dev/github.com/morrisxyang/errors#Wrapf)
- [func WrapWithCode(e error, code int, msg string) error](https://pkg.go.dev/github.com/morrisxyang/errors#WrapWithCode)
- [func WrapWithCodef(e error, code int, format string, args ...interface{}) error](https://pkg.go.dev/github.com/morrisxyang/errors#WrapWithCodef)

### Error Handling

- [func Code(e error) int](https://pkg.go.dev/github.com/morrisxyang/errors#Code)
- [func EffectiveCode(e error) int](https://pkg.go.dev/github.com/morrisxyang/errors#EffectiveCode)
- [func Msg(e error) string](https://pkg.go.dev/github.com/morrisxyang/errors#Msg)
- [func As(err error, target interface{}) bool](https://pkg.go.dev/github.com/morrisxyang/errors#As)
- [func Is(err, target error) bool](https://pkg.go.dev/github.com/morrisxyang/errors#Is)
- [func Cause(e error) error](https://pkg.go.dev/github.com/morrisxyang/errors#Cause)
- [func Unwrap(err error) error](https://pkg.go.dev/github.com/morrisxyang/errors#Unwrap)

### Config

- [type Config](https://pkg.go.dev/github.com/morrisxyang/errors#Config)
- - [func GetCfg() *Config](https://pkg.go.dev/github.com/morrisxyang/errors#GetCfg)

- [func ResetCfg()](https://pkg.go.dev/github.com/morrisxyang/errors#ResetCfg)
- [func SetCfg(c *Config)](https://pkg.go.dev/github.com/morrisxyang/errors#SetCfg)

## FAQ

### Will multiple Wrap errors carry multiple stacks?

You can Wrap multiple times with explanatory information in the calling chain, but only the deepest Wrap operation will set the stack. Continuing to `Wrap`, `return err` and other operations will not affect the stack information.

### If a suitable error code is set for an error in the chain, but not set when continuing to Wrap, how can it be obtained?
It is recommended to set the valid error code at an appropriate and clear time. You can use `EffectiveCode` to obtain the first valid non-zero error code outside the link layer. Due to system calls and other situations, there may be multiple errors carrying error codes in the same link, in which case the error code of the outer layer should be exposed to the outside world by default, shielding the detailed information of the inner layer.
