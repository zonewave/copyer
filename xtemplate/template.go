package xtemplate

import (
	"text/template"

	"github.com/zonewave/copyer/templates"
)

func NewCopyTemplate() (*template.Template, error) {
	funcsMap := template.FuncMap{
		"hasField": HasField,
	}

	return template.New(templates.CopyTmplName).Funcs(funcsMap).Parse(templates.CopyTmpl)
}
