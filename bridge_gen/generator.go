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

// GenerateBridgeCode generates bridge code for a given source path and writes the generated code to the specified output paths.
// Parameters:
//   - src: Source directory containing Go API code to be processed
//   - goOut: Output path for generated Go bridge code
//   - dartOut: Output path for generated Dart bridge code
//
// Returns an error if code generation fails.
func GenerateBridgeCode(src, goOut, dartOut string) error {
	// Handle default output paths if not specified
	if src == "src/api" && goOut == "" && dartOut == "" {
		module, _, _ := ParsePkgPath(src)
		if module != "" {
			goOut = "src/api.go"
			dartOut = "lib/" + module + ".dart"
		}
	}

	// Validate output paths
	if goOut == "" && dartOut == "" {
		return fmt.Errorf("no output path specified")
	}

	// Parse source code
	log.Println("Parsing source code")
	parser := NewGoSrcParser()
	pkg, err := parser.Parse(src)
	if err != nil {
		return fmt.Errorf("failed to parse source code: %w", err)
	}

	// Generate Go code if output path is specified
	if goOut != "" {
		log.Println("Generating Go code")
		goGenerator := NewGoGenerator()
		if err = goGenerator.Generate(goOut, pkg); err != nil {
			return fmt.Errorf("failed to generate Go code: %w", err)
		}
	}

	// Generate Dart code if output path is specified
	if dartOut != "" {
		log.Println("Generating Dart code")
		dartGenerator := NewDartGenerator()
		if err = dartGenerator.Generate(dartOut, pkg); err != nil {
			return fmt.Errorf("failed to generate Dart code: %w", err)
		}
	}

	return nil
}

// BridgeGenerator handles the generation of bridge code between Go and Dart.
// It processes Go structs and functions to create FFI-compatible code.
type BridgeGenerator struct {
	DartClassName string                 // Name of the main Dart class
	PkgPath       string                 // Go package path
	LibName       string                 // Library name for FFI imports
	Structs       []*models.GoStructType // Go struct types to be bridged
	Funcs         []*models.GoFuncType   // Go functions to be bridged
	Slices        []*models.GoField      // Generated wrapper structs for slice types
	Ptrs          []*models.GoField      // Generated wrapper structs for pointer types

	generatedCode []byte // The final generated code
	templatePath  string // Path to the template file
}

// NewGoGenerator creates a new Go bridge code generator.
// It initializes a generator with the Go bridge template.
func NewGoGenerator() *BridgeGenerator {
	return &BridgeGenerator{
		templatePath: "templates/go_bridge.go.tmpl",
	}
}

// NewDartGenerator creates a new Dart bridge code generator.
// It initializes a generator with the Dart bridge template.
func NewDartGenerator() *BridgeGenerator {
	return &BridgeGenerator{
		templatePath: "templates/dart_bridge.go.tmpl",
	}
}

// Generate processes a template file and generates bridge code.
// Parameters:
//   - dest: Destination file path for the generated code
//   - pkg: Package information containing structs and functions
//
// Returns an error if code generation fails.
func (g *BridgeGenerator) Generate(dest string, pkg *models.Package) error {
	// Process package data to prepare for code generation
	g.processPackageData(pkg)

	// Read template file
	templateContent, err := templateFiles.ReadFile(g.templatePath)
	if err != nil {
		return fmt.Errorf("failed to read template file: %w", err)
	}

	// Parse template with custom functions
	tmpl, err := template.New(g.templatePath).
		Funcs(template.FuncMap{
			"makeMap": createMapFromKeyValuePairs,
		}).
		Parse(string(templateContent))
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	// Execute template with generator data
	buffer := bytes.NewBuffer(nil)
	if err = tmpl.Execute(buffer, g); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	// Clean up the generated code by removing excessive empty lines
	g.generatedCode = removeExcessiveEmptyLines(buffer.Bytes())

	// Ensure output directory exists
	dir := filepath.Dir(dest)
	if err = os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Write generated code to file
	if err = os.WriteFile(dest, g.generatedCode, 0644); err != nil {
		return fmt.Errorf("failed to write generated code to file: %w", err)
	}

	log.Printf("Generated code for %s", dest)
	return nil
}

// processPackageData prepares the package data for code generation.
// It processes structs, functions, and special types, and sets package information.
func (g *BridgeGenerator) processPackageData(pkg *models.Package) {
	// Process struct types
	g.processStructTypes(pkg)

	// Process function types
	g.processFunctionTypes(pkg)

	// Process slice and pointer types
	g.processSliceAndPointerTypes()

	// Set package information
	g.PkgPath = pkg.PkgPath
	g.DartClassName = strcase.ToCamel(pkg.Module)
	g.LibName = strings.ToLower(strcase.ToLowerCamel(pkg.Module))
}

