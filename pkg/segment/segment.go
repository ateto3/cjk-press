package segment

import "unicode"

// IsCJK reports whether r belongs to one of the commonly used CJK scripts.
func IsCJK(r rune) bool {
	return unicode.In(r, unicode.Han, unicode.Hangul, unicode.Hiragana, unicode.Katakana)
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
