package upyun

import (
	"io/fs"
	"time"

	"github.com/upyun/go-sdk/v3/upyun"
)

type fileInfo struct {
	*upyun.FileInfo
}

func (f fileInfo) Name() string {
	return f.FileInfo.Name
}

func (f fileInfo) Size() int64 {
	return f.FileInfo.Size
}

func (f fileInfo) Mode() fs.FileMode {
	if f.IsDir() {
		return 0777
	}
	return 0666
}

func (f fileInfo) ModTime() time.Time {
	return f.FileInfo.Time
}

func (f fileInfo) IsDir() bool {
	return f.FileInfo.IsDir
}

func (f fileInfo) Sys() interface{} {
	return nil
}
