package pipeline

import (
	"testing"

	"github.com/atetov/cjk-press/pkg/dictionary"
	"github.com/atetov/cjk-press/pkg/segment"
)

// realChineseText is a representative sample of modern Chinese text.
const realChineseText = `中华人民共和国是世界上人口最多的国家，拥有悠久的历史和灿烂的文化。北京是首都，也是政治、文化中心。上海是国际金融中心，港口吞吐量居世界前列。广州和深圳是南方重要的经济特区，制造业和贸易发达。

中国有56个民族，汉族是人口最多的民族。普通话是全国通用的语言文字，但各地区仍保留着方言，如粤语、闽语、吴语等。春节、清明节、端午节、中秋节是传统四大节日，家家团聚，吃年糕、月饼，放鞭炮。

科技领域，中国在高铁、5G通信、人工智能、量子信息等方面处于世界领先地位。航天方面，嫦娥探月、天问探火、核心空间站建设取得重大进展。经济方面，电子商务、移动支付、共享经济模式在全球范围内产生深远影响。

文化方面，中国书法、绘画、京剧、功夫、茶道、中医等传统文化享誉世界。长城、故宫、兵马俑、莫高窟等世界文化遗产吸引了大量游客。`

const realJapaneseText = `日本国は東アジアに位置する立憲君主国です。東京は首都であり、政治・経済・文化の中心です。日本語が公用語で、文字にはひらがな・カタカナ・漢字が用いられます。

日本は四季がはっきりしており、桜、紅葉、雪景色など自然の美しさが世界的に知られています。温泉、神社、仏閣、茶道、華道、武道など伝統文化も多くの外国人観光客を惹きつけています。

自動車、電子機器、ロボット、アニメ、ゲームなどの産業は世界トップクラスの技術力を持ち、新幹線は安全かつ正確な高速鉄道として有名です。`

const realKoreanText = `대한민국은 동아시아에 위치한 민주주의 국가로, 수도는 서울입니다. 한국어가 공용어이며 한글 문자를 사용합니다.

대한민국은 반도체, 스마트폰, 자동차, 조선, 철강 등 제조업이 발달한 나라이며, K-pop, 드라마, 영화 등 한류 문화가 전 세계적으로 인기를 끌고 있습니다. 경복궁, 창덕궁, 불국사 등 전통 유적지와 한옥 마을은 많은 관광객이 찾는 명소입니다.`

func TestProcessTextPreservesOrder(t *testing.T) {
	text := "abc你好世界defこんにちは안녕"
	parts, dict, err := ProcessText(text, 4, 2, 2)
	if err != nil {
		t.Fatalf("ProcessText returned error: %v", err)
	}

	var reconstructed string
	for _, blockParts := range parts {
		for _, p := range blockParts {
			reconstructed += p
		}
	}
	if reconstructed != text {
		t.Fatalf("parts do not reconstruct input:\n got: %q\n want: %q", reconstructed, text)
	}

	if dict.Size() == 0 {
		t.Error("expected non-empty dictionary")
	}
}

func TestProcessTextSingleBlock(t *testing.T) {
	text := "你好世界"
	parts, dict, err := ProcessText(text, 100, 1, 2)
	if err != nil {
		t.Fatalf("ProcessText returned error: %v", err)
	}

	if len(parts) != 1 {
		t.Fatalf("expected 1 block, got %d", len(parts))
	}

	var reconstructed string
	for _, p := range parts[0] {
		reconstructed += p
	}
	if reconstructed != text {
		t.Fatalf("parts do not reconstruct input: %q", reconstructed)
	}

	if len(parts[0]) != 2 {
		t.Fatalf("expected 2 parts, got %d: %v", len(parts[0]), parts[0])
	}

	if dict.Size() != 2 {
		t.Fatalf("expected 2 dictionary entries, got %d", dict.Size())
	}

	if !dict.Contains("你好") {
		t.Error("dictionary missing '你好'")
	}
	if !dict.Contains("世界") {
		t.Error("dictionary missing '世界'")
	}
}

func TestProcessTextMultipleWorkers(t *testing.T) {
	text := realChineseText
	parts1, dict1, err := ProcessText(text, 500, 1, 8)
	if err != nil {
		t.Fatalf("ProcessText(1 worker) error: %v", err)
	}
	parts4, dict4, err := ProcessText(text, 500, 4, 8)
	if err != nil {
		t.Fatalf("ProcessText(4 workers) error: %v", err)
	}

	var recon1, recon4 string
	for _, bp := range parts1 {
		for _, p := range bp {
			recon1 += p
		}
	}
	for _, bp := range parts4 {
		for _, p := range bp {
			recon4 += p
		}
	}
	if recon1 != text || recon4 != text {
		t.Fatal("reconstructed text does not match input")
	}

	if dict1.Size() != dict4.Size() {
		t.Errorf("dictionary size mismatch: 1 worker=%d, 4 workers=%d", dict1.Size(), dict4.Size())
	}
}

