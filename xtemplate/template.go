package xtemplate

import (
	"text/template"

	ts "github.com/zonewave/copyer/templates"
)

var (
	copyFuncsMap = template.FuncMap{
		"hasField": HasField,
	}
)

func NewCopyTemplate() (*template.Template, error) {
	tmpl, err := template.New(ts.CopyTmplName.String()).Funcs(copyFuncsMap).ParseFS(ts.Fs, ts.CopyTmplName.FileName())
	if err != nil {
		return nil, err
	}
	tmpl2, err := tmpl.New("only out copyFunc").Parse(ts.CopyTmplName.Template())
	if err != nil {
		return nil, err
	}
	return tmpl2, nil
}

func NewOutPutTemplate() (*template.Template, error) {
	funcsMap := template.FuncMap{
		"hasField": HasField,
	}

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
