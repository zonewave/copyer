package main

import (
	"bytes"
	"go/format"
	"io"
	"text/template"

	"github.com/cockroachdb/errors"
	"golang.org/x/tools/go/packages"
)

// Generator is a code generator
type Generator struct {
	Param *TemplateParam
	Tmpl  *template.Template
	Out   io.Writer
}

type GeneratorArg struct {
	FileName       string
	Src            string
	Dst            string
	Line           int
	LoadConfigOpts []func(*packages.Config)
}

// NewGenerator creates a new generator
func NewGenerator(arg *GeneratorArg) (*Generator, error) {

	// load ast.pkgs
	pkgs, err := loadPkgs([]string{"./..."}, arg.LoadConfigOpts...)
	if err != nil {
		return nil, err
	}

	// load template
	tmpl, err := newCopyTemplate()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// parse template paramr
	tmplParam, err := parseTemplateParam(arg.Src, arg.Dst, pkgs)
	if err != nil {
		return nil, err
	}
	g := newGenerate(tmplParam, tmpl, NewOutput(arg.FileName, arg.Line))
	return g, nil
}
func newGenerate(param *TemplateParam, tmpl *template.Template, out io.Writer) *Generator {
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
		return nil, errors.Wrapf(err, "execute template %s", g.Tmpl.Name())
	}
	bs, err := format.Source(buf.Bytes())
	if err != nil {
		return nil, errors.Wrapf(err, "format source %s", g.Tmpl.Name())
	}
	return bs, err
}
