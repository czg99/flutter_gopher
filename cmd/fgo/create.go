package main

import (
	"fmt"
	"os"
	"path/filepath"

	plugingen "github.com/czg99/flutter_gopher/plugin_gen"
	"github.com/spf13/cobra"
)

var (
	projectName string
	outputDir   string
	withExample bool
)

// createCmd 创建Flutter插件的命令
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new Flutter plugin with Go backend.",
	Long: `This command generates a complete Flutter plugin project structure that enables seamless interoperability between Flutter, Go, and Native.

Example usage:
  fgo create -n my_ffi
  fgo create -n my_ffi -o ./output --example`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := validateAndGeneratePlugin(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n\n", err)
			os.Exit(1)
		}
	},
}

// validateAndGeneratePlugin 处理输入验证和插件生成
func validateAndGeneratePlugin() error {
	// 验证项目名称
	if projectName == "" {
		return fmt.Errorf("project name is required (use -n or --name flag)")
	}

	// 如果未指定输出目录则设置默认值
	if outputDir == "" {
		outputDir = projectName
	}

	// 确保输出目录存在
	outputPath, err := filepath.Abs(outputDir)
	if err != nil {
		return fmt.Errorf("failed to resolve output path: %v", err)
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

	if withExample {
		fmt.Println("📱 Example Flutter app has been created in the 'example' subdirectory")
		fmt.Println("   Run 'cd example && flutter run' to test the plugin")
	}

	return nil
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().StringVarP(&projectName, "name", "n", "", "Plugin project name (required)")
	createCmd.Flags().StringVarP(&outputDir, "output", "o", "", "Output directory for the generated plugin project")
	createCmd.Flags().BoolVar(&withExample, "example", false, "Generate example Flutter app that demonstrates the plugin usage")

	// 标记必填标志
	createCmd.MarkFlagRequired("name")
}
