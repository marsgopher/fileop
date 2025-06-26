package minio

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"mime"
	"os"
	"path/filepath"
	"time"

	"github.com/marsgopher/fileop"
	minio "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Client struct {
	*minio.Client
	bucket string
}

func New(cfg Config) (*Client, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AK, cfg.SK, cfg.Token),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("new client: %w", err)
	}
	c := &Client{
		Client: client,
		bucket: cfg.Bucket,
	}
	return c, nil
}

func (c *Client) Bucket(name string) fileop.FileSystemSimpleBucket {
	cp := *c
	cp.bucket = name
	return &cp
}

func (c *Client) Put(localPath, remotePath string) error {
	ctx := context.Background()
	opts := minio.PutObjectOptions{}
	if _, err := c.Client.FPutObject(ctx, c.bucket, remotePath, localPath, opts); err != nil {
		return fmt.Errorf("put %s: %w", remotePath, err)
	}
	return nil
}

func (c *Client) PutStream(rd io.Reader, remotePath string) error {
	ctx := context.Background()
	opts := minio.PutObjectOptions{}
	if opts.ContentType = mime.TypeByExtension(filepath.Ext(remotePath)); opts.ContentType == "" {
		opts.ContentType = "application/octet-stream"
	}
	if _, err := c.Client.PutObject(ctx, c.bucket, remotePath, rd, -1, opts); err != nil {
		return fmt.Errorf("put %s: %w", remotePath, err)
	}
	return nil
}

func (c *Client) PutStreamWithContentType(rd io.Reader, remotePath string, contentType string) error {
	ctx := context.Background()
	opts := minio.PutObjectOptions{}
	if contentType == "" {
		// try fix content type
		contentType = mime.TypeByExtension(filepath.Ext(remotePath))
	}
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	opts.ContentType = contentType
	if _, err := c.Client.PutObject(ctx, c.bucket, remotePath, rd, -1, opts); err != nil {
		return fmt.Errorf("put %s: %w", remotePath, err)
	}
	return nil
}

func (c *Client) PutEmpty(remotePath string) error {
	ctx := context.Background()
	opts := minio.PutObjectOptions{}
	if _, err := c.Client.PutObject(ctx, c.bucket, remotePath, nil, 0, opts); err != nil {
		return fmt.Errorf("put %s: %w", remotePath, err)
	}
	return nil
}

func (c *Client) Exist(remotePath string) bool {
	ctx := context.Background()
	opts := minio.StatObjectOptions{}
	_, err := c.Client.StatObject(ctx, c.bucket, remotePath, opts)
	return err == nil
}

func (c *Client) Close() error {
	return nil
}

func (c *Client) Open(name string) (io.ReadCloser, error) {
	object, err := c.Client.GetObject(context.Background(), c.bucket, name, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}

	return object, nil
}

func (c *Client) Readdir(dirname string, n int) ([]fs.FileInfo, error) {
	if n < 0 {
		n = 0
	}

	objectsCh := c.Client.ListObjects(context.Background(), c.bucket, minio.ListObjectsOptions{
		Prefix:    dirname,
		Recursive: false,
		MaxKeys:   n,
	})

	fileInfos := make([]fs.FileInfo, 0, n)
	for obj := range objectsCh {
		fileInfo := &minioFileInfo{
			name:    obj.Key,
			size:    obj.Size,
			modTime: obj.LastModified,
			isDir:   obj.Key[len(obj.Key)-1] == '/',
		}
		fileInfos = append(fileInfos, fileInfo)
	}

	return fileInfos, nil
}

func (c *Client) Readdirnames(dirname string, n int) ([]string, error) {
	if n < 0 {
		n = 0
	}

	objectsCh := c.Client.ListObjects(context.Background(), c.bucket, minio.ListObjectsOptions{
		Prefix:    dirname,
		Recursive: false,
		MaxKeys:   n,
	})

	names := make([]string, 0, n)
	for obj := range objectsCh {
		names = append(names, obj.Key)
	}

	return names, nil
}

type minioFileInfo struct {
	name    string
	size    int64
	modTime time.Time
	isDir   bool
}

func (f *minioFileInfo) Mode() fs.FileMode {
	if f.IsDir() {
		return os.ModeDir | 0755
	}
	return 0644
}

func (f *minioFileInfo) Name() string       { return f.name }
func (f *minioFileInfo) Size() int64        { return f.size }
func (f *minioFileInfo) ModTime() time.Time { return f.modTime }
func (f *minioFileInfo) IsDir() bool        { return f.isDir }
func (f *minioFileInfo) Sys() interface{}   { return nil }
