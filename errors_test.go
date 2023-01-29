package errors

import (
	"fmt"
	"os"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestErrorPrint(t *testing.T) {
	convey.Convey("test", t, func() {
		convey.Convey("wrap err", func() {
			err := a()
			// %v
			fmt.Printf("test %%v wrap err: %v\n", err)
			s1 := fmt.Sprintf("err: %v\n", err)
			convey.So(s1, convey.ShouldContainSubstring,
				"no such file or directory")
			convey.So(s1, convey.ShouldContainSubstring,
				"Cause by c")
			convey.So(s1, convey.ShouldContainSubstring,
				"c failed reason")
			convey.So(s1, convey.ShouldContainSubstring,
				"Cause by b")
			convey.So(s1, convey.ShouldContainSubstring,
				"b failed reason")
			convey.So(s1, convey.ShouldContainSubstring,
				"a failed reason ")

			//  %+v
			fmt.Printf("test %%+v wrap err: %+v\n", err)
			s2 := fmt.Sprintf("err: %+v\n", err)
			convey.So(s2, convey.ShouldContainSubstring,
				"no such file or directory")
			convey.So(s2, convey.ShouldContainSubstring,
				"github.com/morrisxyang/errors.c")
			convey.So(s2, convey.ShouldContainSubstring,
				"github.com/morrisxyang/errors.b")
			convey.So(s2, convey.ShouldContainSubstring,
				"github.com/morrisxyang/errors.a")
			convey.So(s2, convey.ShouldContainSubstring,
				"c failed reason")
			convey.So(s2, convey.ShouldContainSubstring,
				"b failed reason")
			convey.So(s2, convey.ShouldContainSubstring,
				"a failed reason")
		})
		convey.Convey("new err", func() {
			err := a1()
			// %v
			fmt.Printf("test %%v new err: %v\n", err)
			s1 := fmt.Sprintf("err: %v\n", err)
			convey.So(s1, convey.ShouldContainSubstring,
				"Cause by c1")
			convey.So(s1, convey.ShouldContainSubstring,
				"c1 failed reason")
			convey.So(s1, convey.ShouldContainSubstring,
				"Cause by b1")
			convey.So(s1, convey.ShouldContainSubstring,
				"b1 failed reason")
			convey.So(s1, convey.ShouldContainSubstring,
				"a1 failed reason")

			// %+v
			fmt.Printf("test %%+v new err: %+v\n", err)
			s2 := fmt.Sprintf("err: %+v\n", err)
			convey.So(s2, convey.ShouldContainSubstring,
				"github.com/morrisxyang/errors.c1")
			convey.So(s2, convey.ShouldContainSubstring,
				"github.com/morrisxyang/errors.b1")
			convey.So(s2, convey.ShouldContainSubstring,
				"github.com/morrisxyang/errors.a1")
			convey.So(s2, convey.ShouldContainSubstring,
				"c1 failed reason")
			convey.So(s2, convey.ShouldContainSubstring,
				"b1 failed reason")
			convey.So(s2, convey.ShouldContainSubstring,
				"a1 failed reason")
		})

	})

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
