package models

import "github.com/iancoleman/strcase"

type GoStructType struct {
	Type   GoType
	Fields []*GoField
}

func (t *GoStructType) String() string {
	return t.Type.String()
}

func (t *GoStructType) CType() string {
	return "fg_" + strcase.ToSnake(t.Type.GoType())
}

func (t *GoStructType) GoType() string {
	return t.Type.GoType()
}

func (t *GoStructType) GoCType() string {
	return t.Type.GoCType()
}

func (t *GoStructType) DartType() string {
	return t.Type.DartType()
}

func (t *GoStructType) DartCType() string {
	return t.Type.DartCType()
}

func (t *GoStructType) DartDefault() string {
	return t.Type.DartDefault()
}

func (t *GoStructType) MapName() string {
	return t.Type.MapName()
}

func (t *GoStructType) NeedMap() bool {
	return true
}
