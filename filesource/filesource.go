package filesource

import (
	"fmt"

	"github.com/marsgopher/fileop"
	"github.com/marsgopher/fileop/integration/afero"
	"github.com/marsgopher/fileop/integration/hdfs"
	"github.com/marsgopher/fileop/integration/minio"
	"github.com/marsgopher/fileop/integration/obs"
	"github.com/marsgopher/fileop/integration/upyun"
)

type Config struct {
	Mode  string       `mapstructure:"mode"`
	OBS   obs.Config   `mapstructure:"obs"`
	UPYUN upyun.Config `mapstructure:"upyun"`
	HDFS  hdfs.Config  `mapstructure:"hdfs"`
	MINIO minio.Config `mapstructure:"minio"`
}

func New(c Config) (fileop.ISourceReader, error) {
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
	case "upyun":
		h, err := upyun.New(c.UPYUN)
		if err != nil {
			return nil, fmt.Errorf("new upyun: %w", err)
		}
		return h, nil
	case "obs":
		h, err := obs.New(c.OBS)
		if err != nil {
			return nil, fmt.Errorf("new obs: %w", err)
		}
		return h, nil
	case "minio":
		h, err := minio.New(c.MINIO)
		if err != nil {
			return nil, fmt.Errorf("new minio: %w", err)
		}
		return h, nil
	default:
		return nil, fmt.Errorf("mode %s not support", c.Mode)
	}
}
