package fileop

import (
	"fmt"
	"io"

	"github.com/klauspost/compress/gzip"
	"github.com/klauspost/compress/snappy"
	"github.com/klauspost/compress/zlib"
)

func NewCompressReader(file io.Reader, ct CompressType) (io.ReadCloser, error) {
	switch ct {
	case NONE:
		return &readerNoClose{file}, nil

	case GZIP:
		reader, err := gzip.NewReader(file)
		if err != nil {
			return nil, fmt.Errorf("gzip reader: %w", err)
		}
		return reader, nil

	case ZLIB:
		reader, err := zlib.NewReader(file)
		if err != nil {
			return nil, fmt.Errorf("zlib reader: %w", err)
		}
		return reader, nil

	case SNAPPY:
		return &readerNoClose{snappy.NewReader(file)}, nil

	default:
		return nil, fmt.Errorf("compress type %v not support", ct)
	}
}

type readerNoClose struct {
	io.Reader
}

func (r *readerNoClose) Close() error {
	return nil
}
