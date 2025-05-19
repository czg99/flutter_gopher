package main

import (
	"log"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "fgo",
	Short: "Flutter Gopher - A tool for Flutter and Go integration",
	Long:  "Flutter Gopher is a comprehensive tool for creating Flutter plugins with Go backends and generating FFI bindings between Flutter and Dart.",
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

func main() {
	rootCmd.SetHelpCommand(&cobra.Command{
		Hidden: true,
	})

	createCmd.PersistentFlags().BoolP("help", "h", false, "")
	createCmd.PersistentFlags().MarkHidden("help")

	err := rootCmd.Execute()
	if err != nil {
		log.Fatalln("Error:", err)
	}
}
