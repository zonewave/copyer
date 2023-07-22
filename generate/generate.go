package generate

import (
	"bytes"
	"go/format"
	"text/template"

	"github.com/cockroachdb/errors"
	"github.com/samber/mo"
	"github.com/zonewave/copyer/output"
	"github.com/zonewave/copyer/parser"
	"github.com/zonewave/copyer/xast"
	"github.com/zonewave/copyer/xtemplate"
	"github.com/zonewave/copyer/xutil"
	"golang.org/x/tools/go/packages"
)

// Generator is a code generator
type Generator struct {
	Param *xtemplate.CopyParam
	Tmpl  *template.Template
}

func NewParseTemplateParamArg(arg *GeneratorArg, pkg *packages.Package) *parser.ParseTemplateParamArg {
	return &parser.ParseTemplateParamArg{
		Action:      arg.Action,
		FileName:    arg.GoFile,
		SrcName:     arg.SrcName,
		SrcPkg:      arg.SrcPkg,
		SrcTypeName: arg.SrcType,
		DstName:     arg.DstName,
		DstPkg:      arg.DstPkg,
		DstTypeName: arg.DstType,
		Pkg:         pkg,
	}
}

// NewGenerator creates a new generator
func NewGenerator(arg *GeneratorArg) mo.Result[*Generator] {

	args := xutil.Map2(
		mo.Ok(arg),
		// load ast.pkgs
		xast.LoadLocalPkg(arg.GoPkg, arg.LoadConfigOpts...),
		// set arg
		NewParseTemplateParamArg,
	)

	// parser templates param
	tmplParam := xutil.FlatMap(args, parser.ParseTemplateParam)

	// load templates
	tmpl := xtemplate.NewTmpl(arg.Action)

	return xutil.Map2(tmpl, tmplParam, newGenerate)
}
func newGenerate(tmpl *template.Template, param *xtemplate.CopyParam) *Generator {
	return &Generator{Param: param, Tmpl: tmpl}
}

// OutPut is a wrapper of io.Writer
func OutPut(fileLine int, bs []byte, out output.Writer) mo.Result[bool] {
	return mo.TupleToResult(true, out.LineDataBatchInsert(&output.LinesData{
		StartLine: fileLine,
		Bytes:     bs,
	}))
}

func ProduceCode(g *Generator) mo.Result[[]byte] {
	return xutil.FlatMap(g.BufferExecute(new(bytes.Buffer)), formatCode).MapErr(xutil.MapWrapf[[]byte]("format code %s", g.Tmpl.Name()))
}

func (g *Generator) BufferExecute(value *bytes.Buffer) mo.Result[*bytes.Buffer] {
	err := g.Tmpl.Execute(value, g.Param)
	return mo.TupleToResult(value, errors.Wrapf(err, "execute templates %s", g.Tmpl.Name()))
}

func formatCode(value *bytes.Buffer) mo.Result[[]byte] {
	return mo.TupleToResult(format.Source(value.Bytes()))
}
