package obs

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"mime"
	"os"
	"path/filepath"
	"time"

	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
	"github.com/marsgopher/fileop"
)

var ErrNotSetup = errors.New("not setup")

type GetAclInputFunc func(key string) *obs.SetObjectAclInput

type Client struct {
	*obs.ObsClient
	getAclInput GetAclInputFunc
	bucket      string
}

func New(cfg Config) (*Client, error) {
	if cfg.Endpoint == "" {
		return nil, ErrNotSetup
	}

	obsClient, err := obs.New(cfg.AK, cfg.SK, cfg.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("new obs client: %w", err)
	}

	getAclInput := func(key string) *obs.SetObjectAclInput {
		var aclInput *obs.SetObjectAclInput
		if cfg.ACLControlID != "" {
			aclInput = &obs.SetObjectAclInput{}
			aclInput.Bucket = cfg.Bucket
			aclInput.Key = key

			if cfg.ACLOwnerID != "" {
				aclInput.Owner.ID = cfg.ACLOwnerID
			} else {
				aclInput.Owner.ID = cfg.ACLControlID
			}
			aclInput.Grants = []obs.Grant{
				{
					Grantee: obs.Grantee{
						Type: obs.GranteeUser,
						ID:   cfg.ACLControlID,
					},
					Permission: obs.PermissionFullControl,
				},
			}
		}
		return aclInput
	}

	return &Client{
		ObsClient:   obsClient,
		getAclInput: getAclInput,
		bucket:      cfg.Bucket,
	}, nil
}

func (c *Client) Bucket(name string) fileop.FileSystemSimpleBucket {
	cp := *c
	cp.bucket = name
	return &cp
}

func (c *Client) Put(localPath, remotePath string) error {
	input := &obs.PutFileInput{}
	input.Bucket = c.bucket
	input.Key = remotePath
	input.SourceFile = localPath

	if _, err := c.ObsClient.PutFile(input); err != nil {
		return fmt.Errorf("put %s: %w", remotePath, err)
	}

	if aclInput := c.getAclInput(remotePath); aclInput != nil {
		if _, err := c.ObsClient.SetObjectAcl(aclInput); err != nil {
			return fmt.Errorf("set acl %s: %w", remotePath, err)
		}
	}
	return nil
}

func (c *Client) PutStream(reader io.Reader, remotePath string) error {
	input := &obs.PutObjectInput{}
	input.Bucket = c.bucket
	input.Key = remotePath
	input.Body = reader

	if _, err := c.ObsClient.PutObject(input); err != nil {
		return fmt.Errorf("put %s: %w", remotePath, err)
	}

	if aclInput := c.getAclInput(remotePath); aclInput != nil {
		if _, err := c.ObsClient.SetObjectAcl(aclInput); err != nil {
			return fmt.Errorf("set acl %s: %w", remotePath, err)
		}
	}
	return nil
}

func (c *Client) PutStreamWithContentType(reader io.Reader, remotePath string, contentType string) error {
	if contentType == "" {
		// try fix content type
		contentType = mime.TypeByExtension(filepath.Ext(remotePath))
	}
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	input := &obs.PutObjectInput{}
	input.Bucket = c.bucket
	input.Key = remotePath
	input.Body = reader
	input.ContentType = contentType

	if _, err := c.ObsClient.PutObject(input); err != nil {
		return fmt.Errorf("put %s: %w", remotePath, err)
	}

	if aclInput := c.getAclInput(remotePath); aclInput != nil {
		if _, err := c.ObsClient.SetObjectAcl(aclInput); err != nil {
			return fmt.Errorf("set acl %s: %w", remotePath, err)
		}
	}
	return nil
}

func (c *Client) PutEmpty(remotePath string) error {
	input := &obs.PutFileInput{}
	input.Bucket = c.bucket
	input.Key = remotePath
	aclInput := c.getAclInput(remotePath)

	if _, err := c.ObsClient.PutObject(&obs.PutObjectInput{
		PutObjectBasicInput: input.PutObjectBasicInput,
	}); err != nil {
		return fmt.Errorf("put %s: %w", remotePath, err)
	}

	if aclInput != nil {
		if _, err := c.ObsClient.SetObjectAcl(aclInput); err != nil {
			return fmt.Errorf("set acl %s: %w", remotePath, err)
		}
	}
	return nil
}

// PutFinish (deprecated)
// use PutEmpty(remotePath + ".finish") instead
func (c *Client) PutFinish(remotePath string) error {
	target := remotePath + ".finish"
	return c.PutEmpty(target)
}

func (c *Client) Exist(path string) bool {
	input := &obs.GetObjectMetadataInput{}
	input.Bucket = c.bucket
	input.Key = path
	_, err := c.GetObjectMetadata(input)
	return err == nil
}

func (c *Client) Close() error {
	c.ObsClient.Close()
	return nil
}

func (c *Client) Readdir(dirname string, n int) ([]fs.FileInfo, error) {
	input := &obs.ListObjectsInput{}
	input.Bucket = c.bucket
	input.Prefix = dirname
	input.MaxKeys = n
	input.Delimiter = "/"

	output, err := c.ListObjects(input)
	if err != nil {
		return nil, err
	}

	fileInfos := make([]fs.FileInfo, 0, len(output.Contents)+len(output.CommonPrefixes))
	for _, object := range output.Contents {
		fileInfo := &obsFileInfo{
			name:    object.Key,
			size:    object.Size,
			modTime: object.LastModified,
			isDir:   false,
		}
		fileInfos = append(fileInfos, fileInfo)
	}

	for _, prefix := range output.CommonPrefixes {
		fileInfo := &obsFileInfo{
			name:    prefix,
			isDir:   true,
			modTime: time.Now(),
		}
		fileInfos = append(fileInfos, fileInfo)
	}

	return fileInfos, nil
}

func (c *Client) Readdirnames(dirname string, n int) ([]string, error) {
	input := &obs.ListObjectsInput{}
	input.Bucket = c.bucket
	input.Prefix = dirname
	input.MaxKeys = n
	input.Delimiter = "/"

	output, err := c.ListObjects(input)
	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(output.Contents)+len(output.CommonPrefixes))
	for _, object := range output.Contents {
		names = append(names, object.Key)
	}

	names = append(names, output.CommonPrefixes...)

	return names, nil
}

func (c *Client) Open(name string) (io.ReadCloser, error) {
	input := &obs.GetObjectInput{}
	input.Bucket = c.bucket
	input.Key = name

	output, err := c.GetObject(input)
	if err != nil {
		return nil, err
	}
	return output.Body, nil
}

type obsFileInfo struct {
	name    string
	size    int64
	modTime time.Time
	isDir   bool
}

func (f *obsFileInfo) Name() string       { return f.name }
func (f *obsFileInfo) Size() int64        { return f.size }
func (f *obsFileInfo) ModTime() time.Time { return f.modTime }
func (f *obsFileInfo) IsDir() bool        { return f.isDir }
func (f *obsFileInfo) Mode() fs.FileMode {
	if f.IsDir() {
		return os.ModeDir | 0755
	}
	return 0644
}
func (f *obsFileInfo) Sys() interface{} { return nil }
