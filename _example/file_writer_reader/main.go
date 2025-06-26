package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/marsgopher/fileop"
	"github.com/marsgopher/fileop/filesystem"
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

	if err := write(fs, tmpFilepath); err != nil {
		log.Fatal(err)
	}

	if err := read(fs, tmpFilepath); err != nil {
		log.Fatal(err)
	}
}

func write(fs fileop.FileWriterInterface, path string) error {
	fw, err := fileop.NewFileWriter(fs, path, 0, fileop.NONE)
	if err != nil {
		return fmt.Errorf("new file writer: %w", err)
	}
	defer func() {
		if fw != nil {
			_ = fw.Close()
		}
	}()

	const content = `Hello`
	// writer sth
	if i, err := fw.Write([]byte(content)); err != nil {
		return fmt.Errorf("write: %w", err)
	} else {
		log.Printf("Write %d bytes", i)
	}

	defer func() { fw = nil }()
	return fw.Close()
}

func read(fs fileop.FileReaderInterface, path string) error {
	fr, err := fileop.NewFileReader(fs, path, fileop.NONE)
	if err != nil {
		if fileop.IsUnhandledFileReaderError(err) {
			return nil
		}
		return fmt.Errorf("new file reader: %w", err)
	}
	defer func() {
		_ = fr.Close()
	}()

	buf := make([]byte, 10)
	if i, err := fr.Read(buf); err != nil {
		return fmt.Errorf("read: %w", err)
	} else {
		log.Printf("Read %d bytes: %s", i, buf[:i])
	}

	return nil
}
