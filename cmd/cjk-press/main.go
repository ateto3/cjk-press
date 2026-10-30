package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// version is injected at build time via GoReleaser ldflags:
//
//	-X main.version={{.Version}} -X main.commit={{.Commit}} -X main.date={{.Date}}
var version = "dev"

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

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of cjk-press",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("cjk-press version %s\n", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func main() {
	// Handle --version / -v quickly without requiring a subcommand.
	for _, arg := range os.Args[1:] {
		if arg == "--version" || arg == "-v" {
			fmt.Printf("cjk-press version %s\n", version)
			os.Exit(0)
		}
	}
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
