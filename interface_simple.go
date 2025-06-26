package fileop

import (
	"io"
)

// FileSystemSimple combines basic source reader and target uploader interfaces.
type FileSystemSimple interface {
	ISourceReader
	ITargetUploader
}

// FileSystemSimpleBucket extends FileSystemSimple with bucket operations
// and content type support for object storage systems.
type FileSystemSimpleBucket interface {
	FileSystemSimple
	Bucket(name string) FileSystemSimpleBucket
	PutStreamWithContentType(reader io.Reader, remote string, contentType string) error
}

// ITargetUploader provides upload operations for target file systems.
type ITargetUploader interface {
	io.Closer
	Put(local, remote string) error
	PutStream(reader io.Reader, remote string) error
	PutEmpty(remote string) error
	Exist(remote string) bool
}

// ISourceLister provides directory listing operations for source file systems.
type ISourceLister interface {
	io.Closer
	DirReader
}

// ISourceReader combines source listing and reading operations.
type ISourceReader interface {
	ISourceLister
	Reader
}
