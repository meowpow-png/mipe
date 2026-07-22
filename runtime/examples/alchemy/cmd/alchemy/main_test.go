package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestRunList(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "potion.config"), []byte(`[{"name":"Test"}]`), 0600); err != nil {
		t.Fatal(err)
	}
	old, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(old) }()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	var out bytes.Buffer
	if err := run([]string{"list"}, &out); err != nil {
		t.Fatal(err)
	}
	if out.String() != "Test\n" {
		t.Fatalf("output = %q", out.String())
	}
}
