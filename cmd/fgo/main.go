package main

import (
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "fgo",
	Short: "Flutter Gopher - A tool for Flutter, Go, and Platform integration.",
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
