// Package errors 封装了携带堆栈的统一错误.
package errors

import (
	"runtime"
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
