package fileop

// CompressType represents the type of compression to be used.
type CompressType int

// Supported compression types.
const (
	NONE   CompressType = iota // No compression
	GZIP                       // GZIP compression
	ZLIB                       // ZLIB compression
	SNAPPY                     // SNAPPY compression
)

// String returns the string representation of the compression type.
func (ct CompressType) String() string {
	switch ct {
	case NONE:
		return "none"
	case GZIP:
		return "gzip"
	case ZLIB:
		return "zlib"
	case SNAPPY:
		return "snappy"
	default:
		return "unknown"
	}
}
