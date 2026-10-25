package compress

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"unicode/utf8"

	"github.com/ateto3/cjk-press/pkg/pipeline"
	"github.com/ateto3/cjk-press/pkg/segment"
)

var magic = [4]byte{'C', 'J', 'K', '1'}
var ErrInvalidFormat = errors.New("invalid cjk file format")

// Compress transforms UTF-8 text into the versioned .cjk format.
func Compress(w io.Writer, input []byte) error {
	parts := segment.Segment(string(input), 8)
	dict := make(map[string]uint32)
	words := make([]string, 0)
	for _, p := range parts {
		if len([]rune(p)) > 0 && segment.IsCJK([]rune(p)[0]) {
			if _, ok := dict[p]; !ok {
				dict[p] = uint32(len(words) + 1)
				words = append(words, p)
			}
		}
	}
	bw := bufio.NewWriter(w)
	if _, err := bw.Write(magic[:]); err != nil {
		return err
	}
	if err := binary.Write(bw, binary.BigEndian, uint32(len(words))); err != nil {
		return err
	}
	for _, word := range words {
		b := []byte(word)
		if err := binary.Write(bw, binary.BigEndian, uint32(len(b))); err != nil {
			return err
		}
		if _, err := bw.Write(b); err != nil {
			return err
		}
	}
	if err := binary.Write(bw, binary.BigEndian, uint64(len(parts))); err != nil {
		return err
	}
	for _, p := range parts {
		if id, ok := dict[p]; ok {
			if err := bw.WriteByte(1); err != nil {
				return err
			}
			if err := binary.Write(bw, binary.BigEndian, id); err != nil {
				return err
			}
		} else {
			b := []byte(p)
			if err := bw.WriteByte(0); err != nil {
				return err
			}
			if err := binary.Write(bw, binary.BigEndian, uint32(len(b))); err != nil {
				return err
			}
			if _, err := bw.Write(b); err != nil {
				return err
			}
		}
	}
	return bw.Flush()
}

// Decompress restores the exact original UTF-8 byte sequence.
func Decompress(r io.Reader, w io.Writer) error {
	br := bufio.NewReader(r)
	var got [4]byte
	if _, err := io.ReadFull(br, got[:]); err != nil {
		return err
	}
	if got != magic {
		return ErrInvalidFormat
	}
	var n uint32
	if err := binary.Read(br, binary.BigEndian, &n); err != nil {
		return err
	}
	dict := make([]string, n+1)
	for i := uint32(1); i <= n; i++ {
		var size uint32
		if err := binary.Read(br, binary.BigEndian, &size); err != nil {
			return err
		}
		if size > 1<<30 {
			return ErrInvalidFormat
		}
		b := make([]byte, size)
		if _, err := io.ReadFull(br, b); err != nil {
			return err
		}
		if !utf8Valid(b) {
			return ErrInvalidFormat
		}
		dict[i] = string(b)
	}
	var tokens uint64
	if err := binary.Read(br, binary.BigEndian, &tokens); err != nil {
		return err
	}
	for i := uint64(0); i < tokens; i++ {
		kind, err := br.ReadByte()
		if err != nil {
			return err
		}
		switch kind {
		case 1:
			var id uint32
			if err := binary.Read(br, binary.BigEndian, &id); err != nil {
				return err
			}
			if id == 0 || id > n {
				return fmt.Errorf("%w: dictionary id %d", ErrInvalidFormat, id)
			}
			if _, err := io.WriteString(w, dict[id]); err != nil {
				return err
			}
		case 0:
			var size uint32
			if err := binary.Read(br, binary.BigEndian, &size); err != nil {
				return err
			}
			if size > 1<<30 {
				return ErrInvalidFormat
			}
			b := make([]byte, size)
			if _, err := io.ReadFull(br, b); err != nil {
				return err
			}
			if !utf8Valid(b) {
				return ErrInvalidFormat
			}
			if _, err := w.Write(b); err != nil {
				return err
			}
		default:
			return ErrInvalidFormat
		}
	}
	return nil
}

func utf8Valid(b []byte) bool { return utf8.Valid(b) }

func CompressFile(input, output string) error {
	in, err := os.Open(input)
	if err != nil {
		return err
	}
	defer in.Close()
	data, err := io.ReadAll(in)
	if err != nil {
		return err
	}
	out, err := os.Create(output)
	if err != nil {
		return err
	}
	ok := false
	defer func() {
		if !ok {
			_ = out.Close()
			_ = os.Remove(output)
		}
	}()
	if err = Compress(out, data); err != nil {
		return err
	}
	if err = out.Close(); err != nil {
		return err
	}
	ok = true
	return nil
}
func DecompressFile(input, output string) error {
	in, err := os.Open(input)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(output)
	if err != nil {
		return err
	}
	if err = Decompress(in, out); err != nil {
		_ = out.Close()
		_ = os.Remove(output)
		return err
	}
	return out.Close()
}

