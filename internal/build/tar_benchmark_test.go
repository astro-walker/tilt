package build

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/astro-walker/tilt/internal/ignore"
	"github.com/astro-walker/tilt/internal/testutils/tempdir"
)

func BenchmarkArchivePaths(b *testing.B) {
	f := tempdir.NewTempDirFixture(b)

	fileCount := 10000
	for i := 0; i < fileCount; i++ {
		dir := "dirA"
		if i%2 == 0 {
			dir = "dirB"
		}

		filename := fmt.Sprintf("file%d", i)
		f.WriteFile(filepath.Join(dir, filename), "contents")
	}

	b.ResetTimer()

	run := func() {
		writer := &bytes.Buffer{}
		filter, err := ignore.NewDirectoryMatcher(f.JoinPath(f.Path(), "dirA"))
		assert.NoError(b, err)

		builder := NewArchiveBuilder(writer, filter)
		err = builder.ArchivePathsIfExist(context.Background(), []PathMapping{
			{
				LocalPath:     f.Path(),
				ContainerPath: "/",
			},
		})
		assert.NoError(b, err)
		err = builder.Close()
		assert.NoError(b, err)
	}
	for i := 0; i < b.N; i++ {
		run()
	}
}
