package main

import (
	"go/ast"

	"github.com/cockroachdb/errors"
	"golang.org/x/tools/go/packages"
)

type TemplateStruct struct {
	Name   string
	Fields map[string]struct{}
}

func (t *TemplateStruct) hasField(name string) bool {
	_, ok := t.Fields[name]
	return ok
}

func HasField(src *TemplateStruct, name string) bool {
	return src.hasField(name)
}

type TemplateParam struct {
	Src *TemplateStruct
	Dst *TemplateStruct
}

func NewTemplateParam(src, dst *TemplateStruct) *TemplateParam {
	return &TemplateParam{
		Src: src,
		Dst: dst,
	}
}

// parseTemplateParam parses template param
func parseTemplateParam(srcName, dstName string, pkgs []*packages.Package) (*TemplateParam, error) {
	specs := findTypeSpec(pkgs, []string{srcName, dstName})
	structs, err := parseTypeSecStructs(specs)
	if err != nil {
		return nil, err
	}
	src, ok := structs[srcName]
	if !ok {
		return nil, errors.Errorf("src %s not found", src)
	}
	dst, ok := structs[dstName]
	if !ok {
		return nil, errors.Errorf("dst %s not found", dst)
	}
	return NewTemplateParam(src, dst), nil
}

// parseTypeSecStruct parses ast.TypeSpec to TemplateStruct
func parseTypeSecStruct(spec *ast.TypeSpec) (*TemplateStruct, error) {
	s := &TemplateStruct{
		Name: spec.Name.Name,
	}
	structType, ok := spec.Type.(*ast.StructType)
	if !ok {
		return nil, errors.New("specType isn't  StructType")
	}

	s.Fields = make(map[string]struct{}, structType.Fields.NumFields())
	for _, g := range structType.Fields.List {
		if len(g.Names) == 0 {
			continue
		}
		for _, item := range g.Names {
			s.Fields[item.Name] = struct{}{}
		}
	}

	return s, nil
}

// parseTypeSecStructs parses map[string]ast.TypeSpecs to map[string]TemplateStructs
func parseTypeSecStructs(spec map[string]*ast.TypeSpec) (map[string]*TemplateStruct, error) {
	ret := make(map[string]*TemplateStruct, len(spec))
	for name, ts := range spec {
		s, err := parseTypeSecStruct(ts)
		if err != nil {
			return nil, err
		}
		ret[name] = s
	}

	return ret, nil
}
