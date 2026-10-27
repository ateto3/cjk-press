package library

// Library interface for embedding in other Go applications
type Compressor struct{}

func NewCompressor() *Compressor { return &Compressor{} }
func (c *Compressor) Compress(data []byte) ([]byte, error) { return nil, nil }
func (c *Compressor) Decompress(data []byte) ([]byte, error) { return nil, nil }
