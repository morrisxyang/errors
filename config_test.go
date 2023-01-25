package errors

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestConfig_errorFilterCfg(t *testing.T) {
	convey.Convey("test", t, func() {
		convey.Convey("config default false", func() {
			var cfg Config
			convey.So(cfg.NoStack, convey.ShouldBeFalse)
			convey.So(cfg.NoDetail, convey.ShouldBeFalse)
		})
	})
}

func TestSetCfg(t *testing.T) {
	convey.Convey("test", t, func() {
		SetCfg(&Config{
			NoStack:  false,
			NoDetail: true,
		})
		convey.So(defaultCfg.NoDetail, convey.ShouldBeTrue)
	})
}
