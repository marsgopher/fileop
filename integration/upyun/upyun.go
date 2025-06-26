package upyun

import (
	"io"
	"io/fs"
	"mime"
	"net/http"
	"os"
	"path/filepath"

	"github.com/marsgopher/common/concurrency"
	"github.com/upyun/go-sdk/v3/upyun"
)

type Client struct {
	*upyun.UpYun
}

func New(c Config) (*Client, error) {
	upFS := upyun.NewUpYun(&upyun.UpYunConfig{
		Bucket:    c.Bucket,
		Operator:  c.Operator,
		Password:  c.Password,
		Hosts:     c.Hosts,
		UserAgent: c.UserAgent,
	})
	return &Client{UpYun: upFS}, nil
}

func (w *Client) Put(localPath, remotePath string) error {
	return w.UpYun.Put(&upyun.PutObjectConfig{
		LocalPath: localPath,
		Path:      remotePath,
		Headers:   map[string]string{"Content-Type": "application/octet-stream"},
		UseMD5:    true,
	})
}

func (w *Client) PutStream(reader io.Reader, remotePath string) error {
	return w.UpYun.Put(&upyun.PutObjectConfig{
		Reader:  reader,
		Path:    remotePath,
		Headers: map[string]string{"Content-Type": "application/octet-stream"},
		UseMD5:  true,
	})
}

func (w *Client) PutStreamWithContentType(reader io.Reader, remotePath string, contentType string) error {
	if contentType == "" {
		// try fix content type
		contentType = mime.TypeByExtension(filepath.Ext(remotePath))
	}
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	return w.UpYun.Put(&upyun.PutObjectConfig{
		Reader:  reader,
		Path:    remotePath,
		Headers: map[string]string{"Content-Type": contentType},
		UseMD5:  true,
	})
}

func (w *Client) PutEmpty(remote string) error {
	return w.UpYun.Put(&upyun.PutObjectConfig{
		Path: remote,
	})
}

func (w *Client) PutFinish(remote string) error {
	target := remote + ".finish"
	return w.UpYun.Put(&upyun.PutObjectConfig{
		Path: target,
	})
}

func (w *Client) Exist(remote string) bool {
	_, err := w.UpYun.GetInfo(remote)
	return err == nil
}

func (w *Client) Open(name string) (io.ReadCloser, error) {
	rd, wt := io.Pipe()
	go func() {
		_, err := w.Get(&upyun.GetObjectConfig{
			Path:   name,
			Writer: wt,
		})
		_ = wt.CloseWithError(err)
	}()
	return rd, nil
}

func (w *Client) Readdirnames(name string, n int) ([]string, error) {
	objCh := make(chan *upyun.FileInfo)
	var res []string
	wg := concurrency.NewSemaErrGroup(1)
	wg.Do(func() error {
		if err := w.List(&upyun.GetObjectsConfig{
			Path:           name,
			ObjectsChan:    objCh,
			MaxListObjects: n,
		}); err != nil {
			if e, ok := err.(*upyun.Error); ok {
				if e.StatusCode == http.StatusNotFound {
					return os.ErrNotExist
				}
			}

			return err
		}
		return nil
	})
	for obj := range objCh {
		res = append(res, obj.Name)
	}

	return res, wg.Wait()
}

func (w *Client) Readdir(name string, n int) ([]fs.FileInfo, error) {
	objCh := make(chan *upyun.FileInfo)
	var res []fs.FileInfo
	wg := concurrency.NewSemaErrGroup(1)
	wg.Do(func() error {
		if err := w.List(&upyun.GetObjectsConfig{
			Path:           name,
			ObjectsChan:    objCh,
			MaxListObjects: n,
		}); err != nil {
			if e, ok := err.(*upyun.Error); ok {
				if e.StatusCode == http.StatusNotFound {
					return os.ErrNotExist
				}
			}

			return err
		}
		return nil
	})
	for obj := range objCh {
		res = append(res, fileInfo{obj})
	}

	return res, wg.Wait()
}

func (w *Client) Close() error {
	return nil
}
