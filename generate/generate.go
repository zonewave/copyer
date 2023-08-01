package generate

import (
	"bytes"
	"go/format"

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
	ParseResult *parser.ParseTemplateParamResult
	Tmpls       *xtemplate.Templates
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
		FuncLine:    arg.OutLine,
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
	parseResult := xutil.FlatMap(args, parser.ParseTemplateParam)

	// load templates
	tmpl := xtemplate.NewTmpl()

	return xutil.Map2(tmpl, parseResult, newGenerate)
}
func newGenerate(tmpls *xtemplate.Templates, parseResult *parser.ParseTemplateParamResult) *Generator {
	return &Generator{ParseResult: parseResult, Tmpls: tmpls}
}

// OutPut is a wrapper of io.Writer
func OutPut(fileLine int, bs []byte, out output.Writer) mo.Result[bool] {
	return mo.TupleToResult(true, out.LineDataBatchInsert(&output.LinesData{
		StartLine: fileLine,
		Bytes:     bs,
	}))
}

func ProduceCode(g *Generator) mo.Result[[]byte] {
	return xutil.FlatMap(g.BufferExecute(new(bytes.Buffer)), formatCode).
		MapErr(xutil.MapWrap[[]byte]("format code failed"))
}

func (g *Generator) BufferExecute(value *bytes.Buffer) mo.Result[*bytes.Buffer] {
	err := g.Tmpls.CopyTemplate.Execute(value, g.ParseResult.TmplParam.CopyFunc)
	return mo.TupleToResult(value, err)
}

func formatCode(value *bytes.Buffer) mo.Result[[]byte] {
	return mo.TupleToResult(format.Source(value.Bytes()))
}
