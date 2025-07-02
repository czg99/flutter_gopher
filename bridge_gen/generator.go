package bridgegen

import (
	"bytes"
	"embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/czg99/flutter_gopher/models"
)

//go:embed templates/*
var templateFiles embed.FS

// GenerateBridgeCode 为给定的源路径生成桥接代码，并将生成的代码写入指定的输出目录
// 如果代码生成失败则返回错误
func GenerateBridgeCode(srcDir, goOutDir, dartOutDir, kotlinOutDir string) error {
	// 验证输出路径
	if goOutDir == "" && dartOutDir == "" && kotlinOutDir == "" {
		return fmt.Errorf("no output path specified")
	}

	// 解析源代码
	log.Println("Parsing source code")
	parser := NewGoSrcParser()
	pkg, err := parser.Parse(srcDir)
	if err != nil {
		return fmt.Errorf("failed to parse source code: %w", err)
	}

	var goOut, dartOut, kotlinOut, goJniOut string
	if goOutDir != "" {
		goOut = filepath.Join(goOutDir, "api.go")
	}

	if dartOutDir != "" {
		dartOut = filepath.Join(dartOutDir, pkg.ProjectName+".dart")
	}

	if kotlinOutDir != "" {
		kotlinPackageDir := filepath.Join(strings.Split(pkg.PackageName, ".")...)
		kotlinOut = filepath.Join(kotlinOutDir, kotlinPackageDir, pkg.LibClassName+".kt")
		if goOut != "" {
			goJniOut = filepath.Join(goOutDir, "api_android.go")
		}
	}

	// 如果指定了输出路径则生成Go代码
	if goOut != "" {
		log.Println("Generating Go code")
		if err = NewGoGenerator(*pkg).Generate(goOut); err != nil {
			return fmt.Errorf("failed to generate Go code: %w", err)
		}
	}

	// 如果指定了输出路径则生成Dart代码
	if dartOut != "" {
		log.Println("Generating Dart code")
		if err = NewDartGenerator(*pkg).Generate(dartOut); err != nil {
			return fmt.Errorf("failed to generate Dart code: %w", err)
		}
	}

	// 如果指定了输出路径则生成Kotlin代码
	if kotlinOut != "" {
		log.Println("Generating Kotlin code")
		if err = NewKotlinGenerator(*pkg).Generate(kotlinOut); err != nil {
			return fmt.Errorf("failed to generate Kotlin code: %w", err)
		}
	}

	// 如果指定了输出路径则生成Go JNI代码
	if goJniOut != "" {
		log.Println("Generating Go JNI code")
		if err = NewGoJniGenerator(*pkg).Generate(goJniOut); err != nil {
			return fmt.Errorf("failed to generate Go JNI code: %w", err)
		}
	}
	return nil
}

// BridgeGenerator 处理Go和Dart之间的桥接代码生成
// 处理Go结构体和函数以创建FFI兼容的代码
type BridgeGenerator struct {
	models.Package
	Slices []*models.GoSliceType   // 需要桥接的切片类型
	Ptrs   []*models.GoPointerType // 需要桥接的指针类型
	Chans  []*models.GoChanType    // 需要桥接的通道类型
	Basics []*models.GoBasicType   // 需要桥接的基本类型

	generatedCode []byte // 最终生成的代码
	templatePath  string // 模板文件路径
}

// NewGoGenerator 创建一个新的Go桥接代码生成器
// 使用Go桥接模板初始化生成器
func NewGoGenerator(pkg models.Package) *BridgeGenerator {
	return &BridgeGenerator{
		Package:      pkg,
		templatePath: "templates/go_bridge.go.tmpl",
	}
}

// NewDartGenerator 创建一个新的Dart桥接代码生成器
// 使用Dart桥接模板初始化生成器
func NewDartGenerator(pkg models.Package) *BridgeGenerator {
	return &BridgeGenerator{
		Package:      pkg,
		templatePath: "templates/dart_bridge.go.tmpl",
	}
}

// NewKotlinGenerator 创建一个新的Kotlin桥接代码生成器
// 使用Kotlin桥接模板初始化生成器
func NewKotlinGenerator(pkg models.Package) *BridgeGenerator {
	return &BridgeGenerator{
		Package:      pkg,
		templatePath: "templates/kotlin_bridge.go.tmpl",
	}
}

// NewGoJniGenerator 创建一个新的Go JNI桥接代码生成器
// 使用Go JNI桥接模板初始化生成器
func NewGoJniGenerator(pkg models.Package) *BridgeGenerator {
	return &BridgeGenerator{
		Package:      pkg,
		templatePath: "templates/go_jni_bridge.go.tmpl",
	}
}

