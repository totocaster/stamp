package obsidian

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectNoVault(t *testing.T) {
	dir := t.TempDir()

	result, err := Detect(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.InVault {
		t.Fatalf("expected InVault to be false")
	}
}

func TestDetectWithPlugins(t *testing.T) {
	dir := t.TempDir()
	vault := filepath.Join(dir, "vault")
	obsidianDir := filepath.Join(vault, ".obsidian")
	pluginsDir := filepath.Join(obsidianDir, "plugins", "unique-note-creator")

	if err := os.MkdirAll(pluginsDir, 0o755); err != nil {
		t.Fatalf("setup error: %v", err)
	}

	writeJSON(t, filepath.Join(obsidianDir, "core-plugins.json"), `["daily-notes"]`)
	writeJSON(t, filepath.Join(obsidianDir, "daily-notes.json"), `{"format":"YYYY-MM-DD"}`)

	writeJSON(t, filepath.Join(obsidianDir, "community-plugins.json"), `["unique-note-creator"]`)
	writeJSON(t, filepath.Join(pluginsDir, "data.json"), `{"filenameFormat":"YYYYMMDDHHmm"}`)

	result, err := Detect(filepath.Join(vault, "notes"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !result.InVault {
		t.Fatalf("expected to be inside vault")
	}

	if result.Layouts.Daily != "2006-01-02" {
		t.Fatalf("unexpected daily layout: %q", result.Layouts.Daily)
	}
	if result.Layouts.Default != "200601021504" {
		t.Fatalf("unexpected default layout: %q", result.Layouts.Default)
	}
}

func writeJSON(t *testing.T, path, payload string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir error: %v", err)
	}
	if err := os.WriteFile(path, []byte(payload), 0o644); err != nil {
		t.Fatalf("write error: %v", err)
	}
}
