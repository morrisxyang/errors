package errors

import (
	"errors"
	"fmt"
	"io"
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
		Depth:               0,
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
		Depth:               0,
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
		Depth:               0,
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
