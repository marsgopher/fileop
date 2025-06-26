package fileop

import (
	"fmt"
	"io"

	"github.com/klauspost/compress/gzip"
	"github.com/klauspost/compress/snappy"
	"github.com/klauspost/compress/zlib"
	"github.com/klauspost/pgzip"
)

var (
	UsePGZIP    bool
	PGZIPBlocks = 4
)

func NewCompressWriter(buf io.Writer, ct CompressType) (io.WriteCloser, error) {
	switch ct {
	case NONE:
		return &writerNoClose{buf}, nil

	case GZIP:
		if UsePGZIP {
			w := pgzip.NewWriter(buf)
			if err := w.SetConcurrency(1<<20, PGZIPBlocks); err != nil {
				_ = w.Close()
				return nil, fmt.Errorf("set concurrency: %w", err)
			}
			return w, nil
		} else {
			return gzip.NewWriter(buf), nil
		}

	case ZLIB:
		return zlib.NewWriter(buf), nil

	case SNAPPY:
		return snappy.NewBufferedWriter(buf), nil

	default:
		return nil, fmt.Errorf("compress type %v not support", ct)
	}
}

type writerNoClose struct {
	io.Writer
}

func (r *writerNoClose) Close() error {
	return nil
}
