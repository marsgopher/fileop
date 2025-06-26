package filetarget

import (
	"fmt"
	"io"

	"github.com/marsgopher/fileop"
	"github.com/marsgopher/fileop/integration/afero"
)

type WrapFS struct {
	Source fileop.FileReaderInterface
	Target interface {
		fileop.FileWriterInterface
		fileop.Stater
		io.Closer
	}
}

func newWrapFS(remoteFS fileop.FileSystemWithCloser) *WrapFS {
	localFS, err := afero.New(afero.Disk)
	if err != nil {
		panic(err)
	}

	return &WrapFS{
		Target: remoteFS,
		Source: localFS,
	}
}

func (w *WrapFS) Close() error {
	return w.Target.Close()
}

func (w *WrapFS) Put(local, remote string) error {
	reader, err := fileop.NewFileReader(w.Source, local, fileop.NONE)
	if err != nil {
		return fmt.Errorf("open local file: %w", err)
	}
	defer func() {
		if reader != nil {
			_ = reader.Close()
		}
	}()

	if err := w.PutStream(reader, remote); err != nil {
		return nil
	}
	defer func() { reader = nil }()
	if err := reader.Close(); err != nil {
		return fmt.Errorf("close reader: %w", err)
	}
	return err
}

func (w *WrapFS) PutStream(reader io.Reader, remote string) error {
	writer, err := fileop.NewFileWriter(w.Target, remote, 0, fileop.NONE)
	if err != nil {
		return err
	}
	defer func() {
		if writer != nil {
			_ = writer.Close()
		}
	}()

	if _, err = io.Copy(writer, reader); err != nil {
		return err
	}
	defer func() { writer = nil }()
	if err := writer.Close(); err != nil {
		return fmt.Errorf("close writer: %w", err)
	}
	return err
}

func (w *WrapFS) PutEmpty(remote string) error {
	writer, err := fileop.NewFileWriter(w.Target, remote, 0, fileop.NONE)
	if err != nil {
		return err
	}
	return writer.Close()
}

func (w *WrapFS) PutFinish(remote string) error {
	target := remote + ".finish"

	writer, err := fileop.NewFileWriter(w.Target, target, 0, fileop.NONE)
	if err != nil {
		return err
	}
	return writer.Close()
}

func (w *WrapFS) Exist(remote string) bool {
	_, err := w.Target.Stat(remote)
	return err == nil
}
