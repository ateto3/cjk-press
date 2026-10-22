package compress

import (
	"bufio"
	"io"
	"os"
)

// FileReader wraps bufio.Reader for efficient file reading
type FileReader struct {
	reader *bufio.Reader
	file   *os.File
}

// NewFileReader creates a new FileReader for the given file path
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

// ReadRune reads a single Unicode character from the file
func (fr *FileReader) ReadRune() (r rune, size int, err error) {
	return fr.reader.ReadRune()
}

// ReadBytes reads a slice of bytes from the file
func (fr *FileReader) ReadBytes(delim byte) ([]byte, error) {
	return fr.reader.ReadBytes(delim)
}

// ReadString reads a string until the delimiter
func (fr *FileReader) ReadString(delim byte) (string, error) {
	return fr.reader.ReadString(delim)
}

// Close closes the underlying file
func (fr *FileReader) Close() error {
	return fr.file.Close()
}

// FileWriter wraps bufio.Writer for efficient file writing
type FileWriter struct {
	writer *bufio.Writer
	file   *os.File
}

// NewFileWriter creates a new FileWriter for the given file path
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

// WriteString writes a string to the file
func (fw *FileWriter) WriteString(s string) (int, error) {
	return fw.writer.WriteString(s)
}

// Write writes a byte slice to the file
func (fw *FileWriter) Write(p []byte) (int, error) {
	return fw.writer.Write(p)
}

// WriteRune writes a single Unicode character to the file
func (fw *FileWriter) WriteRune(r rune) (int, error) {
	return fw.writer.WriteRune(r)
}

// Flush flushes the buffer to ensure all data is written
func (fw *FileWriter) Flush() error {
	return fw.writer.Flush()
}

// Close flushes and closes the underlying file
func (fw *FileWriter) Close() error {
	if err := fw.writer.Flush(); err != nil {
		return err
	}
	return fw.file.Close()
}

// ReadAll reads the entire content of a file into a byte slice
func ReadAll(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return io.ReadAll(file)
}

// WriteAll writes a byte slice to a file
func WriteAll(path string, data []byte) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	return err
}
