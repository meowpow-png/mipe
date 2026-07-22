package files

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestCopyContents_CopiesDirectoryContentsRecursively(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	source := filepath.Join(root, "source")
	destination := filepath.Join(root, "destination")
	nested := filepath.Join(source, "nested")

	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatalf("mkdir nested: %v", err)
	}
	if err := os.MkdirAll(destination, 0o755); err != nil {
		t.Fatalf("mkdir destination: %v", err)
	}
	if err := os.WriteFile(filepath.Join(source, "root.txt"), []byte("root"), 0o640); err != nil {
		t.Fatalf("write root file: %v", err)
	}
	if err := os.WriteFile(filepath.Join(nested, "child.txt"), []byte("child"), 0o600); err != nil {
		t.Fatalf("write child file: %v", err)
	}
	if err := CopyContents(source, destination); err != nil {
		t.Fatalf("CopyContents() error = %v", err)
	}
	assertFile(t, filepath.Join(destination, "root.txt"), "root")
	assertFile(t, filepath.Join(destination, "nested", "child.txt"), "child")
	if _, err := os.Stat(filepath.Join(destination, "source")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("destination/source stat error = %v, want not exist", err)
	}
}

func TestCopyContents_ReturnsErrorForMissingSource(t *testing.T) {
	t.Parallel()

	err := CopyContents(filepath.Join(t.TempDir(), "missing"), t.TempDir())
	if err == nil {
		t.Fatal("CopyContents() error = nil, want error")
	}
}

func TestCopyChildren_ReturnsErrorForChildCopyFailure(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	source := filepath.Join(root, "source")
	destination := filepath.Join(root, "destination")
	if err := os.MkdirAll(source, 0o755); err != nil {
		t.Fatalf("mkdir source: %v", err)
	}
	if err := os.WriteFile(filepath.Join(source, "file.txt"), []byte("content"), 0o600); err != nil {
		t.Fatalf("write source file: %v", err)
	}
	if err := os.WriteFile(destination, []byte("not a directory"), 0o600); err != nil {
		t.Fatalf("write destination file: %v", err)
	}
	if err := copyChildren(source, destination); err == nil {
		t.Fatal("copyChildren() error = nil, want error")
	}
}

func TestCopyPath_ReturnsErrorForMissingSource(t *testing.T) {
	t.Parallel()

	err := copyPath(filepath.Join(t.TempDir(), "missing"), filepath.Join(t.TempDir(), "destination"))
	if err == nil {
		t.Fatal("copyPath() error = nil, want error")
	}
}

func TestCopyDirectory_ReturnsErrorForInvalidDestination(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	source := filepath.Join(root, "source")
	destinationParent := filepath.Join(root, "destination-parent")
	destination := filepath.Join(destinationParent, "child")
	if err := os.MkdirAll(source, 0o755); err != nil {
		t.Fatalf("mkdir source: %v", err)
	}
	if err := os.WriteFile(destinationParent, []byte("not a directory"), 0o600); err != nil {
		t.Fatalf("write destination parent: %v", err)
	}
	if err := copyDirectory(source, destination, 0o755); err == nil {
		t.Fatal("copyDirectory() error = nil, want error")
	}
}

func TestCopyDirectory_ReturnsErrorForUnreadableSource(t *testing.T) {
	t.Parallel()

	err := copyDirectory(filepath.Join(t.TempDir(), "missing"), filepath.Join(t.TempDir(), "destination"), 0o755)
	if err == nil {
		t.Fatal("copyDirectory() error = nil, want error")
	}
}

func TestCopyFile_PreservesFileMode(t *testing.T) {
	t.Parallel()

	if runtime.GOOS == "windows" {
		t.Skip("file mode preservation is platform-specific")
	}
	root := t.TempDir()
	source := filepath.Join(root, "source.txt")
	destination := filepath.Join(root, "destination.txt")

	if err := os.WriteFile(source, []byte("content"), 0o640); err != nil {
		t.Fatalf("write source: %v", err)
	}
	if err := copyFile(source, destination, 0o640); err != nil {
		t.Fatalf("copyFile() error = %v", err)
	}
	info, err := os.Stat(destination)
	if err != nil {
		t.Fatalf("stat destination: %v", err)
	}
	if got, want := info.Mode().Perm(), os.FileMode(0o640); got != want {
		t.Fatalf("mode = %v, want %v", got, want)
	}
	assertFile(t, destination, "content")
}

func TestCopyFile_ReturnsErrorForMissingSource(t *testing.T) {
	t.Parallel()

	err := copyFile(filepath.Join(t.TempDir(), "missing"), filepath.Join(t.TempDir(), "destination"), 0o600)
	if err == nil {
		t.Fatal("copyFile() error = nil, want error")
	}
}

func TestCopyFile_ReturnsErrorForInvalidDestination(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	source := filepath.Join(root, "source.txt")
	destination := filepath.Join(root, "missing", "destination.txt")
	if err := os.WriteFile(source, []byte("content"), 0o600); err != nil {
		t.Fatalf("write source: %v", err)
	}
	if err := copyFile(source, destination, 0o600); err == nil {
		t.Fatal("copyFile() error = nil, want error")
	}
}

func assertFile(t *testing.T, path string, want string) {
	t.Helper()

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	if got := string(content); got != want {
		t.Fatalf("%s = %q, want %q", path, got, want)
	}
}
