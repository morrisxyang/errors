// Package errors 封装了携带堆栈的统一错误.
package errors

import (
	"runtime"
	"strconv"
	"strings"
)

// callers 获取堆栈
func callers() *StackTrace {
	// depth 记录的栈深度
	const depth = 10
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:])
	// 转换为 errors.StackTrace
	stack := runtime.CallersFrames(pcs[0:n])
	return &StackTrace{stack}
}

// callerFuncInfo 调用方函数名
func callerFuncInfo() string {
	pc, fileName, line, _ := runtime.Caller(2)
	f := runtime.FuncForPC(pc)
	callerFuncName := f.Name()
	callerFuncName = callerFuncName[strings.LastIndex(callerFuncName, ".")+1:]
	fileName = fileName[strings.LastIndex(fileName, "/")+1:]
	return callerFuncName + "(" + fileName + ":" + strconv.Itoa(line) + ")"
}