// Generate 处理模板文件并生成桥接代码
// 参数:
//   - dest: 生成代码的目标文件路径
//
// 返回生成过程中出现的错误
func (g *BridgeGenerator) Generate(dest string) error {
	// 处理包数据为代码生成做准备
	g.processSpecialTypes()

	// 读取并解析模板文件
	tmpl, err := g.parseTemplate()
	if err != nil {
		return err
	}

	// 使用生成器数据执行模板
	buffer := bytes.NewBuffer(nil)
	if err = tmpl.Execute(buffer, g); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	// 清理生成的代码，移除多余的空行
	g.generatedCode = removeExcessiveEmptyLines(buffer.Bytes())

	// 将生成的代码写入文件
	if err = g.writeToFile(dest); err != nil {
		return err
	}

	log.Println("Generated code for", dest)
	return nil
}

// parseTemplate 读取并解析带有自定义函数的模板文件
func (g *BridgeGenerator) parseTemplate() (*template.Template, error) {
	templateContent, err := templateFiles.ReadFile(g.templatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read template file: %w", err)
	}

	tmpl, err := template.New(g.templatePath).Funcs(template.FuncMap{
		"makeMap": createMapFromKeyValuePairs,
	}).Parse(string(templateContent))

	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	return tmpl, nil
}

// writeToFile 确保输出目录存在并将生成的代码写入文件
func (g *BridgeGenerator) writeToFile(dest string) error {
	// 确保输出目录存在
	dir := filepath.Dir(dest)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// 将生成的代码写入文件
	if err := os.WriteFile(dest, g.generatedCode, 0644); err != nil {
		return fmt.Errorf("failed to write generated code to file: %w", err)
	}

	return nil
}

// processSpecialTypes 处理切片、指针和通道类型
func (g *BridgeGenerator) processSpecialTypes() {
	// 从结构体和函数中收集所有特殊类型
	sliceMap, ptrMap, chanMap, basicMap := g.collectSpecialTypes()

	// 将map转换为slice以便模板处理
	g.Slices = mapToSlice(sliceMap)
	g.Ptrs = mapToSlice(ptrMap)
	g.Chans = mapToSlice(chanMap)
	g.Basics = mapToSlice(basicMap)
}

// mapToSlice 将字段的map转换为slice
func mapToSlice[T models.GoType](fieldMap map[string]T) []T {
	result := make([]T, 0, len(fieldMap))
	for _, v := range fieldMap {
		result = append(result, v)
	}
	return result
}

// collectSpecialTypes 查找结构体和函数中的所有切片、指针、通道和基本类型
func (g *BridgeGenerator) collectSpecialTypes() (sliceMap map[string]*models.GoSliceType, ptrMap map[string]*models.GoPointerType, chanMap map[string]*models.GoChanType, basicMap map[string]*models.GoBasicType) {
	sliceMap = make(map[string]*models.GoSliceType)
	ptrMap = make(map[string]*models.GoPointerType)
	chanMap = make(map[string]*models.GoChanType)
	basicMap = make(map[string]*models.GoBasicType)

	// 处理字段并查找特殊类型的辅助函数
	var processTypes func(t models.GoType)
	processTypes = func(t models.GoType) {
		switch t := t.(type) {
		case *models.GoSliceType:
			sliceMap[t.MapName()] = t
			processTypes(t.Inner)
		case *models.GoPointerType:
			ptrMap[t.MapName()] = t
			processTypes(t.Inner)
		case *models.GoChanType:
			chanMap[t.MapName()] = t
			processTypes(t.Inner)
		case *models.GoBasicType:
			basicMap[t.MapName()] = t
		}
	}

	// 处理所有结构体字段
	for _, structType := range g.Structs {
		for _, field := range structType.Fields {
			processTypes(field.Type)
		}
	}

	// 处理所有函数参数和结果
	for _, funcType := range g.Funcs {
		for _, field := range funcType.Params.Fields {
			processTypes(field.Type)
		}
		for _, field := range funcType.Results.Fields {
			processTypes(field.Type)
		}
	}
	return
}

// removeExcessiveEmptyLines 从生成的代码中移除多余的空行
func removeExcessiveEmptyLines(code []byte) []byte {
	emptyLinePattern := regexp.MustCompile(`(\r\n|\n){3,}`)
	return emptyLinePattern.ReplaceAll(code, []byte("\n\n"))
}

// createMapFromKeyValuePairs 是模板中创建map的辅助函数
// 接受可变数量的键值对参数并返回一个map
func createMapFromKeyValuePairs(values ...any) map[string]any {
	result := make(map[string]any, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		if i+1 < len(values) {
			key := values[i].(string)
			value := values[i+1]
			result[key] = value
		}
	}
	return result
}
