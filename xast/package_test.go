package xast

import (
	"strings"
	"testing"

	c "github.com/smartystreets/goconvey/convey"
	"github.com/zonewave/pkgs/standutil/sliceutil"
	"golang.org/x/tools/go/packages"
)

var loadTestCfg = func(config *packages.Config) {
	config.Mode |= packages.NeedFiles
	config.Tests = true
}

type A struct {
	a int
}
type B struct {
	b int
}

func Test_loadPkgs(t *testing.T) {
	c.Convey("Test_loadPkgs", t, func() {
		c.Convey("Test_loadPkgs for example", func() {
			pkgs, err := loadPkgs([]string{"file=example/struct.go"}, loadTestCfg)
			c.So(err, c.ShouldBeNil)
			c.So(pkgs, c.ShouldNotBeEmpty)
			c.So(pkgs[0].Name, c.ShouldEqual, "example")
			c.So(pkgs[0].GoFiles, c.ShouldNotBeEmpty)

			exist := false
			sliceutil.IterFn(pkgs[0].GoFiles, func(i int, s string) bool {
				if strings.Contains(s, "struct.go") {
					exist = true
					return false
				}
				return true
			})
			c.So(exist, c.ShouldBeTrue)

		})
		c.Convey("Test_loadPkgs failed", func() {
			_, err := loadPkgs([]string{"patterss=tsdfsf/erwer.go"}, loadTestCfg)
			c.So(err, c.ShouldBeError)
		})
	})
}

func Test_findTypeSpec(t *testing.T) {
	c.Convey("Test_findTypeSpec", t, func() {
		pkgs, err := loadPkgs([]string{"./..."}, loadTestCfg)
		c.So(err, c.ShouldBeNil)
		c.So(pkgs, c.ShouldNotBeEmpty)

		ret := findTypeSpec(pkgs, []string{"A", "B"})
		c.So(ret, c.ShouldNotBeEmpty)
		c.So(ret["A"].Name.Name, c.ShouldEqual, "A")
		c.So(ret["B"].Name.Name, c.ShouldEqual, "B")
	})
}