func TestProcessTextEmptyInput(t *testing.T) {
	parts, dict, err := ProcessText("", 100, 2, 8)
	if err != nil {
		t.Fatalf("ProcessText on empty input returned error: %v", err)
	}
	if len(parts) != 1 || len(parts[0]) != 0 {
		t.Fatalf("expected single empty block, got %v", parts)
	}
	if dict.Size() != 0 {
		t.Fatalf("expected empty dictionary, got %d entries", dict.Size())
	}
}

func TestProcessTextOnlyNonCJK(t *testing.T) {
	text := "Hello, World! 12345"
	parts, dict, err := ProcessText(text, 10, 2, 8)
	if err != nil {
		t.Fatalf("ProcessText error: %v", err)
	}

	var reconstructed string
	for _, bp := range parts {
		for _, p := range bp {
			reconstructed += p
		}
	}
	if reconstructed != text {
		t.Fatalf("parts do not reconstruct input: %q", reconstructed)
	}

	if dict.Size() != 0 {
		t.Errorf("expected empty dictionary for non-CJK text, got %d entries", dict.Size())
	}
}

func TestProcessTextMixedScripts(t *testing.T) {
	text := realChineseText + " " + realJapaneseText + " " + realKoreanText
	parts, dict, err := ProcessText(text, 300, 3, 8)
	if err != nil {
		t.Fatalf("ProcessText error: %v", err)
	}

	var reconstructed string
	for _, bp := range parts {
		for _, p := range bp {
			reconstructed += p
		}
	}
	if reconstructed != text {
		t.Fatal("mixed script reconstruction failed")
	}

	if dict.Size() == 0 {
		t.Error("expected non-empty dictionary for mixed CJK text")
	}

	textRunes := []rune(text)
	for _, r := range textRunes {
		if segment.IsCJK(r) {
			found := false
			for _, word := range dict.GetAllWords() {
				for _, wr := range []rune(word) {
					if wr == r {
						found = true
						break
					}
				}
				if found {
					break
				}
			}
			if !found {
				t.Errorf("CJK rune %q not found in any dictionary entry", r)
			}
		}
	}
}

func TestDictionaryMergeConcurrency(t *testing.T) {
	d1 := dictionary.NewDictionary()
	d1.Add("你好")
	d1.Add("世界")

	d2 := dictionary.NewDictionary()
	d2.Add("世界")
	d2.Add("北京")

	d1.Merge(d2)

	if d1.Size() != 3 {
		t.Errorf("expected 3 entries after merge, got %d", d1.Size())
	}
	if !d1.Contains("你好") || !d1.Contains("世界") || !d1.Contains("北京") {
		t.Error("dictionary missing entries after merge")
	}
	if d1.GetID("世界") == 0 {
		t.Error("'世界' should have a valid ID")
	}
}

// Benchmarks

func generateRepeatedChineseText(n int) []byte {
	base := "北京是中国的首都。上海是国际金融中心。广州是南方大门。深圳是科技创新之都。"
	var result string
	for i := 0; i < n; i++ {
		result += base
		if i%10 == 0 {
			result += "\n"
		}
	}
	return []byte(result)
}

func BenchmarkProcessTextSequential(b *testing.B) {
	text := string(generateRepeatedChineseText(1000))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := ProcessText(text, 1000, 1, 8)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkProcessTextConcurrent(b *testing.B) {
	text := string(generateRepeatedChineseText(1000))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := ProcessText(text, 1000, 4, 8)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkProcessTextLargeConcurrent(b *testing.B) {
	text := string(generateRepeatedChineseText(10000))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := ProcessText(text, 5000, 8, 8)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDictionaryMerge(b *testing.B) {
	dicts := make([]*dictionary.Dictionary, 8)
	for i := range dicts {
		dicts[i] = dictionary.NewDictionary()
		for j := 0; j < 1000; j++ {
			dicts[i].Add(string(rune('一' + j%100)))
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		global := dictionary.NewDictionary()
		for _, d := range dicts {
			global.Merge(d)
		}
		if global.Size() == 0 {
			b.Fatal("empty merged dictionary")
		}
	}
}