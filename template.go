package main

import (
	"text/template"

	"github.com/cockroachdb/errors"
	mTemplate "github.com/zonewave/copyer/template"
)

func newCopyTemplate() (*template.Template, error) {
	funcsMap := template.FuncMap{
		"hasField": HasField,
	}
	t, err := template.New(mTemplate.CopyTmplName).Funcs(funcsMap).Parse(mTemplate.CopyTmpl)
	if err != nil {
		return nil, errors.Wrapf(err, "parse template %s", mTemplate.CopyTmplName)
	}

	return t, nil
}
