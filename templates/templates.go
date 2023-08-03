package templates

import (
	"embed"
	"fmt"
)

//go:embed "*.gotmpl"
var Fs embed.FS

type TmplName string

const (
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
	return fmt.Sprintf(`{{template "%s" . }}`, t.String())
}
