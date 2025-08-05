package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	plugingen "github.com/czg99/flutter_gopher/plugin_gen"
	"github.com/spf13/cobra"
)

var withExample bool

// createCmd 创建Flutter插件的命令
var createCmd = &cobra.Command{
	Use:   "create <project_name>",
	Short: "Create a new Flutter plugin with Go backend.",
	Long: `This command generates a complete Flutter plugin project structure that enables seamless interoperability between Flutter, Go, and Platform.

Example usage:
  fgo create my_ffi
  fgo create my_ffi --example`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := validateAndGeneratePlugin(args[0]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n\n", err)
			os.Exit(1)
		}
	},
}

// isValidProjectName 检查项目名称是否合法
func isValidProjectName(name string) bool {
	if name == "" {
		return false
	}
	// 正则：只能包含字母、数字和下划线，且必须以字母开头，不能以下划线结尾
	pattern := `^[a-zA-Z][a-zA-Z0-9_]*[a-zA-Z0-9]+$`
	return regexp.MustCompile(pattern).MatchString(name)
}

// validateAndGeneratePlugin 处理输入验证和插件生成
func validateAndGeneratePlugin(projectName string) error {
	// 确保输出目录存在
	outputPath, err := filepath.Abs(projectName)
	if err != nil {
		return fmt.Errorf("failed to resolve output path: %v", err)
	}

	// 检查项目名称是否合法
	projectName = filepath.Base(outputPath)
	if !isValidProjectName(projectName) {
		return fmt.Errorf("invalid project name: %s", projectName)
	}

	// 如果输出目录不存在则创建
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		fmt.Println("Creating output directory:", outputPath)
		if err = os.MkdirAll(outputPath, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %v", err)
		}
	} else if err != nil {
		return fmt.Errorf("error accessing output directory: %v", err)
	}

	// 初始化插件生成器
	fmt.Printf("Initializing plugin generator for '%s'...\n", projectName)
	generator := plugingen.NewPluginGenerator(projectName)

	// 生成插件项目结构
	fmt.Println("Generating plugin project structure...")
	if err := generator.Generate(outputPath); err != nil {
		return fmt.Errorf("failed to generate plugin project: %v", err)
	}

	// 切换到输出目录
	if err := os.Chdir(outputPath); err != nil {
		return fmt.Errorf("failed to change directory: %v", err)
	}

	if withExample {
		fmt.Println()
		// 生成 example 应用
		if err := generator.GeneratorFlutterExample("example"); err != nil {
			return fmt.Errorf("failed to generate flutter example: %v", err)
		}
	}

	fmt.Println("\n✅ Plugin project created successfully!")
	fmt.Println("📁 Location:", outputPath)
	fmt.Println("📦 Plugin name:", projectName)
	return nil
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().BoolVar(&withExample, "example", false, "Generate example Flutter app that demonstrates the plugin usage")
}
