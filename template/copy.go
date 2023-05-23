package template

const CopyTmplName = "copy"
const CopyTmpl = `
func Copy{{.Src.Name}}To{{.Dst.Name}}(src *{{.Src.Name}}, dst *{{.Dst.Name}}) {
{{- range $name,$value := .Dst.Fields}}
    {{- if hasField $.Src $name -}}
        dst.{{$name}} = src.{{$name -}}
    {{end}}
{{end -}}
}
`
