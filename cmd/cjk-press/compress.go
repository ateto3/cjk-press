package main

import (
	"fmt"

	"github.com/atetov/cjk-press/pkg/compress"

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
		concurrent, _ := cmd.Flags().GetBool("concurrent")
		workers, _ := cmd.Flags().GetInt("workers")
		blockSize, _ := cmd.Flags().GetInt("block-size")

		fmt.Printf("Compressing %s to %s...", inputFile, outputFile)
		if concurrent {
			fmt.Printf(" (concurrent, %d workers, block size %d runes)", workers, blockSize)
		}
		fmt.Println()

		var err error
		if concurrent {
			err = compress.CompressFileConcurrent(inputFile, outputFile, blockSize, workers)
		} else {
			err = compress.CompressFile(inputFile, outputFile)
		}
		if err != nil {
			fmt.Printf("compression failed: %v\n", err)
			return
		}
		fmt.Println("Compression complete")
	},
}

func init() {
	compressCmd.Flags().StringP("output", "o", "", "Output file path")
	compressCmd.Flags().Bool("concurrent", false, "Use concurrent block processing with goroutines")
	compressCmd.Flags().Int("workers", 4, "Number of worker goroutines for concurrent compression")
	compressCmd.Flags().Int("block-size", 10000, "Block size in runes for concurrent processing")
	rootCmd.AddCommand(compressCmd)
}
