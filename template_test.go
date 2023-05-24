package main

import (
	"testing"

	c "github.com/smartystreets/goconvey/convey"
)

func Test_newCopyTemplate(t *testing.T) {
	c.Convey("Test_newCopyTemplate", t, func() {
		tmpl, err := newCopyTemplate()
		c.So(err, c.ShouldBeNil)
		c.So(tmpl.ParseName, c.ShouldEqual, "copy")
		c.So(tmpl.Tree.Name, c.ShouldEqual, "copy")
	})
}
