package compress

type CompressionLevel int

const (
	Fast CompressionLevel = iota
	Balanced
	Maximum
)

func SetLevel(l CompressionLevel) {}
