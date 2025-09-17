package plugingen

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/czg99/flutter_gopher/models"
)

//go:embed templates/*
var templateFiles embed.FS

// PluginGenerator 保存用于 Flutter 插件生成的配置信息
type PluginGenerator struct {
	models.ProjectNaming
}

// NewPluginGenerator 根据提供的项目名创建一个新的插件生成器
func NewPluginGenerator(projectName string) *PluginGenerator {
	return &PluginGenerator{
		ProjectNaming: models.NewProjectNaming(projectName),
	}
}

// Generate 在指定的目标目录下创建一个新的 Flutter 插件项目
func (g *PluginGenerator) Generate(destDir string) error {
	// 确保目标目录存在
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}

	// 创建 .timestamp 文件
	if err := g.CreateTimestampFile(destDir); err != nil {
		return fmt.Errorf("failed to create .timestamp file: %w", err)
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

	return nil
}

// GeneratorFlutterExample 生成一个 example 应用
func (g *PluginGenerator) GeneratorFlutterExample(destDir string) error {
	// 创建 example 目录
	if err := os.RemoveAll(destDir); err != nil {
		return fmt.Errorf("failed to remove example directory: %w", err)
	}

	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to recreate example directory: %w", err)
	}

	// 执行 flutter create 命令创建示例项目
	cmd := exec.Command("flutter", "create", ".", "--no-pub", "--offline")
	cmd.Dir = destDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Println("Creating Flutter example project...")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create flutter example project: %w", err)
	}

	pubspecFile := filepath.Join(destDir, "pubspec.yaml")
	// 读取 pubspec.yaml 文件内容
	content, err := os.ReadFile(pubspecFile)
	if err != nil {
		return fmt.Errorf("failed to read pubspec.yaml file: %w", err)
	}

	// 加入项目依赖
	projectDep := fmt.Sprintf("  %s:\n    path: ../\n\n", g.ProjectName)
	content = bytes.ReplaceAll(content, []byte("dev_dependencies:"), []byte(projectDep+"dev_dependencies:"))

	// 写入修改后的内容
	if err = os.WriteFile(pubspecFile, content, 0644); err != nil {
		return fmt.Errorf("failed to write pubspec.yaml file: %w", err)
	}

	// 从模板复制 example 文件
	err = fs.WalkDir(templateFiles, "templates/example", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		return g.processTemplateFile(path, filepath.Join(destDir, "../"), d.IsDir())
	})

	if err != nil {
		return fmt.Errorf("failed to process example template files: %w", err)
	}
	return nil
}

// processTemplateFile 处理模板文件或目录
func (g *PluginGenerator) processTemplateFile(templatePath, destDir string, isDir bool) error {
	// 获取相对于 templates 目录的路径
	var relPath string
	relPath, err := filepath.Rel("templates", templatePath)
	if err != nil {
		return err
	}

	if relPath == "." {
		return nil
	}

	// 处理 PackageName 占位符
	if strings.Contains(relPath, "PackageName") {
		packageDir := filepath.Join(strings.Split(g.PackageName, ".")...)
		relPath = strings.ReplaceAll(relPath, "PackageName", packageDir)
	}

	if isDir {
		// 在目标位置创建目录
		destPath := filepath.Join(destDir, relPath)
		if err = os.MkdirAll(destPath, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", destPath, err)
		}
		return nil
	}

	// 替换占位符
	relPath = strings.ReplaceAll(relPath, "PluginClassName", g.PluginClassName)
	relPath = strings.ReplaceAll(relPath, "ProjectName", g.ProjectName)
	relPath = strings.ReplaceAll(relPath, "LibName", g.LibName)
	relPath = strings.TrimSuffix(relPath, ".tmpl")
	relPath = filepath.FromSlash(relPath)

	log.Println("Processing template file:", relPath)

	// 读取模板文件内容
	content, err := templateFiles.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("failed to read template file %s: %w", relPath, err)
	}

	// 处理包含模板变量的文件
	if bytes.Contains(content, []byte("{{.")) {
		var tmpl *template.Template
		tmpl, err = template.New(relPath).Parse(string(content))
		if err != nil {
			return fmt.Errorf("failed to parse template %s: %w", relPath, err)
		}

		buffer := bytes.NewBuffer(nil)
		if err = tmpl.Execute(buffer, g); err != nil {
			return fmt.Errorf("failed to execute template %s: %w", relPath, err)
		}

		content = buffer.Bytes()
	}

	// 写入处理后的内容到目标文件
	destPath := filepath.Join(destDir, relPath)
	err = os.WriteFile(destPath, content, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file %s: %w", destPath, err)
	}
	return nil
}
