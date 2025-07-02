package models

type GoSliceType struct {
	Inner GoType
}

func (t *GoSliceType) String() string {
	return "[]" + t.Inner.String()
}

func (t *GoSliceType) CType() string {
	return "FgArray"
}

func (t *GoSliceType) GoType() string {
	return "[]" + t.Inner.GoType()
}

func (t *GoSliceType) GoCType() string {
	return "C.FgArray"
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
	return t.Inner.MapName() + "List"
}

func (t *GoSliceType) NeedMap() bool {
	return true
}
