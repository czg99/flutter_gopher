package bridgegen

import (
	"fmt"
	"go/ast"
	"go/token"
	"log"
	"os"
	"reflect"
	"unicode"

	"github.com/czg99/flutter_gopher/models"
	"github.com/iancoleman/strcase"
	"golang.org/x/tools/go/packages"
)

// GoSrcParser 实现了 Go 代码的 Parser 接口
type GoSrcParser struct {
	models.ProjectNaming

	typeNodes []*ast.TypeSpec
	funcNodes []*ast.FuncDecl

	structs []*models.GoStructType
	funcs   []*models.GoFuncType
}

// NewGoSrcParser 创建一个新的 GoParser 实例
func NewGoSrcParser() *GoSrcParser {
	return &GoSrcParser{
		typeNodes: make([]*ast.TypeSpec, 0),
		funcNodes: make([]*ast.FuncDecl, 0),
		structs:   make([]*models.GoStructType, 0),
		funcs:     make([]*models.GoFuncType, 0),
	}
}

// Parse 实现了 GoParser 的 Parser 接口
func (p *GoSrcParser) Parse(path string) (*models.Package, error) {
	// 验证路径
	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("failed to access path: %w", err)
	}

	// 获取包路径
	module, pkgPath, err := ParsePkgPath(path)
	if err != nil {
		return nil, err
	}
	p.ProjectNaming = models.NewProjectNaming(module)

	// 加载包
	pkgs, err := p.loadPackages(path, fileInfo)
	if err != nil {
		return nil, err
	}

	// 开始解析
	log.Println("Starting package parsing")
	if err := p.parsePackages(pkgs); err != nil {
		return nil, err
	}

	// 处理收集的节点
	if err := p.processNodes(); err != nil {
		return nil, err
	}

	return &models.Package{
		ProjectNaming: p.ProjectNaming,
		Module:        module,
		PkgPath:       pkgPath,
		Structs:       p.structs,
		Funcs:         p.funcs,
	}, nil
}

// loadPackages 从指定路径加载 Go 包
func (p *GoSrcParser) loadPackages(path string, fileInfo os.FileInfo) ([]*packages.Package, error) {
	// 配置包加载
	config := &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles | packages.NeedImports |
			packages.NeedTypes | packages.NeedSyntax | packages.NeedTypesInfo | packages.NeedTypesSizes,
	}

	var pkgs []*packages.Package
	var err error

	if fileInfo.IsDir() {
		config.Dir = path
		pkgs, err = packages.Load(config)
	} else {
		pkgs, err = packages.Load(config, "file="+path)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to load packages: %w", err)
	}

	if len(pkgs) == 0 {
		return nil, fmt.Errorf("no packages found in the specified path")
	}

	return pkgs, nil
}

// parsePackages 处理所有包及其文件
func (p *GoSrcParser) parsePackages(pkgs []*packages.Package) error {
	for _, pkg := range pkgs {
		log.Println(" - Processing package:", pkg.Name)

		for i, file := range pkg.CompiledGoFiles {
			log.Println("   - Analyzing file:", file)
			syntax := pkg.Syntax[i]
			if err := p.collectNodes(syntax); err != nil {
				return err
			}
		}
	}
	return nil
}

// collectNodes 从文件中收集 AST 节点
func (p *GoSrcParser) collectNodes(file *ast.File) error {
	for _, decl := range file.Decls {
		switch node := decl.(type) {
		case *ast.GenDecl:
			// 处理通用声明（IMPORT, CONST, TYPE, VAR）
			if err := p.handleGenDecl(node); err != nil {
				return err
			}

		case *ast.FuncDecl:
			// 收集函数声明
			p.funcNodes = append(p.funcNodes, node)

		default:
			return fmt.Errorf("unexpected declaration type: %v", reflect.TypeOf(decl))
		}
	}

	return nil
}

