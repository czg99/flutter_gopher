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

// GoSrcParser implements the Parser interface for Go code
type GoSrcParser struct {
	// Collected AST nodes
	typeNodes []*ast.TypeSpec
	funcNodes []*ast.FuncDecl

	// Result data
	structs []*models.GoStructType
	funcs   []*models.GoFuncType
}

// NewGoSrcParser creates a new instance of GoParser
func NewGoSrcParser() *GoSrcParser {
	return &GoSrcParser{
		typeNodes: make([]*ast.TypeSpec, 0),
		funcNodes: make([]*ast.FuncDecl, 0),
		structs:   make([]*models.GoStructType, 0),
		funcs:     make([]*models.GoFuncType, 0),
	}
}

// Parse implements the Parser interface for GoParser
func (p *GoSrcParser) Parse(path string) (*models.Package, error) {
	// Validate path
	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("failed to access path: %w", err)
	}

	// Load packages
	pkgs, err := p.loadPackages(path, fileInfo)
	if err != nil {
		return nil, err
	}

	// Get package path
	module, pkgPath, err := ParsePkgPath(path)
	if err != nil {
		return nil, err
	}

	// Parse all files
	log.Println("Starting package parsing")
	if err := p.parsePackages(pkgs); err != nil {
		return nil, err
	}

	// Process collected nodes
	if err := p.processNodes(); err != nil {
		return nil, err
	}

	// Return results
	return &models.Package{
		Module:  module,
		PkgPath: pkgPath,
		Structs: p.structs,
		Funcs:   p.funcs,
	}, nil
}

// loadPackages loads Go packages from the specified path
func (p *GoSrcParser) loadPackages(path string, fileInfo os.FileInfo) ([]*packages.Package, error) {
	// Configure package loading
	config := &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles | packages.NeedImports |
			packages.NeedTypes | packages.NeedSyntax | packages.NeedTypesInfo | packages.NeedTypesSizes,
	}

	// Load packages
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

// parsePackages processes all packages and their files
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

// collectNodes collects AST nodes from a file
func (p *GoSrcParser) collectNodes(file *ast.File) error {
	for _, decl := range file.Decls {
		switch node := decl.(type) {
		case *ast.GenDecl:
			// Handle general declarations (IMPORT, CONST, TYPE, VAR)
			if err := p.handleGenDecl(node); err != nil {
				return err
			}

		case *ast.FuncDecl:
			// Collect function declarations
			p.funcNodes = append(p.funcNodes, node)

		default:
			return fmt.Errorf("unexpected declaration type: %v", reflect.TypeOf(decl))
		}
	}

	return nil
}

// handleGenDecl processes general declarations
func (p *GoSrcParser) handleGenDecl(decl *ast.GenDecl) error {
	switch decl.Tok {
	case token.IMPORT, token.CONST, token.VAR:
		// Ignore these declarations
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

// processNodes processes all collected AST nodes
func (p *GoSrcParser) processNodes() error {
	log.Println("Processing collected nodes")

	// Process types
	for _, spec := range p.typeNodes {
		if err := p.processTypeNode(spec); err != nil {
			return err
		}
	}

	// Process functions
	for _, decl := range p.funcNodes {
		if err := p.processFunctionNode(decl); err != nil {
			return err
		}
	}

	return nil
}

// processTypeNode processes a type declaration
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

// processFunctionNode processes a function declaration
func (p *GoSrcParser) processFunctionNode(funcDecl *ast.FuncDecl) error {
	// Skip methods (functions with receivers)
	if funcDecl.Recv != nil {
		return nil
	}

	name := funcDecl.Name.Name
	// Skip unexported functions
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

	p.funcs = append(p.funcs, funcType)
	return nil
}

// parseTypeExpr parses a Go type expression
func (p *GoSrcParser) parseTypeExpr(name string, expr ast.Expr) (models.GoType, error) {
	switch e := expr.(type) {
	case *ast.Ident:
		// Handle basic types and identifiers
		if basicType := models.BasicTypeMap[e.Name]; basicType != nil {
			return basicType, nil
		}

		if e.Name == "any" {
			return nil, fmt.Errorf("unsupported type: any")
		}

		return &models.GoIdentType{
			Name: e.Name,
		}, nil

	case *ast.SelectorExpr:
		return nil, fmt.Errorf("imported types are not supported: %v.%v", e.X, e.Sel)

	case *ast.StarExpr:
		// Handle pointer types
		inner, err := p.parseTypeExpr("", e.X)
		if err != nil {
			return nil, err
		}
		return &models.GoPointerType{
			Inner: inner,
		}, nil

	case *ast.ArrayType:
		// Handle slice types
		if e.Len != nil {
			return nil, fmt.Errorf("fixed-size arrays are not supported")
		}

		inner, err := p.parseTypeExpr("", e.Elt)
		if err != nil {
			return nil, err
		}
		return &models.GoSliceType{
			Inner: inner,
		}, nil

	case *ast.StructType:
		// Handle struct types
		fields, err := p.parseFields(e.Fields, true)
		if err != nil {
			return nil, err
		}
		return &models.GoStructType{
			Type: &models.GoIdentType{
				Name: name,
			},
			Fields: fields,
		}, nil

	case *ast.FuncType:
		// Handle function types
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
			Name: name,
			Params: &models.GoStructType{
				Type:   &models.GoIdentType{Name: strcase.ToLowerCamel(name) + "Params"},
				Fields: params,
			},
			Results: &models.GoStructType{
				Type:   &models.GoIdentType{Name: strcase.ToLowerCamel(name) + "Results"},
				Fields: results,
			},
		}, nil
	case *ast.ChanType:
		inner, err := p.parseTypeExpr("", e.Value)
		if err != nil {
			return nil, err
		}
		return &models.GoChanType{
			Inner: inner,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported type expression: %v (%T)", expr, expr)
	}
}

// parseFields parses field lists for structs, function parameters, and results
func (p *GoSrcParser) parseFields(list *ast.FieldList, isStruct bool) ([]*models.GoField, error) {
	if list == nil {
		return nil, nil
	}

	var fields []*models.GoField

	for _, field := range list.List {
		fieldType, err := p.parseTypeExpr("", field.Type)
		if err != nil {
			return nil, err
		}

		if field.Names == nil {
			// Anonymous field
			fields = append(fields, &models.GoField{
				Type: fieldType,
			})
			continue
		}

		for _, name := range field.Names {
			// Skip unexported struct fields
			if isStruct && unicode.IsLower(rune(name.Name[0])) {
				continue
			}

			fields = append(fields, &models.GoField{
				Name: name.Name,
				Type: fieldType,
			})
		}
	}

	return fields, nil
}
