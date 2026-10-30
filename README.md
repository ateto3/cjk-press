# cjk-press

[![GitHub release](https://img.shields.io/github/release/atetov/cjk-press.svg?include_prereleases&sort=semver)](https://github.com/atetov/cjk-press/releases)
[![CI](https://img.shields.io/github/actions/workflow/status/atetov/cjk-press/ci.yml?branch=main)](https://github.com/atetov/cjk-press/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/actions/workflow/status/atetov/cjk-press/release.yml?branch=main)](https://github.com/atetov/cjk-press/actions/workflows/release.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/ateto3/cjk-press)](https://goreportcard.com/report/github.com/ateto3/cjk-press)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

`cjk-press` is an Open Source text pre-processor and compressor written in Go, specifically designed to optimize storage for text files containing CJK (Chinese, Japanese, and Korean) languages.

## The Problem

Traditional compressors (Gzip, Zip, DEFLATE, LZ77) were built for the Latin alphabet: 1 byte per ASCII character and spaces between words. CJK text breaks those assumptions:

* **Heavy encoding:** Hanzi, Kanji, Hiragana, Katakana, and Hangul characters usually occupy **3 bytes in UTF-8**.
* **No word delimiters:** Chinese and Japanese do not use spaces between words, so standard compressors struggle to find repeating word patterns.
* **Low short-byte entropy:** Thousands of unique characters instead of a 26-letter alphabet means identical byte sequences are rarer inside standard dictionary windows.

## The Solution

`cjk-press` is an intelligent transformation engine. Instead of processing text as blind byte streams, it analyzes the text flow using Go's native `rune` type and applies linguistic techniques:

1. **Word Segmentation:** Groups characters into actual words using open-source statistical dictionaries.
2. **Dynamic Dictionary Mapping:** Replaces repeating multi-byte character sequences (6, 9, or more bytes) with compact 1- or 2-byte identifiers.
3. **Input Optimization:** Produces an intermediate file with drastically improved repetition frequency, ready for final packaging or transmission.

## Features

* **High Performance:** Pure Go with native concurrency (goroutines) for large files.
* **Native UTF-8 Handling:** Zero risk of character corruption thanks to Go's safe rune handling.
* **Cross-platform:** Linux, macOS, and Windows binaries (amd64, arm64).
* **Docker-ready:** Official multi-arch Docker image.
* **Streaming:** Compress/decompress through stdin/stdout pipes.
* **Modular:** Usable as a standalone CLI or as a Go library (`pkg/library`).

## Installation

### Option A: One-line installer (Linux / macOS)

```bash
curl -fsSL https://raw.githubusercontent.com/atetov/cjk-press/main/install.sh | bash
```

Descarga el último release para tu arquitectura (`linux amd64`/`arm64`, `darwin`, etc.) e instala `cjk-press` en `/usr/local/bin` (o `$INSTALL_DIR`).

### Option B: From release binaries

Download the tarball/zip for your platform from [GitHub Releases](https://github.com/atetov/cjk-press/releases):

```bash
VERSION=v0.1.0
ARCH=linux-amd64   # usa darwin-amd64, windows-amd64.zip, linux-arm64, etc.
curl -fsSL -o cjk-press.tar.gz \
  "https://github.com/atetov/cjk-press/releases/download/${VERSION}/cjk-press_${VERSION#v}_${ARCH}.tar.gz"
tar -xzf cjk-press.tar.gz
sudo install -m 0755 cjk-press /usr/local/bin/
cjk-press --version
```

### Option C: From source (Go 1.21+)

```bash
go install github.com/ateto3/cjk-press/cmd/cjk-press@latest
```

### Option D: Homebrew (macOS / Linux)

```bash
brew tap atetov/homebrew-tap
brew install cjk-press
```

### Option E: Scoop (Windows)

```powershell
scoop bucket add atetov/scoop-bucket
scoop install cjk-press
```

### Option F: Docker

```bash
# Construye la imagen tú mismo (recomendado):
docker build -t cjk-press .

docker run --rm -v "$(pwd):/data" cjk-press compress /data/texto_chino.txt
docker run --rm -v "$(pwd):/data" cjk-press decompress /data/texto_chino.txt.cjk -o /data/restaurado.txt
```

La imagen se genera con **distroless** (mínima, segura) y soporta `linux/amd64` y `linux/arm64` vía Docker Buildx.

## Basic Usage

### Compress a file

```bash
# Output file will be input.cjk
cjk-press compress chinese_file.txt

# Specify output file
cjk-press compress chinese_file.txt -o compressed.cjk

# Use more workers for large files
cjk-press compress big_japanese.txt --workers 8 --block-size 20000
```

### Decompress a file

```bash
cjk-press decompress compressed.cjk
cjk-press decompress compressed.cjk -o restored_file.txt
```

### Streaming via stdin/stdout

```bash
cat texto_coreano.txt | cjk-press compress -o output.cjk
cjk-press decompress output.cjk -o -
```

### Version

```bash
cjk-press --version
cjk-press version
```

## Example with real CJK text

```bash
echo "東京は日本の首都です。北京是中国的首都。서울은 대한민국의 수도입니다." > sample.txt
cjk-press compress sample.txt -o sample.cjk
ls -lh sample.txt sample.cjk
cjk-press decompress sample.cjk -o sample_restored.txt
diff sample.txt sample_restored.txt && echo "OK: lossless round-trip"
```

## Benchmarks

Run the built-in Go benchmarks:

```bash
go test -bench=. -benchmem ./pkg/compress/ ./pkg/segment/ ./pkg/pipeline/
```

**Sugerencia para contribuir:** publica una tabla comparativa real (cjk-press vs gzip vs zstd) sobre corpus de Wikipedia en chino, japonés y coreano, midiendo ratio de compresión y tiempo. Los números convencen más que las palabras.

## Project Structure

```
cjk-press/
├── cmd/cjk-press/       # CLI (Cobra): compress, decompress, version
├── pkg/
│   ├── compress/        # Codecs, levels, file I/O
│   ├── dictionary/      # Dynamic dictionary
│   ├── segment/         # CJK text segmentation
│   ├── morphseg/        # Morphological segmentation
│   ├── pipeline/        # Concurrent processing pipeline
│   ├── api/             # HTTP API (web service)
│   ├── simd/            # SIMD optimizations
│   ├── mmap/            # Memory-mapped I/O
│   ├── crypto/          # Encryption for .cjk files
│   └── ...
├── py/cjkpress/         # Python bindings (pendiente de implementar con cgo)
├── Dockerfile           # Imagen multi-stage para Linux
├── install.sh           # Instalador de un solo comando para Linux/macOS
├── .goreleaser.yml      # Builds multiplataforma + Homebrew + Scoop
└── .github/workflows/   # CI y releases automáticas
```

## Roadmap

See [ROADMAP.md](ROADMAP.md). Current phase highlights:

* [x] Core CLI and project structure
* [ ] Real lossless compression/decompression logic (currently stubbed in `cmd/cjk-press/compress.go`)
* [x] CI/CD with GitHub Actions and GoReleaser
* [ ] Docker image published on Docker Hub / GHCR
* [ ] Python bindings (`py/cjkpress`)
* [ ] VS Code extension for `.cjk` files

## How to help spread the word (marketing)

If you want to promote cjk-press:

1. **Star el repositorio** y compártelo en redes (LinkedIn, X/Twitter, Mastodon) en español e inglés.
2. **Publica en Product Hunt, Hacker News y Reddit** (r/programming, r/linux, r/selfhosted, y comunidades de Japón/China/Corea).
3. **Escribe un artículo técnico** (dev.to, Medium, tu blog): *"¿Por qué gzip sufre comprimiendo texto CJK?"* e incluye benchmarks reales.
4. **Graba un vídeo demo corto** (YouTube Shorts, TikTok, Reels) comprimiendo un archivo japonés o chino.
5. **Añádelo a listas:** awesome-go, awesome-compression, y directorios de software libre.
6. **Etiqueta el repo:** `cjk`, `compression`, `nlp`, `chinese`, `japanese`, `korean`, `go`.
7. **Traduce el README** a otros idiomas (chino, japonés, coreano) para alcanzar a los usuarios objetivo.

## Contributing

Contributions are welcome! See [CONTRIBUTING.md](CONTRIBUTING.md). Before implementing compression logic, please check existing issues and open a proposal first.

## License

This project is licensed under the MIT License.
