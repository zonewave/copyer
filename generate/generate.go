package generate

import (
	"bytes"
	"go/format"
	"io"
	"os"
	"text/template"

	"github.com/cockroachdb/errors"
	"github.com/zonewave/copyer/output"
	"github.com/zonewave/copyer/parser"
	"github.com/zonewave/copyer/xast"
	"github.com/zonewave/copyer/xtemplate"
)

// Generator is a code generator
type Generator struct {
	Param *xtemplate.CopyParam
	Tmpl  *template.Template
	Out   io.Writer
}

// NewGenerator creates a new generator
func NewGenerator(arg *GeneratorArg) (*Generator, error) {

	// load ast.pkgs
	pkg, err := xast.LoadLocalPkg(arg.LoadConfigOpts...)
	if err != nil {
		return nil, err
	}

	// load templates
	tmpl, err := xtemplate.NewTmpl(arg.Action)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// parser templates param
	tmplParam, err := parser.ParseTemplateParam(&parser.ParseTemplateParamArg{
		Action:      arg.Action,
		FileName:    arg.GoFile,
		SrcName:     arg.SrcName,
		SrcPkg:      arg.SrcPkg,
		SrcTypeName: arg.SrcType,
		DstName:     arg.DstName,
		DstPkg:      arg.DstPkg,
		DstTypeName: arg.DstType,
		Pkg:         pkg,
	})
	if err != nil {
		return nil, err
	}
	// output
	var outPut io.Writer
	if arg.Print {
		outPut = os.Stdout
	} else {
		outPut = output.NewOutput(arg.OutFile, arg.OutLine)
	}
	g := newGenerate(tmplParam, tmpl, outPut)
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
