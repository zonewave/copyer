package parser

import (
	"go/ast"

	"github.com/samber/mo"
	"github.com/zonewave/copyer/common"
	"github.com/zonewave/copyer/xast"
	"github.com/zonewave/copyer/xtemplate"
	"github.com/zonewave/copyer/xutil"
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
}

func NewParseTemplateParamResult(tmplParam *xtemplate.FileParam, importLine, funcLine int) *ParseTemplateParamResult {
	return &ParseTemplateParamResult{TmplParam: tmplParam, ImportLine: importLine, FuncLine: funcLine}
}

// ParseTemplateParam parses templates param
func ParseTemplateParam(arg *ParseTemplateParamArg) mo.Result[*ParseTemplateParamResult] {
	file := xast.FindAstFile(arg.Pkg, arg.FileName)

	varDataSpec := xutil.FlatMap(
		file,
		func(file *ast.File) mo.Result[map[string]*xast.VarDataSpec] {
			return xast.VarSpecParseTry(arg.Pkg, file,
				xast.NewFindVarDataSpecPair(arg.SrcName, arg.SrcPkg, arg.SrcTypeName),
				xast.NewFindVarDataSpecPair(arg.DstName, arg.DstPkg, arg.DstTypeName)).
				Map(xutil.ChecksItems[string, *xast.VarDataSpec](arg.SrcName, arg.DstName))
		},
	)

	return xutil.Map(varDataSpec, func(varDataSpec map[string]*xast.VarDataSpec) *ParseTemplateParamResult {
		src := varDataSpec[arg.SrcName]
		dst := varDataSpec[arg.DstName]
		varParam := xtemplate.NewTemplateParam(parseVar(arg.Pkg, src), parseVar(arg.Pkg, dst))
		fileParam := xtemplate.NewFileParam(varParam)
		return NewParseTemplateParamResult(fileParam, -1, arg.FuncLine)
	})

}
