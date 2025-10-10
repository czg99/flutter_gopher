package models

import (
	"strings"

	"github.com/iancoleman/strcase"
)

type GoFuncType struct {
	Name               string        //函数名
	Params             *GoStructType //参数
	Results            *GoStructType //返回值
	HasParams          bool          //是否存在参数
	HasResults         bool          //是否存在返回
	ResultCount        int           //返回数量
	IsAnonymousResults bool          //是否匿名的返回
	HasErr             bool          //是否存在错误字段
}

func (t *GoFuncType) String() string {
	return t.Name
}

func (t *GoFuncType) CType() string {
	return "fg_" + strcase.ToSnake(t.Name)
}

func (t *GoFuncType) GoType() string {
	return t.Name
}

func (t *GoFuncType) GoCType() string {
	return t.CType()
}

func (t *GoFuncType) DartType() string {
	return strcase.ToLowerCamel(t.Name)
}

func (t *GoFuncType) DartCType() string {
	return "_fg" + strcase.ToCamel(t.Name)
}

func (t *GoFuncType) DartDefault() string {
	return "null"
}

func (t *GoFuncType) DartResultType() string {
	if t.ResultCount > 1 {
		builder := strings.Builder{}
		builder.WriteString("(")
		for i, field := range t.Results.Fields {
			if i > 0 {
				builder.WriteString(", ")
			}
			builder.WriteString(field.DartType())
			if !t.IsAnonymousResults {
				builder.WriteString(" ")
				builder.WriteString(field.DartName())
			}
		}
		builder.WriteString(")")
		return builder.String()
	} else if t.ResultCount == 1 {
		return t.Results.Fields[0].DartType()
	}
	return "void"
}

func (t *GoFuncType) MapName() string {
	return ""
}

func (t *GoFuncType) NeedMap() bool {
	return false
}
