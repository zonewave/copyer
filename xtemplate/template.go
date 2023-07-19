package xtemplate

import (
	"text/template"

	"github.com/cockroachdb/errors"
	"github.com/samber/mo"
	"github.com/zonewave/copyer/common"
	ts "github.com/zonewave/copyer/templates"
	"github.com/zonewave/copyer/xutil"
)

var (
	funcsMap = template.FuncMap{
		"hasField": HasField,
	}
)

func NewTmpl(tmplType common.ActionType) mo.Result[*template.Template] {
	switch tmplType {
	case common.Local:
		return NewCopyTemplate()
	case common.Outfile:
		return mo.TupleToResult(NewOutPutFileTemplate())
	default:
		return mo.Err[*template.Template](errors.Errorf("tmplType %s not found", tmplType))
	}
}

func NewCopyTemplate() mo.Result[*template.Template] {
	mapErrWrap := xutil.MapWrap[*template.Template]
	return mo.TupleToResult(
		template.New(ts.CopyTmplName.String()).Funcs(funcsMap).ParseFS(ts.Fs, ts.CopyTmplName.FileName()),
	).
		MapErr(mapErrWrap("load CopyTmplName template error")).
		Map(func(tmpl *template.Template) (*template.Template, error) {
			return tmpl.New("only out copyFunc").Parse(ts.CopyTmplName.Template())
		}).
		MapErr(mapErrWrap("parse copyTemplate error"))

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
