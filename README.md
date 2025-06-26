# fileop

Collection of file operation helpers.

## Use Cases

This repository contains interface abstraction implementations for several commonly used file systems. It makes it easy to quickly adapt various file system operations in your business code.

## Basic Interfaces

[interface.go](interface.go)

**Read-only Operations**

```
type Reader interface {
	Open(name string) (io.ReadCloser, error)
}

type DirReader interface {
	Readdir(dirname string, n int) ([]fs.FileInfo, error)
	Readdirnames(dirname string, n int) ([]string, error)
}

type Stater interface {
	Stat(name string) (fs.FileInfo, error)
}

type Walker interface {
	Walk(root string, walkFn filepath.WalkFunc) error
}
```

**Create/Modify Operations**

```
type Writer interface {
	DirCreator
	Create(name string) (io.WriteCloser, error)
	Rename(oldPath, newPath string) error
}

type DirCreator interface {
	Mkdir(dirname string, perm fs.FileMode) error
	MkdirAll(dirname string, perm fs.FileMode) error
}

type Cleaner interface {
	Remove(name string) error
	RemoveAll(name string) error
}
```

File system operation interface needs to implement all the above methods

```
type FileSystem interface {
	Reader
	DirReader
	Writer
	Cleaner
	Stater
	Walker
}
```

### File System

Instantiation

```
filesystem.New(c Config) (fileop.FileSystemWithCloser, error)
```

Supported file systems:

- disk (local)
- hdfs (HDFS)

### Simple File System

[interface_simple.go](interface_simple.go)

Instantiation

```
simplefs.New(c Config) (fileop.FileSystemSimpleBucket, error) 
```

Supported file systems:

- disk (local)
- hdfs (HDFS)
- obs (Huawei OBS)
- minio (minio, S3 compatible)

## Helper Functions

### Compression

Supported compression file types:

- NONE
- GZIP
- ZLIB
- SNAPPY

```
NewCompressReader(file io.Reader, ct CompressType) (io.ReadCloser, error)
NewCompressWriter(buf io.Writer, ct CompressType) (io.WriteCloser, error)
```

#### Global Configuration for CompressWriter

- `fileop.UsePGZIP=[true|false]` 
  - Use "github.com/klauspost/pgzip" library for writing gzip files.
  - Default is "github.com/klauspost/compress/gzip" library.
- `fileop.PGZIPBlocks=[number]`
  - Set the number of PGZIP blocks.
  - Default is 4.

### File

File read/write

- [example/file_writer_reader](example/file_writer_reader/main.go)
- [example/fileutil_write_read](example/fileutil_write_read/main.go)

## How to Extend

- Extend a new data source [example/extend_filesource](example/extend_filesource/main.go) 
