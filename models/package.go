package models

type Package struct {
	ProjectNaming
	PkgPath string
	Structs []*GoStructType
	Funcs   []*GoFuncType
}
