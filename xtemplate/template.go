package xtemplate

import (
	"text/template"

	"github.com/cockroachdb/errors"
	"github.com/zonewave/copyer/common"
	ts "github.com/zonewave/copyer/templates"
)

var (
	funcsMap = template.FuncMap{
		"hasField": HasField,
	}
)

func NewTmpl(tmplType common.ActionType) (*template.Template, error) {
	switch tmplType {
	case common.Local:
		return NewCopyTemplate()
	case common.Outfile:
		return NewOutPutFileTemplate()
	default:
		return nil, errors.Errorf("tmplType %s not found", tmplType)
	}
}

func NewCopyTemplate() (*template.Template, error) {
	tmpl, err := template.New(ts.CopyTmplName.String()).Funcs(funcsMap).ParseFS(ts.Fs, ts.CopyTmplName.FileName())
	if err != nil {
		return nil, err
	}
	tmpl2, err := tmpl.New("only out copyFunc").Parse(ts.CopyTmplName.Template())
	if err != nil {
		return nil, err
	}
	return tmpl2, err
}

func NewOutPutFileTemplate() (*template.Template, error) {
	tmpl, err := template.New(ts.OutFileTmplName.String()).Funcs(funcsMap).ParseFS(ts.Fs,
		ts.CopyTmplName.FileName(),
		ts.ImportTmplName.FileName(),
		ts.OutFileTmplName.FileName())
	if err != nil {
		return nil, err
	}
	tmpl2, err := tmpl.New("generate file").Funcs(funcsMap).Parse(ts.OutFileTmplName.Template())
	if err != nil {
		return nil, err
	}
	return tmpl2, nil
}
