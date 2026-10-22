# Project Roadmap: cjk-press

This document outlines the development roadmap for cjk-press. The project is divided into incremental phases to ensure a rock-solid core before introducing complex optimizations.

## Phase 1: Core Architecture and CLI (MVP)

Objective: Establish the Go project structure and the command-line interface.
- [x] Set up the Go module (go.mod) and directory structure (/cmd, /pkg).
- [x] Implement the CLI using the flag package or the cobra library.
- [x] Develop file readers and writers using bufio to guarantee efficient memory management.
- [x] Design the basic data structures for the in-memory dictionary (Go Maps).

## Phase 2: Segmentation and Simple CJK Mapping

Objective: Implement the logic to detect and replace repeated character sequences.
- [ ] Implement a rune-based text analyzer to identify characters within the CJK Unicode ranges (Unified Chinese, Hangul, Kana).
- [ ] Develop an initial segmentation algorithm (Max Match or based on a basic static dictionary embedded via //go:embed).
- [ ] Design the binary format for the compressed file (.cjk), defining the header, the dictionary table, and the indexed data.
- [ ] Program the symmetrical Compress() and Decompress() functions, ensuring the compression is 100% lossless.

## Phase 3: Concurrency and Compression Pipelines

Objective: Speed ​​up the processing of large files by distributing the workload.
- [ ] Implement text block segmentation using Goroutines and Channels.
- [ ] Synchronize and merge local dictionaries from each block into a single global dictionary.
- [ ] Create an automated test suite (Go unit tests and built-in benchmarks) using real-world Chinese text datasets (e.g., Chinese Wikipedia).

## Phase 4: Advanced Optimization and Language-Specific Fine-Tuning

Objective: Improve the compression ratio by integrating advanced linguistic tools.
- [ ] Integrate an advanced morphological segmenter for Chinese (e.g., porting or wrapping HMM/Viterbi-based segmenters).
- [ ] Adapt the dynamic dictionary to support the grammatical nuances of Korean (particle usage) and Japanese (mixed Kanji and Kana).
- [ ] Evaluate and implement hybrid pre-processing support (e.g., applying the CJK transformation and then piping the output into Go's standard compress/zlib or compress/gzip algorithms).