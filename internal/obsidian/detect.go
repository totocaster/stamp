package obsidian

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// Layouts captures Go time layouts that can be applied when the CLI
// executes inside an Obsidian vault.
type Layouts struct {
	Default string
	Daily   string
}

// Result describes detected Obsidian metadata for the current working directory.
type Result struct {
	InVault   bool
	VaultPath string
	Layouts   Layouts
}

// Detect attempts to discover whether startPath is within an Obsidian vault and,
// when it is, extracts time formats from known plugins (Daily Notes and Unique Note Creator).
func Detect(startPath string) (*Result, error) {
	absStart, err := filepath.Abs(startPath)
	if err != nil {
		return nil, err
	}

	vaultPath, err := findVault(absStart)
	if err != nil {
		return nil, err
	}

	if vaultPath == "" {
		return &Result{InVault: false}, nil
	}

	res := &Result{
		InVault:   true,
		VaultPath: vaultPath,
	}

	layouts, err := collectLayouts(vaultPath)
	res.Layouts = layouts

	if err != nil {
		return res, err
	}

	return res, nil
}

func findVault(start string) (string, error) {
	current := start
	for {
		obsidianDir := filepath.Join(current, ".obsidian")
		info, err := os.Stat(obsidianDir)
		if err == nil && info.IsDir() {
			return current, nil
		}

		if !errors.Is(err, os.ErrNotExist) && err != nil {
			return "", err
		}

		parent := filepath.Dir(current)
		if parent == current {
			return "", nil
		}
		current = parent
	}
}

func collectLayouts(vaultPath string) (Layouts, error) {
	var layouts Layouts
	var firstErr error

	if format, err := detectDailyNotesFormat(vaultPath); err == nil && format != "" {
		if goLayout, ok := momentToGoLayout(format); ok {
			layouts.Daily = goLayout
		}
	} else if err != nil && firstErr == nil {
		firstErr = err
	}

	if format, err := detectUniqueNoteCreatorFormat(vaultPath); err == nil && format != "" {
		if goLayout, ok := momentToGoLayout(format); ok {
			layouts.Default = goLayout
		}
	} else if err != nil && firstErr == nil {
		firstErr = err
	}

	return layouts, firstErr
}

func detectDailyNotesFormat(vaultPath string) (string, error) {
	if !isCorePluginEnabled(vaultPath, "daily-notes") {
		return "", nil
	}

	if format, err := loadDailyNotesJSON(filepath.Join(vaultPath, ".obsidian", "daily-notes.json")); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return "", err
		}
	} else if format != "" {
		return format, nil
	}

	if format, err := loadDailyNotesFromApp(filepath.Join(vaultPath, ".obsidian", "app.json")); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return "", err
		}
	} else if format != "" {
		return format, nil
	}

	return "", nil
}

func detectUniqueNoteCreatorFormat(vaultPath string) (string, error) {
	if !isCommunityPluginEnabled(vaultPath, "unique-note-creator") && !pluginDirectoryExists(vaultPath, "unique-note-creator") {
		return "", nil
	}

	path := filepath.Join(vaultPath, ".obsidian", "plugins", "unique-note-creator", "data.json")
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", nil
		}
		return "", err
	}

	return findFormatInJSON(data), nil
}

func isCorePluginEnabled(vaultPath, pluginID string) bool {
	ids, err := loadPluginList(filepath.Join(vaultPath, ".obsidian", "core-plugins.json"))
	if err != nil {
		return false
	}
	for _, id := range ids {
		if id == pluginID {
			return true
		}
	}
	return false
}

func isCommunityPluginEnabled(vaultPath, pluginID string) bool {
	ids, err := loadPluginList(filepath.Join(vaultPath, ".obsidian", "community-plugins.json"))
	if err != nil {
		return false
	}
	for _, id := range ids {
		if id == pluginID {
			return true
		}
	}
	return false
}

func pluginDirectoryExists(vaultPath, pluginID string) bool {
	path := filepath.Join(vaultPath, ".obsidian", "plugins", pluginID)
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func loadPluginList(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var ids []string
	if err := json.Unmarshal(data, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

func loadDailyNotesJSON(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	var payload struct {
		Format string `json:"format"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		return "", err
	}
	return payload.Format, nil
}

func loadDailyNotesFromApp(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	var payload struct {
		DailyNotes struct {
			Format string `json:"format"`
		} `json:"dailyNotes"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		return "", err
	}
	return payload.DailyNotes.Format, nil
}

func findFormatInJSON(data []byte) string {
	var parsed interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		return ""
	}
	if format, ok := searchFormat(parsed); ok {
		return format
	}
	return ""
}

func searchFormat(node interface{}) (string, bool) {
	switch v := node.(type) {
	case map[string]interface{}:
		for key, val := range v {
			if isFormatKey(key) {
				if str, ok := val.(string); ok && looksLikeMomentFormat(str) {
					return str, true
				}
			}
		}
		for _, val := range v {
			if str, ok := searchFormat(val); ok {
				return str, true
			}
		}
	case []interface{}:
		for _, item := range v {
			if str, ok := searchFormat(item); ok {
				return str, true
			}
		}
	case string:
		if looksLikeMomentFormat(v) {
			return v, true
		}
	}
	return "", false
}

func isFormatKey(key string) bool {
	lower := strings.ToLower(key)
	switch lower {
	case "format", "dateformat", "filenameformat", "fileformat", "fileNameFormat":
		return true
	}
	if strings.Contains(lower, "format") {
		return true
	}
	return false
}

func looksLikeMomentFormat(value string) bool {
	return strings.ContainsAny(value, "YMDHhms")
}
