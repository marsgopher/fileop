package hdfs

import (
	"time"
)

type Config struct {
	// for client setup
	User         string   `mapstructure:"user"`
	NameNodes    []string `mapstructure:"namenodes"`
	OldNameNodes []string `mapstructure:"name_nodes"` // for backward compatibility

	// for create file
	Replication int   `mapstructure:"replication"`
	BlockSize   int64 `mapstructure:"block_size"`

	Timeout time.Duration `mapstructure:"timeout"`
}
