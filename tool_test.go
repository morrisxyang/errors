// Package errors ...
package errors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDepth(t *testing.T) {
	SetCfg(&Config{
		StackDepth:          1,
		ErrorConnectionFlag: ": ",
	})

	fmt.Printf("%+v\n", New("callers"))
	assert.Regexp(t, "^callers\ngithub.com/morrisxyang/errors.TestDepth\n\t.*tool_test.go:19$",
		fmt.Sprintf("%+v", New("callers")))
}
