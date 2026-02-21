package sequential

import (
	"os"
	"path/filepath"
	"testing"
)

func TestHighestAndNext(t *testing.T) {
	dir := t.TempDir()

	files := []string{
		"P0005 Alpha.md",
		"p0100 Beta.txt",
		"notes.txt",
		"X0010.txt",
		"P0099",
	}

	for _, name := range files {
		path := filepath.Join(dir, name)
		if err := os.WriteFile(path, []byte("test"), 0o644); err != nil {
			t.Fatalf("failed to write %s: %v", name, err)
		}
	}

	highest, err := Highest(dir, Spec{Prefix: "P", Width: 4})
	if err != nil {
		t.Fatalf("Highest() error = %v", err)
	}
	if highest != 100 {
		t.Fatalf("Highest() = %d, want 100", highest)
	}

	code, value, err := Next(dir, Spec{Prefix: "P", Width: 4})
	if err != nil {
		t.Fatalf("Next() error = %v", err)
	}
	if value != 101 {
		t.Fatalf("Next() value = %d, want 101", value)
	}
	if code != "P0101" {
		t.Fatalf("Next() code = %s, want P0101", code)
	}
}

func TestNextDefaults(t *testing.T) {
	dir := t.TempDir()

	code, value, err := Next(dir, Spec{})
	if err != nil {
		t.Fatalf("Next() error = %v", err)
	}
	if value != 1 {
		t.Fatalf("Next() value = %d, want 1", value)
	}
	if code != "P0001" {
		t.Fatalf("Next() code = %s, want P0001", code)
	}
}

func TestFormatWidth(t *testing.T) {
	code := Format(Spec{Prefix: "X", Width: 2}, 7)
	if code != "X07" {
		t.Fatalf("Format() = %s, want X07", code)
	}
}

func TestHighestMultiCharPrefix(t *testing.T) {
	dir := t.TempDir()

	entries := []string{"Proj0001 One", "proj0003 Two", "PRJ0002"}
	for _, name := range entries {
		if err := os.Mkdir(filepath.Join(dir, name), 0o755); err != nil {
			t.Fatalf("failed to mkdir %s: %v", name, err)
		}
	}

	highest, err := Highest(dir, Spec{Prefix: "Proj", Width: 4})
	if err != nil {
		t.Fatalf("Highest() error = %v", err)
	}
	if highest != 3 {
		t.Fatalf("Highest() = %d, want 3", highest)
	}
}
