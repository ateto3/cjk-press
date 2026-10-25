package segment

import (
	"strings"
	"testing"
)

func TestSegmentPreservesText(t *testing.T) {
	input := "abc你好世界defこんにちは안녕"
	parts := Segment(input, 2)
	var got string
	for _, p := range parts {
		got += p
	}
	if got != input {
		t.Fatalf("got %q", got)
	}
}

func TestIsCJK(t *testing.T) {
	for _, r := range []rune("中あア한") {
		if !IsCJK(r) {
			t.Errorf("%q not detected", r)
		}
	}
	if IsCJK('a') {
		t.Error("ASCII incorrectly detected")
	}
}

func TestIsKorean(t *testing.T) {
	hangul := "한글"
	for _, r := range hangul {
		if !IsKorean(r) {
			t.Errorf("%q not detected as Korean", r)
		}
	}
}

func TestIsJapanese(t *testing.T) {
	ja := "日本語"
	for _, r := range ja {
		if !IsJapanese(r) {
			t.Errorf("%q not detected as Japanese", r)
		}
	}
}

func TestIsChinese(t *testing.T) {
	zh := "中文"
	for _, r := range zh {
		if !IsChinese(r) {
			t.Errorf("%q not detected as Chinese", r)
		}
	}
}

func TestIsKanji(t *testing.T) {
	kanji := []rune("中")
	if !IsKanji(kanji[0]) {
		t.Error("kanji not detected")
	}
}

func TestIsKana(t *testing.T) {
	hira := []rune("あ")
	kata := []rune("ア")
	if !IsKana(hira[0]) {
		t.Error("hiragana not detected")
	}
	if !IsKana(kata[0]) {
		t.Error("katakana not detected")
	}
}

func TestIsKoreanParticle(t *testing.T) {
	particles := []string{"이", "가", "은", "는", "을", "를", "의", "에", "로", "도", "만", "에서", "부터", "까지", "처럼", "같이", "께서"}
	for _, p := range particles {
		if !IsKoreanParticle(p) {
			t.Errorf("%q not detected as particle", p)
		}
	}
}

func TestIsKoreanParticleNot(t *testing.T) {
	notParticles := []string{"서울", "부산", "안녕하세요", "compressor", ""}
	for _, w := range notParticles {
		if IsKoreanParticle(w) {
			t.Errorf("%q incorrectly detected as particle", w)
		}
	}
}

func TestSegmentKorean(t *testing.T) {
	// Hello + particle 가
	result := SegmentKorean("서울가")
	if len(result) != 2 {
		t.Fatalf("expected 2 segments, got %d: %v", len(result), result)
	}
	reconstructed := ""
	for _, s := range result {
		reconstructed += s
	}
	if reconstructed != "서울가" {
		t.Errorf("reconstruction failed: %q", reconstructed)
	}
}

func TestSegmentKoreanWithParticles(t *testing.T) {
	// 대한민국 + 은 (topic particle)
	result := SegmentKorean("대한민국은")
	reconstructed := ""
	for _, s := range result {
		reconstructed += s
	}
	if reconstructed != "대한민국은" {
		t.Errorf("reconstruction failed: %q", reconstructed)
	}
}

func TestSegmentJapaneseKanjiKana(t *testing.T) {
	// All-Kanji text should be 1 segment (no Kana boundary present)
	result := SegmentJapanese("漢字語")
	if len(result) != 1 {
		t.Fatalf("expected 1 segment for all-Kanji text, got %d: %v", len(result), result)
	}
	reconstructed := ""
	for _, s := range result {
		reconstructed += s
	}
	if reconstructed != "漢字語" {
		t.Errorf("reconstruction failed: %q", reconstructed)
	}
}

func TestSegmentJapaneseKanaKana(t *testing.T) {
	// ひらがな + ひらがな (all Kana, no boundary)
	result := SegmentJapanese("くろねこ")
	if len(result) != 1 {
		t.Fatalf("expected 1 segment, got %d: %v", len(result), result)
	}
	reconstructed := ""
	for _, s := range result {
		reconstructed += s
	}
	if reconstructed != "くろねこ" {
		t.Errorf("reconstruction failed: %q", reconstructed)
	}
}

func TestSegmentJapaneseMixed(t *testing.T) {
	// Mixed Kanji, Kana, and punctuation
	result := SegmentJapanese("こんにちは世界!")
	if len(result) < 2 {
		t.Fatalf("expected at least 2 segments, got %d: %v", len(result), result)
	}
	reconstructed := ""
	for _, s := range result {
		reconstructed += s
	}
	// Non-Japanese parts pass through unchanged
	if !strings.HasSuffix(reconstructed, "世界!") {
		t.Errorf("punctuation not preserved: %q", reconstructed)
	}
}

func TestSegmentLanguage(t *testing.T) {
	chinese := "北京是中国的首都"
	ko := "서울은한국의수도"

	// Default (zh) uses segmenter
	zhResult := SegmentLanguage(chinese, "", 8)
	if len(zhResult) == 0 {
		t.Error("expected segments for Chinese")
	}
	reconstructed := ""
	for _, s := range zhResult {
		reconstructed += s
	}
	if reconstructed != chinese {
		t.Errorf("Chinese reconstruction failed: %q", reconstructed)
	}

	// Korean keeps particles attached
	koResult := SegmentLanguage(ko, "ko", 8)
	reconstructed = ""
	for _, s := range koResult {
		reconstructed += s
	}
	if reconstructed != ko {
		t.Errorf("Korean reconstruction failed: %q", reconstructed)
	}

	// Japanese splits at Kanji-Kana boundary
	jaResult := SegmentLanguage("漢字語を", "ja", 8)
	reconstructed = ""
	for _, s := range jaResult {
		reconstructed += s
	}
	if reconstructed != "漢字語を" {
		t.Errorf("Japanese reconstruction failed: %q", reconstructed)
	}
}
