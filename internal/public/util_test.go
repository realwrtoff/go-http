package public

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestIsValidMobile(t *testing.T) {
	Convey("TestIsValidMobile", t, func() {
		mobileShort := "1774458194"
		mobile := "17744581949"
		mobileLong := "177445819499"
		mobileStr := "177hell1949"

		Convey("case short", func() {
			So(IsValidMobile(mobileShort), ShouldBeFalse)
		})

		Convey("case ok", func() {
			So(IsValidMobile(mobile), ShouldBeTrue)
		})

		Convey("case long", func() {
			So(IsValidMobile(mobileLong), ShouldBeFalse)
		})

		Convey("case str", func() {
			So(IsValidMobile(mobileStr), ShouldBeFalse)
		})
	})
}
