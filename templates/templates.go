package templates

import (
	"embed"
)

//go:embed "*.gotmpl"
var Fs embed.FS

type TmplName string

const (
	AllTmplName     TmplName = "copyer"
	CopyTmplName    TmplName = "copyFunc"
	ImportTmplName  TmplName = "import"
	OutFileTmplName TmplName = "outfile"
)

func (t TmplName) String() string {
	return string(t)
}
func (t TmplName) FileName() string {
	return t.String() + ".gotmpl"
}

func (t TmplName) Template() string {
	return `{{template "copyFunc" . }}`
}
