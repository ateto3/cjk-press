package pipeline

import (
	"sync"

	"github.com/ateto3/cjk-press/pkg/dictionary"
	"github.com/ateto3/cjk-press/pkg/segment"
)

// BlockResult contains the segmented parts and local dictionary for one block
type BlockResult struct {
	ID    int
	Parts []string
	Dict  *dictionary.Dictionary
}

// processBlock segments one block of text and builds a local dictionary of CJK terms
func processBlock(id int, text string, maxRunes int) BlockResult {
	parts := segment.Segment(text, maxRunes)
	dict := dictionary.NewDictionary()
	for _, p := range parts {
		runes := []rune(p)
		if len(runes) > 0 && segment.IsCJK(runes[0]) {
			dict.Add(p)
		}
	}
	return BlockResult{ID: id, Parts: parts, Dict: dict}
}

// ProcessText splits text into rune-based blocks, processes each block
// concurrently using a worker pool and channels, then merges all local
// dictionaries into a single global dictionary.
//
// Returns the parts for each block (in original order) and the merged
// global dictionary. An error is returned only for invalid parameters.
func ProcessText(text string, blockSizeRunes int, numWorkers int, maxRunes int) ([][]string, *dictionary.Dictionary, error) {
	if blockSizeRunes < 1 {
		blockSizeRunes = 10000
	}
	if numWorkers < 1 {
		numWorkers = 1
	}

	runes := []rune(text)
	totalBlocks := (len(runes) + blockSizeRunes - 1) / blockSizeRunes
	if totalBlocks == 0 {
		totalBlocks = 1
	}

	blocks := make([]string, totalBlocks)
	for b := 0; b < totalBlocks; b++ {
		start := b * blockSizeRunes
		end := start + blockSizeRunes
		if end > len(runes) {
			end = len(runes)
		}
		blocks[b] = string(runes[start:end])
	}

	jobs := make(chan int, totalBlocks)
	results := make(chan BlockResult, totalBlocks)

	var wg sync.WaitGroup
	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for id := range jobs {
				results <- processBlock(id, blocks[id], maxRunes)
			}
		}()
	}

	for i := 0; i < totalBlocks; i++ {
		jobs <- i
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	blockResults := make([]BlockResult, totalBlocks)
	for result := range results {
		blockResults[result.ID] = result
	}

	globalDict := dictionary.NewDictionary()
	for _, br := range blockResults {
		globalDict.Merge(br.Dict)
	}

	allParts := make([][]string, totalBlocks)
	for _, br := range blockResults {
		allParts[br.ID] = br.Parts
	}

	return allParts, globalDict, nil
}
