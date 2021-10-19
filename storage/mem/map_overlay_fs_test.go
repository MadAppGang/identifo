package mem_test

import (
	"io/fs"
	"path"
	"testing"

	"github.com/madappgang/identifo/storage/mem"
	"github.com/spf13/afero"
)

func TestNewMapOverlayFS(t *testing.T) {
	base := afero.NewMemMapFs()
	afero.WriteFile(base, "./test1.txt", []byte("file1"), 644)
	afero.WriteFile(base, "./path/test2.txt", []byte("file2"), 644)

	files := map[string][]byte{
		"./path/test3.txt": []byte("file3"),
		"test4.txt":        []byte("file4"),
	}

	m := mem.NewMapOverlayFS(afero.NewIOFS(base), files)

	d1, _ := fs.ReadFile(m, "test1.txt")
	if string(d1) != "file1" {
		t.Fatalf("Error getting data, got %s, expected: %s", d1, "file1")
	}

	d2, _ := fs.ReadFile(m, path.Clean("./path/test2.txt"))
	if string(d2) != "file2" {
		t.Fatalf("Error getting data, got %s, expected: %s", d2, "file2")
	}

	d3, _ := fs.ReadFile(m, path.Clean("path/test3.txt"))
	if string(d3) != "file3" {
		t.Fatalf("Error getting data, got %s, expected: %s", d3, "file3")
	}

	d4, _ := fs.ReadFile(m, path.Clean("test4.txt"))
	if string(d4) != "file4" {
		t.Fatalf("Error getting data, got %s, expected: %s", d4, "file4")
	}
}
