package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	bridgegen "github.com/czg99/flutter_gopher/bridge_gen"
	"github.com/spf13/cobra"
)

// generateCmd 桥接代码生成命令
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate Go and Dart FFI code from Go source files",
	Long: `Generate Go and Dart FFI code from Go source files.

This command parses Go source files and generates the corresponding FFI code
for both Go and Dart to enable Flutter-Go interoperability.

Example usage:
  fgo generate`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := validateAndProcess(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n\n", err)
			os.Exit(1)
		}
	},
}

// validateAndProcess 处理输入验证和源文件处理
func validateAndProcess() error {
	log.Println("Starting code generation process...")

	// 查找工程根目录
	projectRoot, err := findProjectRoot()
	if err != nil {
		return fmt.Errorf("failed to find project root: %v", err)
	}
	log.Println("Found project at: ", projectRoot)

	// 解析pubspec.yaml中的name字段
	projectName, err := parseProjectName(filepath.Join(projectRoot, "pubspec.yaml"))
	if err != nil {
		return fmt.Errorf("failed to parse project name: %v", err)
	}

	if err = os.Chdir(projectRoot); err != nil {
		return fmt.Errorf("failed to change directory: %v", err)
	}

	if err := bridgegen.GenerateBridgeCode("src/api", "src/api.go", "lib/"+projectName+".dart"); err != nil {
		return fmt.Errorf("failed to generate FFI code: %v", err)
	}
	return nil
}

// 解析pubspec.yaml中的name字段
func parseProjectName(pubspecPath string) (string, error) {
	// 读取pubspec.yaml文件
	data, err := os.ReadFile(pubspecPath)
	if err != nil {
		return "", fmt.Errorf("failed to read pubspec.yaml: %v", err)
	}

	// 查找name字段
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "name:") {
			return strings.TrimSpace(line[5:]), nil
		}
	}
	return "", fmt.Errorf("name field is missing in pubspec.yaml")
}

// 查找pubspec.yaml的工程目录，并且目录中存在src/api目录
func findProjectRoot() (string, error) {
	// 从当前目录开始向上查找
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %v", err)
	}

	// 检查文件是否存在的辅助函数
	fileExists := func(path string) (bool, error) {
		_, err := os.Stat(path)
		if err == nil {
			return true, nil
		}
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	// 向上查找pubspec.yaml
	for {
		// 检查pubspec.yaml是否存在
		pubspecPath := filepath.Join(currentDir, "pubspec.yaml")
		pubspecExists, err := fileExists(pubspecPath)
		if err != nil {
			return "", fmt.Errorf("error checking pubspec.yaml: %v", err)
		}

		if pubspecExists {
			// 找到pubspec.yaml，检查src/api目录是否存在
			apiDir := filepath.Join(currentDir, "src/api")
			apiDirExists, err := fileExists(apiDir)
			if err != nil {
				return "", fmt.Errorf("error checking src/api directory: %v", err)
			}

			if apiDirExists {
				// 找到符合条件的目录，返回pubspec.yaml路径
				return currentDir, nil
			}
		}

		// 向上移动一层目录
		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			// 到达根目录仍未找到
			break
		}
		currentDir = parentDir
	}

	return "", fmt.Errorf("pubspec.yaml with src/api directory not found in any parent directory")
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
