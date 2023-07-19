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
}

// ParseTemplateParam parses templates param
func ParseTemplateParam(arg *ParseTemplateParamArg) mo.Result[*xtemplate.CopyParam] {
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

	return xutil.Map(varDataSpec, func(varDataSpec map[string]*xast.VarDataSpec) *xtemplate.CopyParam {
		src := varDataSpec[arg.SrcName]
		dst := varDataSpec[arg.DstName]
		return xtemplate.NewTemplateParam(parseVar(arg.Pkg, src), parseVar(arg.Pkg, dst))
	})

}
