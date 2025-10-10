package models

import (
	"unicode"

	"github.com/iancoleman/strcase"
)

type GoIdentType struct {
	Name string
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
	if unicode.IsLower(rune(t.Name[0])) {
		return "_" + strcase.ToLowerCamel(t.Name)
	}
	return strcase.ToCamel(t.Name)
}

func (t *GoIdentType) DartCType() string {
	return "_fg" + strcase.ToCamel(t.Name)
}

func (t *GoIdentType) DartDefault() string {
	return t.DartType() + "()"
}

func (t *GoIdentType) MapName() string {
	return strcase.ToCamel(t.Name)
}

func (t *GoIdentType) NeedMap() bool {
	return true
}
