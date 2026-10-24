package segment

import "testing"

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
