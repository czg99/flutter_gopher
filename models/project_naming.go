package models

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/iancoleman/strcase"
)

type ProjectNaming struct {
	ProjectName     string // 蛇形命名的项目名（例如 "my_api"）
	PackageName     string // 插件包名（例如 "com.flutter_gopher.my_api"）
	PluginClassName string // 原生插件类名（例如 "MyApiPlugin"）
	LibClassName    string // 库的类名（例如 "MyApi"）
	LibName         string // 用于导入的库名（例如 "myapi"）
	Timestamp       int64  // 用于标识导出函数唯一性
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
		Timestamp:       time.Now().UnixMilli(),
	}
}

// CreateTimestampFile 创建 .timestamp 文件
func (p *ProjectNaming) CreateTimestampFile(destDir string) error {
	timestampFile := filepath.Join(destDir, ".timestamp")
	// 读取文件内容
	content, _ := os.ReadFile(timestampFile)
	if len(content) > 0 {
		timestamp, _ := strconv.ParseInt(string(content), 10, 64)
		if timestamp > 0 {
			p.Timestamp = timestamp
			return nil
		}
	}
	// 创建文件
	if err := os.WriteFile(timestampFile, fmt.Append(nil, p.Timestamp), 0644); err != nil {
		return err
	}
	return nil
}
