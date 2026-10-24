package main

import (
	"fmt"
	"github.com/atetov/cjk-press/pkg/compress"

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
		if err := compress.DecompressFile(inputFile, outputFile); err != nil {
			fmt.Printf("decompression failed: %v\n", err)
			return
		}
		fmt.Println("Decompression complete")
	},
}

func init() {
	decompressCmd.Flags().StringP("output", "o", "", "Output file path")
	rootCmd.AddCommand(decompressCmd)
}
