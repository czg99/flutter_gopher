package models

import (
	"strings"
	"time"

	"github.com/iancoleman/strcase"
)

type ProjectNaming struct {
	ProjectName     string // 蛇形命名的项目名（例如 "my_api"）
	PackageName     string // 插件包名（例如 "com.flutter_gopher.myapi"）
	PluginClassName string // 原生插件类名（例如 "MyApiPlugin"）
	LibClassName    string // 库的类名（例如 "MyApi"）
	LibName         string // 用于导入的库名（例如 "myapi"）
	Timestamp       int64  // 用于标识导出函数唯一性
}

func NewProjectNaming(projectName string) ProjectNaming {
	snake := strcase.ToSnake(projectName)
	flatName := strings.ToLower(strcase.ToLowerCamel(projectName))
	pkgName := "com.flutter_gopher." + flatName

	return ProjectNaming{
		ProjectName:     snake,
		PackageName:     pkgName,
		PluginClassName: strcase.ToCamel(projectName) + "Plugin",
		LibClassName:    strcase.ToCamel(projectName),
		LibName:         flatName,
		Timestamp:       time.Now().UnixMilli(),
	}
}
