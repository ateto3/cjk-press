package morphseg

import (
	"strings"
	"testing"
)

const sampleChinese = `中华人民共和国是世界上人口最多的国家，拥有悠久的历史和灿烂的文化。北京是首都，也是政治、文化中心。上海是国际金融中心。`

func TestNewSegmenter(t *testing.T) {
	s := NewSegmenter()
	if s == nil {
		t.Fatal("expected non-nil segmenter")
	}
}

func TestSegmentBasic(t *testing.T) {
	s := NewSegmenter()
	result := s.Segment("北京是中国的首都")
	if len(result) == 0 {
		t.Fatal("expected non-empty segmentation")
	}
	reconstructed := strings.Join(result, "")
	if reconstructed != "北京是中国的首都" {
		t.Errorf("segment reconstruction failed: %q", reconstructed)
	}
}

func TestSegmentEmpty(t *testing.T) {
	s := NewSegmenter()
	result := s.Segment("")
	if result != nil {
		t.Errorf("expected nil for empty input, got %v", result)
	}
}

func TestSegmentNonChinese(t *testing.T) {
	s := NewSegmenter()
	result := s.Segment("Hello World 123")
	if result == nil {
		t.Fatal("expected non-nil result for non-Chinese text")
	}
	reconstructed := strings.Join(result, "")
	if reconstructed != "Hello World 123" {
		t.Errorf("non-Chinese reconstruction failed: %q", reconstructed)
	}
}

func TestIsChineseText(t *testing.T) {
	if !IsChineseText(sampleChinese) {
		t.Error("expected sample to be detected as Chinese")
	}
	if IsChineseText("Hello World") {
		t.Error("expected English text not to be detected as Chinese")
	}
	if IsChineseText("") {
		t.Error("expected empty text not to be detected as Chinese")
	}
}

func TestTrain(t *testing.T) {
	s := NewSegmenter()
	corpus := []string{
		"北京 上海 广州 深圳",
		"中国 人民 共和国",
		"北京 是 中国 的 首都",
		"上海 是 国际 金融 中心",
	}
	s.Train(corpus)
	result := s.Segment("北京上海")
	reconstructed := strings.Join(result, "")
	if reconstructed != "北京上海" {
		t.Errorf("trailing should be preserved: %q", reconstructed)
	}
}

func TestSegmentReconstructsOriginal(t *testing.T) {
	s := NewSegmenter()
	result := s.Segment(sampleChinese)
	var buf strings.Builder
	for _, w := range result {
		buf.WriteString(w)
	}
	if buf.String() != sampleChinese {
		t.Errorf("reconstruction mismatch:\n got: %q\n want: %q",
			buf.String(), sampleChinese)
	}
}

func BenchmarkSegmentChinese(b *testing.B) {
	s := NewSegmenter()
	text := sampleChinese
	for i := 0; i < 100; i++ {
		text = text + sampleChinese
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Segment(text)
	}
}
