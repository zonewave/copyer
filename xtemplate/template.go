package xtemplate

import (
	"text/template"

	ts "github.com/zonewave/copyer/templates"
)

func NewCopyTemplate() (*template.Template, error) {
	funcsMap := template.FuncMap{
		"hasField": HasField,
	}

	tmpl, err := template.New(ts.CopyTmplName.String()).Funcs(funcsMap).ParseFS(ts.Fs, ts.CopyTmplName.FileName())
	if err != nil {
		return nil, err
	}
	tmpl2, err := tmpl.New("only out copyFunc").Funcs(funcsMap).Parse(ts.CopyTmplName.Template())
	if err != nil {
		return nil, err
	}
	return tmpl2, nil
}
