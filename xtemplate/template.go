package xtemplate

import (
	"text/template"

	"github.com/samber/mo"
	ts "github.com/zonewave/copyer/templates"
	"github.com/zonewave/copyer/xutil"
)

var (
	funcsMap = template.FuncMap{
		"hasField": HasField,
	}
)

type Templates struct {
	CopyTemplate   *template.Template
	ImportTemplate *template.Template
}

func NewTemplates(copyTemplate *template.Template, importTemplate *template.Template) *Templates {
	return &Templates{CopyTemplate: copyTemplate, ImportTemplate: importTemplate}
}

func NewTmpl() mo.Result[*Templates] {
	return xutil.Map2(NewCopyTemplate(), NewImportTemplate(), NewTemplates)
}

func NewCopyTemplate() mo.Result[*template.Template] {
	mapErrWrap := xutil.MapWrap[*template.Template]
	return mo.TupleToResult(
		template.New(ts.CopyTmplName.String()).Funcs(funcsMap).ParseFS(ts.Fs, ts.CopyTmplName.FileName()),
	).
		MapErr(mapErrWrap("load CopyTmplName template error")).
		FlatMap(func(tmpl *template.Template) mo.Result[*template.Template] {
			return mo.TupleToResult(tmpl.New("only out copyFunc").Parse(ts.CopyTmplName.Template())).MapErr(mapErrWrap("parse copyTemplate error"))
		})

}
func NewImportTemplate() mo.Result[*template.Template] {
	mapErrWrap := xutil.MapWrap[*template.Template]
	return mo.TupleToResult(
		template.New(ts.CopyTmplName.String()).Funcs(funcsMap).ParseFS(ts.Fs, ts.CopyTmplName.FileName()),
	).
		MapErr(mapErrWrap("load import template error")).
		FlatMap(func(tmpl *template.Template) mo.Result[*template.Template] {
			return mo.TupleToResult(tmpl.New(" out import").Parse(ts.ImportTmplName.Template())).MapErr(mapErrWrap("parse ImportTmpl error"))
		})

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
