package main

import (
	"log"

	"github.com/czg99/flutter_gopher/locales"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "fgo",
	Short: locales.MustLocalizeMessage(&i18n.Message{
		ID:    "fgo.main.desc",
		Other: "Flutter Gopher - 一个 Flutter、Go、Platform 的桥接代码生成工具",
	}),
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

	ffiCmd.PersistentFlags().BoolP("help", "h", false, "")
	ffiCmd.PersistentFlags().MarkHidden("help")

	rootCmd.Flags().BoolP("help", "h", false, locales.MustLocalizeMessage(&i18n.Message{
		ID:    "fgo.main.help",
		Other: "fgo的帮助",
	}))

	if err := rootCmd.Execute(); err != nil {
		log.Fatalln("Error:", err)
	}
}
