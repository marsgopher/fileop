package afero

import (
	"fmt"
	"io"
	"io/fs"
	"path/filepath"

	"github.com/spf13/afero"
)

type FileSystemType int

const (
	Memory = iota
	Disk
)

func (ft FileSystemType) String() string {
	switch ft {
	case Memory:
		return "memory"
	case Disk:
		return "disk"
	default:
		return "unknown"
	}
}

type Handler struct {
	afero.Fs
}

func New(t FileSystemType) (*Handler, error) {
	switch t {
	case Memory:
		return &Handler{
			Fs: afero.NewMemMapFs(),
		}, nil
	case Disk:
		return &Handler{
			Fs: afero.NewOsFs(),
		}, nil
	default:
		return nil, fmt.Errorf("filesystem type %v not support", t)
	}
}

func (h *Handler) Open(name string) (io.ReadCloser, error) {
	return h.Fs.Open(name)
}

func (h *Handler) Create(name string) (io.WriteCloser, error) {
	return h.Fs.Create(name)
}

func (h *Handler) Walk(root string, walkFn filepath.WalkFunc) error {
	return afero.Walk(h.Fs, root, walkFn)
}

func (h *Handler) Readdirnames(dirname string, n int) ([]string, error) {
	dir, err := h.Fs.Open(dirname)
	if err != nil {
		return nil, err
	}
	defer func() { _ = dir.Close() }()
	return dir.Readdirnames(n)
}

func (h *Handler) Readdir(dirname string, n int) ([]fs.FileInfo, error) {
	dir, err := h.Fs.Open(dirname)
	if err != nil {
		return nil, err
	}
	defer func() { _ = dir.Close() }()
	return dir.Readdir(n)
}

func (h *Handler) Close() error {
	return nil
}
