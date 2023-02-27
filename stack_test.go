package errors

import (
	"fmt"
	"regexp"
	"runtime"
	"testing"
)

func TestStackTrace(t *testing.T) {
	tests := []struct {
		err  error
		want string
	}{
		{
			New("ooh"),
			"github.com/morrisxyang/errors.TestStackTrace" +
				"\n\t.+errors/stack_test.go:16",
		},
		{
			Wrap(New("ooh"), "ahh"),
			"github.com/morrisxyang/errors.TestStackTrace" +
				"\n\t.+errors/stack_test.go:21", // this is the stack of Wrap, not New
		},
		{
			Cause(Wrap(New("ooh"), "ahh")),
			"github.com/morrisxyang/errors.TestStackTrace" +
				"\n\t.+errors/stack_test.go:26", // this is the stack of New
		},
		{
			func() error { return New("ooh") }(),
			`github.com/morrisxyang/errors.TestStackTrace.func1` +
				"\n\t.+errors/stack_test.go:31" + "\n" + // this is the stack of New
				"github.com/morrisxyang/errors.TestStackTrace" +
				"\n\t.+errors/stack_test.go:31", // this is the stack of New's caller
		},
		{
			Cause(func() error {
				return func() error {
					return Errorf("hello %s", fmt.Sprintf("world: %s", "ooh"))
				}()
			}()),
			`github.com/morrisxyang/errors.TestStackTrace.func2.1` +
				"\n\t.+errors/stack_test.go:40" + "\n" + // this is the stack of Errorf
				`github.com/morrisxyang/errors.TestStackTrace.func2` +
				"\n\t.+errors/stack_test.go:41" + "\n" + // this is the stack of Errorf's caller
				"github.com/morrisxyang/errors.TestStackTrace" +
				"\n\t.+errors/stack_test.go:42", // this is the stack of Errorf's caller's caller
		},
	}
	for i, tt := range tests {
		x, ok := tt.err.(interface {
			StackTrace() StackTrace
		})
		if !ok {
			t.Errorf("expected %#v to implement StackTrace() StackTrace", tt.err)
			continue
		}
		st := x.StackTrace()
		testFormatRegexp(t, i, st, "%+v", tt.want)

	}
}

var initCallers = func() *StackTrace {
	// maxDepth 记录的栈深度
	const maxDepth = 64
	var pcs [maxDepth]uintptr
	n := runtime.Callers(2, pcs[:])

	var stack *runtime.Frames
	cfg := GetCfg()
	if cfg.StackDepth > 0 && cfg.StackDepth < n {
		stack = runtime.CallersFrames(pcs[0:cfg.StackDepth])
	} else {
		// 转换为 errors.StackTrace
		stack = runtime.CallersFrames(pcs[0:n])
	}

	return &StackTrace{*stack}
}()

func TestStackTraceFormat(t *testing.T) {
	tests := []struct {
		*StackTrace
		format string
		want   string
	}{
		{
			&StackTrace{},
			"%s",
			`\[\]`,
		},
		{
			&StackTrace{},
			"%q",
			`unsupported format: %!q, use %s: \[\]`,
		},
		{
			&StackTrace{},
			"%v",
			`\[\]`,
		},
		{
			&StackTrace{},
			"%+v",
			"",
		},
		{
			&StackTrace{},
			"%#v",
			`\[\]`,
		},
		{
			initCallers,
			"%s",
			`\[github.com/morrisxyang/errors.init runtime.doInit runtime.doInit runtime.main runtime.goexit\]`,
		},
		{
			initCallers,
			"%q",
			`unsupported format: %!q, use %s: \[github.com/morrisxyang/errors.init runtime.doInit runtime.doInit ` +
				`runtime.main runtime.goexit\]`,
		},
		{
			initCallers,
			"%v",
			`\[github.com/morrisxyang/errors.init runtime.doInit runtime.doInit runtime.main runtime.goexit\]`,
		},
		{
			initCallers,
			"%+v",
			`\ngithub.com/morrisxyang/errors.init\n\t.+errors/stack_test.go:81\nruntime.doInit\n\t.*`,
		},
		{
			initCallers,
			"%#v",
			`\[github.com/morrisxyang/errors.init runtime.doInit runtime.doInit runtime.main runtime.goexit\]`,
		},
		{
			(*StackTrace)(nil),
			"%v",
			`<nil>`,
		},
	}

	for i, tt := range tests {
		testFormatRegexp(t, i, tt.StackTrace, tt.format, tt.want)
	}
}

func testFormatRegexp(t *testing.T, n int, arg interface{}, format, want string) {
	t.Helper()
	got := fmt.Sprintf(format, arg)

	match, err := regexp.MatchString(want, got)
	if err != nil {
		t.Fatal(err)
	}
	if !match {
		t.Errorf("failed test %d: fmt.Sprintf(%q, err):\n got: %q\nwant: %q", n+1, format, got, want)
	}
	t.Logf("success test %d: fmt.Sprintf(%q, err):\n got: %q\nwant: %q", n+1, format, got, want)
}
