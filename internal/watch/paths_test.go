package watch

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/astro-walker/tilt/internal/testutils/tempdir"
)

func TestGreatestExistingAncestor(t *testing.T) {
	f := tempdir.NewTempDirFixture(t)

	p, err := greatestExistingAncestor(f.Path())
	assert.NoError(t, err)
	assert.Equal(t, f.Path(), p)

	p, err = greatestExistingAncestor(f.JoinPath("missing"))
	assert.NoError(t, err)
	assert.Equal(t, f.Path(), p)

	missingTopLevel := "/missingDir/a/b/c"
	if runtime.GOOS == "windows" {
		missingTopLevel = "C:\\missingDir\\a\\b\\c"
	}
	_, err = greatestExistingAncestor(missingTopLevel)
	assert.Contains(t, err.Error(), "cannot watch root directory")
}