// processStructTypes processes struct types from the package.
// This ensures that dependent structs are processed after their dependencies.
func (g *BridgeGenerator) processStructTypes(pkg *models.Package) {
	g.Structs = pkg.Structs
}

// processFunctionTypes processes function types for code generation.
// It ensures that return values are properly named and categorized.
func (g *BridgeGenerator) processFunctionTypes(pkg *models.Package) {
	funcs := make([]*models.GoFuncType, 0, len(pkg.Funcs))
	for _, value := range pkg.Funcs {
		log.Printf(" - Processing function: %s", value.Name)
		g.processFunctionReturnValues(value)
		funcs = append(funcs, value)
	}
	g.Funcs = funcs
}

// processFunctionReturnValues ensures that function return values have proper names.
// It also identifies error return values and calculates the result count.
func (g *BridgeGenerator) processFunctionReturnValues(funcType *models.GoFuncType) {
	fields := funcType.Results.Fields
	if len(fields) == 0 {
		return
	}

	// Name unnamed return values and ensure error has a name
	for idx, field := range fields {
		// If this is the last field and it's an error type, ensure it has a name
		if idx+1 == len(fields) && isErrorType(field.Type) {
			if field.Name == "" {
				field.Name = "err"
			}
		}

		// If field has no name, give it a default name
		if field.Name == "" {
			field.Name = fmt.Sprintf("res%d", idx)
		}
	}

	// Calculate result count excluding error
	resultCount := len(fields)
	errorName := ""

	// Check if the last return value is an error
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

// isErrorType checks if a type is an error type.
// Returns true if the Go type is "error".
func isErrorType(t models.GoType) bool {
	return t.GoType() == "error"
}

// processSliceAndPointerTypes creates wrapper structs for slice and pointer types.
// These wrapper structs are used to pass slices and pointers between Go and Dart.
func (g *BridgeGenerator) processSliceAndPointerTypes() {
	sliceMap := make(map[string]*models.GoField)
	ptrMap := make(map[string]*models.GoField)

	// Process all struct fields and function parameters/results to find slice and pointer types
	g.collectSpecialTypes(sliceMap, ptrMap)

	// Convert maps to slices for template processing
	g.Slices = mapToSlice(sliceMap)
	g.Ptrs = mapToSlice(ptrMap)
}

// mapToSlice converts a map of fields to a slice.
// This is a helper function to simplify the code in processSliceAndPointerTypes.
func mapToSlice(fieldMap map[string]*models.GoField) []*models.GoField {
	result := make([]*models.GoField, 0, len(fieldMap))
	for _, v := range fieldMap {
		result = append(result, v)
	}
	return result
}

// collectSpecialTypes finds all slice and pointer types in structs and functions.
// It creates wrapper structs for each unique type.
func (g *BridgeGenerator) collectSpecialTypes(sliceMap map[string]*models.GoField, ptrMap map[string]*models.GoField) {
	// Helper function to process fields and find special types
	processFields := func(fields []*models.GoField) {
		for _, field := range fields {
			if field.IsSlice() {
				key := field.MapName()
				if _, exists := sliceMap[key]; !exists {
					sliceMap[key] = field
				}
			}

			if field.IsPtr() {
				key := field.MapName()
				if _, exists := ptrMap[key]; !exists {
					ptrMap[key] = field
				}
			}
		}
	}

	// Process all struct fields
	for _, structType := range g.Structs {
		processFields(structType.Fields)
	}

	// Process all function parameters and results
	for _, funcType := range g.Funcs {
		processFields(funcType.Params.Fields)
		processFields(funcType.Results.Fields)
	}
}

// removeExcessiveEmptyLines removes excessive empty lines from generated code
// to improve readability.
func removeExcessiveEmptyLines(code []byte) []byte {
	emptyLinePattern := regexp.MustCompile(`(\r\n|\n){3,}`)
	return emptyLinePattern.ReplaceAll(code, []byte("\n\n"))
}

// createMapFromKeyValuePairs is a helper function to create maps in templates.
// It takes a variable number of arguments as key-value pairs and returns a map.
func createMapFromKeyValuePairs(values ...any) map[string]any {
	result := make(map[string]any)
	for i := 0; i < len(values); i += 2 {
		if i+1 < len(values) {
			key := values[i].(string)
			value := values[i+1]
			result[key] = value
		}
	}
	return result
}
