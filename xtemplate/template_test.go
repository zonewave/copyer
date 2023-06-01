package xtemplate

import (
	"bytes"
	"go/format"
	"testing"

	c "github.com/smartystreets/goconvey/convey"
)

func Test_CopyTemplate(t *testing.T) {
	c.Convey("Test_CopyTemplate", t, func() {
		tmpl, err := NewCopyTemplate()
		c.So(err, c.ShouldBeNil)

		var buf bytes.Buffer
		tmplStructS := &TmplStruct{
			Name: "A",
			Fields: map[string]*TmplVar{
				"Name": {
					Name:          "Name",
					TypeNameNoDot: "int",
					Type:          "int",
					Exported:      true,
				},
				"Age": {
					Name:          "Age",
					TypeNameNoDot: "int",
					Type:          "int",
					Exported:      false,
				},
			},
		}
		tmplStructD := &TmplStruct{}
		CopyTmplStructToTmplStruct(tmplStructS, tmplStructD)
		err = tmpl.Execute(&buf, &CopyParam{
			Src: &TmplVar{
				Name:          "src",
				TypeNameNoDot: "A",
				Type:          "A",
				StructType:    tmplStructS,
				Exported:      true,
			},
			Dst: &TmplVar{
				Name:          "dst",
				TypeNameNoDot: "A",
				Type:          "A",
				StructType:    tmplStructD,
				Exported:      false,
			},
		})
		bs, err := format.Source(buf.Bytes())
		c.So(err, c.ShouldBeNil)
		expect :=
			`func CopyAToA(src *A, dst *A) {
	dst.Name = src.Name
}`
		c.So(string(bs), c.ShouldEqual, expect)
	})
}
