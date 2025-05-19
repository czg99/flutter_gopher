package models

import "github.com/iancoleman/strcase"

type GoSliceType struct {
	Inner GoType
}

func (t *GoSliceType) IsInnerPtr() bool {
	_, ok := t.Inner.(*GoPointerType)
	return ok
}

func (t *GoSliceType) String() string {
	return "[]" + t.Inner.String()
}

func (t *GoSliceType) CType() string {
	return "fg_array"
}

func (t *GoSliceType) GoType() string {
	return "[]" + t.Inner.GoType()
}

func (t *GoSliceType) GoCType() string {
	return "C.fg_array"
}

func (t *GoSliceType) DartType() string {
	return "List<" + t.Inner.DartType() + ">"
}

func (t *GoSliceType) DartCType() string {
	return "_fgArray"
}

func (t *GoSliceType) DartDefault() string {
	return "[]"
}

func (t *GoSliceType) MapName() string {
	if t.IsInnerPtr() {
		return strcase.ToCamel("Null" + t.Inner.DartType() + "List")
	}
	return strcase.ToCamel(t.Inner.DartType() + "List")
}
