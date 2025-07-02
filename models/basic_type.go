package models

import "github.com/iancoleman/strcase"

var BasicTypeMap = map[string]*GoBasicType{
	"bool": {cType: "bool", goType: "bool", goCType: "C.bool", dartCType: "ffi.Bool", dartType: "bool", dartDefault: "false", kotlinType: "Boolean", kotlinCType: "Boolean", kotlinDefault: "false", kotlinPackagePath: "java/lang/Boolean", kotlinGetValue: "booleanValue"},

	"string": {cType: "void*", goType: "string", goCType: "unsafe.Pointer", dartCType: "ffi.Pointer<ffi.Void>", dartType: "String", dartDefault: "''", needMap: true, kotlinType: "String", kotlinCType: "String", kotlinDefault: `""`, kotlinPackagePath: "java/lang/String", kotlinGetValue: "getBytes"},
	"error":  {cType: "void*", goType: "error", goCType: "unsafe.Pointer", dartCType: "ffi.Pointer<ffi.Void>", dartType: "Error", dartDefault: "null", needMap: true, kotlinType: "Error", kotlinCType: "String", kotlinDefault: "null", kotlinPackagePath: "java/lang/String", kotlinGetValue: "getBytes"},

	"int8":  {cType: "int8_t", goType: "int8", goCType: "C.int8_t", dartCType: "ffi.Int8", dartType: "int", dartDefault: "0", kotlinType: "Int", kotlinCType: "Int", kotlinDefault: "0", kotlinPackagePath: "java/lang/Integer", kotlinGetValue: "intValue"},
	"int16": {cType: "int16_t", goType: "int16", goCType: "C.int16_t", dartCType: "ffi.Int16", dartType: "int", dartDefault: "0", kotlinType: "Int", kotlinCType: "Int", kotlinDefault: "0", kotlinPackagePath: "java/lang/Integer", kotlinGetValue: "intValue"},
	"int32": {cType: "int32_t", goType: "int32", goCType: "C.int32_t", dartCType: "ffi.Int32", dartType: "int", dartDefault: "0", kotlinType: "Int", kotlinCType: "Int", kotlinDefault: "0", kotlinPackagePath: "java/lang/Integer", kotlinGetValue: "intValue"},
	"int64": {cType: "int64_t", goType: "int64", goCType: "C.int64_t", dartCType: "ffi.Int64", dartType: "int", dartDefault: "0", kotlinType: "Long", kotlinCType: "Long", kotlinDefault: "0", kotlinPackagePath: "java/lang/Long", kotlinGetValue: "longValue"},

	"byte":   {cType: "uint8_t", goType: "byte", goCType: "C.uint8_t", dartCType: "ffi.Uint8", dartType: "int", dartDefault: "0", kotlinType: "Byte", kotlinCType: "Byte", kotlinDefault: "0", kotlinPackagePath: "java/lang/Byte", kotlinGetValue: "byteValue"},
	"uint8":  {cType: "uint8_t", goType: "uint8", goCType: "C.uint8_t", dartCType: "ffi.Uint8", dartType: "int", dartDefault: "0", kotlinType: "Int", kotlinCType: "Int", kotlinDefault: "0", kotlinPackagePath: "java/lang/Integer", kotlinGetValue: "intValue"},
	"uint16": {cType: "uint16_t", goType: "uint16", goCType: "C.uint16_t", dartCType: "ffi.Uint16", dartType: "int", dartDefault: "0", kotlinType: "Int", kotlinCType: "Int", kotlinDefault: "0", kotlinPackagePath: "java/lang/Integer", kotlinGetValue: "intValue"},
	"uint32": {cType: "uint32_t", goType: "uint32", goCType: "C.uint32_t", dartCType: "ffi.Uint32", dartType: "int", dartDefault: "0", kotlinType: "Int", kotlinCType: "Int", kotlinDefault: "0", kotlinPackagePath: "java/lang/Integer", kotlinGetValue: "intValue"},
	"uint64": {cType: "uint64_t", goType: "uint64", goCType: "C.uint64_t", dartCType: "ffi.Uint64", dartType: "int", dartDefault: "0", kotlinType: "Long", kotlinCType: "Long", kotlinDefault: "0", kotlinPackagePath: "java/lang/Long", kotlinGetValue: "longValue"},

	"float32": {cType: "float", goType: "float32", goCType: "C.float", dartCType: "ffi.Float", dartType: "double", dartDefault: "0", kotlinType: "Float", kotlinCType: "Float", kotlinDefault: "0.0", kotlinPackagePath: "java/lang/Float", kotlinGetValue: "floatValue"},
	"float64": {cType: "double", goType: "float64", goCType: "C.double", dartCType: "ffi.Double", dartType: "double", dartDefault: "0", kotlinType: "Double", kotlinCType: "Double", kotlinDefault: "0.0", kotlinPackagePath: "java/lang/Double", kotlinGetValue: "doubleValue"},

	"int":  {cType: "int", goType: "int", goCType: "C.int", dartCType: "ffi.Int", dartType: "int", dartDefault: "0", kotlinType: "Int", kotlinCType: "Int", kotlinDefault: "0", kotlinPackagePath: "java/lang/Integer", kotlinGetValue: "intValue"},
	"uint": {cType: "unsigned int", goType: "uint", goCType: "C.uint", dartCType: "ffi.UnsignedInt", dartType: "int", dartDefault: "0", kotlinType: "Int", kotlinCType: "Int", kotlinDefault: "0", kotlinPackagePath: "java/lang/Integer", kotlinGetValue: "intValue"},

	"uintptr": {cType: "uintptr_t", goType: "uintptr", goCType: "C.uintptr_t", dartCType: "ffi.UintPtr", dartType: "int", dartDefault: "0", kotlinType: "Long", kotlinCType: "Long", kotlinDefault: "0", kotlinPackagePath: "java/lang/Long", kotlinGetValue: "longValue"},
}

type GoBasicType struct {
	cType string

	goType  string
	goCType string

	dartType    string
	dartCType   string
	dartDefault string

	kotlinType        string
	kotlinCType       string
	kotlinDefault     string
	kotlinPackagePath string
	kotlinGetValue    string

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
func (t *GoBasicType) KotlinType() string {
	return t.kotlinType
}
func (t *GoBasicType) KotlinCType() string {
	return t.kotlinCType
}
func (t *GoBasicType) KotlinDefault() string {
	return t.kotlinDefault
}
func (t *GoBasicType) KotlinPackagePath() string {
	return t.kotlinPackagePath
}
func (t *GoBasicType) KotlinGetValue() string {
	return t.kotlinGetValue
}
func (t *GoBasicType) MapName() string {
	return strcase.ToCamel(t.goType)
}
func (t *GoBasicType) NeedMap() bool {
	return t.needMap
}
