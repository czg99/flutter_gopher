package plugingen

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
)

//go:embed templates/*
var templateFiles embed.FS

// PluginGenerator 保存用于 Flutter 插件生成的配置信息
type PluginGenerator struct {
	ProjectName     string // 蛇形命名的项目名（例如 "my_api"）
	PackageName     string // 插件包名（例如 "com.flutter_gopher.myapi"）
	PluginClassName string // 原生插件类名（例如 "MyApiPlugin"）
	DartClassName   string // Dart API 类名（例如 "MyApi"）
	LibName         string // 用于导入的库名（例如 "myapi"）
}

// NewPluginGenerator 根据提供的项目名创建一个新的插件生成器
func NewPluginGenerator(projectName string) *PluginGenerator {
	snake := strcase.ToSnake(projectName)
	flatName := strings.ToLower(strcase.ToLowerCamel(projectName))
	return &PluginGenerator{
		ProjectName:     snake,
		PackageName:     "com.flutter_gopher." + flatName,
		PluginClassName: strcase.ToCamel(projectName) + "Plugin",
		DartClassName:   strcase.ToCamel(projectName),
		LibName:         flatName,
	}
}

// Generate 在指定的目标目录下创建一个新的 Flutter 插件项目
// 如果 example 为 true，则还会创建一个使用该插件的 example 应用
func (g *PluginGenerator) Generate(destDir string, example bool) error {
	// 确保目标目录存在
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}

	// 遍历嵌入的模板文件
	err := fs.WalkDir(templateFiles, "templates", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// 跳过 example 文件
		if strings.Contains(path, "example") {
			return nil
		}

		return g.processTemplateFile(path, destDir, d.IsDir())
	})

	if err != nil {
		return fmt.Errorf("failed to process template files: %w", err)
	}

	// 如果需要则创建 example 应用
	if example {
		if err := g.createFlutterExample(destDir); err != nil {
			return fmt.Errorf("failed to create flutter example: %w", err)
		}
	}

	return nil
}

// processTemplateFile 处理模板文件或目录
func (g *PluginGenerator) processTemplateFile(path, destDir string, isDir bool) error {
	// 获取相对于 templates 目录的路径
	var relPath string
	relPath, err := filepath.Rel("templates", path)
	if err != nil {
		return err
	}

	if isDir {
		// 在目标位置创建目录
		if relPath != "." {
			destPath := filepath.Join(destDir, relPath)
			if err = os.MkdirAll(destPath, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", destPath, err)
			}
		}

		return nil
	}

	// 处理文件
	destPath := filepath.Join(destDir, relPath)

	// 读取模板文件内容
	content, err := templateFiles.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read template file %s: %w", relPath, err)
	}

	// 处理包含模板变量的文件
	if bytes.Contains(content, []byte("{{.")) {
		var tmpl *template.Template
		tmpl, err = template.New(filepath.Base(relPath)).Parse(string(content))
		if err != nil {
			return fmt.Errorf("failed to parse template %s: %w", relPath, err)
		}

		buffer := bytes.NewBuffer(nil)
		if err = tmpl.Execute(buffer, g); err != nil {
			return fmt.Errorf("failed to execute template %s: %w", relPath, err)
		}

		content = buffer.Bytes()
	}

	// 去除目标路径中的 .tmpl 后缀
	destPath = strings.TrimSuffix(destPath, ".tmpl")

	// 替换占位符文件名为实际名称
	fileName := filepath.Base(destPath)
	dir := filepath.Dir(destPath)
	if strings.HasPrefix(fileName, "PluginClassName") {
		fileName = strings.Replace(fileName, "PluginClassName", g.PluginClassName, 1)
	} else if strings.HasPrefix(fileName, "ProjectName") {
		fileName = strings.Replace(fileName, "ProjectName", g.ProjectName, 1)
	}

	// 写入处理后的内容到目标文件
	destPath = filepath.Join(dir, fileName)
	err = os.WriteFile(destPath, content, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file %s: %w", destPath, err)
	}
	return nil
}

// createFlutterExample 创建一个 example 应用
func (g *PluginGenerator) createFlutterExample(destDir string) error {
	// 创建 example 目录
	exampleDir := filepath.Join(destDir, "example")
	if err := os.RemoveAll(exampleDir); err != nil {
		return fmt.Errorf("failed to remove example directory: %w", err)
	}

	if err := os.MkdirAll(exampleDir, 0755); err != nil {
		return fmt.Errorf("failed to recreate example directory: %w", err)
	}

	// 执行 flutter create 命令创建示例项目
	cmd := exec.Command("flutter", "create", ".")
	cmd.Dir = exampleDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("Creating Flutter example project...")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create flutter example project: %w", err)
	}

	// 添加主插件项目依赖
	cmd = exec.Command("flutter", "pub", "add", g.ProjectName, "--path", "..")
	cmd.Dir = exampleDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("Adding project dependency...")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to add project dependency: %w", err)
	}

	// 从模板复制 example 文件
	err := fs.WalkDir(templateFiles, "templates/example", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		return g.processTemplateFile(path, destDir, d.IsDir())
	})

	if err != nil {
		return fmt.Errorf("failed to process example template files: %w", err)
	}
	return nil
}
