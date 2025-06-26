package fileop

import (
	"bufio"
	"fmt"
	"io"
	"path/filepath"
	"sync"
)

var fwFree = sync.Pool{
	New: func() interface{} {
		return new(FileWriter)
	},
}

// constraint
var (
	_ io.WriteCloser = &FileWriter{}

	_ NewWriterFunc = func(fwi FileWriterInterface, path string, bufSize int, ct CompressType) (io.WriteCloser, error) {
		return NewFileWriter(fwi, path, bufSize, ct)
	}
)

type FileWriter struct {
	Path   string
	file   io.WriteCloser
	buf    *bufio.Writer
	writer io.WriteCloser
}

func (fw *FileWriter) free() {
	fw.Path = ""
	if fw.writer != nil {
		_ = fw.writer.Close()
	}
	fw.writer = nil
	if fw.buf != nil {
		_ = fw.buf.Flush()
	}
	fw.buf = nil
	if fw.file != nil {
		_ = fw.file.Close()
	}
	fw.file = nil
}

func NewFileWriter(fwi FileWriterInterface, dstPath string, bufSize int, ct CompressType) (*FileWriter, error) {
	fw := fwFree.Get().(*FileWriter)
	fw.Path = dstPath

	// auto create dir
	dstDir := filepath.Dir(dstPath)
	if err := fwi.MkdirAll(dstDir, 0755); err != nil {
		return nil, fmt.Errorf("mkdir: %w", err)
	}

	file, err := fwi.Create(dstPath)
	if err != nil {
		return nil, fmt.Errorf("create dst: %w", err)
	}
	fw.file = file

	var buf *bufio.Writer
	if size := bufSize; size > 0 {
		buf = bufio.NewWriterSize(file, size)
	} else {
		buf = bufio.NewWriter(file)
	}
	fw.buf = buf

	writer, err := NewCompressWriter(buf, ct)
	if err != nil {
		_ = fw.Close()
		return nil, fmt.Errorf("compress writer: %w", err)
	}
	fw.writer = writer

	return fw, nil
}

func (fw *FileWriter) Close() error {
	defer func() {
		fw.free()
		fwFree.Put(fw)
	}()

	defer func() { fw.writer = nil }()
	if wt := fw.writer; wt != nil {
		if err := wt.Close(); err != nil {
			return fmt.Errorf("close writer: %w", err)
		}
	}

	defer func() { fw.buf = nil }()
	if buf := fw.buf; buf != nil {
		if err := buf.Flush(); err != nil {
			return fmt.Errorf("flush buf: %w", err)
		}
	}

	defer func() { fw.file = nil }()
	if wt := fw.file; wt != nil {
		if err := wt.Close(); err != nil {
			return fmt.Errorf("close file: %w", err)
		}
	}

	return nil
}

func (fw *FileWriter) Write(p []byte) (int, error) {
	return fw.writer.Write(p)
}

func (fw *FileWriter) WriteLine(line []byte) (int, error) {
	if len(line) == 0 {
		return 0, nil
	}
	return fw.writer.Write(append(line, '\n'))
}
