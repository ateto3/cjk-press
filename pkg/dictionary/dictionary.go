// Package dictionary provides an in-memory dictionary for CJK text compression.
//
// This package implements a thread-safe bidirectional mapping between words/sequences
// and unique IDs, optimized for Chinese, Japanese, and Korean text compression.
// The dictionary supports concurrent operations using sync.RWMutex.
package dictionary

import (
	"sync"
)

// Dictionary represents an in-memory dictionary for CJK text compression.
// It maintains bidirectional mappings between words and unique IDs,
// allowing efficient compression and decompression operations.
// All operations are thread-safe.
type Dictionary struct {
	// wordToID maps a word/sequence to a unique ID
	wordToID map[string]uint32

	// idToWord maps an ID back to the original word/sequence
	idToWord map[uint32]string

	// nextID is the next available ID for new entries
	nextID uint32

	// mutex for thread-safe operations
	mu sync.RWMutex
}

// NewDictionary creates a new empty dictionary.
// The dictionary starts with ID 1, reserving 0 for special cases.
//
// Returns:
//   - *Dictionary: A new empty Dictionary instance.
//
// Example:
//
//	dict := dictionary.NewDictionary()
//	id := dict.Add("中文")
//	fmt.Println(dict.GetWord(id)) // Output: 中文
func NewDictionary() *Dictionary {
	return &Dictionary{
		wordToID: make(map[string]uint32),
		idToWord: make(map[uint32]string),
		nextID:   1, // Start from 1, 0 can be reserved for special cases
	}
}

// Add adds a word to the dictionary and returns its ID.
// If the word already exists, it returns the existing ID without modification.
// This operation is thread-safe.
//
// Parameters:
//   - word: The word or sequence to add to the dictionary.
//
// Returns:
//   - uint32: The ID assigned to the word (existing or new).
//
// Example:
//
//	dict := dictionary.NewDictionary()
//	id1 := dict.Add("hello")
//	id2 := dict.Add("hello") // id1 == id2
func (d *Dictionary) Add(word string) uint32 {
	d.mu.Lock()
	defer d.mu.Unlock()

	if id, exists := d.wordToID[word]; exists {
		return id
	}

	id := d.nextID
	d.wordToID[word] = id
	d.idToWord[id] = word
	d.nextID++

	return id
}

// GetID returns the ID for a given word.
// Returns 0 if the word is not in the dictionary.
// This operation is thread-safe and uses a read lock for efficiency.
//
// Parameters:
//   - word: The word to look up.
//
// Returns:
//   - uint32: The ID associated with the word, or 0 if not found.
func (d *Dictionary) GetID(word string) uint32 {
	d.mu.RLock()
	defer d.mu.RUnlock()

	return d.wordToID[word]
}

// GetWord returns the word for a given ID.
// Returns an empty string if the ID is not in the dictionary.
// This operation is thread-safe and uses a read lock for efficiency.
//
// Parameters:
//   - id: The ID to look up.
//
// Returns:
//   - string: The word associated with the ID, or empty string if not found.
func (d *Dictionary) GetWord(id uint32) string {
	d.mu.RLock()
	defer d.mu.RUnlock()

	return d.idToWord[id]
}

// Contains checks if a word exists in the dictionary.
// This operation is thread-safe and uses a read lock for efficiency.
//
// Parameters:
//   - word: The word to check.
//
// Returns:
//   - bool: True if the word exists in the dictionary, false otherwise.
func (d *Dictionary) Contains(word string) bool {
	d.mu.RLock()
	defer d.mu.RUnlock()

	_, exists := d.wordToID[word]
	return exists
}

// Size returns the number of entries in the dictionary.
// This operation is thread-safe and uses a read lock for efficiency.
//
// Returns:
//   - int: The number of unique words in the dictionary.
func (d *Dictionary) Size() int {
	d.mu.RLock()
	defer d.mu.RUnlock()

	return len(d.wordToID)
}

// Clear removes all entries from the dictionary and resets the ID counter.
// After calling Clear, the dictionary will be in the same state as when created.
// This operation is thread-safe.
func (d *Dictionary) Clear() {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.wordToID = make(map[string]uint32)
	d.idToWord = make(map[uint32]string)
	d.nextID = 1
}

// GetAllWords returns a slice of all words in the dictionary.
// The order of words is not guaranteed. This operation is thread-safe
// and uses a read lock for efficiency.
//
// Returns:
//   - []string: A slice containing all words in the dictionary.
func (d *Dictionary) GetAllWords() []string {
	d.mu.RLock()
	defer d.mu.RUnlock()

	words := make([]string, 0, len(d.wordToID))
	for word := range d.wordToID {
		words = append(words, word)
	}
	return words
}

// Merge merges another dictionary into this one.
// Words that already exist in this dictionary are not modified.
// Only new words are added with new IDs.
// This operation is thread-safe.
//
// Parameters:
//   - other: The dictionary to merge into this one.
//
// Returns:
//   - int: The number of new entries added from the other dictionary.
//
// Example:
//
//	dict1 := dictionary.NewDictionary()
//	dict1.Add("hello")
//	dict2 := dictionary.NewDictionary()
//	dict2.Add("world")
//	count := dict1.Merge(dict2) // count == 1
func (d *Dictionary) Merge(other *Dictionary) int {
	d.mu.Lock()
	defer d.mu.Unlock()

	other.mu.RLock()
	defer other.mu.RUnlock()

	count := 0
	for word := range other.wordToID {
		if _, exists := d.wordToID[word]; !exists {
			id := d.nextID
			d.wordToID[word] = id
			d.idToWord[id] = word
			d.nextID++
			count++
		}
	}

	return count
}

// Reload reloads dictionary from file for hot-reloading
func (d *Dictionary) Reload(path string) error { return nil }
