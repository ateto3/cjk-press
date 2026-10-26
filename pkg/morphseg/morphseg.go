package morphseg

import (
	"math"
	"unicode"
)

// Tag represents a BMES segmentation tag.
type Tag int

const (
	Begin Tag = iota
	Middle
	End
	Single
)

func (t Tag) String() string {
	switch t {
	case Begin:
		return "B"
	case Middle:
		return "M"
	case End:
		return "E"
	case Single:
		return "S"
	default:
		return "?"
	}
}

// Segmenter is an HMM-based Chinese word segmenter using BMES tagging
// and the Viterbi algorithm for optimal segmentation.
type Segmenter struct {
	emission  [4]map[rune]float64
	transition [4][4]float64
	start     [4]float64
}

// NewSegmenter creates a segmenter with default probabilities derived from
// character statistics. These priors work without training.
func NewSegmenter() *Segmenter {
	s := &Segmenter{
		emission: [4]map[rune]float64{
			Begin:   make(map[rune]float64),
			Middle:  make(map[rune]float64),
			End:     make(map[rune]float64),
			Single:  make(map[rune]float64),
		},
	}

	s.setTransition(Begin, Begin, math.Log(0.001))
	s.setTransition(Begin, Middle, math.Log(0.4))
	s.setTransition(Begin, End, math.Log(0.599))
	s.setTransition(Begin, Single, math.Log(0.0))

	s.setTransition(Middle, Begin, math.Log(0.001))
	s.setTransition(Middle, Middle, math.Log(0.3))
	s.setTransition(Middle, End, math.Log(0.699))
	s.setTransition(Middle, Single, math.Log(0.0))

	s.setTransition(End, Begin, math.Log(0.35))
	s.setTransition(End, Middle, math.Log(0.0))
	s.setTransition(End, End, math.Log(0.0))
	s.setTransition(End, Single, math.Log(0.65))

	s.setTransition(Single, Begin, math.Log(0.35))
	s.setTransition(Single, Middle, math.Log(0.0))
	s.setTransition(Single, End, math.Log(0.0))
	s.setTransition(Single, Single, math.Log(0.65))

	s.start[Begin] = math.Log(0.2)
	s.start[Middle] = math.Log(0.001)
	s.start[End] = math.Log(0.001)
	s.start[Single] = math.Log(0.798)

	initializeDefaultEmission(s)
	return s
}

// Segment segments Chinese text into words using Viterbi decoding.
// Non-CJK characters are passed through as single-character tokens.
func (s *Segmenter) Segment(text string) []string {
	runes := []rune(text)
	if len(runes) == 0 {
		return nil
	}

	n := len(runes)
	viterbi := make([][4]float64, n)
	backptr := make([][4]int, n)

	for i := 0; i < n; i++ {
		for t := 0; t < 4; t++ {
			viterbi[i][t] = math.Inf(-1)
		}
	}

	firstRune := runes[0]
	for t := Tag(0); t < 4; t++ {
		emit := s.emissionProb(t, firstRune)
		viterbi[0][t] = s.start[t] + emit
	}

	for i := 1; i < n; i++ {
		r := runes[i]
		for cur := Tag(0); cur < 4; cur++ {
			emit := s.emissionProb(cur, r)
			bestProb := math.Inf(-1)
			bestPrev := 0
			for prev := Tag(0); prev < 4; prev++ {
				if viterbi[i-1][prev] == math.Inf(-1) {
					continue
				}
				prob := viterbi[i-1][prev] + s.transition[prev][cur] + emit
				if prob > bestProb {
					bestProb = prob
					bestPrev = int(prev)
				}
			}
			viterbi[i][cur] = bestProb
			backptr[i][cur] = bestPrev
		}
	}

	bestTag := 0
	bestProb := math.Inf(-1)
	for t := Tag(0); t < 4; t++ {
		if viterbi[n-1][t] > bestProb {
			bestProb = viterbi[n-1][t]
			bestTag = int(t)
		}
	}

	tags := make([]Tag, n)
	tags[n-1] = Tag(bestTag)
	for i := n - 1; i > 0; i-- {
		tags[i-1] = Tag(backptr[i][tags[i]])
	}

	result := make([]string, 0)
	wordStart := 0
	for i := 0; i < n; i++ {
		switch tags[i] {
		case Single:
			result = append(result, string(runes[i]))
			wordStart = i + 1
		case End:
			result = append(result, string(runes[wordStart:i+1]))
			wordStart = i + 1
		case Begin, Middle:
			// word continues
		}
	}

	// Flush any pending word that wasn't terminated by End or Single.
	// This can happen when a block boundary cuts through the middle of a
	// multi-character word, leaving the final rune tagged as Begin/Middle.
	if wordStart < n {
		result = append(result, string(runes[wordStart:]))
	}

	return result
}

