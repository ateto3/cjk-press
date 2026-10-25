package segment

import "unicode"

// IsCJK reports whether r belongs to one of the commonly used CJK scripts.
func IsCJK(r rune) bool {
	return unicode.In(r, unicode.Han, unicode.Hangul, unicode.Hiragana, unicode.Katakana)
}

// IsKorean reports whether r is a Hangul character.
func IsKorean(r rune) bool {
	return unicode.In(r, unicode.Hangul)
}

// IsJapanese reports whether r is a Japanese character (Kanji, Hiragana, or Katakana).
func IsJapanese(r rune) bool {
	return unicode.In(r, unicode.Han, unicode.Hiragana, unicode.Katakana)
}

// IsChinese reports whether r is a Han (Chinese) character.
func IsChinese(r rune) bool {
	return unicode.In(r, unicode.Han)
}

// IsKanji reports whether r is a Kanji character (Han used in Japanese).
func IsKanji(r rune) bool {
	return unicode.In(r, unicode.Han)
}

// IsKana reports whether r is Hiragana or Katakana.
func IsKana(r rune) bool {
	return unicode.In(r, unicode.Hiragana, unicode.Katakana)
}

// IsKoreanParticle reports whether the given Hangul syllable is a common
// Korean particle (조사). Particles attach to preceding nouns and should
// generally be kept attached for better compression and linguistic accuracy.
func IsKoreanParticle(word string) bool {
	particles := []string{
		// Subject particles
		"이", "가",
		// Topic particles
		"은", "는",
		// Object particles
		"을", "를",
		// Possessive/genitive
		"의",
		// Location/direction
		"에", "에서",
		// Instrumental/directional
		"로", "으로",
		// Also/only
		"도", "만",
		// From
		"부터", "까지",
		// Like/as
		"처럼", "같이",
		// Honorific subject
		"께서",
	}
	for _, p := range particles {
		if word == p {
			return true
		}
	}
	return false
}

// SegmentKorean splits Korean text, keeping particles attached to their
// preceding noun where possible. It returns segments that preserve the
// original text when concatenated.
func SegmentKorean(text string) []string {
	runes := []rune(text)
	if len(runes) == 0 {
		return nil
	}

	result := make([]string, 0)
	i := 0
	for i < len(runes) {
		// Skip non-Korean
		if !IsKorean(runes[i]) {
			start := i
			for i < len(runes) && !IsKorean(runes[i]) {
				i++
			}
			result = append(result, string(runes[start:i]))
			continue
		}

		// Collect a Hangul run
		start := i
		for i < len(runes) && IsKorean(runes[i]) {
			i++
		}
		runesPart := runes[start:i]
		run := string(runesPart)

		// Try to split off a trailing particle. We work with rune counts
		// because Korean particles attach at run boundaries.
		if len(runesPart) > 1 {
			split := false
			for r := 2; r >= 1; r-- {
				if IsKoreanParticle(string(runesPart[len(runesPart)-r:])) {
					result = append(result, string(runesPart[:len(runesPart)-r]))
					result = append(result, string(runesPart[len(runesPart)-r:]))
					split = true
					break
				}
			}
			if split {
				continue
			}
		}
		result = append(result, run)
	}
	return result
}

// SegmentJapanese splits Japanese text at Kanji-Kana boundaries and
// between Kanji compounds, producing linguistically meaningful segments.
// Non-Japanese text is passed through unchanged.
func SegmentJapanese(text string) []string {
	runes := []rune(text)
	if len(runes) == 0 {
		return nil
	}

	result := make([]string, 0)
	i := 0
	for i < len(runes) {
		if !IsJapanese(runes[i]) {
			start := i
			for i < len(runes) && !IsJapanese(runes[i]) {
				i++
			}
			result = append(result, string(runes[start:i]))
			continue
		}

		// Collect a Japanese run
		start := i
		prevKana := false
		for i < len(runes) && IsJapanese(runes[i]) {
			curIsKana := IsKana(runes[i])
			// Split at Kanji-Kana or Kana-Kanji transitions
			if i > start && curIsKana != prevKana {
				break
			}
			prevKana = curIsKana
			i++
		}
		result = append(result, string(runes[start:i]))
	}
	return result
}

// SegmentLanguage applies language-specific segmentation for CJK text.
// For Chinese, it uses the provided segmenter. For Korean, it keeps
// particles attached. For Japanese, it splits at Kanji-Kana boundaries.
// The maxRunes parameter limits segment size for Chinese max-match fallback.
func SegmentLanguage(text string, lang string, maxRunes int) []string {
	if maxRunes < 1 {
		maxRunes = 8
	}

	switch lang {
	case "ko":
		return SegmentKorean(text)
	case "ja":
		return SegmentJapanese(text)
	case "zh", "":
		return Segment(text, maxRunes)
	default:
		return Segment(text, maxRunes)
	}
}

// Segment splits text into CJK runs and non-CJK runs. CJK runs are further
// split into chunks of up to maxRunes runes, which gives the codec useful
// repeated phrases without changing non-CJK text at all.
func Segment(text string, maxRunes int) []string {
	if maxRunes < 1 {
		maxRunes = 8
	}
	runes := []rune(text)
	result := make([]string, 0)
	for i := 0; i < len(runes); {
		start := i
		cjk := IsCJK(runes[i])
		for i < len(runes) && IsCJK(runes[i]) == cjk {
			i++
		}
		if !cjk {
			result = append(result, string(runes[start:i]))
			continue
		}
		for start < i {
			end := start + maxRunes
			if end > i {
				end = i
			}
			result = append(result, string(runes[start:end]))
			start = end
		}
	}
	return result
}
