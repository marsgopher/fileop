package afero

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMemFsMkdirInConcurrency(t *testing.T) {
	t.Parallel()
	assert := require.New(t)

	const dir = "test_dir"
	const n = 1000

	for i := 0; i < n; i++ {
		fs, err := New(Memory)
		assert.NoError(err)
		c1 := make(chan error, 1)
		c2 := make(chan error, 1)

		go func() {
			c1 <- fs.Mkdir(dir, 0755)
		}()
		go func() {
			c2 <- fs.Mkdir(dir, 0755)
		}()

		// Only one attempt of creating the directory should succeed.
		err1 := <-c1
		err2 := <-c2
		assert.NotEqualf(err1, err2, "run #%v, more than one success", i)
	}
}
