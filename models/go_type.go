package models

type GoType interface {
	String() string

	CType() string
	GoType() string
	GoCType() string
	DartType() string
	DartCType() string
	DartDefault() string
	MapName() string
}
