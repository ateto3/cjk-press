package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var compressCmd = &cobra.Command{
	Use:   "compress [input file]",
	Short: "Compress a CJK text file",
	Long:  `Compress a text file containing CJK characters into .cjk format.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		inputFile := args[0]
		outputFile, _ := cmd.Flags().GetString("output")
		if outputFile == "" {
			outputFile = inputFile + ".cjk"
		}
		fmt.Printf("Compressing %s to %s...\n", inputFile, outputFile)
		fmt.Println("Compression logic will be implemented in Phase 2")
	},
}

func init() {
	compressCmd.Flags().StringP("output", "o", "", "Output file path")
	rootCmd.AddCommand(compressCmd)
}
