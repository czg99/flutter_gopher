package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/czg99/flutter_gopher/locales"
	plugingen "github.com/czg99/flutter_gopher/plugin_gen"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/spf13/cobra"
)

var withExample bool

// createCmd 创建Flutter插件的命令
var createCmd = &cobra.Command{
	Use: "create <project_name>",
	Short: locales.MustLocalizeMessage(&i18n.Message{
		ID:    "fgo.create.short",
		Other: "创建一个带有Go绑定的 Flutter 插件项目",
	}),
	Long: locales.MustLocalizeMessage(&i18n.Message{
		ID: "fgo.create.long",
		Other: `此命令生成一个完整的 Flutter 插件项目结构，使 Flutter、Go、Platform 之间的数据交互变得简单

使用示例:
fgo create my_ffi
fgo create my_ffi --example`,
	}),
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := validateAndGeneratePlugin(args[0]); err != nil {
			fmt.Fprintf(os.Stderr, "\n%v", err)
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
		return fmt.Errorf(locales.MustLocalizeMessage(&i18n.Message{
			ID:    "fgo.create.resolvepath.error",
			Other: "解析输出路径失败: %w",
		}), err)
	}

	// 检查项目名称是否合法
	projectName = filepath.Base(outputPath)
	if !isValidProjectName(projectName) {
		return fmt.Errorf(locales.MustLocalizeMessage(&i18n.Message{
			ID:    "fgo.create.invalidname.error",
			Other: "无效的项目名称: %s",
		}), projectName)
	}

	// 如果输出目录不存在则创建
	if _, err = os.Stat(outputPath); os.IsNotExist(err) {
		log.Println(locales.MustLocalizeMessage(&i18n.Message{
			ID:    "fgo.create.createdir.info",
			Other: "创建输出目录:",
		}), outputPath)

		if err = os.MkdirAll(outputPath, 0755); err != nil {
			return fmt.Errorf(locales.MustLocalizeMessage(&i18n.Message{
				ID:    "fgo.create.createdir.error",
				Other: "创建输出目录失败: %w",
			}), err)
		}
	} else if err != nil {
		return fmt.Errorf(locales.MustLocalizeMessage(&i18n.Message{
			ID:    "fgo.create.accessdir.error",
			Other: "访问输出目录失败: %w",
		}), err)
	}

	// 初始化插件生成器
	log.Println(locales.MustLocalizeMessage(&i18n.Message{
		ID:    "fgo.create.initgen.info",
		Other: "初始化插件生成器:",
	}), projectName)
	generator := plugingen.NewPluginGenerator(projectName)

	// 生成插件项目结构
	log.Println(locales.MustLocalizeMessage(&i18n.Message{
		ID:    "fgo.create.genstruct.info",
		Other: "生成插件项目结构...",
	}))
	if err = generator.Generate(outputPath); err != nil {
		return fmt.Errorf(locales.MustLocalizeMessage(&i18n.Message{
			ID:    "fgo.create.genstruct.error",
			Other: "生成插件项目结构失败: %w",
		}), err)
	}

	// 切换到输出目录
	if err = os.Chdir(outputPath); err != nil {
		return fmt.Errorf(locales.MustLocalizeMessage(&i18n.Message{
			ID:    "fgo.create.chdir.error",
			Other: "切换到输出目录失败: %w",
		}), err)
	}

	// 运行自身的ffi命令
	ffiCmd.Run(nil, nil)

	if withExample {
		fmt.Println()
		// 生成 example 应用
		if err := generator.GeneratorFlutterExample("example"); err != nil {
			return fmt.Errorf(locales.MustLocalizeMessage(&i18n.Message{
				ID:    "fgo.create.genexample.error",
				Other: "生成 Flutter 示例应用失败: %w",
			}), err)
		}
	}

	fmt.Println()
	fmt.Println(locales.MustLocalizeMessage(&i18n.Message{
		ID:    "fgo.create.success.info",
		Other: "✅ 插件项目创建成功!",
	}))
	fmt.Println(locales.MustLocalizeMessage(&i18n.Message{
		ID:    "fgo.create.pluginloc.info",
		Other: "📁 项目位置:",
	}), outputPath)
	fmt.Println(locales.MustLocalizeMessage(&i18n.Message{
		ID:    "fgo.create.pluginname.info",
		Other: "📦 插件名称:",
	}), projectName)
	return nil
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().BoolVar(&withExample, "example", false, locales.MustLocalizeMessage(&i18n.Message{
		ID:    "fgo.create.example.flag",
		Other: "生成一个演示 Flutter 插件使用的示例应用",
	}))
}
