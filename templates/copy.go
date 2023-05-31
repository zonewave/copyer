package templates

import _ "embed"

const CopyTmplName = "copy"

//go:embed copy.gotmpl
var CopyTmpl string
