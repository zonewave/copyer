package main

import (
	"text/template"

	mTemplate "github.com/zonewave/copyer/template"
)

func newCopyTemplate() (*template.Template, error) {
	funcsMap := template.FuncMap{
		"hasField": HasField,
	}

	return template.New(mTemplate.CopyTmplName).Funcs(funcsMap).Parse(mTemplate.CopyTmpl)
}
