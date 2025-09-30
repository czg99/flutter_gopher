package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	ffigen "github.com/czg99/flutter_gopher/ffi_gen"
	"github.com/czg99/flutter_gopher/locales"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/spf13/cobra"
)

// ffiCmd 桥接代码生成命令
var ffiCmd = &cobra.Command{
	Use: "ffi",
	Short: locales.MustLocalizeMessage(&i18n.Message{
		ID:    "fgo.ffi.short",
		Other: "解析gosrc/ffi目录并生成CGO和Dart FFI代码",
	}),
	Long: locales.MustLocalizeMessage(&i18n.Message{
		ID: "fgo.ffi.long",
		Other: `此命令解析gosrc/ffi目录的源文件并生成对应的FFI代码，使Dart可以直接调用Go函数

使用示例:
fgo ffi
`,
	}),
	Run: func(cmd *cobra.Command, args []string) {
		if err := validateAndProcess(); err != nil {
			fmt.Fprintf(os.Stderr, "\n%v", err)
			os.Exit(1)
		}
	},
}

// validateAndProcess 处理输入验证和源文件处理
func validateAndProcess() error {
	log.Println(locales.MustLocalizeMessage(&i18n.Message{
		ID:    "fgo.ffi.gen.start",
		Other: "开始生成FFI代码...",
	}))

	// 查找工程根目录
	projectRoot, err := findProjectRoot()
	if err != nil {
		return fmt.Errorf(locales.MustLocalizeMessage(&i18n.Message{
			ID:    "fgo.ffi.gen.findproject.error",
			Other: "查找项目根目录失败: %w",
		}), err)
	}

	log.Println(locales.MustLocalizeMessage(&i18n.Message{
		ID:    "fgo.ffi.gen.findproject.info",
		Other: "找到项目根目录:",
	}), projectRoot)

	if err = os.Chdir(projectRoot); err != nil {
		return fmt.Errorf(locales.MustLocalizeMessage(&i18n.Message{
			ID:    "fgo.ffi.gen.chdir.error",
			Other: "切换到项目根目录失败: %w",
		}), err)
	}

	goffiDir := "gosrc/ffi"
	if err := ffigen.GenerateFfiCode(goffiDir, "lib/src/ffi"); err != nil {
		return fmt.Errorf(locales.MustLocalizeMessage(&i18n.Message{
			ID:    "fgo.ffi.gen.error",
			Other: "生成FFI代码失败: %w",
		}), err)
	}
	return nil
}

// findProjectRoot 查找pubspec.yaml的工程目录，并且目录中存在gosrc/ffi目录
func findProjectRoot() (root string, err error) {
	// 从当前目录开始向上查找
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf(locales.MustLocalizeMessage(&i18n.Message{
			ID:    "fgo.ffi.gen.findproject.getwd.error",
			Other: "获取当前目录失败: %w",
		}), err)
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
			return "", fmt.Errorf(locales.MustLocalizeMessage(&i18n.Message{
				ID:    "fgo.ffi.gen.findproject.check.pubspec.error",
				Other: "未找到pubspec.yaml文件: %w",
			}), err)
		}

		if pubspecExists {
			// 找到pubspec.yaml，检查gosrc目录是否存在
			apiDir := filepath.Join(currentDir, "gosrc")
			apiDirExists, err := fileExists(apiDir)
			if err != nil {
				return "", fmt.Errorf(locales.MustLocalizeMessage(&i18n.Message{
					ID:    "fgo.ffi.gen.findproject.check.gosrc.error",
					Other: "未找到gosrc目录: %w",
				}), err)
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

	return "", errors.New(locales.MustLocalizeMessage(&i18n.Message{
		ID:    "fgo.ffi.gen.findproject.notfound.error",
		Other: "未找到pubspec.yaml文件与gosrc目录在任何父目录中",
	}))
}

func init() {
	rootCmd.AddCommand(ffiCmd)
}
