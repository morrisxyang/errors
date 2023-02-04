package errors

import (
	"fmt"
	"io"
	"runtime"
)

// StackTrace is stack of Frames from innermost (newest) to outermost (oldest).
type StackTrace struct {
	*runtime.Frames
}

// Format formats the stack of Frames according to the fmt.Formatter interface.
//
//	%s	lists source files for each Frame in the stack
//	%v	lists the source file and line number for each Frame in the stack
//
// Format accepts flags that alter the printing of some verbs, as follows:
//
//	%+v   Prints filename, function, and line number for each Frame in the stack.
func (st StackTrace) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		switch {
		case s.Flag('+'):
			io.WriteString(s, "\nStack Info: ")
			for {
				frame, more := st.Frames.Next()
				if !more {
					break
				}
				fmt.Fprintf(s, "\n%s", frame.Function)
				fmt.Fprintf(s, "\n\t%s:%d", frame.File, frame.Line)
			}
		default:
			st.formatSlice(s, verb)
		}
	case 's':
		st.formatSlice(s, verb)
	default:
		_, _ = fmt.Fprintf(s, "unsupported format: %%!%c, use %%s: ", verb)
		st.formatSlice(s, verb)
	}
}

// formatSlice will format this StackTrace into the given buffer as a slice of
// Frame, only valid when called with '%s' or '%v'.
func (st StackTrace) formatSlice(s fmt.State, verb rune) {
	io.WriteString(s, "[")
	for {
		frame, more := st.Frames.Next()
		if !more {
			break
		}
		io.WriteString(s, " ")
		fmt.Fprintf(s, "%s", frame.Function)
	}
	io.WriteString(s, "]")
}
