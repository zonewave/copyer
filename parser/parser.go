package parser

import (
	"github.com/cockroachdb/errors"
	"github.com/zonewave/copyer/xast"
	"github.com/zonewave/copyer/xtemplate"
	"golang.org/x/tools/go/packages"
)

// ParseTemplateParam parses templates param
func ParseTemplateParam(fileName, srcPkg, srcTypeName, dstPkg, dstTypeName string, pkg *packages.Package) (*xtemplate.CopyParam, error) {
	file, err := xast.FindAstFile(pkg, fileName)
	if err != nil {
		return nil, err
	}
	varSpecs, err := xast.VarSpecLocalParseMust(pkg, file,
		xast.NewFindVarDataSpecPair("src", srcPkg, srcTypeName),
		xast.NewFindVarDataSpecPair("dst", dstPkg, dstTypeName))
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
	return xtemplate.NewTemplateParam(parseVar(pkg, src), parseVar(pkg, dst)), nil
}

// parseTypeSecStructs parses map[string]ast.TypeSpecs to map[string]TemplateStructs
