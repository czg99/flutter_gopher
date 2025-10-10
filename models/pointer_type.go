package models

type GoPointerType struct {
	Inner GoType
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

func (t *GoPointerType) MapName() string {
	return "Nullable" + t.Inner.MapName()
}

func (t *GoPointerType) NeedMap() bool {
	return true
}
