package errors

import (
	"runtime"
)

// callers function retrieves the stack trace of the current goroutine.
func callers() *StackTrace {
	// constant to limit the depth of the stack trace
	const maxDepth = 64
	var pcs [maxDepth]uintptr
	n := runtime.Callers(3, pcs[:])

	cfg := GetCfg()
	var stack *runtime.Frames
	// if StackDepth is set and less than total number of frames then limit stack trace depth
	if cfg.StackDepth > 0 && cfg.StackDepth < n {
		stack = runtime.CallersFrames(pcs[0:cfg.StackDepth])
	} else {
		// otherwise, return all frames
		stack = runtime.CallersFrames(pcs[0:n])
	}

	return &StackTrace{*stack}
}
