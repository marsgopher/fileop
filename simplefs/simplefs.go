package simplefs

import (
	"fmt"
	"io"
	"io/fs"
	"path/filepath"

	"github.com/marsgopher/fileop"
	"github.com/marsgopher/fileop/filetarget"
	"github.com/marsgopher/fileop/integration/afero"
	"github.com/marsgopher/fileop/integration/hdfs"
	"github.com/marsgopher/fileop/integration/minio"
	"github.com/marsgopher/fileop/integration/obs"
	"github.com/marsgopher/fileop/integration/upyun"
)

type Config struct {
	Mode  string       `mapstructure:"mode"`
	OBS   obs.Config   `mapstructure:"obs"`
	UPYUN upyun.Config `mapstructure:"upyun"`
	HDFS  hdfs.Config  `mapstructure:"hdfs"`
	MINIO minio.Config `mapstructure:"minio"`
}

func NewV1(c Config) (fileop.FileSystemSimple, error) {
	switch c.Mode {
	case "disk":
		h, err := afero.New(afero.Disk)
		if err != nil {
			return nil, fmt.Errorf("new disk: %w", err)
		}
		return newWrapFS(h), nil
	case "hdfs":
		h, err := hdfs.New(c.HDFS)
		if err != nil {
			return nil, fmt.Errorf("new hdfs: %w", err)
		}
		return newWrapFS(h), nil
	case "upyun":
		h, err := upyun.New(c.UPYUN)
		if err != nil {
			return nil, fmt.Errorf("new upyun: %w", err)
		}
		return h, nil
	case "obs":
		h, err := obs.New(c.OBS)
		if err != nil {
			return nil, fmt.Errorf("new obs: %w", err)
		}
		return h, nil
	case "minio", "s3":
		h, err := minio.New(c.MINIO)
		if err != nil {
			return nil, fmt.Errorf("new minio: %w", err)
		}
		return h, nil
	default:
		return nil, fmt.Errorf("mode %s not support", c.Mode)
	}
}

func New(c Config) (fileop.FileSystemSimpleBucket, error) {
	switch c.Mode {
	case "disk":
		h, err := afero.New(afero.Disk)
		if err != nil {
			return nil, fmt.Errorf("new disk: %w", err)
		}
		return newWrapFS(h), nil
	case "hdfs":
		h, err := hdfs.New(c.HDFS)
		if err != nil {
			return nil, fmt.Errorf("new hdfs: %w", err)
		}
		return newWrapFS(h), nil
	case "obs":
		h, err := obs.New(c.OBS)
		if err != nil {
			return nil, fmt.Errorf("new obs: %w", err)
		}
		return h, nil
	case "minio", "s3":
		h, err := minio.New(c.MINIO)
		if err != nil {
			return nil, fmt.Errorf("new minio: %w", err)
		}
		return h, nil
	default:
		return nil, fmt.Errorf("mode %s not support", c.Mode)
	}
}

type WrapFS struct {
	Source   fileop.FileSystemWithCloser
	Target   *filetarget.WrapFS
	BasePath string
}

func (w *WrapFS) Bucket(name string) fileop.FileSystemSimpleBucket {
	cp := *w
	cp.BasePath = name
	return &cp
}

func (w *WrapFS) Readdir(dirname string, n int) ([]fs.FileInfo, error) {
	dirname = filepath.Join(w.BasePath, dirname)
	return w.Source.Readdir(dirname, n)
}

func (w *WrapFS) Readdirnames(dirname string, n int) ([]string, error) {
	dirname = filepath.Join(w.BasePath, dirname)
	return w.Source.Readdirnames(dirname, n)
}

func (w *WrapFS) Open(name string) (io.ReadCloser, error) {
	name = filepath.Join(w.BasePath, name)
	return w.Source.Open(name)
}

func (w *WrapFS) Put(local, remote string) error {
	remote = filepath.Join(w.BasePath, remote)
	return w.Target.Put(local, remote)
}

func (w *WrapFS) PutStream(reader io.Reader, remote string) error {
	remote = filepath.Join(w.BasePath, remote)
	return w.Target.PutStream(reader, remote)
}

func (w *WrapFS) PutStreamWithContentType(reader io.Reader, remote string, _ string) error {
	remote = filepath.Join(w.BasePath, remote)
	return w.Target.PutStream(reader, remote)
}

func (w *WrapFS) PutEmpty(remote string) error {
	remote = filepath.Join(w.BasePath, remote)
	return w.Target.PutEmpty(remote)
}

func (w *WrapFS) Exist(remote string) bool {
	remote = filepath.Join(w.BasePath, remote)
	return w.Target.Exist(remote)
}

func (w *WrapFS) Close() error {
	return w.Source.Close()
}

func newWrapFS(remoteFS fileop.FileSystemWithCloser) *WrapFS {
	localFS, err := afero.New(afero.Disk)
	if err != nil {
		panic(err)
	}
	return &WrapFS{
		Source: remoteFS,
		Target: &filetarget.WrapFS{
			Target: remoteFS,
			Source: localFS,
		},
	}
}
