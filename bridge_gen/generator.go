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
	"github.com/iancoleman/strcase"
)

//go:embed templates/*
var templateFiles embed.FS

// GenerateBridgeCode 为给定的源路径生成桥接代码，并将生成的代码写入指定的输出路径
// 参数:
//   - src: 包含待处理Go API代码的源目录
//   - goOut: 生成的Go桥接代码的输出路径
//   - dartOut: 生成的Dart桥接代码的输出路径
//
// 如果代码生成失败则返回错误
func GenerateBridgeCode(src, goOut, dartOut string) error {
	// 验证输出路径
	if goOut == "" && dartOut == "" {
		return fmt.Errorf("no output path specified")
	}

	// 解析源代码
	log.Println("Parsing source code")
	parser := NewGoSrcParser()
	pkg, err := parser.Parse(src)
	if err != nil {
		return fmt.Errorf("failed to parse source code: %w", err)
	}

	// 如果指定了输出路径则生成Go代码
	if goOut != "" {
		log.Println("Generating Go code")
		if err = NewGoGenerator().Generate(goOut, pkg); err != nil {
			return fmt.Errorf("failed to generate Go code: %w", err)
		}
	}

	// 如果指定了输出路径则生成Dart代码
	if dartOut != "" {
		log.Println("Generating Dart code")
		if err = NewDartGenerator().Generate(dartOut, pkg); err != nil {
			return fmt.Errorf("failed to generate Dart code: %w", err)
		}
	}

	return nil
}

// BridgeGenerator 处理Go和Dart之间的桥接代码生成
// 处理Go结构体和函数以创建FFI兼容的代码
type BridgeGenerator struct {
	DartClassName string                  // Dart主类名
	PkgPath       string                  // Go包路径
	LibName       string                  // FFI导入的库名
	Structs       []*models.GoStructType  // 需要桥接的Go结构体类型
	Funcs         []*models.GoFuncType    // 需要桥接的Go函数
	Slices        []*models.GoSliceType   // 需要桥接的切片类型
	Ptrs          []*models.GoPointerType // 需要桥接的指针类型
	Chans         []*models.GoChanType    // 需要桥接的通道类型

	generatedCode []byte // 最终生成的代码
	templatePath  string // 模板文件路径
}

// NewGoGenerator 创建一个新的Go桥接代码生成器
// 使用Go桥接模板初始化生成器
func NewGoGenerator() *BridgeGenerator {
	return &BridgeGenerator{
		templatePath: "templates/go_bridge.go.tmpl",
	}
}

// NewDartGenerator 创建一个新的Dart桥接代码生成器
// 使用Dart桥接模板初始化生成器
func NewDartGenerator() *BridgeGenerator {
	return &BridgeGenerator{
		templatePath: "templates/dart_bridge.go.tmpl",
	}
}

// Generate 处理模板文件并生成桥接代码
// 参数:
//   - dest: 生成代码的目标文件路径
//   - pkg: 包含结构体和函数的包信息
//
// 返回生成过程中出现的错误
func (g *BridgeGenerator) Generate(dest string, pkg *models.Package) error {
	// 处理包数据为代码生成做准备
	g.processPackageData(pkg)

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

	log.Printf("Generated code for %s", dest)
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

// processPackageData 为代码生成准备包数据
// 处理结构体、函数和特殊类型，并设置包信息
func (g *BridgeGenerator) processPackageData(pkg *models.Package) {
	// 处理结构体类型
	g.Structs = pkg.Structs

	// 处理函数类型
	g.processFunctionTypes(pkg)

	// 处理特殊类型
	g.processSpecialTypes()

	// 设置包信息
	g.PkgPath = pkg.PkgPath
	g.DartClassName = strcase.ToCamel(pkg.Module)
	g.LibName = strings.ToLower(strcase.ToLowerCamel(pkg.Module))
}

// processFunctionTypes 为代码生成处理函数类型
func (g *BridgeGenerator) processFunctionTypes(pkg *models.Package) {
	funcs := make([]*models.GoFuncType, 0, len(pkg.Funcs))
	for _, value := range pkg.Funcs {
		log.Printf(" - Processing function: %s", value.Name)
		g.processFunctionReturnValues(value)
		funcs = append(funcs, value)
	}
	g.Funcs = funcs
}

// processFunctionReturnValues 确保函数返回值有正确的名称
// 同时识别错误返回值并计算结果数量
func (g *BridgeGenerator) processFunctionReturnValues(funcType *models.GoFuncType) {
	fields := funcType.Results.Fields
	if len(fields) == 0 {
		return
	}

	// 为未命名的返回值命名并确保错误有名称
	for idx, field := range fields {
		// 如果是最后一个字段且是错误类型，确保它有名称
		if idx+1 == len(fields) && isErrorType(field.Type) {
			if field.Name == "" {
				field.Name = "err"
			}
		}

		// 如果字段没有名称，给它一个默认名称
		if field.Name == "" {
			field.Name = fmt.Sprintf("res%d", idx)
		}
	}

	// 计算结果数量（不包括错误）
	resultCount := len(fields)
	errorName := ""

	// 检查最后一个返回值是否是错误
	if resultCount > 0 {
		lastField := fields[resultCount-1]
		if isErrorType(lastField.Type) {
			errorName = lastField.DartName()
			resultCount--
		}
	}

	funcType.ResultCount = resultCount
	funcType.DartErrorName = errorName
}

// isErrorType 检查类型是否为错误类型
// 如果Go类型是"error"则返回true
func isErrorType(t models.GoType) bool {
	return t.GoType() == "error"
}

// processSpecialTypes 处理切片、指针和通道类型
func (g *BridgeGenerator) processSpecialTypes() {
	// 从结构体和函数中收集所有特殊类型
	sliceMap, ptrMap, chanMap := g.collectSpecialTypes()

	// 将map转换为slice以便模板处理
	g.Slices = mapToSlice(sliceMap)
	g.Ptrs = mapToSlice(ptrMap)
	g.Chans = mapToSlice(chanMap)
}

// mapToSlice 将字段的map转换为slice
func mapToSlice[T models.GoType](fieldMap map[string]T) []T {
	result := make([]T, 0, len(fieldMap))
	for _, v := range fieldMap {
		result = append(result, v)
	}
	return result
}

// collectSpecialTypes 查找结构体和函数中的所有切片、指针和通道类型
func (g *BridgeGenerator) collectSpecialTypes() (sliceMap map[string]*models.GoSliceType, ptrMap map[string]*models.GoPointerType, chanMap map[string]*models.GoChanType) {
	sliceMap = make(map[string]*models.GoSliceType)
	ptrMap = make(map[string]*models.GoPointerType)
	chanMap = make(map[string]*models.GoChanType)

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
