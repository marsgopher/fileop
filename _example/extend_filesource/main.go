package main

import (
	"fmt"
	"io"
	"io/fs"

	"github.com/marsgopher/fileop"
	"github.com/marsgopher/fileop/filesource"
)

type Config struct {
	// ",squash" makes extended config at the same level as base config
	// https://pkg.go.dev/github.com/mitchellh/mapstructure#hdr-Embedded_Structs_and_Squashing
	filesource.Config `mapstructure:",squash"`

	// Extended data source configuration
	ExtCfg ExtSourceConfig `mapstructure:"extend_source"`
}

type ExtSourceConfig struct{}
type ExtSource struct{}

func (e ExtSource) Close() error {
	fmt.Println("call Close")
	return nil
}
func (e ExtSource) Readdir(_ string, _ int) ([]fs.FileInfo, error) {
	fmt.Println("call Readdir")
	return nil, nil
}
func (e ExtSource) Readdirnames(_ string, _ int) ([]string, error) {
	fmt.Println("call Readdirnames")
	return nil, nil
}
func (e ExtSource) Open(_ string) (io.ReadCloser, error) {
	fmt.Println("call Open")
	return nil, nil
}
func newExtSource(_ ExtSourceConfig) (*ExtSource, error) {
	return &ExtSource{}, nil
}

func New(c Config) (fileop.ISourceReader, error) {
	switch c.Mode {
	case "extend_source":
		return newExtSource(c.ExtCfg)
	default:
		return filesource.New(c.Config)
	}
}

func main() {
	cfg := Config{}
	cfg.Mode = "extend_source"
	cfg.ExtCfg = ExtSourceConfig{}

	src, err := New(cfg)
	if err != nil {
		panic(err)
	}
	defer func() { _ = src.Close() }()

	_, _ = src.Open("")
	_, _ = src.Readdir("", 0)
	_, _ = src.Readdirnames("", 0)
}
