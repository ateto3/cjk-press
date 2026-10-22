package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var decompressCmd = &cobra.Command{
	Use:   "decompress [input file]",
	Short: "Decompress a .cjk file",
	Long:  `Decompress a .cjk file back to its original text format.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		inputFile := args[0]
		outputFile, _ := cmd.Flags().GetString("output")
		if outputFile == "" {
			outputFile = inputFile + ".txt"
		}
		fmt.Printf("Decompressing %s to %s...\n", inputFile, outputFile)
		fmt.Println("Decompression logic will be implemented in Phase 2")
	},
}

func init() {
	decompressCmd.Flags().StringP("output", "o", "", "Output file path")
	rootCmd.AddCommand(decompressCmd)
}
