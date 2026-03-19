package models

import (
	"encoding/base32"
	"strings"

	"github.com/iancoleman/strcase"
)

type ProjectNaming struct {
	ProjectName     string // 蛇形命名的项目名（例如 "my_api"）
	PackageName     string // 插件包名（例如 "com.flutter_gopher.my_api"）
	PluginClassName string // 原生插件类名（例如 "MyApiPlugin"）
	LibClassName    string // 库的类名（例如 "MyApi"）
	LibName         string // 用于导入的库名（例如 "myapi"）
	ID              string // 用于标识导出函数唯一性

	JavaPackagePath string
	JniPackagePath  string
}

func NewProjectNaming(projectName string) ProjectNaming {
	snake := strcase.ToSnake(projectName)
	camel := strcase.ToCamel(projectName)
	pkgName := "com.flutter_gopher." + snake

	return ProjectNaming{
		ProjectName:     snake,
		PackageName:     pkgName,
		PluginClassName: camel + "Plugin",
		LibClassName:    camel,
		LibName:         strings.ToLower(camel),
		ID:              strings.ToLower(strings.TrimRight(base32.StdEncoding.EncodeToString([]byte(snake)), "=")),

		JavaPackagePath: strings.ReplaceAll(pkgName, ".", "/"),
		JniPackagePath:  strings.ReplaceAll(strings.ReplaceAll(pkgName, "_", "_1"), ".", "_"),
	}
}
