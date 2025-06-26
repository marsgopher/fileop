package filesystem

import (
	"fmt"

	"github.com/marsgopher/fileop"
	"github.com/marsgopher/fileop/integration/afero"
	"github.com/marsgopher/fileop/integration/hdfs"
)

type Config struct {
	Mode string      `mapstructure:"mode"`
	HDFS hdfs.Config `mapstructure:"hdfs"`
}

func New(c Config) (fileop.FileSystemWithCloser, error) {
	switch c.Mode {
	case "disk":
		h, err := afero.New(afero.Disk)
		if err != nil {
			return nil, fmt.Errorf("new disk: %w", err)
		}
		return h, nil
	case "hdfs":
		h, err := hdfs.New(c.HDFS)
		if err != nil {
			return nil, fmt.Errorf("new hdfs: %w", err)
		}
		return h, nil
	default:
		return nil, fmt.Errorf("mode %s not support", c.Mode)
	}
}
