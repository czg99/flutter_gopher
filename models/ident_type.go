package models

import (
	"strings"

	"github.com/iancoleman/strcase"
)

type GoIdentType struct {
	PackageName string
	Name        string
}

func (t *GoIdentType) String() string {
	return t.Name
}

func (t *GoIdentType) CType() string {
	return "struct Fg" + strcase.ToCamel(t.Name)
}

func (t *GoIdentType) GoType() string {
	return t.Name
}

func (t *GoIdentType) GoCType() string {
	return "C.Fg" + strcase.ToCamel(t.Name)
}

func (t *GoIdentType) DartType() string {
	return strcase.ToCamel(t.Name)
}

func (t *GoIdentType) DartCType() string {
	return "_fg" + strcase.ToCamel(t.Name)
}

func (t *GoIdentType) DartDefault() string {
	return t.DartType() + "()"
}

func (t *GoIdentType) KotlinType() string {
	return strcase.ToCamel(t.Name)
}

func (t *GoIdentType) KotlinCType() string {
	return "fg" + strcase.ToCamel(t.Name)
}

func (t *GoIdentType) KotlinDefault() string {
	return t.KotlinType() + "()"
}

func (t *GoIdentType) KotlinPackagePath() string {
	return strings.ReplaceAll(t.PackageName+"."+t.KotlinType(), ".", "/")
}

func (t *GoIdentType) MapName() string {
	return strcase.ToCamel(t.Name)
}

func (t *GoIdentType) NeedMap() bool {
	return true
}
