package main

import (
	"fmt"
	"os"
	"path/filepath"

	bridgegen "github.com/czg99/flutter_gopher/bridge_gen"
	"github.com/spf13/cobra"
)

var (
	srcPath     string
	goOutPath   string
	dartOutPath string
)

// generateCmd 桥接代码生成命令
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate Go and Dart FFI code from Go source files",
	Long: `Generate Go and Dart FFI code from Go source files.

This command parses Go source files and generates the corresponding FFI code
for both Go and Dart to enable Flutter-Go interoperability.

Example usage:
  fgo generate -s src/api -g output_go.go -d output_dart.dart`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := validateAndProcess(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n\n", err)
			os.Exit(1)
		}
	},
}

// validateAndProcess 处理输入验证和源文件处理
func validateAndProcess() error {
	// 如果未指定源路径，则使用默认路径
	if srcPath == "" {
		srcPath = "src/api"
	}

	// 解析绝对路径
	absPath, err := filepath.Abs(srcPath)
	if err != nil {
		return fmt.Errorf("failed to resolve source path: %v", err)
	}

	// 检查源路径是否存在
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return fmt.Errorf("source path does not exist: %s", absPath)
	} else if err != nil {
		return fmt.Errorf("error accessing source path: %v", err)
	}

	if err := bridgegen.GenerateBridgeCode(srcPath, goOutPath, dartOutPath); err != nil {
		return fmt.Errorf("failed to generate FFI code: %v", err)
	}
	return nil
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().StringVarP(&srcPath, "src", "s", "src/api", "Source path containing Go API files to parse")
	generateCmd.Flags().StringVarP(&goOutPath, "go_out", "g", "", "Output path for generated Go FFI code")
	generateCmd.Flags().StringVarP(&dartOutPath, "dart_out", "d", "", "Output path for generated Dart FFI code")
}
