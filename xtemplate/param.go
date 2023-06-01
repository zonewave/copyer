package xtemplate

type ImplImportParam struct {
	Alias   string
	PkgPath string
}
type TmplImports struct {
	Imports []*ImplImportParam
}
type TmplVar struct {
	Name          string
	TypeNameNoDot string
	Type          string
	StructType    *TmplStruct
	Exported      bool
}

type TmplStruct struct {
	// typeName
	Name string
	// key: field name
	Fields map[string]*TmplVar
}

func (t *TmplStruct) NameGet() string {
	return t.Name
}
func (t *TmplStruct) HasField(name string) bool {
	field, ok := t.Fields[name]
	if !ok {
		return false
	}
	if field.Exported {
		return true
	}
	return false
}

//go:generate copyer -s TmplStruct -d TmplStruct
func CopyTmplStructToTmplStruct(src *TmplStruct, dst *TmplStruct) {
	dst.Fields = src.Fields
	dst.Name = src.Name
}

func HasField(src *TmplVar, name string) bool {
	return src.StructType.HasField(name)
}

type CopyParam struct {
	Src *TmplVar
	Dst *TmplVar
}

type FileParam struct {
	PackageName string
	Imports     *TmplImports
	Param       *CopyParam
}

func NewTemplateParam(src, dst *TmplVar) *CopyParam {
	return &CopyParam{
		Src: src,
		Dst: dst,
	}
}
