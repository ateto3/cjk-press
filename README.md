# cjk-press

`cjk-press` is an Open Source text pre-processor and compressor written in Go, specifically designed to optimize storage for text files containing CJK (Chinese, Japanese, and Korean) languages.

## The Problem

Traditional compressors (such as Gzip or Zip) and their underlying algorithms (DEFLATE, LZ77) were built with the Latin alphabet in mind, where each character occupies 1 byte in ASCII and words are separated by spaces. 

In CJK languages, the scenario changes drastically:
* **Heavy encoding:** Hanzi, Kanji, Hiragana, Katakana, and Hangul characters usually occupy 3 bytes in UTF-8.
* **No word delimiters:** Languages like Chinese or Japanese do not use spaces between words, making it difficult for standard compressors to efficiently identify repeating word patterns.
* **Low short-byte entropy:** Since there are thousands of unique characters instead of a short 26-letter alphabet, identical byte sequences occur much less frequently within standard dictionary windows.

## The Solution

`cjk-press` acts as an intelligent transformation engine. Instead of processing text as blind byte streams, it analyzes the text flow using Go's native `rune` type and applies linguistic techniques:

1. **Word Segmentation:** Groups characters into actual words using open-source statistical dictionaries.
2. **Dynamic Dictionary Mapping:** Replaces repeating multi-byte character sequences (which can take up 6, 9, or more bytes) with compact 1- or 2-byte identifiers.
3. **Input Optimization:** Generates an intermediate file with a drastically optimized repetition frequency, ready for final packaging or transmission.

## Features

* **High Performance:** Written in pure Go, leveraging native concurrency via Goroutines to process massive files in parallel.
* **Native UTF-8 Handling:** Zero risk of character corruption thanks to Go's safe handling of runes.
* **Modular Architecture:** Designed as a pluggable CLI tool that can be used standalone or integrated into data and AI (NLP) pipelines.

## Installation

### From Source
To install the CLI tool from source, make sure you have Go 1.21 or higher installed:

```bash
go install github.com/atetov/cjk-press/cmd/cjk-press@latest
```

### From Release Binaries
Download the latest release from the [GitHub Releases](https://github.com/atetov/cjk-press/releases) page for your platform (Linux, macOS, or Windows).

### Using Homebrew (macOS/Linux)
```bash
brew tap atetov/homebrew-tap
brew install cjk-press
```

### Using Scoop (Windows)
```powershell
scoop bucket add atetov/scoop-bucket
scoop install cjk-press
```

### Using Docker
```bash
docker pull atetov/cjk-press:latest
docker run --rm -v $(pwd):/data atetov/cjk-press compress /data/input.txt
```

## Basic Usage

### Compress a file
```bash
# Basic compression (output file will be input.cjk)
cjk-press compress chinese_file.txt

# Specify output file
cjk-press compress chinese_file.txt -o compressed.cjk
```

### Decompress a file
```bash
# Basic decompression (output file will be input.txt)
cjk-press decompress compressed.cjk

# Specify output file
cjk-press decompress compressed.cjk -o restored_file.txt
```

### Using stdin/stdout for streaming
```bash
# Compress from stdin
cat chinese_file.txt | cjk-press compress -o output.cjk

# Decompress to stdout
cjk-press decompress compressed.cjk -o -
```

## Contributing

Contributions are welcome! If you want to improve the segmentation dictionaries or propose optimizations to the compression algorithm, please open an Issue or submit a Pull Request.

## License

This project is licensed under the MIT License.
