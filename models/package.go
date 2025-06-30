package models

type Package struct {
	ProjectNaming
	Module  string
	PkgPath string
	Structs []*GoStructType
	Funcs   []*GoFuncType
}
