package compress

import (
"bytes"
"testing"
)

func TestDebugRoundTrip(t *testing.T) {
input := "北京是中国的首都。上海是国际金融中心。广州是南方大门。深圳是科技创新之都。\n"
for i := 0; i < 500; i++ {
input += "北京是中国的首都。上海是国际金融中心。广州是南方大门。深圳是科技创新之都。\n"
}

inputBytes := []byte(input)
t.Logf("Input length: %d", len(inputBytes))

var encoded bytes.Buffer
if err := Compress(&encoded, inputBytes, "zh"); err != nil {
t.Fatal(err)
}
t.Logf("Encoded length: %d", encoded.Len())

var decoded bytes.Buffer
if err := Decompress(&encoded, &decoded); err != nil {
t.Fatalf("Decompress failed: %v", err)
}
t.Logf("Decoded length: %d", decoded.Len())

t.Logf("Input[:100]: %q", inputBytes[:min(100, len(inputBytes))])
t.Logf("Decoded[:100]: %q", decoded.Bytes()[:min(100, decoded.Len())])

if !bytes.Equal(inputBytes, decoded.Bytes()) {
t.Logf("DIFF: input=%d, decoded=%d", len(inputBytes), decoded.Len())
t.Logf("Decoded differs at position %d: got %q\n want %q",
findDiff(inputBytes, decoded.Bytes()),
decoded.Bytes()[findDiff(inputBytes, decoded.Bytes()):findDiff(inputBytes, decoded.Bytes())+10],
inputBytes[findDiff(inputBytes, decoded.Bytes()):findDiff(inputBytes, decoded.Bytes())+10])
}
}

func findDiff(a, b []byte) int {
n := len(a)
if len(b) < n {
n = len(b)
}
for i := 0; i < n; i++ {
if a[i] != b[i] {
return i
}
}
return n
}

func min(a, b int) int {
if a < b {
return a
}
return b
}
