package models

import (
	"strings"

	"github.com/iancoleman/strcase"
)

type ProjectNaming struct {
	ProjectName     string // 蛇形命名的项目名（例如 "my_api"）
	PackageName     string // 插件包名（例如 "com.flutter_gopher.myapi"）
	PluginClassName string // 原生插件类名（例如 "MyApiPlugin"）
	LibClassName    string // 库的类名（例如 "MyApi"）
	LibName         string // 用于导入的库名（例如 "myapi"）
}

func NewProjectNaming(projectName string) ProjectNaming {
	snake := strcase.ToSnake(projectName)
	flatName := strings.ToLower(strcase.ToLowerCamel(projectName))
	return ProjectNaming{
		ProjectName:     snake,
		PackageName:     "com.flutter_gopher." + flatName,
		PluginClassName: strcase.ToCamel(projectName) + "Plugin",
		LibClassName:    strcase.ToCamel(projectName),
		LibName:         flatName,
	}
}
