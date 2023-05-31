package main

import (
	"bytes"
	"go/format"
	"io"
	"text/template"

	"github.com/cockroachdb/errors"
	"github.com/zonewave/copyer/parser"
	"github.com/zonewave/copyer/xast"
	"github.com/zonewave/copyer/xtemplate"
	"golang.org/x/tools/go/packages"
)

// Generator is a code generator
type Generator struct {
	Param *xtemplate.CopyParam
	Tmpl  *template.Template
	Out   io.Writer
}

type GeneratorArg struct {
	FileName       string
	Src            string
	SrcPkg         string
	DstPkg         string
	Dst            string
	Line           int
	LoadConfigOpts []func(*packages.Config)
}

// NewGenerator creates a new generator
func NewGenerator(arg *GeneratorArg) (*Generator, error) {

	// load ast.pkgs
	pkg, err := xast.LoadLocalPkg(arg.LoadConfigOpts...)
	if err != nil {
		return nil, err
	}

	// load templates
	tmpl, err := xtemplate.NewCopyTemplate()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// parser templates param
	tmplParam, err := parser.ParseTemplateParam(arg.FileName, arg.SrcPkg, arg.Src, arg.DstPkg, arg.Dst, pkg)
	if err != nil {
		return nil, err
	}
	g := newGenerate(tmplParam, tmpl, NewOutput(arg.FileName, arg.Line))
	return g, nil
}
func newGenerate(param *xtemplate.CopyParam, tmpl *template.Template, out io.Writer) *Generator {
	return &Generator{Param: param, Tmpl: tmpl, Out: out}
}

// OutPut is a wrapper of io.Writer
func (g *Generator) OutPut(bs []byte) error {
	_, e := g.Out.Write(bs)
	return e
}

// Generate generates code
func (g *Generator) Generate() ([]byte, error) {
	var buf bytes.Buffer
	err := g.Tmpl.Execute(&buf, g.Param)
	if err != nil {
		return nil, errors.Wrapf(err, "execute templates %s", g.Tmpl.Name())
	}
	bs, err := format.Source(buf.Bytes())
	if err != nil {
		return nil, errors.Wrapf(err, "format source %s", g.Tmpl.Name())
	}
	return bs, err
}
