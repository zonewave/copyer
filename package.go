package main

import (
	"go/ast"
	"go/token"

	"github.com/cockroachdb/errors"
	"golang.org/x/tools/go/packages"
)

func findTypeSpec(pkgs []*packages.Package, typeNames []string) map[string]*ast.TypeSpec {
	ret := make(map[string]*ast.TypeSpec, len(typeNames))
	for _, pkg := range pkgs {
		for _, file := range pkg.Syntax {
			ast.Inspect(file, func(node ast.Node) bool {
				genDecl, ok := node.(*ast.GenDecl)
				if !ok || genDecl.Tok != token.TYPE {
					return true
				}
				if len(genDecl.Specs) == 0 {
					return true
				}
				for _, spec := range genDecl.Specs {
					for _, name := range typeNames {
						if spec.(*ast.TypeSpec).Name.Name == name {
							ret[name] = spec.(*ast.TypeSpec)
							break
						}
					}
				}
				return true
			})
		}
	}
	return ret
}

func newPackagesConfig(opts ...func(config *packages.Config)) *packages.Config {
	cfg := &packages.Config{
		Mode:  packages.NeedName | packages.NeedSyntax,
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