// handleGenDecl 处理通用声明
func (p *GoSrcParser) handleGenDecl(decl *ast.GenDecl) error {
	switch decl.Tok {
	case token.IMPORT, token.CONST, token.VAR:
		// 忽略这些声明
		return nil
	case token.TYPE:
		if len(decl.Specs) != 1 {
			return fmt.Errorf("expected exactly one type specification, got %d", len(decl.Specs))
		}

		typeSpec, ok := (decl.Specs[0]).(*ast.TypeSpec)
		if !ok {
			return fmt.Errorf("type specification is not *ast.TypeSpec")
		}

		p.typeNodes = append(p.typeNodes, typeSpec)
		return nil
	default:
		return fmt.Errorf("unexpected declaration token: %v", decl.Tok)
	}
}

// processNodes 处理所有收集的 AST 节点
func (p *GoSrcParser) processNodes() error {
	log.Println("Processing collected nodes")

	// 处理类型
	for _, spec := range p.typeNodes {
		if err := p.processTypeNode(spec); err != nil {
			return err
		}
	}

	// 处理函数
	for _, decl := range p.funcNodes {
		if err := p.processFunctionNode(decl); err != nil {
			return err
		}
	}

	return nil
}

// processTypeNode 处理类型声明
func (p *GoSrcParser) processTypeNode(typeSpec *ast.TypeSpec) error {
	name := typeSpec.Name.Name
	log.Println(" - Processing type:", name)

	if typeSpec.TypeParams != nil {
		return fmt.Errorf("generic types with type parameters are not supported")
	}

	goType, err := p.parseTypeExpr(name, typeSpec.Type)
	if err != nil {
		return err
	}

	structType, ok := goType.(*models.GoStructType)
	if !ok {
		return fmt.Errorf("expected struct type, got %v", reflect.TypeOf(goType))
	}

	p.structs = append(p.structs, structType)
	return nil
}

// processFunctionNode 处理函数声明
func (p *GoSrcParser) processFunctionNode(funcDecl *ast.FuncDecl) error {
	// 跳过方法（带接收者的函数）
	if funcDecl.Recv != nil {
		return nil
	}

	name := funcDecl.Name.Name
	// 跳过未导出的函数
	if unicode.IsLower(rune(name[0])) {
		return nil
	}

	log.Println(" - Processing function:", name)

	goType, err := p.parseTypeExpr(name, funcDecl.Type)
	if err != nil {
		return err
	}

	funcType, ok := goType.(*models.GoFuncType)
	if !ok {
		return fmt.Errorf("expected function type, got %v", reflect.TypeOf(goType))
	}

	// 处理函数返回值
	processFunctionReturnValues(funcType)

	p.funcs = append(p.funcs, funcType)
	return nil
}