// defaultBlockSizeRunes is the default number of runes per block for concurrent processing
const defaultBlockSizeRunes = 10000

// defaultMaxRunes is the default max chunk size for segmentation
const defaultMaxRunes = 8

// CompressConcurrent compresses input using concurrent block processing.
// Text is split into rune-based blocks, each block is segmented and assigned
// a local dictionary in a separate goroutine, then all local dictionaries are
// merged into a single global dictionary. The output is binary-compatible with
// the sequential Compress() function and can be decompressed with Decompress().
func CompressConcurrent(w io.Writer, input []byte, blockSizeRunes int, numWorkers int) error {
	parts, dict, err := pipeline.ProcessText(string(input), blockSizeRunes, numWorkers, defaultMaxRunes)
	if err != nil {
		return err
	}

	// Build word->ID lookup from the merged global dictionary
	wordToID := make(map[string]uint32, dict.Size())
	allWords := dict.GetAllWords()
	for _, word := range allWords {
		wordToID[word] = dict.GetID(word)
	}

	bw := bufio.NewWriter(w)
	if _, err := bw.Write(magic[:]); err != nil {
		return err
	}
	if err := binary.Write(bw, binary.BigEndian, uint32(dict.Size())); err != nil {
		return err
	}

	// Write dictionary entries sorted by ID for deterministic output
	type dictEntry struct {
		id   uint32
		word string
	}
	entries := make([]dictEntry, 0, len(allWords))
	for _, word := range allWords {
		entries = append(entries, dictEntry{id: dict.GetID(word), word: word})
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].id < entries[j].id })
	for _, e := range entries {
		b := []byte(e.word)
		if err := binary.Write(bw, binary.BigEndian, uint32(len(b))); err != nil {
			return err
		}
		if _, err := bw.Write(b); err != nil {
			return err
		}
	}

	// Flatten all parts in original block order
	var allParts []string
	for _, blockParts := range parts {
		allParts = append(allParts, blockParts...)
	}

	if err := binary.Write(bw, binary.BigEndian, uint64(len(allParts))); err != nil {
		return err
	}
	for _, p := range allParts {
		if id, ok := wordToID[p]; ok {
			if err := bw.WriteByte(1); err != nil {
				return err
			}
			if err := binary.Write(bw, binary.BigEndian, id); err != nil {
				return err
			}
		} else {
			b := []byte(p)
			if err := bw.WriteByte(0); err != nil {
				return err
			}
			if err := binary.Write(bw, binary.BigEndian, uint32(len(b))); err != nil {
				return err
			}
			if _, err := bw.Write(b); err != nil {
				return err
			}
		}
	}
	return bw.Flush()
}

// CompressFileConcurrent compresses a file using concurrent block processing.
func CompressFileConcurrent(input, output string, blockSizeRunes int, numWorkers int) error {
	in, err := os.Open(input)
	if err != nil {
		return err
	}
	defer in.Close()
	data, err := io.ReadAll(in)
	if err != nil {
		return err
	}
	out, err := os.Create(output)
	if err != nil {
		return err
	}
	ok := false
	defer func() {
		if !ok {
			_ = out.Close()
			_ = os.Remove(output)
		}
	}()
	if err = CompressConcurrent(out, data, blockSizeRunes, numWorkers); err != nil {
		return err
	}
	if err = out.Close(); err != nil {
		return err
	}
	ok = true
	return nil
}

// DecompressConcurrent decompresses a .cjk file. Decompression is inherently
// sequential because the global dictionary must be read before tokens can be
// resolved, but this wrapper uses buffered I/O for efficiency and accepts a
// worker count parameter for API symmetry with CompressConcurrent.
func DecompressConcurrent(r io.Reader, w io.Writer, numWorkers int) error {
	_ = numWorkers // decompression is sequential by format design
	return Decompress(r, w)
}

// DecompressFileConcurrent decompresses a .cjk file using buffered I/O.
func DecompressFileConcurrent(input, output string, numWorkers int) error {
	in, err := os.Open(input)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(output)
	if err != nil {
		return err
	}
	if err = DecompressConcurrent(in, out, numWorkers); err != nil {
		_ = out.Close()
		_ = os.Remove(output)
		return err
	}
	return out.Close()
}

var (
	magicZlib = [4]byte{'C', 'J', 'K', 'Z'}
	magicGzip = [4]byte{'C', 'J', 'K', 'G'}
)

