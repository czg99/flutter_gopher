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

	"github.com/czg99/flutter_gopher/locales"
	"github.com/czg99/flutter_gopher/models"
	"github.com/nicksnyder/go-i18n/v2/i18n"
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
		return fmt.Errorf(locales.MustLocalizeMessage(&i18n.Message{
			ID:    "plugingen.target.createdir.error",
			Other: "创建目标目录失败: %w",
		}), err)
	}

	// 创建 .timestamp 文件
	if err := g.CreateTimestampFile(destDir); err != nil {
		return fmt.Errorf(locales.MustLocalizeMessage(&i18n.Message{
			ID:    "plugingen.target.createtimestamp.error",
			Other: "创建.timestamp文件失败: %w",
		}), err)
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
		return fmt.Errorf(locales.MustLocalizeMessage(&i18n.Message{
			ID:    "plugingen.target.template.error",
			Other: "处理模板文件失败: %w",
		}), err)
	}

	return nil
}

// GeneratorFlutterExample 生成一个 example 应用
func (g *PluginGenerator) GeneratorFlutterExample(destDir string) error {
	// 创建 example 目录
	if err := os.RemoveAll(destDir); err != nil {
		return fmt.Errorf(locales.MustLocalizeMessage(&i18n.Message{
			ID:    "plugingen.example.remove.error",
			Other: "删除example目录失败: %w",
		}), err)
	}

	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf(locales.MustLocalizeMessage(&i18n.Message{
			ID:    "plugingen.example.createdir.error",
			Other: "创建example目录失败: %w",
		}), err)
	}

	// 执行 flutter create 命令创建示例项目
	cmd := exec.Command("flutter", "create", ".", "--no-pub", "--offline")
	cmd.Dir = destDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Println(locales.MustLocalizeMessage(&i18n.Message{
		ID:    "plugingen.example.create.info",
		Other: "正在创建Flutter example项目...",
	}))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf(locales.MustLocalizeMessage(&i18n.Message{
			ID:    "plugingen.example.create.error",
			Other: "创建Flutter example项目失败: %w",
		}), err)
	}

	pubspecFile := filepath.Join(destDir, "pubspec.yaml")
	// 读取 pubspec.yaml 文件内容
	content, err := os.ReadFile(pubspecFile)
	if err != nil {
		return fmt.Errorf(locales.MustLocalizeMessage(&i18n.Message{
			ID:    "plugingen.example.readpubspec.error",
			Other: "读取example项目的pubspec.yaml文件失败: %w",
		}), err)
	}

	// 加入项目依赖
	projectDep := fmt.Sprintf("  %s:\n    path: ../\n\n", g.ProjectName)
	content = bytes.ReplaceAll(content, []byte("dev_dependencies:"), []byte(projectDep+"dev_dependencies:"))

	// 写入修改后的内容
	if err = os.WriteFile(pubspecFile, content, 0644); err != nil {
		return fmt.Errorf(locales.MustLocalizeMessage(&i18n.Message{
			ID:    "plugingen.example.writepubspec.error",
			Other: "写入example项目的pubspec.yaml文件失败: %w",
		}), err)
	}

	// 从模板复制 example 文件
	err = fs.WalkDir(templateFiles, "templates/example", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		return g.processTemplateFile(path, filepath.Join(destDir, "../"), d.IsDir())
	})

	if err != nil {
		return fmt.Errorf(locales.MustLocalizeMessage(&i18n.Message{
			ID:    "plugingen.example.template.error",
			Other: "处理example模板文件失败: %w",
		}), err)
	}
	return nil
}

// processTemplateFile 处理模板文件或目录
func (g *PluginGenerator) processTemplateFile(templatePath, destDir string, isDir bool) error {
	// 获取相对于 templates 目录的路径
	var relPath string
	relPath, err := filepath.Rel("templates", templatePath)
	if err != nil {
		return fmt.Errorf(locales.MustLocalizeMessage(&i18n.Message{
			ID:    "plugingen.template.relpath.error",
			Other: "获取模板文件相对路径失败: %w",
		}), err)
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
			return fmt.Errorf(locales.MustLocalizeMessage(&i18n.Message{
				ID:    "plugingen.template.createdir.error",
				Other: "创建目录 %s 失败: %w",
			}), destPath, err)
		}
		return nil
	}

	// 替换占位符
	relPath = strings.ReplaceAll(relPath, "PluginClassName", g.PluginClassName)
	relPath = strings.ReplaceAll(relPath, "ProjectName", g.ProjectName)
	relPath = strings.ReplaceAll(relPath, "LibName", g.LibName)
	relPath = strings.TrimSuffix(relPath, ".tmpl")
	relPath = filepath.FromSlash(relPath)

	log.Println(locales.MustLocalizeMessage(&i18n.Message{
		ID:    "plugingen.template.process.file",
		Other: "处理模板文件:",
	}), relPath)

	// 读取模板文件内容
	content, err := templateFiles.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf(locales.MustLocalizeMessage(&i18n.Message{
			ID:    "plugingen.template.read.error",
			Other: "读取模板文件 %s 失败: %w",
		}), relPath, err)
	}

	// 处理包含模板变量的文件
	if bytes.Contains(content, []byte("{{.")) {
		var tmpl *template.Template
		tmpl, err = template.New(relPath).Parse(string(content))
		if err != nil {
			return fmt.Errorf(locales.MustLocalizeMessage(&i18n.Message{
				ID:    "plugingen.template.parse.error",
				Other: "解析模板文件 %s 失败: %w",
			}), relPath, err)
		}

		buffer := bytes.NewBuffer(nil)
		if err = tmpl.Execute(buffer, g); err != nil {
			return fmt.Errorf(locales.MustLocalizeMessage(&i18n.Message{
				ID:    "plugingen.template.execute.error",
				Other: "执行模板文件 %s 失败: %w",
			}), relPath, err)
		}

		content = buffer.Bytes()
	}

	// 写入处理后的内容到目标文件
	destPath := filepath.Join(destDir, relPath)
	err = os.WriteFile(destPath, content, 0644)
	if err != nil {
		return fmt.Errorf(locales.MustLocalizeMessage(&i18n.Message{
			ID:    "plugingen.template.writefile.error",
			Other: "写入文件 %s 失败: %w",
		}), destPath, err)
	}
	return nil
}