// parseTypeExpr 解析 Go 类型表达式
func (p *GoSrcParser) parseTypeExpr(name string, expr ast.Expr) (models.GoType, error) {
	switch e := expr.(type) {
	case *ast.Ident:
		// 处理基本类型
		if basicType := models.BasicTypeMap[e.Name]; basicType != nil {
			return basicType, nil
		}

		if e.Name == "any" {
			return nil, fmt.Errorf("unsupported type: any")
		}

		return &models.GoIdentType{
			PackageName: p.PackageName,
			Name:        e.Name,
		}, nil

	case *ast.SelectorExpr:
		return nil, fmt.Errorf("imported types are not supported: %v.%v", e.X, e.Sel)

	case *ast.StarExpr:
		// 处理指针类型
		inner, err := p.parseTypeExpr("", e.X)
		if err != nil {
			return nil, err
		}

		switch inner.(type) {
		case *models.GoPointerType:
			return nil, fmt.Errorf("double pointer types are not supported")
		}

		return &models.GoPointerType{
			PackageName: p.PackageName,
			Inner:       inner,
		}, nil

	case *ast.ArrayType:
		// 处理切片类型
		if e.Len != nil {
			return nil, fmt.Errorf("fixed-size arrays are not supported")
		}

		inner, err := p.parseTypeExpr("", e.Elt)
		if err != nil {
			return nil, err
		}
		return &models.GoSliceType{
			PackageName: p.PackageName,
			Inner:       inner,
		}, nil

	case *ast.StructType:
		// 处理结构类型
		fields, err := p.parseFields(e.Fields, true)
		if err != nil {
			return nil, err
		}

		if len(fields) == 0 {
			return nil, fmt.Errorf("struct with no fields is not supported")
		}

		return &models.GoStructType{
			Type: &models.GoIdentType{
				PackageName: p.PackageName,
				Name:        name,
			},
			Fields: fields,
		}, nil

	case *ast.FuncType:
		// 处理函数类型
		if e.TypeParams != nil {
			return nil, fmt.Errorf("generic functions with type parameters are not supported")
		}

		params, err := p.parseFields(e.Params, false)
		if err != nil {
			return nil, err
		}

		results, err := p.parseFields(e.Results, false)
		if err != nil {
			return nil, err
		}

		return &models.GoFuncType{
			PackageName: p.PackageName,
			ClassName:   p.LibClassName,
			Name:        name,
			Params: &models.GoStructType{
				Type: &models.GoIdentType{
					PackageName: p.PackageName,
					Name:        strcase.ToLowerCamel(name) + "Params",
				},
				Fields: params,
			},
			Results: &models.GoStructType{
				Type: &models.GoIdentType{
					PackageName: p.PackageName,
					Name:        strcase.ToLowerCamel(name) + "Results",
				},
				Fields: results,
			},
		}, nil
	case *ast.ChanType:
		// 处理通道类型
		if e.Arrow != token.NoPos {
			return nil, fmt.Errorf("chan<- or <-chan are not supported, please use chan instead")
		}

		inner, err := p.parseTypeExpr("", e.Value)
		if err != nil {
			return nil, err
		}
		return &models.GoChanType{
			PackageName: p.PackageName,
			Inner:       inner,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported type expression: %v (%T)", expr, expr)
	}
}

// parseFields 解析结构体、函数参数和结果的字段列表
func (p *GoSrcParser) parseFields(list *ast.FieldList, isStruct bool) ([]*models.GoField, error) {
	if list == nil {
		return nil, nil
	}

	fields := make([]*models.GoField, 0, len(list.List))
	for _, field := range list.List {
		var names []string
		if field.Names == nil {
			//结构体中过滤匿名字段
			if isStruct {
				continue
			}
			names = []string{""}
		} else {
			names = make([]string, 0, len(field.Names))
			for _, name := range field.Names {
				// 跳过私有的结构体字段
				if isStruct && unicode.IsLower(rune(name.Name[0])) {
					continue
				}
				names = append(names, name.Name)
			}
		}

		if len(names) == 0 {
			continue
		}

		fieldType, err := p.parseTypeExpr("", field.Type)
		if err != nil {
			return nil, err
		}

		for _, name := range names {
			fields = append(fields, &models.GoField{
				Name: name,
				Type: fieldType,
			})
		}
	}

	return fields, nil
}

// isErrorType 检查类型是否为错误类型
// 如果Go类型是"error"则返回true
func isErrorType(t models.GoType) bool {
	return t.GoType() == "error"
}

// processFunctionReturnValues 确保函数返回值有正确的名称
// 同时识别错误返回值并计算结果数量
func processFunctionReturnValues(funcType *models.GoFuncType) {
	fields := funcType.Results.Fields
	if len(fields) == 0 {
		return
	}

	hasErr := false
	resultCount := len(fields)
	// 为未命名的返回值命名并确保错误类型有名称
	for idx, field := range fields {
		// 如果是最后一个字段且是错误类型，确保它有名称
		if idx+1 == resultCount && isErrorType(field.Type) {
			if field.Name == "" {
				field.Name = "err"
			}

			if field.Name == "err" {
				//只有错误类型名称为err时才为错误信息
				hasErr = true
			}
		}

		// 如果字段没有名称，给它一个默认名称
		if field.Name == "" {
			field.Name = fmt.Sprintf("res%d", idx)
		}
	}

	// 结果数量不含错误
	if hasErr {
		resultCount--
	}

	funcType.ResultCount = resultCount
	funcType.HasErr = hasErr
}