// CompressWithZlib applies CJK dictionary compression followed by zlib
// compression. The output is not compatible with the plain .cjk format and
// must be decompressed with DecompressWithZlib.
func CompressWithZlib(w io.Writer, input []byte) error {
	var cjkBuf bytes.Buffer
	if err := Compress(&cjkBuf, input); err != nil {
		return err
	}

	bw := bufio.NewWriter(w)
	if _, err := bw.Write(magicZlib[:]); err != nil {
		return err
	}
	zw := zlib.NewWriter(bw)
	if _, err := zw.Write(cjkBuf.Bytes()); err != nil {
		return err
	}
	if err := zw.Close(); err != nil {
		return err
	}
	return bw.Flush()
}

// DecompressWithZlib reverses CompressWithZlib: first zlib-decompresses,
// then applies CJK dictionary decompression.
func DecompressWithZlib(r io.Reader, w io.Writer) error {
	br := bufio.NewReader(r)
	var got [4]byte
	if _, err := io.ReadFull(br, got[:]); err != nil {
		return err
	}
	if got != magicZlib {
		return ErrInvalidFormat
	}
	zr, err := zlib.NewReader(br)
	if err != nil {
		return err
	}
	defer zr.Close()

	var cjkBuf bytes.Buffer
	if _, err := cjkBuf.ReadFrom(zr); err != nil {
		return err
	}
	return Decompress(&cjkBuf, w)
}

// CompressWithGzip applies CJK dictionary compression followed by gzip
// compression. The output is not compatible with the plain .cjk format and
// must be decompressed with DecompressWithGzip.
func CompressWithGzip(w io.Writer, input []byte) error {
	var cjkBuf bytes.Buffer
	if err := Compress(&cjkBuf, input); err != nil {
		return err
	}

	bw := bufio.NewWriter(w)
	if _, err := bw.Write(magicGzip[:]); err != nil {
		return err
	}
	gw := gzip.NewWriter(bw)
	if _, err := gw.Write(cjkBuf.Bytes()); err != nil {
		return err
	}
	if err := gw.Close(); err != nil {
		return err
	}
	return bw.Flush()
}

// DecompressWithGzip reverses CompressWithGzip: first gzip-decompresses,
// then applies CJK dictionary decompression.
func DecompressWithGzip(r io.Reader, w io.Writer) error {
	br := bufio.NewReader(r)
	var got [4]byte
	if _, err := io.ReadFull(br, got[:]); err != nil {
		return err
	}
	if got != magicGzip {
		return ErrInvalidFormat
	}
	gr, err := gzip.NewReader(br)
	if err != nil {
		return err
	}
	defer gr.Close()

	var cjkBuf bytes.Buffer
	if _, err := cjkBuf.ReadFrom(gr); err != nil {
		return err
	}
	return Decompress(&cjkBuf, w)
}

// CompressFileWithZlib compresses a file using CJK + zlib hybrid compression.
func CompressFileWithZlib(input, output string) error {
	in, err := os.Open(input)
	if err != nil {
		return err
	}
	defer in.Close()
	data, err := io.ReadAll(in)
	if err != nil {
		return err
	}
	out, err := os.Create(output)
	if err != nil {
		return err
	}
	ok := false
	defer func() {
		if !ok {
			_ = out.Close()
			_ = os.Remove(output)
		}
	}()
	if err = CompressWithZlib(out, data); err != nil {
		return err
	}
	if err = out.Close(); err != nil {
		return err
	}
	ok = true
	return nil
}

// DecompressFileWithZlib decompresses a file compressed with CompressWithZlib.
func DecompressFileWithZlib(input, output string) error {
	in, err := os.Open(input)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(output)
	if err != nil {
		return err
	}
	if err = DecompressWithZlib(in, out); err != nil {
		_ = out.Close()
		_ = os.Remove(output)
		return err
	}
	return out.Close()
}

// CompressFileWithGzip compresses a file using CJK + gzip hybrid compression.
func CompressFileWithGzip(input, output string) error {
	in, err := os.Open(input)
	if err != nil {
		return err
	}
	defer in.Close()
	data, err := io.ReadAll(in)
	if err != nil {
		return err
	}
	out, err := os.Create(output)
	if err != nil {
		return err
	}
	ok := false
	defer func() {
		if !ok {
			_ = out.Close()
			_ = os.Remove(output)
		}
	}()
	if err = CompressWithGzip(out, data); err != nil {
		return err
	}
	if err = out.Close(); err != nil {
		return err
	}
	ok = true
	return nil
}

// DecompressFileWithGzip decompresses a file compressed with CompressWithGzip.
func DecompressFileWithGzip(input, output string) error {
	in, err := os.Open(input)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(output)
	if err != nil {
		return err
	}
	if err = DecompressWithGzip(in, out); err != nil {
		_ = out.Close()
		_ = os.Remove(output)
		return err
	}
	return out.Close()
}