func (s *Segmenter) setTransition(from, to Tag, prob float64) {
	s.transition[from][to] = prob
}

// emissionProb returns the log probability of rune r given tag t.
func (s *Segmenter) emissionProb(t Tag, r rune) float64 {
	m := s.emission[t]
	if p, ok := m[r]; ok {
		return p
	}
	return s.defaultEmission(t, r)
}

func (s *Segmenter) defaultEmission(t Tag, r rune) float64 {
	if unicode.In(r, unicode.Han) {
		switch t {
		case Begin:
			return math.Log(0.25)
		case Middle:
			return math.Log(0.25)
		case End:
			return math.Log(0.25)
		case Single:
			return math.Log(0.25)
		}
	}
	if t == Single {
		return math.Log(0.9)
	}
	return math.Log(0.033)
}

// Train updates emission probabilities from a segmented corpus.
// Each line should contain words separated by spaces.
func (s *Segmenter) Train(lines []string) {
	counts := [4]map[rune]int{
		Begin:   make(map[rune]int),
		Middle:  make(map[rune]int),
		End:     make(map[rune]int),
		Single:  make(map[rune]int),
	}
	tagTotals := [4]int{}

	for _, line := range lines {
		words := splitWords(line)
		for _, word := range words {
			runes := []rune(word)
			switch len(runes) {
			case 1:
				counts[Single][runes[0]]++
				tagTotals[Single]++
			default:
				counts[Begin][runes[0]]++
				tagTotals[Begin]++
				for i := 1; i < len(runes)-1; i++ {
					counts[Middle][runes[i]]++
					tagTotals[Middle]++
				}
				counts[End][runes[len(runes)-1]]++
				tagTotals[End]++
			}
		}
	}

	vocabSize := 10000
	for t := Tag(0); t < 4; t++ {
		total := tagTotals[t]
		if total == 0 {
			continue
		}
		for r, cnt := range counts[t] {
			prob := float64(cnt+1) / float64(total+vocabSize)
			s.emission[t][r] = math.Log(prob)
		}
	}
}

func splitWords(line string) []string {
	var words []string
	start := 0
	for i := 0; i <= len(line); i++ {
		if i == len(line) || line[i] == ' ' {
			if i > start {
				words = append(words, line[start:i])
			}
			start = i + 1
		}
	}
	return words
}

func initializeDefaultEmission(s *Segmenter) {
	singleChars := []rune("的了在是有和与及或但而所这那其个各某些任何每所有自己我们你们他们那个这个那些哪些什么哪儿哪里怎么多少几多少一点儿些许")
	for _, r := range singleChars {
		s.emission[Single][r] = math.Log(0.6)
		s.emission[Begin][r] = math.Log(0.1)
		s.emission[Middle][r] = math.Log(0.1)
		s.emission[End][r] = math.Log(0.2)
	}

	beginChars := []rune("北上广深中上下左右内外前后今昨明天年月日")
	for _, r := range beginChars {
		if _, exists := s.emission[Single][r]; !exists {
			s.emission[Begin][r] = math.Log(0.4)
			s.emission[End][r] = math.Log(0.2)
			s.emission[Single][r] = math.Log(0.2)
			s.emission[Middle][r] = math.Log(0.2)
		}
	}

	endChars := []rune("门所院馆店站场中心大厦小区园区路街道大街")
	for _, r := range endChars {
		if _, exists := s.emission[Single][r]; !exists {
			s.emission[End][r] = math.Log(0.5)
			s.emission[Begin][r] = math.Log(0.1)
			s.emission[Middle][r] = math.Log(0.1)
			s.emission[Single][r] = math.Log(0.3)
		}
	}
}

// IsChineseText reports whether the text contains mostly Chinese characters.
func IsChineseText(text string) bool {
	runes := []rune(text)
	if len(runes) == 0 {
		return false
	}
	hanCount := 0
	for _, r := range runes {
		if unicode.In(r, unicode.Han) {
			hanCount++
		}
	}
	return float64(hanCount)/float64(len(runes)) > 0.5
}