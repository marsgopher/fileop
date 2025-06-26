package fileop

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/marsgopher/fileop/integration/afero"
)

func TestListDirAfterConcurrencyFileWrite(t *testing.T) {
	t.Parallel()
	assert := require.New(t)

	const dir = "test_dir"
	const n = 10000
	mfs, err := afero.New(afero.Memory)
	assert.NoError(err)

	allFilePaths := make([]string, 0, n)

	// run concurrency test
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		fp := filepath.Join(
			dir,
			fmt.Sprintf("%02d", n%10),
			fmt.Sprintf("%d.txt", i),
		)
		allFilePaths = append(allFilePaths, fp)

		wg.Add(1)
		go func() {
			defer wg.Done()

			wt, err := NewFileWriter(mfs, fp, 0, NONE)
			assert.NoError(err)
			defer func() {
				assert.NoError(wt.Close())
			}()

			// write 30 bytes
			for j := 0; j < 10; j++ {
				_, err := wt.Write([]byte("000"))
				assert.NoError(err)
			}
		}()
	}
	wg.Wait()

	// Test1: find all files by full path access
	for _, fp := range allFilePaths {
		info, err := mfs.Stat(fp)
		assert.NoErrorf(err, "stat file %s", fp)
		assert.Equalf(int64(30), info.Size(), "file %s size unmatch", fp)
	}

	// Test2: find all files by walk
	foundFiles := make([]string, 0, n)
	wErr := mfs.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		assert.NoError(err)
		if info.IsDir() {
			return nil // skip dir
		}
		if strings.HasSuffix(info.Name(), ".txt") {
			foundFiles = append(foundFiles, path)
		}
		return nil
	})
	assert.NoError(wErr, "walk")
	assert.Equal(n, len(foundFiles), "missing files")
}
