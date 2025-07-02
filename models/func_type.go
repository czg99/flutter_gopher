package models

import (
	"strings"

	"github.com/iancoleman/strcase"
)

type GoFuncType struct {
	PackageName string
	Name        string
	ClassName   string
	Params      *GoStructType
	Results     *GoStructType
	ResultCount int
	HasErr      bool
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

func (t *GoFuncType) KotlinType() string {
	return strcase.ToCamel(t.Name)
}

func (t *GoFuncType) KotlinCType() string {
	return "fg" + strcase.ToCamel(t.Name)
}

func (t *GoFuncType) KotlinDefault() string {
	return "null"
}

func (t *GoFuncType) KotlinPackagePath() string {
	return strings.ReplaceAll(t.PackageName+"."+t.KotlinType(), ".", "/")
}

func (t *GoFuncType) GoJniType() string {
	name := "Java." + t.PackageName + "._" + t.ClassName + "." + t.KotlinCType()
	return strings.ReplaceAll(strings.ReplaceAll(name, "_", "_1"), ".", "_")
}

func (t *GoFuncType) DartResultType() string {
	if t.ResultCount > 1 {
		return t.Results.DartType()
	} else if t.ResultCount == 1 {
		return t.Results.Fields[0].DartType()
	}
	return "void"
}

func (t *GoFuncType) KotlinResultType() string {
	if t.ResultCount > 1 {
		return t.Results.KotlinType()
	} else if t.ResultCount == 1 {
		return t.Results.Fields[0].KotlinType()
	}
	return "Unit"
}

func (t *GoFuncType) MapName() string {
	return ""
}

func (t *GoFuncType) NeedMap() bool {
	return false
}
