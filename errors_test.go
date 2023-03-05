package errors

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorPrint(t *testing.T) {
	ResetCfg()
	err := a()
	assert.Equal(t, fmt.Sprintf("%s", err), fmt.Sprintf("%v", err))
	fmt.Printf("%s\n", err)

	s1 := fmt.Sprintf("err: %v\n", err)
	assert.Contains(t, s1, "no such file or directory")
	assert.Regexp(t, GetCfg().ErrorConnectionFlag+".*c failed reason", s1)
	assert.Regexp(t, GetCfg().ErrorConnectionFlag+".*b failed reason", s1)
	assert.Contains(t, s1, "c failed reason")
	assert.Contains(t, s1, "b failed reason")
	assert.Contains(t, s1, "a failed reason")

	SetCfg(&Config{
		StackDepth:          100,
		ErrorConnectionFlag: ": ",
	})
	assert.Contains(t, s1, "no such file or directory")
	assert.Regexp(t, GetCfg().ErrorConnectionFlag+".*c failed reason", s1)
	assert.Regexp(t, GetCfg().ErrorConnectionFlag+".*b failed reason", s1)
	assert.Contains(t, s1, "c failed reason")
	assert.Contains(t, s1, "b failed reason")
	assert.Contains(t, s1, "a failed reason")

	ResetCfg()
	//  %+v print stack
	fmt.Printf("test %%+v wrap err: %+v\n", err)
	s2 := fmt.Sprintf("err: %+v\n", err)
	assert.Contains(t, s2, "no such file or directory")
	assert.Contains(t, s2, "github.com/morrisxyang/errors.c")
	assert.Contains(t, s2, "github.com/morrisxyang/errors.b")
	assert.Contains(t, s2, "github.com/morrisxyang/errors.a")
	assert.Contains(t, s2, "github.com/morrisxyang/errors.TestErrorPrint")
	assert.Contains(t, s2, "c failed reason")
	assert.Contains(t, s2, "b failed reason")
	assert.Contains(t, s2, "a failed reason")
}

func TestErrorPrintWithCode(t *testing.T) {
	ResetCfg()
	err := a1()
	// %v
	fmt.Printf("test %%v new err: %v\n", err)
	s1 := fmt.Sprintf("err: %v\n", err)
	assert.Regexp(t, GetCfg().ErrorConnectionFlag+".*123, c1 failed reason", s1)
	assert.Regexp(t, GetCfg().ErrorConnectionFlag+".*456, b1 failed reason", s1)
	assert.Contains(t, s1, "a1 failed reason")

	// %+v
	fmt.Printf("test %%+v new err: %+v\n", err)
	s2 := fmt.Sprintf("err: %+v\n", err)
	assert.Contains(t, s2, "github.com/morrisxyang/errors.c1")
	assert.Contains(t, s2, "github.com/morrisxyang/errors.b1")
	assert.Contains(t, s2, "github.com/morrisxyang/errors.a1")
	assert.Contains(t, s2, "c1 failed reason")
	assert.Contains(t, s2, "b1 failed reason")
	assert.Contains(t, s2, "a1 failed reason")

	err = a2()
	// %v
	fmt.Printf("test %%v new err: %v\n", err)
	s3 := fmt.Sprintf("err: %v\n", err)
	assert.Regexp(t, GetCfg().ErrorConnectionFlag+".*123, c2 failed reason", s3)
	assert.Regexp(t, GetCfg().ErrorConnectionFlag+".*b2 failed reason", s3)
	assert.Contains(t, s3, "789, a2 failed reason")
}

func a() error {
	err := b()
	err = Wrap(err, "a failed reason")
	return err
}

func b() error {
	err := c()
	err = Wrap(err, "b failed reason")
	return err
}

func c() error {
	_, err := os.Open("test")
	if err != nil {
		return WrapWithCode(err, 123, "c failed reason")
	}
	return nil
}

func a1() error {
	err := b1()
	err = Wrap(err, "a1 failed reason")
	return err
}

func b1() error {
	err := c1()
	err = WrapWithCode(err, 456, "b1 failed reason")
	return err
}
func c1() error {
	return NewWithCode(123, "c1 failed reason")
}

func a2() error {
	err := b2()
	err = WrapWithCode(err, 789, "a2 failed reason")
	return err
}

func b2() error {
	err := c2()
	err = Wrap(err, "b2 failed reason")
	return err
}
func c2() error {
	return NewWithCode(123, "c2 failed reason")
}
