package dictionary

import (
	"sync"
)

// Dictionary represents an in-memory dictionary for CJK text compression
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

// NewDictionary creates a new empty dictionary
func NewDictionary() *Dictionary {
	return &Dictionary{
		wordToID: make(map[string]uint32),
		idToWord: make(map[uint32]string),
		nextID:   1, // Start from 1, 0 can be reserved for special cases
	}
}

// Add adds a word to the dictionary and returns its ID
// If the word already exists, returns the existing ID
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

// GetID returns the ID for a given word
// Returns 0 if the word is not in the dictionary
func (d *Dictionary) GetID(word string) uint32 {
	d.mu.RLock()
	defer d.mu.RUnlock()

	return d.wordToID[word]
}

// GetWord returns the word for a given ID
// Returns empty string if the ID is not in the dictionary
func (d *Dictionary) GetWord(id uint32) string {
	d.mu.RLock()
	defer d.mu.RUnlock()

	return d.idToWord[id]
}

// Contains checks if a word exists in the dictionary
func (d *Dictionary) Contains(word string) bool {
	d.mu.RLock()
	defer d.mu.RUnlock()

	_, exists := d.wordToID[word]
	return exists
}

// Size returns the number of entries in the dictionary
func (d *Dictionary) Size() int {
	d.mu.RLock()
	defer d.mu.RUnlock()

	return len(d.wordToID)
}

// Clear removes all entries from the dictionary
func (d *Dictionary) Clear() {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.wordToID = make(map[string]uint32)
	d.idToWord = make(map[uint32]string)
	d.nextID = 1
}

// GetAllWords returns a slice of all words in the dictionary
func (d *Dictionary) GetAllWords() []string {
	d.mu.RLock()
	defer d.mu.RUnlock()

	words := make([]string, 0, len(d.wordToID))
	for word := range d.wordToID {
		words = append(words, word)
	}
	return words
}

// Merge merges another dictionary into this one
// Returns the number of new entries added
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
