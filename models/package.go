package models

type Package struct {
	Module  string
	PkgPath string
	Structs []*GoStructType
	Funcs   []*GoFuncType
}
