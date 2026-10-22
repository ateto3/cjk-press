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

To install the CLI tool, make sure you have Go installed and run:

```bash
go install ://github.com
```

## Basic Usage

### Compress a file
```bash
cjk-press compress -i chinese_file.txt -o file.cjk
```

### Decompress a file
```bash
cjk-press decompress -i file.cjk -o restored_file.txt
```

## Contributing

Contributions are welcome! If you want to improve the segmentation dictionaries or propose optimizations to the compression algorithm, please open an Issue or submit a Pull Request.

## License

This project is licensed under the MIT License.
