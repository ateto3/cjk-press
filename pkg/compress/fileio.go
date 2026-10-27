// Package compress provides efficient file I/O operations for CJK text compression.
//
// This package implements buffered readers and writers using bufio to ensure
// efficient memory management when processing large text files containing
// Chinese, Japanese, and Korean characters.
package compress

import (
	"bufio"
	"io"
	"os"
)

// FileReader wraps bufio.Reader for efficient file reading.
// It provides buffered I/O operations optimized for reading Unicode text.
type FileReader struct {
	reader *bufio.Reader
	file   *os.File
}

// NewFileReader creates a new FileReader for the given file path.
// It opens the file and wraps it with a buffered reader for efficient I/O.
//
// Parameters:
//   - path: The path to the file to read.
//
// Returns:
//   - *FileReader: A new FileReader instance.
//   - error: An error if the file cannot be opened.
//
// Example:
//
//	reader, err := compress.NewFileReader("input.txt")
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer reader.Close()
func NewFileReader(path string) (*FileReader, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return &FileReader{
		reader: bufio.NewReader(file),
		file:   file,
	}, nil
}

// ReadRune reads a single Unicode character from the file.
// This method is particularly useful for processing CJK characters which
// may be multi-byte in UTF-8 encoding.
//
// Returns:
//   - r: The rune (Unicode character) read.
//   - size: The number of bytes consumed by the rune.
//   - err: An error if the read fails or EOF is reached.
func (fr *FileReader) ReadRune() (r rune, size int, err error) {
	return fr.reader.ReadRune()
}

// ReadBytes reads until the first occurrence of delim in the input,
// returning a slice containing the data up to and including the delimiter.
//
// Parameters:
//   - delim: The delimiter byte to read until.
//
// Returns:
//   - []byte: The bytes read including the delimiter.
//   - error: An error if the read fails or EOF is reached before the delimiter.
func (fr *FileReader) ReadBytes(delim byte) ([]byte, error) {
	return fr.reader.ReadBytes(delim)
}

// ReadString reads until the first occurrence of delim in the input,
// returning a string containing the data up to and including the delimiter.
//
// Parameters:
//   - delim: The delimiter byte to read until.
//
// Returns:
//   - string: The string read including the delimiter.
//   - error: An error if the read fails or EOF is reached before the delimiter.
func (fr *FileReader) ReadString(delim byte) (string, error) {
	return fr.reader.ReadString(delim)
}

// Close closes the underlying file and releases associated resources.
// It is safe to call Close multiple times.
//
// Returns:
//   - error: An error if the file cannot be closed.
func (fr *FileReader) Close() error {
	return fr.file.Close()
}

// FileWriter wraps bufio.Writer for efficient file writing.
// It provides buffered I/O operations optimized for writing Unicode text.
type FileWriter struct {
	writer *bufio.Writer
	file   *os.File
}

// NewFileWriter creates a new FileWriter for the given file path.
// It creates the file (truncating if it exists) and wraps it with a
// buffered writer for efficient I/O.
//
// Parameters:
//   - path: The path to the file to write.
//
// Returns:
//   - *FileWriter: A new FileWriter instance.
//   - error: An error if the file cannot be created.
//
// Example:
//
//	writer, err := compress.NewFileWriter("output.txt")
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer writer.Close()
func NewFileWriter(path string) (*FileWriter, error) {
	file, err := os.Create(path)
	if err != nil {
		return nil, err
	}

	return &FileWriter{
		writer: bufio.NewWriter(file),
		file:   file,
	}, nil
}

// WriteString writes a string to the file buffer.
// The data is not flushed to disk until Flush or Close is called.
//
// Parameters:
//   - s: The string to write.
//
// Returns:
//   - int: The number of bytes written.
//   - error: An error if the write fails.
func (fw *FileWriter) WriteString(s string) (int, error) {
	return fw.writer.WriteString(s)
}

// Write writes a byte slice to the file buffer.
// The data is not flushed to disk until Flush or Close is called.
//
// Parameters:
//   - p: The byte slice to write.
//
// Returns:
//   - int: The number of bytes written.
//   - error: An error if the write fails.
func (fw *FileWriter) Write(p []byte) (int, error) {
	return fw.writer.Write(p)
}

// WriteRune writes a single Unicode character to the file buffer.
// This method handles UTF-8 encoding automatically, making it ideal
// for CJK characters.
//
// Parameters:
//   - r: The rune (Unicode character) to write.
//
// Returns:
//   - int: The number of bytes written.
//   - error: An error if the write fails.
func (fw *FileWriter) WriteRune(r rune) (int, error) {
	return fw.writer.WriteRune(r)
}

// Flush writes any buffered data to the underlying file.
// This should be called to ensure data is persisted before Close.
//
// Returns:
//   - error: An error if the flush fails.
func (fw *FileWriter) Flush() error {
	return fw.writer.Flush()
}

// Close flushes any buffered data and closes the underlying file.
// It is safe to call Close multiple times.
//
// Returns:
//   - error: An error if the flush or close fails.
func (fw *FileWriter) Close() error {
	if err := fw.writer.Flush(); err != nil {
		return err
	}
	return fw.file.Close()
}

// ReadAll reads the entire content of a file into a byte slice.
// This is a convenience function for small to medium-sized files.
// For large files, consider using FileReader for streaming.
//
// Parameters:
//   - path: The path to the file to read.
//
// Returns:
//   - []byte: The complete file contents.
//   - error: An error if the file cannot be read.
func ReadAll(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return io.ReadAll(file)
}

// WriteAll writes a byte slice to a file, creating or truncating it.
// This is a convenience function for small to medium-sized files.
// For large files, consider using FileWriter for streaming.
//
// Parameters:
//   - path: The path to the file to write.
//   - data: The data to write.
//
// Returns:
//   - error: An error if the file cannot be written.
func WriteAll(path string, data []byte) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	return err
}
