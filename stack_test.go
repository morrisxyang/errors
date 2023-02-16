package errors

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
)

// var initpc = caller()

// a version of runtime.Caller that returns a Frame, not a uintptr.
// func caller() Frame {
// 	var pcs [3]uintptr
// 	n := runtime.Callers(2, pcs[:])
// 	frames := runtime.CallersFrames(pcs[:n])
// 	frame, _ := frames.Next()
// 	return Frame(frame.PC)
// }

func TestStackTrace(t *testing.T) {
	tests := []struct {
		err  error
		want string
	}{
		{
			New("ooh"),
			"github.com/morrisxyang/errors.TestStackTrace" +
				"\n\t.+errors/stack_test.go:27",
		},
		{
			Wrap(New("ooh"), "ahh"),
			"github.com/morrisxyang/errors.TestStackTrace" +
				"\n\t.+errors/stack_test.go:32", // this is the stack of Wrap, not New
		},
		{
			Cause(Wrap(New("ooh"), "ahh")),
			"github.com/morrisxyang/errors.TestStackTrace" +
				"\n\t.+errors/stack_test.go:37", // this is the stack of New
		},
		{
			func() error { return New("ooh") }(),
			`github.com/morrisxyang/errors.TestStackTrace.func1` +
				"\n\t.+errors/stack_test.go:42" + "\n" + // this is the stack of New
				"github.com/morrisxyang/errors.TestStackTrace" +
				"\n\t.+errors/stack_test.go:42", // this is the stack of New's caller
		},
		{
			Cause(func() error {
				return func() error {
					return Errorf("hello %s", fmt.Sprintf("world: %s", "ooh"))
				}()
			}()),
			`github.com/morrisxyang/errors.TestStackTrace.func2.1` +
				"\n\t.+errors/stack_test.go:51" + "\n" + // this is the stack of Errorf
				`github.com/morrisxyang/errors.TestStackTrace.func2` +
				"\n\t.+errors/stack_test.go:52" + "\n" + // this is the stack of Errorf's caller
				"github.com/morrisxyang/errors.TestStackTrace" +
				"\n\t.+errors/stack_test.go:53", // this is the stack of Errorf's caller's caller
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

func testFormatRegexp(t *testing.T, n int, arg interface{}, format, want string) {
	t.Helper()
	got := fmt.Sprintf(format, arg)
	gotLines := strings.SplitN(got, "\n", -1)
	gotLines = gotLines[2:]
	wantLines := strings.SplitN(want, "\n", -1)

	if len(wantLines) > len(gotLines) {
		t.Errorf("test %d: wantLines(%d) > gotLines(%d):\n got: %q\nwant: %q", n+1, len(wantLines), len(gotLines), got, want)
		return
	}

	for i, w := range wantLines {
		match, err := regexp.MatchString(w, gotLines[i])
		if err != nil {
			t.Fatal(err)
		}
		if !match {
			t.Errorf("test %d: line %d: fmt.Sprintf(%q, err):\n got: %q\nwant: %q", n+1, i+1, format, got, want)
		}
	}
}
