package models

import "github.com/iancoleman/strcase"

type GoFuncType struct {
	Name          string
	Params        *GoStructType
	Results       *GoStructType
	ResultCount   int
	DartErrorName string
}

func (t *GoFuncType) String() string {
	return t.Name
}

func (t *GoFuncType) CType() string {
	return "fg_" + strcase.ToSnake(t.Name)
}

func (t *GoFuncType) GoType() string {
	return t.Name
}

func (t *GoFuncType) GoCType() string {
	return t.CType()
}

func (t *GoFuncType) DartType() string {
	return strcase.ToLowerCamel(t.Name)
}

func (t *GoFuncType) DartCType() string {
	return "_fg" + strcase.ToCamel(t.Name)
}

func (t *GoFuncType) DartDefault() string {
	return "null"
}

func (t *GoFuncType) DartResultType() string {
	if t.ResultCount > 1 {
		return t.Results.DartType()
	} else if t.ResultCount == 1 {
		return t.Results.Fields[0].DartType()
	}
	return "void"
}

func (t *GoFuncType) MapName() string {
	return ""
}

func (t *GoFuncType) NeedMap() bool {
	return false
}
