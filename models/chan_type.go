package models

import "strings"

type GoChanType struct {
	PackageName string
	Inner       GoType
}

func (t *GoChanType) String() string {
	return "chan " + t.Inner.String()
}

func (t *GoChanType) CType() string {
	return "FgChan"
}

func (t *GoChanType) GoType() string {
	return "chan " + t.Inner.GoType()
}

func (t *GoChanType) GoCType() string {
	return "C.FgChan"
}

func (t *GoChanType) DartType() string {
	return "FgChan<" + t.Inner.DartType() + ">"
}

func (t *GoChanType) DartCType() string {
	return "_fgChan"
}

func (t *GoChanType) DartDefault() string {
	return "FgChan<" + t.Inner.DartType() + ">()"
}

func (t *GoChanType) KotlinType() string {
	return "FgChan<" + t.Inner.KotlinType() + ">"
}

func (t *GoChanType) KotlinCType() string {
	return "_fgChan"
}

func (t *GoChanType) KotlinDefault() string {
	return "FgChan<" + t.Inner.KotlinType() + ">()"
}

func (t *GoChanType) KotlinPackagePath() string {
	return strings.ReplaceAll(t.PackageName+"."+t.KotlinType(), ".", "/")
}

func (t *GoChanType) MapName() string {
	return t.Inner.MapName() + "Chan"
}

func (t *GoChanType) NeedMap() bool {
	return true
}
