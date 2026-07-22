package storage

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad(t *testing.T) {
	path := filepath.Join(t.TempDir(), "potion.config")
	data := `[{"name":"Test","description":"A potion","ingredients":["water"]}]`
	if err := os.WriteFile(path, []byte(data), 0600); err != nil {
		t.Fatal(err)
	}
	book, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if got := book.All()[0].Name; got != "Test" {
		t.Fatalf("recipe name = %q", got)
	}
}
