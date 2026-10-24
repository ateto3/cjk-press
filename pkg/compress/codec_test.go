package compress

import (
	"bytes"
	"testing"
)

func TestRoundTrip(t *testing.T) {
	input := []byte("标题：你好，世界！ Hello こんにちは 안녕하세요\n你好，世界！")
	var encoded bytes.Buffer
	if err := Compress(&encoded, input); err != nil {
		t.Fatal(err)
	}
	var decoded bytes.Buffer
	if err := Decompress(&encoded, &decoded); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(input, decoded.Bytes()) {
		t.Fatalf("round trip changed input: %q", decoded.Bytes())
	}
}

func TestRejectsBadMagic(t *testing.T) {
	if err := Decompress(bytes.NewReader([]byte("nope")), &bytes.Buffer{}); err == nil {
		t.Fatal("expected invalid format error")
	}
}

func TestConcurrentRoundTrip(t *testing.T) {
	input := []byte("标题：你好，世界！ Hello こんにちは 안녕하세요\n你好，世界！")
	var encoded bytes.Buffer
	if err := CompressConcurrent(&encoded, input, 1000, 4); err != nil {
		t.Fatal(err)
	}
	var decoded bytes.Buffer
	if err := Decompress(&encoded, &decoded); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(input, decoded.Bytes()) {
		t.Fatalf("concurrent round trip changed input: %q", decoded.Bytes())
	}
}

func TestConcurrentLargeTextRoundTrip(t *testing.T) {
	// Build a larger Chinese text with repeated phrases
	base := "北京是中国的首都。上海是国际金融中心。广州是南方大门。深圳是科技创新之都。\n"
	var input string
	for i := 0; i < 500; i++ {
		input += base
	}
	data := []byte(input)

	var encoded bytes.Buffer
	if err := CompressConcurrent(&encoded, data, 500, 8); err != nil {
		t.Fatal(err)
	}

	var decoded bytes.Buffer
	if err := Decompress(&encoded, &decoded); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(data, decoded.Bytes()) {
		t.Fatalf("concurrent large text round trip failed: input=%d bytes, output=%d bytes", len(data), decoded.Len())
	}
}

func TestConcurrentMixedScriptsRoundTrip(t *testing.T) {
	input := []byte(realChineseText + " " + realJapaneseText + " " + realKoreanText)
	var encoded bytes.Buffer
	if err := CompressConcurrent(&encoded, input, 300, 3); err != nil {
		t.Fatal(err)
	}
	var decoded bytes.Buffer
	if err := Decompress(&encoded, &decoded); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(input, decoded.Bytes()) {
		t.Fatalf("mixed scripts concurrent round trip failed")
	}
}

func TestConcurrentEmptyInput(t *testing.T) {
	var encoded bytes.Buffer
	if err := CompressConcurrent(&encoded, []byte{}, 1000, 4); err != nil {
		t.Fatal(err)
	}
	var decoded bytes.Buffer
	if err := Decompress(&encoded, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.Len() != 0 {
		t.Fatalf("expected empty output, got %q", decoded.String())
	}
}

func TestConcurrentSequentialCompatibility(t *testing.T) {
	// Verify that output from CompressConcurrent can be decompressed by Decompress
	input := []byte("标题：你好，世界！ Hello こんにちは 안녕하세요\n你好，世界！")
	var encoded bytes.Buffer
	if err := CompressConcurrent(&encoded, input, 100, 2); err != nil {
		t.Fatal(err)
	}
	var decoded bytes.Buffer
	if err := Decompress(&encoded, &decoded); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(input, decoded.Bytes()) {
		t.Fatalf("consequential compatibility failed")
	}
}

// Real-world Chinese text samples for testing
const realChineseText = `中华人民共和国是世界上人口最多的国家，拥有悠久的历史和灿烂的文化。北京是首都，也是政治、文化中心。上海是国际金融中心，港口吞吐量居世界前列。广州和深圳是南方重要的经济特区，制造业和贸易发达。

中国有56个民族，汉族是人口最多的民族。普通话是全国通用的语言文字，但各地区仍保留着方言，如粤语、闽语、吴语等。春节、清明节、端午节、中秋节是传统四大节日，家家团聚，吃年糕、月饼，放鞭炮。

科技领域，中国在高铁、5G通信、人工智能、量子信息等方面处于世界领先地位。航天方面，嫦娥探月、天问探火、核心空间站建设取得重大进展。经济方面，电子商务、移动支付、共享经济模式在全球范围内产生深远影响。

文化方面，中国书法、绘画、京剧、功夫、茶道、中医等传统文化享誉世界。长城、故宫、兵马俑、莫高窟等世界文化遗产吸引了大量游客。`

const realJapaneseText = `日本国は東アジアに位置する立憲君主国です。東京は首都であり、政治・経済・文化の中心です。日本語が公用語で、文字にはひらがな・カタカナ・漢字が用いられます。

日本は四季がはっきりしており、桜、紅葉、雪景色など自然の美しさが世界的に知られています。温泉、神社、仏閣、茶道、華道、武道など伝統文化も多くの外国人観光客を惹きつけています。

自動車、電子機器、ロボット、アニメ、ゲームなどの産業は世界トップクラスの技術力を持ち、新幹線は安全かつ正確な高速鉄道として有名です。`

const realKoreanText = `대한민국은 동아시아에 위치한 민주주의 국가로, 수도는 서울입니다. 한국어가 공용어이며 한글 문자를 사용합니다.

대한민국은 반도체, 스마트폰, 자동차, 조선, 철강 등 제조업이 발달한 나라이며, K-pop, 드라마, 영화 등 한류 문화가 전 세계적으로 인기를 끌고 있습니다. 경복궁, 창덕궁, 불국사 등 전통 유적지와 한옥 마을은 많은 관광객이 찾는 명소입니다.`
