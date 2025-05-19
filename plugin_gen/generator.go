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

// PluginGenerator holds configuration for Flutter plugin generation.
// It contains naming conventions and identifiers used across generated files.
type PluginGenerator struct {
	ProjectName     string // Snake case project name (e.g. "my_api")
	PackageName     string // Java package name (e.g. "com.flutter_gopher.myapi")
	PluginClassName string // Native plugin class name (e.g. "MyApiPlugin")
	DartClassName   string // Dart API class name (e.g. "MyApi")
	LibName         string // Library name for imports (e.g. "myapi")
}

// NewPluginGenerator creates a new plugin generator with standardized naming
// based on the provided project name.
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

// Generate creates a new Flutter plugin project in the specified destination directory.
// If example is true, it also creates an example Flutter app that uses the plugin.
func (g *PluginGenerator) Generate(destDir string, example bool) error {
	// Ensure target directory exists
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}

	// Walk through embedded template files
	err := fs.WalkDir(templateFiles, "templates", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip example files
		if strings.Contains(path, "example") {
			return nil
		}

		return g.processTemplateFile(path, destDir, d.IsDir())
	})

	if err != nil {
		return fmt.Errorf("failed to process template files: %w", err)
	}

	// Create example Flutter app if requested
	if example {
		if err := g.createFlutterExample(destDir); err != nil {
			return fmt.Errorf("failed to create flutter example: %w", err)
		}
	}

	return nil
}

// processTemplateFile processes a single template file or directory.
// For directories, it creates the corresponding directory in the destination.
// For files, it processes templates, replaces placeholders, and writes to destination.
func (g *PluginGenerator) processTemplateFile(path, destDir string, isDir bool) error {
	// Get path relative to templates directory
	var relPath string
	relPath, err := filepath.Rel("templates", path)
	if err != nil {
		return err
	}

	if isDir {
		// Create directory in destination
		if relPath != "." {
			destPath := filepath.Join(destDir, relPath)
			if err = os.MkdirAll(destPath, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", destPath, err)
			}
		}

		return nil
	}

	// Process file
	destPath := filepath.Join(destDir, relPath)

	// Read template file content
	content, err := templateFiles.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read template file %s: %w", relPath, err)
	}

	// Process template files (containing template variables)
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

	// Remove .tmpl extension from destination path
	destPath = strings.TrimSuffix(destPath, ".tmpl")

	// Replace placeholder filenames with actual names
	fileName := filepath.Base(destPath)
	dir := filepath.Dir(destPath)
	if strings.HasPrefix(fileName, "PluginClassName") {
		fileName = strings.Replace(fileName, "PluginClassName", g.PluginClassName, 1)
	} else if strings.HasPrefix(fileName, "ProjectName") {
		fileName = strings.Replace(fileName, "ProjectName", g.ProjectName, 1)
	}

	// Write processed content to destination file
	destPath = filepath.Join(dir, fileName)
	err = os.WriteFile(destPath, content, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file %s: %w", destPath, err)
	}
	return nil
}

// createFlutterExample creates a Flutter example app that demonstrates
// how to use the generated plugin.
func (g *PluginGenerator) createFlutterExample(destDir string) error {
	// Create example directory
	exampleDir := filepath.Join(destDir, "example")
	if err := os.RemoveAll(exampleDir); err != nil {
		return fmt.Errorf("failed to remove example directory: %w", err)
	}

	if err := os.MkdirAll(exampleDir, 0755); err != nil {
		return fmt.Errorf("failed to recreate example directory: %w", err)
	}

	// Execute flutter create command to create example project
	cmd := exec.Command("flutter", "create", ".")
	cmd.Dir = exampleDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("Creating Flutter example project...")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create flutter example project: %w", err)
	}

	// Add dependency to the main plugin project
	cmd = exec.Command("flutter", "pub", "add", g.ProjectName, "--path", "..")
	cmd.Dir = exampleDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("Adding project dependency...")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to add project dependency: %w", err)
	}

	// Copy example files from templates
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
