package main

import (
	"fmt"
	"os"
	"path/filepath"

	bridgegen "github.com/czg99/flutter_gopher/bridge_gen"
	plugingen "github.com/czg99/flutter_gopher/plugin_gen"
	"github.com/spf13/cobra"
)

var (
	projectName string
	outputDir   string
	withExample bool
)

// createCmd represents the Flutter plugin creation command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new Flutter plugin with Go backend",
	Long: `Create a new Flutter plugin with Go backend integration.

This command generates a complete Flutter plugin project structure with
all necessary Go and Dart code to enable seamless Flutter-Go interoperability.
The generated plugin includes:
  - Go API structure for implementing native functionality
  - Dart API for calling Go code from Flutter
  - Platform-specific integration code
  - Bridge code for communication between Flutter and Go

Example usage:
  fg create -n my_api
  fg create -n my_api -o ./output --example`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := validateAndGeneratePlugin(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n\n", err)
			os.Exit(1)
		}
	},
}

// validateAndGeneratePlugin handles the validation of inputs and generation of the plugin
func validateAndGeneratePlugin() error {
	// Validate project name
	if projectName == "" {
		return fmt.Errorf("project name is required (use -n or --name flag)")
	}

	// Set default output directory if not specified
	if outputDir == "" {
		outputDir = projectName
	}

	// Ensure output directory exists
	outputPath, err := filepath.Abs(outputDir)
	if err != nil {
		return fmt.Errorf("failed to resolve output path: %v", err)
	}

	// Create output directory if it doesn't exist
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		fmt.Printf("Creating output directory: %s\n", outputPath)
		if err = os.MkdirAll(outputPath, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %v", err)
		}
	} else if err != nil {
		return fmt.Errorf("error accessing output directory: %v", err)
	}

	// Create plugin generator
	fmt.Printf("Initializing plugin generator for '%s'...\n", projectName)
	generator := plugingen.NewPluginGenerator(projectName)

	// Generate plugin project
	fmt.Printf("Generating plugin project structure...\n")
	if err := generator.Generate(outputPath, withExample); err != nil {
		return fmt.Errorf("failed to generate plugin project: %v", err)
	}

	// Generate bridge code
	fmt.Printf("Generating Go-Dart bridge code...\n")
	// Change to the output directory
	if err := os.Chdir(outputPath); err != nil {
		return fmt.Errorf("failed to change directory: %v", err)
	}

	// Generate bridge code from the API directory
	if err := bridgegen.GenerateBridgeCode("src/api", "", ""); err != nil {
		return fmt.Errorf("failed to generate bridge code: %v", err)
	}

	// Success message
	fmt.Println("\n✅ Plugin project created successfully!")
	fmt.Printf("📁 Location: %s\n", outputPath)
	fmt.Printf("📦 Plugin name: %s\n", projectName)

	if withExample {
		fmt.Println("📱 Example Flutter app has been created in the 'example' subdirectory")
		fmt.Println("   Run 'cd example && flutter run' to test the plugin")
	}

	fmt.Println("\n📝 Next steps:")
	fmt.Println("  1. Implement your Go API in the 'src/api' directory")
	fmt.Println("  2. Run 'fg generate' to regenerate bridge code after API changes")
	fmt.Println("  3. Use the plugin in your Flutter app with 'flutter pub add <plugin_name> --path <plugin_path>'")

	return nil
}

func init() {
	rootCmd.AddCommand(createCmd)

	// Define command line flags
	createCmd.Flags().StringVarP(&projectName, "name", "n", "", "Plugin project name (required)")
	createCmd.Flags().StringVarP(&outputDir, "output", "o", "", "Output directory for the generated plugin project")
	createCmd.Flags().BoolVar(&withExample, "example", false, "Generate example Flutter app that demonstrates the plugin usage")

	// Mark required flags
	createCmd.MarkFlagRequired("name")
}
