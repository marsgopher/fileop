package main

import (
	"bytes"
	"log"
	"os"
	"path/filepath"

	"github.com/marsgopher/fileop/filesystem"
	"github.com/marsgopher/fileop/fileutil"
)

func main() {
	cfg := filesystem.Config{
		Mode: "disk",
	}
	fs, err := filesystem.New(cfg)
	if err != nil {
		log.Fatalf("new filesystem: %v", err)
	}

	tmpFilepath := filepath.Join(os.TempDir(), "fileop_test")

	rd := bytes.NewReader([]byte(`Hello`))

	if err := fileutil.WriteFile(fs, tmpFilepath, rd); err != nil {
		log.Fatal(err)
	}

	if content, err := fileutil.ReadFile(fs, tmpFilepath); err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Read: %s", content)
	}
}
