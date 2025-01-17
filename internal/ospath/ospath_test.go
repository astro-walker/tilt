package ospath

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/astro-walker/tilt/internal/testutils/tempdir"
)

func TestChild(t *testing.T) {
	f := NewOspathFixture(t)

	paths := []string{
		filepath.Join("parent", "fileA"),
		filepath.Join("parent", "child", "fileB"),
		filepath.Join("parent", "child", "grandchild", "fileC"),
		filepath.Join("sibling", "fileD"),
	}
	f.TouchFiles(paths)

	f.assertChild("parent", filepath.Join("sibling", "fileD"), "")
	f.assertChild(filepath.Join("parent", "child"), filepath.Join("parent", "fileA"), "")
	f.assertChild("parent", filepath.Join("parent", "fileA"), "fileA")
	f.assertChild("parent", filepath.Join("parent", "child", "fileB"), filepath.Join("child", "fileB"))
	f.assertChild("parent", filepath.Join("parent", "child", "grandchild", "fileC"), filepath.Join("child", "grandchild", "fileC"))

	f.assertChild(filepath.Join("parent", "child"), filepath.Join("parent", "child", "fileB"), "fileB")

	f.assertChild("parent", "parent", ".")
}

func TestCaseInsensitiveFileSystem(t *testing.T) {
	f := NewOspathFixture(t)

	fileA := filepath.Join("parent", "fileA")
	f.TouchFiles([]string{fileA})

	// Assume that macOS and Windows are case-insensitive, and other operating
	// systems are not. This isn't strictly accurate, but is good enough for
	// testing.
	if runtime.GOOS == "darwin" || runtime.GOOS == "windows" {
		f.assertChild("Parent", fileA, "fileA")
	} else {
		f.assertChild("Parent", fileA, "")
	}
}

func TestIsBrokenSymlink(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("windows does not support user-land symlinks")
	}
	f := NewOspathFixture(t)

	f.TouchFiles([]string{
		"fileA",
		"child/fileB",
		"child/grandchild/fileC",
	})

	f.symlink("fileA", "symlinkFileA")
	f.symlink("fileB", "symlinkFileB")

	f.assertBrokenSymlink("fileA", false)
	f.assertBrokenSymlink("fileB", false)
	f.assertBrokenSymlink("child/fileB", false)
	f.assertBrokenSymlink("child/grandchild/fileC", false)
	f.assertBrokenSymlink("symlinkFileA", false)
	f.assertBrokenSymlink("symlinkFileB", true)
}

func TestInvalidDir(t *testing.T) {
	f := NewOspathFixture(t)

	// Passing "" as dir used to make Child hang forever. Let's make sure it doesn't do that.
	_, isChild := Child("", "random")
	if isChild {
		f.t.Fatalf("Expected file 'random' to NOT be a child of empty dir")
	}
}

func TestDirTrailingSlash(t *testing.T) {
	f := NewOspathFixture(t)

	f.TouchFiles([]string{filepath.Join("some", "dir", "file")})

	// Should work regardless of whether directory has trailing slash
	f.assertChild(filepath.Join("some", "dir"),
		filepath.Join("some", "dir", "file"), "file")
	f.assertChild(filepath.Join("some", "dir")+string(filepath.Separator),
		filepath.Join("some", "dir", "file"), "file")
}

func TestTryAsCwdChildren(t *testing.T) {
	f := NewOspathFixture(t)
	f.Chdir()

	results := TryAsCwdChildren([]string{f.Path()})

	if len(results) == 0 {
		t.Fatal("Expected 1 result, got 0")
	}

	r := results[0]

	if r != "." {
		t.Errorf("Expected %s to equal \".\"", r)
	}
}

type OspathFixture struct {
	*tempdir.TempDirFixture
	t *testing.T
}

func NewOspathFixture(t *testing.T) *OspathFixture {
	ret := &OspathFixture{
		TempDirFixture: tempdir.NewTempDirFixture(t),
		t:              t,
	}
	return ret
}

// pass `expectedRelative` = "" to indicate that `file` is NOT a child of `dir`
func (f *OspathFixture) assertChild(dir, file, expectedRel string) {
	dir = f.JoinPath(dir)
	file = f.JoinPath(file)
	rel, isChild := Child(dir, file)
	if expectedRel == "" {
		if isChild {
			f.t.Fatalf("Expected file '%s' to NOT be a child of dir '%s'", file, dir)
		}
		return
	}

	if !isChild {
		f.t.Fatalf("Expected file '%s' to be a child of dir '%s', but got !isChild", file, dir)
	}

	if rel != expectedRel {
		f.t.Fatalf("Expected relative path of '%s' to dir '%s' to be: '%s'. Actual: '%s'.", file, dir, expectedRel, rel)
	}
}

func (f *OspathFixture) symlink(oldPath, newPath string) {
	oldPath = filepath.Join(f.Path(), oldPath)
	newPath = filepath.Join(f.Path(), newPath)
	err := os.Symlink(oldPath, newPath)
	if err != nil {
		f.t.Fatal(err)
	}
}

func (f *OspathFixture) assertBrokenSymlink(file string, expected bool) {
	broken, err := IsBrokenSymlink(filepath.Join(f.Path(), file))
	if err != nil {
		f.t.Fatal(err)
	}

	if broken != expected {
		if broken {
			f.t.Fatalf("Expected a regular file or working symlink: %s", file)
		} else {
			f.t.Fatalf("Expected a broken symlink: %s", file)
		}
	}
}
