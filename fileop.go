package fileop

import (
	"io"
	"io/fs"
)

// FileReaderInterface is an alias for Reader interface.
type FileReaderInterface Reader

// NewReaderFunc is a function type for creating new readers with compression support.
type NewReaderFunc func(fri FileReaderInterface, path string, ct CompressType) (io.ReadCloser, error)

// FileWriterInterface provides basic file writing operations.
type FileWriterInterface interface {
	MkdirAll(dirname string, perm fs.FileMode) error
	Create(name string) (io.WriteCloser, error)
}

// NewWriterFunc is a function type for creating new writers with compression and buffer support.
type NewWriterFunc func(fwi FileWriterInterface, path string, bufSize int, ct CompressType) (io.WriteCloser, error)
