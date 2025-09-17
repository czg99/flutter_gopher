package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	ffigen "github.com/czg99/flutter_gopher/ffi_gen"
	"github.com/spf13/cobra"
)

// ffiCmd 桥接代码生成命令
var ffiCmd = &cobra.Command{
	Use:   "ffi",
	Short: "Generate Go and Dart FFI code from Go source files",
	Long: `This command parses Go source files and generates the corresponding FFI code.

Example usage:
  fgo ffi`,
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

	gosrcDir := "gosrc"
	// 查找工程根目录
	projectRoot, err := findProjectRoot(gosrcDir)
	if err != nil {
		return fmt.Errorf("failed to find project root: %v", err)
	}
	log.Println("Found project at:", projectRoot)

	if err = os.Chdir(projectRoot); err != nil {
		return fmt.Errorf("failed to change directory: %v", err)
	}

	goffiDir := gosrcDir + "/ffi"
	if err := ffigen.GenerateFfiCode(goffiDir, "lib/src/ffi"); err != nil {
		return fmt.Errorf("failed to generate FFI code: %v", err)
	}
	return nil
}

// 查找pubspec.yaml的工程目录，并且目录中存在src/ffi目录
func findProjectRoot(gosrcDir string) (root string, err error) {
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
			// 找到pubspec.yaml，检查gosrc目录是否存在
			apiDir := filepath.Join(currentDir, gosrcDir)
			apiDirExists, err := fileExists(apiDir)
			if err != nil {
				return "", fmt.Errorf("error checking %s directory: %v", gosrcDir, err)
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

	return "", fmt.Errorf("pubspec.yaml with %s directory not found in any parent directory", gosrcDir)
}

func init() {
	rootCmd.AddCommand(ffiCmd)
}
