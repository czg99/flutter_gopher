package models

import "github.com/iancoleman/strcase"

type GoField struct {
	Name string
	Type GoType
}

func (f *GoField) InnerMost() GoType {
	var most func(t GoType) GoType
	most = func(t GoType) GoType {
		if p, ok := t.(*GoPointerType); ok {
			return most(p.Inner)
		}

		if s, ok := t.(*GoSliceType); ok {
			return most(s.Inner)
		}

		if s, ok := t.(*GoChanType); ok {
			return most(s.Inner)
		}
		return t
	}

	return most(f.Type)
}

func (f *GoField) CName() string {
	return strcase.ToSnake(f.Name)
}

func (f *GoField) GoName() string {
	return f.Name
}

func (f *GoField) DartName() string {
	return strcase.ToLowerCamel(f.Name)
}

func (f *GoField) String() string {
	return f.Type.String()
}

func (f *GoField) CType() string {
	return f.Type.CType()
}

func (f *GoField) GoType() string {
	return f.Type.GoType()
}

func (f *GoField) GoCType() string {
	return f.Type.GoCType()
}

func (f *GoField) DartType() string {
	return f.Type.DartType()
}

func (f *GoField) DartCType() string {
	return f.Type.DartCType()
}

func (f *GoField) DartDefault() string {
	return f.Type.DartDefault()
}

func (f *GoField) MapName() string {
	return f.Type.MapName()
}

func (f *GoField) NeedMap() bool {
	return f.Type.NeedMap()
}
