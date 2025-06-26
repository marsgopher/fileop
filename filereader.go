package fileop

import (
	"bufio"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
)

var frFree = sync.Pool{
	New: func() interface{} {
		return new(FileReader)
	},
}

// constraint
var (
	_ io.ReadCloser = &FileReader{}

	_ NewReaderFunc = func(fri FileReaderInterface, path string, ct CompressType) (io.ReadCloser, error) {
		return NewFileReader(fri, path, ct)
	}
)

type FileReader struct {
	Path   string
	file   io.ReadCloser
	reader io.ReadCloser

	EOF     bool
	scanner *bufio.Scanner
}

func (fr *FileReader) free() {
	fr.Path = ""
	fr.EOF = false
	fr.scanner = nil
	if fr.reader != nil {
		_ = fr.reader.Close()
	}
	fr.reader = nil
	if fr.file != nil {
		_ = fr.file.Close()
	}
	fr.file = nil
}

// NewFileReader create *FileReader on any FileReaderInterface.
// NOTE: you can call IsUnhandledFileReaderError judge errors can not solve by retry.
func NewFileReader(fri FileReaderInterface, srcPath string, ct CompressType) (*FileReader, error) {
	fr := frFree.Get().(*FileReader)
	fr.Path = srcPath

	file, err := fri.Open(srcPath)
	if err != nil {
		return nil, fmt.Errorf("open file %s: %w", srcPath, err)
	}
	fr.file = file

	reader, err := NewCompressReader(file, ct)
	if err != nil {
		_ = fr.Close()
		return nil, fmt.Errorf("compress reader: %w", err)
	}
	fr.reader = reader

	return fr, nil
}

func IsUnhandledFileReaderError(err error) bool {
	return errors.Is(err, io.EOF) ||
		errors.Is(err, os.ErrNotExist) ||
		errors.Is(err, gzip.ErrHeader)
}

func (fr *FileReader) Close() error {
	defer func() {
		fr.free()
		frFree.Put(fr)
	}()

	defer func() { fr.reader = nil }()
	if rd := fr.reader; rd != nil {
		if err := rd.Close(); err != nil {
			return fmt.Errorf("close reader: %w", err)
		}
	}

	defer func() { fr.file = nil }()
	if rd := fr.file; rd != nil {
		if err := rd.Close(); err != nil {
			return fmt.Errorf("close file: %w", err)
		}
	}

	return nil
}

func (fr *FileReader) Read(p []byte) (int, error) {
	return fr.reader.Read(p)
}

func (fr *FileReader) ReadLine() ([]byte, error) {
	if fr.EOF || fr.reader == nil {
		return nil, io.EOF
	}
	if fr.scanner == nil {
		fr.scanner = bufio.NewScanner(fr.reader)
	}

	var line []byte
	var hasOne bool
	for fr.scanner.Scan() {
		line = fr.scanner.Bytes()
		hasOne = true
		break
	}
	if err := fr.scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan %s: %w", fr.Path, err)
	}

	if !hasOne {
		fr.EOF = true
	}
	return line, nil
}
