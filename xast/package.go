package xast

import (
	"go/ast"
	"go/token"
	"path/filepath"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/duke-git/lancet/v2/maputil"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/samber/mo"
	"github.com/zonewave/copyer/xutil"
	"github.com/zonewave/copyer/xutil/xslice"
	"golang.org/x/tools/go/packages"
)

type AstSpec struct {
	Pkg      *packages.Package
	File     *ast.File
	TypeSpec *ast.TypeSpec // option
}

type VarDataSpec struct {
	Name string
	//   astSpec file/pkg info where the variable name is declared
	AstSpec         *AstSpec
	TypePackageName string
	TypeSpec        *TypeSpec
}

func NewVarDataSpec(name string, astSpec *AstSpec, typePackageName string, typeSpec *TypeSpec) *VarDataSpec {
	return &VarDataSpec{Name: name, AstSpec: astSpec, TypePackageName: typePackageName, TypeSpec: typeSpec}
}
func NewVarDataSpecFromPair(astSpec *AstSpec, pair *VariableDef) *VarDataSpec {
	return NewVarDataSpec(pair.varName, astSpec, pair.pkgName, nil)
}
func (v *VarDataSpec) IsPublic() bool {

	return ast.IsExported(v.Name)
}

type TypeSpec struct {
	Name    string
	AstSpec *AstSpec
}

type VariableDef struct {
	varName  string
	pkgName  string
	pkgId    int
	typeName string
}

func NewFindVarDataSpecPair(varName string, pkgName string, typeName string) *VariableDef {
	return &VariableDef{varName: varName, pkgName: pkgName, typeName: typeName}
}

func VarSpecParseTry(pkg *packages.Package, file *ast.File, pairs ...*VariableDef) mo.Result[map[string]*VarDataSpec] {

	pairsByPkg := slice.GroupWith(pairs, func(item *VariableDef) string {
		return item.pkgName
	})
	// search import pkgs
	importPkgsNames := slice.Filter(maputil.Keys(pairsByPkg), func(_ int, item string) bool {
		return item != ""
	})
	pkgs := mo.Ok(importedPkgFind(pkg, file, importPkgsNames...)).
		Map(xutil.ChecksItems[string, *packages.Package](importPkgsNames...)).
		Map(func(value map[string]*packages.Package) (map[string]*packages.Package, error) {
			value[""] = pkg
			return value, nil
		})

	return xutil.FlatMap(pkgs, func(pkgs map[string]*packages.Package) mo.Result[map[string]*VarDataSpec] {
		// search typeSpec
		ret := make(map[string]*VarDataSpec)
		astSpec := &AstSpec{
			Pkg:      pkg,
			File:     file,
			TypeSpec: nil,
		}
		for _, pair := range pairs {
			ret[pair.varName] = &VarDataSpec{
				Name:            pair.varName,
				TypePackageName: pair.pkgName,
				AstSpec:         astSpec,
			}
		}
		for pkgName, fPkg := range pkgs {
			gPairs := pairsByPkg[pkgName]
			typeNames := slice.Map(gPairs, func(_ int, item *VariableDef) string {
				return item.typeName
			})
			typeSpec := TypeSpecGet(fPkg, typeNames...).Map(xutil.ChecksItems[string, *TypeSpec](typeNames...))
			if typeSpec.IsError() {
				return mo.Err[map[string]*VarDataSpec](typeSpec.Error())
			}
			for _, pair := range gPairs {
				ret[pair.varName].TypeSpec = typeSpec.MustGet()[pair.typeName]
			}
		}
		return mo.Ok(ret)
	})

}

func TypeSpecGet(pkg *packages.Package, typeNames ...string) mo.Result[map[string]*TypeSpec] {
	ret := make(map[string]*TypeSpec, len(typeNames))

	for _, file := range pkg.Syntax {

		ast.Inspect(file, func(node ast.Node) bool {
			genDecl, ok := node.(*ast.GenDecl)
			// only find type
			if !ok || genDecl.Tok != token.TYPE {
				return true
			}
			if len(genDecl.Specs) == 0 {
				return true
			}
			for _, spec := range genDecl.Specs {
				for _, name := range typeNames {
					typeSpec := spec.(*ast.TypeSpec)
					if typeSpec.Name.Name != name {
						continue
					}
					ret[name] = &TypeSpec{
						Name: name,
						AstSpec: &AstSpec{
							Pkg:      pkg,
							File:     file,
							TypeSpec: typeSpec,
						},
					}
					break
				}
			}
			return true
		})
	}
	return mo.Ok(ret)
}

func importedPkgFind(pkg *packages.Package, file *ast.File, importPkgNames ...string) map[string]*packages.Package {

	ret := make(map[string]*packages.Package)
	slice.ForEach(file.Imports, func(_ int, importSpec *ast.ImportSpec) {
		importPath := strings.Trim(importSpec.Path.Value, "\"")
		var importName string
		if importSpec.Name != nil {
			importName = importSpec.Name.Name
		} else {
			importName = filepath.Base(importPath)
		}
		if !slice.Contain(importPkgNames, importName) {
			return
		}
		ret[importName] = pkg.Imports[importPath]
	})

	return ret
}
func FindAstFile(pkg *packages.Package, fileName string) mo.Result[*ast.File] {
	return xslice.FindByR(pkg.Syntax, func(_ int, file *ast.File) bool {
		return getFileName(pkg, file) == fileName
	}, errors.Newf("file %s not found", fileName))
}
func newPackagesConfig(opts ...func(config *packages.Config)) *packages.Config {
	cfg := &packages.Config{
		Mode:  packages.NeedName | packages.NeedImports | packages.NeedTypes | packages.NeedSyntax | packages.NeedDeps,
		Tests: false,
	}

	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}
func loadPkgs(patterns []string, opts ...func(config *packages.Config)) mo.Result[[]*packages.Package] {
	cfg := newPackagesConfig(opts...)
	allPkgs, err := packages.Load(cfg, patterns...)
	return mo.TupleToResult(allPkgs, errors.Wrap(err, "packages load failed"))
}

func LoadLocalPkg(pkgName string, opts ...func(config *packages.Config)) mo.Result[*packages.Package] {
	return xutil.FlatMap(
		loadPkgs([]string{"./"}, opts...),
		func(pkgs []*packages.Package) mo.Result[*packages.Package] {
			return xslice.FindByR(pkgs, func(index int, item *packages.Package) bool {
				return item.Name == pkgName
			}, errors.New("must load local pkg"))
		},
	)

}

func getFileName(pkg *packages.Package, file *ast.File) string {
	// 获取文件路径
	return pkg.Fset.Position(file.Pos()).Filename
}
