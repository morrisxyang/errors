package errors

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestSetCfg(t *testing.T) {
	convey.Convey("test", t, func() {
		SetCfg(&Config{
			ErrorConnectionFlag: ":",
		})
		convey.So(defaultCfg.ErrorConnectionFlag, convey.ShouldEqual, ":")
	})
}
