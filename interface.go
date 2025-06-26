// Package fileop provides a collection of file operation helpers with
// abstraction interfaces for various file systems.
package fileop

import (
	"io"
	"io/fs"
	"path/filepath"
)

// Reader provides read operations for files.
type Reader interface {
	Open(name string) (io.ReadCloser, error)
}

// DirReader provides directory reading operations.
type DirReader interface {
	Readdir(dirname string, n int) ([]fs.FileInfo, error)
	Readdirnames(dirname string, n int) ([]string, error)
}

// DirCreator provides directory creation operations.
type DirCreator interface {
	Mkdir(dirname string, perm fs.FileMode) error
	MkdirAll(dirname string, perm fs.FileMode) error
}

// Stater provides file stat operations.
type Stater interface {
	Stat(name string) (fs.FileInfo, error)
}

// Walker provides file tree walking operations.
type Walker interface {
	Walk(root string, walkFn filepath.WalkFunc) error
}

// Writer provides write operations for files and directories.
type Writer interface {
	DirCreator
	Create(name string) (io.WriteCloser, error)
	Rename(oldPath, newPath string) error
}

// Cleaner provides file and directory removal operations.
type Cleaner interface {
	Remove(name string) error
	RemoveAll(name string) error
}

// FileSystem combines all file system operations.
type FileSystem interface {
	Reader
	DirReader
	Writer
	Cleaner
	Stater
	Walker
}

// FileSystemWithCloser extends FileSystem with a Close method.
type FileSystemWithCloser interface {
	FileSystem
	io.Closer
}
