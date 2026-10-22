package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cjk-press",
	Short: "A compression tool for CJK (Chinese, Japanese, Korean) text",
	Long: `cjk-press is a specialized compression tool designed for CJK text.
It uses dictionary-based compression optimized for Chinese, Japanese, and Korean characters.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("cjk-press - CJK Text Compression Tool")
		fmt.Println("Use 'cjk-press --help' to see available commands")
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
