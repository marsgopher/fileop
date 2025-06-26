package hdfs

import (
	"fmt"
	"io"
	"io/fs"
	"net"

	hdfs "github.com/colinmarc/hdfs/v2"
)

const userDefault = "root"

type Handler struct {
	*hdfs.Client
	cfg Config
}

func New(c Config) (*Handler, error) {
	user := userDefault
	if u := c.User; u != "" {
		user = u
	}

	if len(c.OldNameNodes) > 0 {
		c.NameNodes = append(c.NameNodes, c.OldNameNodes...)
	}
	opts := hdfs.ClientOptions{
		Addresses: c.NameNodes,
		User:      user,
	}
	if t := c.Timeout; t != 0 {
		dialFunc := (&net.Dialer{
			Timeout:   t,
			KeepAlive: t,
		}).DialContext
		opts.NamenodeDialFunc = dialFunc
		opts.DatanodeDialFunc = dialFunc
	}
	client, err := hdfs.NewClient(opts)
	if err != nil {
		return nil, err
	}

	h := &Handler{
		cfg:    c,
		Client: client,
	}
	return h, nil
}

func (h *Handler) Create(name string) (io.WriteCloser, error) {
	// NOTE: overwrite exist file
	if _, err := h.Stat(name); err == nil {
		if err := h.Remove(name); err != nil {
			return nil, fmt.Errorf("remove old file: %w", err)
		}
	}

	var fd *hdfs.FileWriter
	var err error

	c := h.cfg
	if c.Replication == 0 || c.BlockSize == 0 {
		fd, err = h.Client.Create(name)
		if err != nil {
			return nil, fmt.Errorf("create file: %w", err)
		}
	} else {
		fd, err = h.Client.CreateFile(name, c.Replication, c.BlockSize, 0644)
		if err != nil {
			return nil, fmt.Errorf("create file: %w", err)
		}
	}

	return fd, nil
}

func (h *Handler) Open(name string) (io.ReadCloser, error) {
	f, err := h.Client.Open(name)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func (h *Handler) Readdirnames(dirname string, n int) ([]string, error) {
	dir, err := h.Client.Open(dirname)
	if err != nil {
		return nil, err
	}
	defer func() { _ = dir.Close() }()
	return dir.Readdirnames(n)
}

func (h *Handler) Readdir(dirname string, n int) ([]fs.FileInfo, error) {
	dir, err := h.Client.Open(dirname)
	if err != nil {
		return nil, err
	}
	defer func() { _ = dir.Close() }()
	return dir.Readdir(n)
}
