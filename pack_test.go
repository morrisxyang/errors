package errors

import (
	"errors"
	"fmt"
	"io"
	"math"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		err  string
		want error
	}{
		{"", fmt.Errorf("")},
		{"foo", fmt.Errorf("foo")},
		{"foo", New("foo")},
		{"foo", Newf("%s", "foo")},
		{"foo", Errorf("foo")},
		{"string with format specifiers: %v", errors.New("string with format specifiers: %v")},
	}

	for _, tt := range tests {
		got := New(tt.err)
		if got.Error() != tt.want.Error() {
			t.Errorf("New.Error(): got: %q, want %q", got, tt.want)
		}
	}
}

func TestWrapNil(t *testing.T) {
	tests := []struct {
		err  error
		want error
	}{
		{Wrap(nil, "no error"), nil},
		{Wrapf(nil, "no %s", "error"), nil},
	}
	for _, tt := range tests {
		got := tt.err
		if got != tt.want {
			t.Errorf("Wrap nil: got: %q, want %q", got, tt.want)
		}
	}
}

func TestWrap(t *testing.T) {
	tests := []struct {
		err     error
		message string
		want    string
	}{
		{io.EOF, "read error", "read error: EOF"},
		{Wrap(io.EOF, "read error"), "client error", "client error: read error: EOF"},
	}
	SetCfg(&Config{
		StackDepth:          0,
		ErrorConnectionFlag: ": ",
	})
	defer ResetCfg()
	for _, tt := range tests {
		got := Wrap(tt.err, tt.message).Error()
		if got != tt.want {
			t.Errorf("Wrap(%v, %q): got: %v, want %v", tt.err, tt.message, got, tt.want)
		}
	}
}

func TestWrapf(t *testing.T) {
	tests := []struct {
		err     error
		message string
		want    string
	}{
		{io.EOF, "read error", "read error: EOF"},
		{Wrapf(io.EOF, "read error without format specifiers"), "client error",
			"client error: read error without format specifiers: EOF"},
		{Wrapf(io.EOF, "read error with %d format specifier", 1), "client error",
			"client error: read error with 1 format specifier: EOF"},
	}
	SetCfg(&Config{
		StackDepth:          0,
		ErrorConnectionFlag: ": ",
	})
	defer ResetCfg()
	for _, tt := range tests {
		got := Wrapf(tt.err, tt.message).Error()
		if got != tt.want {
			t.Errorf("Wrapf(%v, %q): got: %v, want %v", tt.err, tt.message, got, tt.want)
		}
	}
}

func TestErrorf(t *testing.T) {
	tests := []struct {
		err  error
		want string
	}{
		{Errorf("read error without format specifiers"), "read error without format specifiers"},
		{Errorf("read error with %d format specifier", 1), "read error with 1 format specifier"},
	}
	SetCfg(&Config{
		StackDepth:          0,
		ErrorConnectionFlag: ": ",
	})
	defer ResetCfg()
	for _, tt := range tests {
		got := tt.err.Error()
		if got != tt.want {
			t.Errorf("Errorf(%v): got: %q, want %q", tt.err, got, tt.want)
		}
	}
}

type nilError struct{}

func (nilError) Error() string { return "nil error" }

func TestCause(t *testing.T) {
	x := New("error")
	tests := []struct {
		err  error
		want error
	}{{
		// nil error is nil
		err:  nil,
		want: nil,
	}, {
		// explicit nil error is nil
		err:  (error)(nil),
		want: nil,
	}, {
		// typed nil is nil
		err:  (*nilError)(nil),
		want: (*nilError)(nil),
	}, {
		// uncaused error is unaffected
		err:  io.EOF,
		want: io.EOF,
	}, {
		// caused error returns cause
		err:  Wrap(io.EOF, "ignored"),
		want: io.EOF,
	}, {
		err:  x, // return from errors.New
		want: x,
	}, {
		Wrap(nil, "whoops"),
		nil,
	}, {
		Wrap(io.EOF, "whoops"),
		io.EOF,
	}, {
		Wrap(nil, ""),
		nil,
	}, {
		Wrap(io.EOF, ""),
		io.EOF,
	}}

	for i, tt := range tests {
		got := Cause(tt.err)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("test %d: got %#v, want %#v", i+1, got, tt.want)
		}
	}
}

// errors.New, etc values are not expected to be compared by value
// but the change in errors#27 made them incomparable. Assert that
// various kinds of errors have a functional equality operator, even
// if the result of that equality is always false.
func TestErrorEquality(t *testing.T) {
	vals := []error{
		nil,
		io.EOF,
		errors.New("EOF"),
		New("EOF"),
		Errorf("EOF"),
		Wrap(io.EOF, "EOF"),
		Wrapf(io.EOF, "EOF%d", 2),
		Wrap(nil, "whoops"),
		Wrap(io.EOF, "whoops"),
		Wrap(io.EOF, ""),
		Wrap(nil, ""),
	}

	for i := range vals {
		for j := range vals {
			_ = vals[i] == vals[j] // mustn't panic
		}
	}
}

func TestCode(t *testing.T) {
	// Test when error is nil
	code := Code(nil)
	if code != 0 {
		t.Errorf("Expected Code(nil) to return 0, but got %d", code)
	}

	// Test when error is not of type *baseError
	err := errors.New("some error")
	code = Code(err)
	if code != math.MinInt32 {
		t.Errorf("Expected Code(%v) to return %d, but got %d", err, math.MinInt32, code)
	}

	// Test when error is of type *baseError
	baseErr := &baseError{code: 123, msg: "some error"}
	code = Code(baseErr)
	if code != baseErr.code {
		t.Errorf("Expected Code(%v) to return %d, but got %d", baseErr, baseErr.code, code)
	}
}

func TestMsg(t *testing.T) {
	// Test when error is nil
	msg := Msg(nil)
	if msg != "" {
		t.Errorf("Expected Msg(nil) to return empty string, but got %s", msg)
	}

	// Test when error is not of type *baseError
	err := errors.New("some error")
	msg = Msg(err)
	if msg != "unknown error" {
		t.Errorf("Expected Msg(%v) to return 'unknown error', but got %s", err, msg)
	}

	// Test when error is of type *baseError
	baseErr := &baseError{code: 123, msg: "some error"}
	msg = Msg(baseErr)
	if msg != baseErr.msg {
		t.Errorf("Expected Msg(%v) to return '%s', but got '%s'", baseErr, baseErr.msg, msg)
	}
}
