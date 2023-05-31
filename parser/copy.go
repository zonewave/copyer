package parser

import (
	"go/ast"

	"github.com/zonewave/copyer/xast"
	"github.com/zonewave/copyer/xtemplate"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"golang.org/x/tools/go/packages"
)

func parseVar(localPkg *packages.Package, v *xast.VarDataSpec) *xtemplate.TmplVar {
	tVar := &xtemplate.TmplVar{
		Name:          v.Name,
		TypeNameNoDot: cases.Title(language.Und).String(v.TypePackageName) + v.TypeSpec.Name,
		Exported:      isExported(localPkg, v.AstSpec.Pkg, v.Name),
	}
	if v.TypePackageName == "" {
		tVar.Type = v.TypeSpec.Name
	} else {
		tVar.Type = v.TypePackageName + "." + v.TypeSpec.Name
	}

	tVar.StructType = parseTypeSecStruct(localPkg, v.TypeSpec)
	return tVar
}

func isExported(localPkg *packages.Package, pkg *packages.Package, name string) bool {
	if pkg.ID == localPkg.ID {
		return true
	}

	return ast.IsExported(name)

}
func parseTypeSecStruct(localPkg *packages.Package, spec *xast.TypeSpec) *xtemplate.TmplStruct {
	structType := spec.AstSpec.TypeSpec.Type.(*ast.StructType)
	s := &xtemplate.TmplStruct{
		Name:   spec.Name,
		Fields: make(map[string]*xtemplate.TmplVar, structType.Fields.NumFields()),
	}
	for _, g := range structType.Fields.List {
		if len(g.Names) == 0 {
			continue
		}
		name := typeExprToName(g.Type)

		for _, item := range g.Names {
			s.Fields[item.Name] = &xtemplate.TmplVar{
				Name:          item.Name,
				TypeNameNoDot: name,
				Type:          name,
				Exported:      isExported(localPkg, spec.AstSpec.Pkg, item.Name),
			}
		}
	}

	return s
}
func typeExprToName(typeExpr ast.Expr) string {
	switch typeExpr.(type) {
	case *ast.Ident:
		return typeExpr.(*ast.Ident).Name
	case *ast.SelectorExpr:
		return typeExpr.(*ast.SelectorExpr).Sel.Name
	case *ast.StarExpr:
		return "*" + typeExprToName(typeExpr.(*ast.StarExpr).X)
	case *ast.ArrayType:
		return "[]" + typeExprToName(typeExpr.(*ast.ArrayType).Elt)
	case *ast.MapType:
		mapType := typeExpr.(*ast.MapType)
		return "map[" + typeExprToName(mapType.Key) + "]" + typeExprToName(mapType.Value)
	}
	return ""
}
