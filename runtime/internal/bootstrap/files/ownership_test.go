package files

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"
)

func TestChownRecursive_AppliesOwnershipToRootAndChildren(t *testing.T) {
	restore := replaceChownSeam(t)
	defer restore()

	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, "nested"), 0o755); err != nil {
		t.Fatalf("mkdir nested: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "nested", "file.txt"), []byte("content"), 0o600); err != nil {
		t.Fatalf("write file: %v", err)
	}
	var paths []string
	chown = func(path string, uid int, gid int) error {
		if uid != 1000 || gid != 1001 {
			t.Fatalf("uid/gid = %d/%d, want 1000/1001", uid, gid)
		}
		paths = append(paths, path)
		return nil
	}
	if err := ChownRecursive(root, 1000, 1001); err != nil {
		t.Fatalf("ChownRecursive() error = %v", err)
	}
	sort.Strings(paths)
	want := []string{
		root,
		filepath.Join(root, "nested"),
		filepath.Join(root, "nested", "file.txt"),
	}
	sort.Strings(want)
	if !reflect.DeepEqual(paths, want) {
		t.Fatalf("paths = %#v, want %#v", paths, want)
	}
}

func TestChownRecursive_ReturnsWalkError(t *testing.T) {
	t.Parallel()

	err := ChownRecursive(filepath.Join(t.TempDir(), "missing"), 1000, 1001)
	if err == nil {
		t.Fatal("ChownRecursive() error = nil, want error")
	}
}

func TestChownRecursive_ReturnsChownError(t *testing.T) {
	restore := replaceChownSeam(t)
	defer restore()

	sentinel := errors.New("chown failed")
	chown = func(path string, uid int, gid int) error {
		return sentinel
	}
	if err := ChownRecursive(t.TempDir(), 1000, 1001); !errors.Is(err, sentinel) {
		t.Fatalf("ChownRecursive() error = %v, want sentinel", err)
	}
}

func replaceChownSeam(t *testing.T) func() {
	t.Helper()

	originalChown := chown
	return func() {
		chown = originalChown
	}
}
