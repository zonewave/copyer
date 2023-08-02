package generate

import (
	"bytes"
	"go/format"
	"text/template"

	"github.com/samber/mo"
	"github.com/zonewave/copyer/common"
	"github.com/zonewave/copyer/output"
	"github.com/zonewave/copyer/parser"
	"github.com/zonewave/copyer/xast"
	"github.com/zonewave/copyer/xtemplate"
	"github.com/zonewave/copyer/xutil/xmo"
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

	args := xmo.Map2(
		mo.Ok(arg), xast.LoadLocalPkg(arg.GoPkg, arg.LoadConfigOpts...),
		// set arg
		NewParseTemplateParamArg,
	)

	// parser templates param
	parseResult := xmo.FlatMap(
		args,
		parser.ParseTemplateParam,
	)

	// load templates
	tmpl := xtemplate.NewTmpl()

	return xmo.Map2(
		tmpl, parseResult,
		newGenerate,
	)
}
func newGenerate(tmpls *xtemplate.Templates, parseResult *parser.ParseTemplateParamResult) *Generator {
	return &Generator{ParseResult: parseResult, Tmpls: tmpls}
}

// OutPut is a wrapper of io.Writer
func OutPut(data []*output.LinesData, out output.Writer) mo.Result[bool] {
	return mo.TupleToResult(true, out.LineDataBatchInsert(data...))
}

func ProduceCode(g *Generator) mo.Result[[]*output.LinesData] {
	switch g.ParseResult.Action {
	case common.Local:
		var ret []mo.Result[*output.LinesData]
		if g.ParseResult.TmplParam.Imports != nil {
			ret = append(ret, g.TemplateImportExecute())
		}
		ret = append(ret, g.TemplateFuncExecute())
		return xmo.ToSlice(ret...)
	}
	return mo.Ok([]*output.LinesData{})
}

func (g *Generator) TemplateFuncExecute() mo.Result[*output.LinesData] {
	return xmo.Map(
		templateExecuteAndFormat(g.Tmpls.CopyTemplate, g.ParseResult.TmplParam.CopyFunc),
		func(bs []byte) *output.LinesData {
			return &output.LinesData{
				StartLine: g.ParseResult.FuncLine,
				Bytes:     bs,
			}
		})
}

func (g *Generator) TemplateImportExecute() mo.Result[*output.LinesData] {
	return xmo.Map(
		templateExecuteAndFormat(g.Tmpls.ImportTemplate, g.ParseResult.TmplParam.Imports),
		func(bs []byte) *output.LinesData {
			return &output.LinesData{
				StartLine: g.ParseResult.ImportLine,
				Bytes:     bs,
			}
		})
}
func templateExecuteAndFormat(tmpl *template.Template, tmplParam any) mo.Result[[]byte] {
	return xmo.FlatMap(
		templateExecute(tmpl, tmplParam),
		formatCode,
	)
}
func templateExecute(tmpl *template.Template, tmplParam any) mo.Result[*bytes.Buffer] {
	buffer := new(bytes.Buffer)
	err := tmpl.Execute(buffer, tmplParam)
	return mo.TupleToResult(buffer, err)
}

func formatCode(value *bytes.Buffer) mo.Result[[]byte] {
	return mo.TupleToResult(format.Source(value.Bytes())).MapErr(xmo.MapWrapf[[]byte]("format code failed,src:%s", string(value.Bytes())))
}
