package models

type GoPointerType struct {
	PackageName string
	Inner       GoType
}

func (t *GoPointerType) String() string {
	return "*" + t.Inner.String()
}

func (t *GoPointerType) CType() string {
	return t.Inner.CType() + "*"
}

func (t *GoPointerType) GoType() string {
	return "*" + t.Inner.GoType()
}

func (t *GoPointerType) GoCType() string {
	return "*" + t.Inner.GoCType()
}

func (t *GoPointerType) DartType() string {
	return t.Inner.DartType() + "?"
}

func (t *GoPointerType) DartCType() string {
	return "ffi.Pointer<" + t.Inner.DartCType() + ">"
}

func (t *GoPointerType) DartDefault() string {
	return "null"
}

func (t *GoPointerType) KotlinType() string {
	return t.Inner.KotlinType() + "?"
}

func (t *GoPointerType) KotlinCType() string {
	return t.Inner.KotlinType() + "?"
}

func (t *GoPointerType) KotlinDefault() string {
	return "null"
}

func (t *GoPointerType) KotlinPackagePath() string {
	return t.Inner.KotlinPackagePath()
}

func (t *GoPointerType) MapName() string {
	return "Null" + t.Inner.MapName()
}

func (t *GoPointerType) NeedMap() bool {
	return true
}
