package xast

import (
	"go/ast"
	"go/token"
	"path/filepath"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/hashicorp/go-set"
	"github.com/zonewave/pkgs/standutil/sliceutil"
	"golang.org/x/tools/go/packages"
)

type AstSpec struct {
	Pkg      *packages.Package
	File     *ast.File
	TypeSpec *ast.TypeSpec // option
}

type VarDataSpec struct {
	Name            string
	AstSpec         *AstSpec
	TypePackageName string
	TypeSpec        *TypeSpec
}

func (v *VarDataSpec) IsPublic() bool {

	return ast.IsExported(v.Name)
}

type TypeSpec struct {
	Name    string
	AstSpec *AstSpec
}

type FindVarDataSpecPair struct {
	varName  string
	pkgName  string
	typeName string
}

func NewFindVarDataSpecPair(varName string, pkgName string, typeName string) *FindVarDataSpecPair {
	return &FindVarDataSpecPair{varName: varName, pkgName: pkgName, typeName: typeName}
}

func VarSpecLocalParseMust(pkg *packages.Package, file *ast.File, pairs ...*FindVarDataSpecPair) (map[string]*VarDataSpec, error) {
	ret, err := VarSpecLocalParse(pkg, file, pairs...)
	if err != nil {
		return nil, err
	}
	if len(ret) != len(pairs) {
		return nil, errors.Errorf("pairs  :%v not all found", pairs)
	}
	return ret, nil

}

func VarSpecLocalParse(pkg *packages.Package, file *ast.File, pairs ...*FindVarDataSpecPair) (map[string]*VarDataSpec, error) {
	ret := make(map[string]*VarDataSpec, len(pairs))
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

	localNamesFilter := func(pair *FindVarDataSpecPair) bool { return pair.pkgName == "" }
	// local
	localPairs := sliceutil.Filter(pairs, localNamesFilter)

	for _, pair := range localPairs {
		typeSpecs, err := TypeSpecMustGet(pkg, pair.typeName)
		if err != nil {
			return nil, err
		}
		ret[pair.varName].TypeSpec = typeSpecs[pair.typeName]
	}

	// import
	importPairs := sliceutil.Filter(pairs, func(pair *FindVarDataSpecPair) bool {
		return pair.pkgName != ""
	})
	pkgNamesGet := func(pair *FindVarDataSpecPair) string { return pair.pkgName }
	importPkgsNames := sliceutil.Map(importPairs, pkgNamesGet)
	pkgs, err := mustFindImportPkg(pkg, file, importPkgsNames...)
	if err != nil {
		return nil, err
	}
	for _, pair := range importPairs {
		importPkg := pkgs[pair.pkgName]
		typeSpecs, gErr := TypeSpecMustGet(importPkg, pair.typeName)
		if gErr != nil {
			return nil, gErr
		}
		ret[pair.varName].TypeSpec = typeSpecs[pair.typeName]
	}

	return ret, nil
}

func TypeSpecMustGet(pkg *packages.Package, typeNames ...string) (map[string]*TypeSpec, error) {
	ret := TypeSpecGet(pkg, typeNames...)
	if len(ret) != len(typeNames) {
		return nil, errors.Errorf("type spec :%v not all found", typeNames)
	}
	return ret, nil
}
func TypeSpecGet(pkg *packages.Package, typeNames ...string) map[string]*TypeSpec {
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
					if typeSpec.Name.Name == name {
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
			}
			return true
		})
	}
	return ret
}

func RemoveDuplicate[T comparable](s []T) []T {
	return set.From(s).Slice()
}
func mustFindImportPkg(pkg *packages.Package, file *ast.File, importPkgNames ...string) (map[string]*packages.Package, error) {
	importPkgNames = RemoveDuplicate(importPkgNames)
	ret := importedPkgSearch(pkg, file, importPkgNames...)
	if len(ret) != len(importPkgNames) {
		return nil, errors.Errorf("import pkg :%v not all found", importPkgNames)
	}
	return ret, nil

}
func importedPkgSearch(pkg *packages.Package, file *ast.File, importPkgNames ...string) map[string]*packages.Package {

	ret := make(map[string]*packages.Package)
	for _, importSpec := range file.Imports {
		importName := ""
		importPath := strings.Trim(importSpec.Path.Value, "\"")
		if importSpec.Name != nil {
			importName = importSpec.Name.Name
		} else {
			importName = filepath.Base(importPath)
		}
		if !sliceutil.Contain(importName, importPkgNames) {
			continue
		}
		ret[importName] = pkg.Imports[importPath]
	}

	return ret
}
func FindAstFile(pkg *packages.Package, fileName string) (*ast.File, error) {
	for _, file := range pkg.Syntax {
		if getFileName(pkg, file) == fileName {
			return file, nil
		}
	}
	return nil, errors.Newf("file %s not found", fileName)
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
func loadPkgs(patterns []string, opts ...func(config *packages.Config)) ([]*packages.Package, error) {

	cfg := newPackagesConfig(opts...)
	allPkgs, err := packages.Load(cfg, patterns...)
	if err != nil {
		return nil, errors.Wrap(err, "packages load failed")
	}
	return allPkgs, nil
}

func LoadLocalPkg(pkgName string, opts ...func(config *packages.Config)) (*packages.Package, error) {
	pkgs, err := loadPkgs([]string{"./"}, opts...)
	if err != nil {
		return nil, err
	}
	var localPkg *packages.Package
	for _, pkg := range pkgs {
		if pkg.Name == pkgName {
			localPkg = pkg
			break
		}
	}
	if localPkg == nil {
		return nil, errors.New("must load local pkg")
	}
	return localPkg, nil
}

func getFileName(pkg *packages.Package, file *ast.File) string {
	// 获取文件路径
	return pkg.Fset.Position(file.Pos()).Filename
}
