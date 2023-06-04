package parser

import (
	"github.com/cockroachdb/errors"
	"github.com/zonewave/copyer/common"
	"github.com/zonewave/copyer/xast"
	"github.com/zonewave/copyer/xtemplate"
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
func ParseTemplateParam(arg *ParseTemplateParamArg) (*xtemplate.CopyParam, error) {
	file, err := xast.FindAstFile(arg.Pkg, arg.FileName)
	if err != nil {
		return nil, err
	}
	varSpecs, err := xast.VarSpecLocalParseMust(arg.Pkg, file,
		xast.NewFindVarDataSpecPair(arg.SrcName, arg.SrcPkg, arg.SrcTypeName),
		xast.NewFindVarDataSpecPair(arg.DstName, arg.DstPkg, arg.DstTypeName))
	if err != nil {
		return nil, err
	}

	src, ok := varSpecs["src"]
	if !ok {
		return nil, errors.Errorf("src %s not found", src)
	}
	dst, ok := varSpecs["dst"]
	if !ok {
		return nil, errors.Errorf("dst %s not found", dst)
	}
	return xtemplate.NewTemplateParam(parseVar(arg.Pkg, src), parseVar(arg.Pkg, dst)), nil
}
