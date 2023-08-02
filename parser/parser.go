package parser

import (
	"go/ast"

	"github.com/cockroachdb/errors"
	"github.com/samber/mo"
	"github.com/zonewave/copyer/common"
	"github.com/zonewave/copyer/xast"
	"github.com/zonewave/copyer/xtemplate"
	"github.com/zonewave/copyer/xutil"
	"github.com/zonewave/copyer/xutil/xmo"
	"golang.org/x/tools/go/packages"
)

type ParseTemplateParamArg struct {
	Action      common.ActionType
	FileName    string
	SrcName     string
	SrcPkg      string
	SrcTypeName string
	DstName     string
	DstPkg      string
	DstTypeName string
	Pkg         *packages.Package
	FuncLine    int
}
type ParseTemplateParamResult struct {
	TmplParam  *xtemplate.FileParam
	ImportLine int
	FuncLine   int
	Action     common.ActionType
}

func NewParseTemplateParamResult(tmplParam *xtemplate.FileParam, importLine int, arg *ParseTemplateParamArg) *ParseTemplateParamResult {
	return &ParseTemplateParamResult{TmplParam: tmplParam, ImportLine: importLine, FuncLine: arg.FuncLine, Action: arg.Action}
}

// ParseTemplateParam parses templates param
func ParseTemplateParam(arg *ParseTemplateParamArg) mo.Result[*ParseTemplateParamResult] {
	file := xast.FindAstFile(arg.Pkg, arg.FileName)
	switch arg.Action {
	case common.Local:
		return xmo.FlatMap2(file, mo.Ok(arg), ParseLocalTemplateParam)
	case common.Outfile:
		// TODO
	}
	return mo.Err[*ParseTemplateParamResult](errors.New("unknown action"))
}

func ParseLocalTemplateParam(file *ast.File, arg *ParseTemplateParamArg) mo.Result[*ParseTemplateParamResult] {
	varDataSpec := xast.VarSpecParseTry(arg.Pkg, file,
		xast.NewFindVarDataSpecPair(arg.SrcName, arg.SrcPkg, arg.SrcTypeName),
		xast.NewFindVarDataSpecPair(arg.DstName, arg.DstPkg, arg.DstTypeName)).
		Map(xutil.ChecksItems[string, *xast.VarDataSpec](arg.SrcName, arg.DstName))

	return xmo.Map(
		varDataSpec,
		func(varDataSpec map[string]*xast.VarDataSpec) *ParseTemplateParamResult {
			src := varDataSpec[arg.SrcName]
			dst := varDataSpec[arg.DstName]
			varParam := xtemplate.NewTemplateParam(parseVar(arg.Pkg, src), parseVar(arg.Pkg, dst))
			fileParam := xtemplate.NewFileParam(varParam)
			// todo: parse external add import pkg
			return NewParseTemplateParamResult(fileParam, -1, arg)
		})

}
