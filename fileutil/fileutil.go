package fileutil

import (
	"io"

	"github.com/marsgopher/fileop"
)

func WriteFile(h fileop.FileSystemWithCloser, path string, content io.Reader) error {
	wt, err := fileop.NewFileWriter(h, path, 0, fileop.NONE)
	if err != nil {
		return err
	}
	defer func() {
		if wt != nil {
			_ = wt.Close()
		}
	}()
	if _, err := io.Copy(wt, content); err != nil {
		return err
	}
	defer func() { wt = nil }()
	return wt.Close()
}

func ReadFile(h fileop.FileSystemWithCloser, path string) ([]byte, error) {
	rd, err := fileop.NewFileReader(h, path, fileop.NONE)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rd.Close()
	}()
	return io.ReadAll(rd)
}
