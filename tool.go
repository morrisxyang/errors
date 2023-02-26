// Package errors 封装了携带堆栈的统一错误.
package errors

import (
	"runtime"
)

// callers 获取堆栈
func callers() *StackTrace {
	// maxDepth 记录的栈深度
	const maxDepth = 64
	var pcs [maxDepth]uintptr
	n := runtime.Callers(3, pcs[:])

	var stack *runtime.Frames
	cfg := GetCfg()
	if cfg.StackDepth > 0 && cfg.StackDepth < n {
		stack = runtime.CallersFrames(pcs[0:cfg.StackDepth])
	} else {
		// 转换为 errors.StackTrace
		stack = runtime.CallersFrames(pcs[0:n])
	}

	return &StackTrace{*stack}
}
