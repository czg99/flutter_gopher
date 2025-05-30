package models

type GoChanType struct {
	Inner GoType
}

func (t *GoChanType) String() string {
	return "chan " + t.Inner.String()
}

func (t *GoChanType) CType() string {
	return "fg_chan"
}

func (t *GoChanType) GoType() string {
	return "chan " + t.Inner.GoType()
}

func (t *GoChanType) GoCType() string {
	return "C.fg_chan"
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

func (t *GoChanType) MapName() string {
	return t.Inner.MapName() + "Chan"
}

func (t *GoChanType) NeedMap() bool {
	return true
}
