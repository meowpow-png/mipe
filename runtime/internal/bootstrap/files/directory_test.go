package files

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCreateDirectory_CreatesMissingParents(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "nested", "directory")

	if err := CreateDirectory(path); err != nil {
		t.Fatalf("CreateDirectory() error = %v", err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat path: %v", err)
	}
	if !info.IsDir() {
		t.Fatalf("%s is not a directory", path)
	}
}
