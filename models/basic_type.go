package models

import (
	"github.com/iancoleman/strcase"
)

var BasicTypeMap = map[string]*GoBasicType{
	"bool": {cType: "bool", goType: "bool", goCType: "C.bool", dartCType: "ffi.Bool", dartType: "bool", dartDefault: "false"},

	"string": {cType: "FgData", goType: "string", goCType: "C.FgData", dartCType: "_fgData", dartType: "String", dartDefault: "''", needMap: true},
	"error":  {cType: "FgData", goType: "error", goCType: "C.FgData", dartCType: "_fgData", dartType: "String?", dartDefault: "null", needMap: true},
	"[]byte": {cType: "FgData", goType: "[]byte", goCType: "C.FgData", dartCType: "_fgData", dartType: "Uint8List", dartDefault: "Uint8List(0)", needMap: true, mapName: "Bytes"},

	"int8":  {cType: "int8_t", goType: "int8", goCType: "C.int8_t", dartCType: "ffi.Int8", dartType: "int", dartDefault: "0"},
	"int16": {cType: "int16_t", goType: "int16", goCType: "C.int16_t", dartCType: "ffi.Int16", dartType: "int", dartDefault: "0"},
	"int32": {cType: "int32_t", goType: "int32", goCType: "C.int32_t", dartCType: "ffi.Int32", dartType: "int", dartDefault: "0"},
	"int64": {cType: "int64_t", goType: "int64", goCType: "C.int64_t", dartCType: "ffi.Int64", dartType: "int", dartDefault: "0"},

	"byte":   {cType: "uint8_t", goType: "byte", goCType: "C.uint8_t", dartCType: "ffi.Uint8", dartType: "int", dartDefault: "0"},
	"uint8":  {cType: "uint8_t", goType: "uint8", goCType: "C.uint8_t", dartCType: "ffi.Uint8", dartType: "int", dartDefault: "0"},
	"uint16": {cType: "uint16_t", goType: "uint16", goCType: "C.uint16_t", dartCType: "ffi.Uint16", dartType: "int", dartDefault: "0"},
	"uint32": {cType: "uint32_t", goType: "uint32", goCType: "C.uint32_t", dartCType: "ffi.Uint32", dartType: "int", dartDefault: "0"},
	"uint64": {cType: "uint64_t", goType: "uint64", goCType: "C.uint64_t", dartCType: "ffi.Uint64", dartType: "int", dartDefault: "0"},

	"float32": {cType: "float", goType: "float32", goCType: "C.float", dartCType: "ffi.Float", dartType: "double", dartDefault: "0"},
	"float64": {cType: "double", goType: "float64", goCType: "C.double", dartCType: "ffi.Double", dartType: "double", dartDefault: "0"},

	"int":  {cType: "int", goType: "int", goCType: "C.int", dartCType: "ffi.Int", dartType: "int", dartDefault: "0"},
	"uint": {cType: "unsigned int", goType: "uint", goCType: "C.uint", dartCType: "ffi.UnsignedInt", dartType: "int", dartDefault: "0"},

	"uintptr": {cType: "uintptr_t", goType: "uintptr", goCType: "C.uintptr_t", dartCType: "ffi.UintPtr", dartType: "int", dartDefault: "0"},
}

type GoBasicType struct {
	cType string

	goType  string
	goCType string

	dartType    string
	dartCType   string
	dartDefault string

	mapName string
	needMap bool
}

func (t *GoBasicType) String() string {
	return t.goType
}

func (t *GoBasicType) CType() string {
	return t.cType
}
func (t *GoBasicType) GoType() string {
	return t.goType
}
func (t *GoBasicType) GoCType() string {
	return t.goCType
}
func (t *GoBasicType) DartType() string {
	return t.dartType
}
func (t *GoBasicType) DartCType() string {
	return t.dartCType
}
func (t *GoBasicType) DartDefault() string {
	return t.dartDefault
}
func (t *GoBasicType) MapName() string {
	if t.mapName != "" {
		return t.mapName
	}
	return strcase.ToCamel(t.goType)
}
func (t *GoBasicType) NeedMap() bool {
	return t.needMap
}
