package benchmarks

import "testing"

func BenchmarkCompression(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// TODO: implement benchmark
	}
}
