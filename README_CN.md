# errors
[![Go Reference](https://pkg.go.dev/badge/github.com/morrisxyang/errors.svg)](https://pkg.go.dev/github.com/morrisxyang/errors)
![Static Badge](https://img.shields.io/badge/License-BSD2-Green)
[![Coverage Status](https://coveralls.io/repos/github/morrisxyang/errors/badge.svg?branch=master)](https://coveralls.io/github/morrisxyang/errors?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/morrisxyang/errors)](https://goreportcard.com/report/github.com/morrisxyang/errors)
![Static Badge](https://img.shields.io/badge/go%20verion-%3E%3D1.15-blue)

简单的支持**错误堆栈**, **错误码**, **错误链**的工具库:

- 支持携带堆栈, 嵌套构造错误链

- 支持携带错误码, 方便接口返回

- 支持自定义堆栈打印深度和错误链打印格式

- 使用 CallersFrames 替代 FuncForPC 生成堆栈, 避免特殊情况`line number错误`等问题, 详见[runtime: strongly encourage using CallersFrames over FuncForPC with Callers result](https://github.com/golang/go/issues/19426)

- 简化堆栈信息, 一条链路只保留最深层堆栈, 只打印一次.

## 安装和文档

安装使用 `go get github.com/morrisxyang/errors`

文档地址是 https://pkg.go.dev/github.com/morrisxyang/errors


## 快速开始

构造错误链

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

打印错误信息, `%+v`会打印堆栈, `%v`只打印错误信息

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
....堆栈信息省略
```

## 核心方法

### 错误封装

- [func New(msg string) error](https://pkg.go.dev/github.com/morrisxyang/errors#New)
- [func Newf(format string, args ...interface{}) error](https://pkg.go.dev/github.com/morrisxyang/errors#Newf)
- [func NewWithCode(code int, msg string) error](https://pkg.go.dev/github.com/morrisxyang/errors#NewWithCode)
- [func NewWithCodef(code int, format string, args ...interface{}) error](https://pkg.go.dev/github.com/morrisxyang/errors#NewWithCodef)
- [func Wrap(e error, msg string) error](https://pkg.go.dev/github.com/morrisxyang/errors#Wrap)
- [func Wrapf(e error, format string, args ...interface{}) error](https://pkg.go.dev/github.com/morrisxyang/errors#Wrapf)
- [func WrapWithCode(e error, code int, msg string) error](https://pkg.go.dev/github.com/morrisxyang/errors#WrapWithCode)
- [func WrapWithCodef(e error, code int, format string, args ...interface{}) error](https://pkg.go.dev/github.com/morrisxyang/errors#WrapWithCodef)

### 错误解析

- [func Code(e error) int](https://pkg.go.dev/github.com/morrisxyang/errors#Code)
- [func EffectiveCode(e error) int](https://pkg.go.dev/github.com/morrisxyang/errors#EffectiveCode)
- [func Msg(e error) string](https://pkg.go.dev/github.com/morrisxyang/errors#Msg)
- [func As(err error, target interface{}) bool](https://pkg.go.dev/github.com/morrisxyang/errors#As)
- [func Is(err, target error) bool](https://pkg.go.dev/github.com/morrisxyang/errors#Is)
- [func Cause(e error) error](https://pkg.go.dev/github.com/morrisxyang/errors#Cause)
- [func Unwrap(err error) error](https://pkg.go.dev/github.com/morrisxyang/errors#Unwrap)

### 配置

- [type Config](https://pkg.go.dev/github.com/morrisxyang/errors#Config)
- - [func GetCfg() *Config](https://pkg.go.dev/github.com/morrisxyang/errors#GetCfg)

- [func ResetCfg()](https://pkg.go.dev/github.com/morrisxyang/errors#ResetCfg)
- [func SetCfg(c *Config)](https://pkg.go.dev/github.com/morrisxyang/errors#SetCfg)

## FAQ

1. 多次 Wrap 错误会携带多次堆栈吗?

   可在调用链路上多次Wrap, 添加说明信息, 但只有最深层的Wrap操作会设置堆栈, 继续 `Wrap`, `return err` 等操作不会影响堆栈信息

2. 在链路中某个错误设置了合适的错误码, 然后继续Wrap时没有设置, 如何获取?

   建议在合适的清晰的时机设置有效的错误码, 可以使用`EffectiveCode`获取链路中外层第一个有效的非0错误码, 由于系统调用等情况, 同一链路中可能有多个错误携带错误码, 此时默认外层的错误码应该对外暴露, 屏蔽了内层的详细信息.
