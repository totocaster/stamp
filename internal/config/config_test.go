package config

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func setEnv(key, value string) {
	if value == "" {
		os.Unsetenv(key)
	} else {
		os.Setenv(key, value)
	}
}

func withHomeDir(t *testing.T, dir string) {
	t.Helper()

	originalHome := os.Getenv("HOME")
	originalUserProfile := os.Getenv("USERPROFILE")
	originalHomeDrive := os.Getenv("HOMEDRIVE")
	originalHomePath := os.Getenv("HOMEPATH")

	t.Cleanup(func() {
		setEnv("HOME", originalHome)
		setEnv("USERPROFILE", originalUserProfile)
		setEnv("HOMEDRIVE", originalHomeDrive)
		setEnv("HOMEPATH", originalHomePath)
	})

	setEnv("HOME", dir)

	// On Windows, os.UserHomeDir falls back to USERPROFILE or HOMEDRIVE+HOMEPATH.
	if runtime.GOOS == "windows" {
		setEnv("USERPROFILE", dir)

		volume := filepath.VolumeName(dir)
		path := strings.TrimPrefix(dir, volume)
		if volume != "" {
			setEnv("HOMEDRIVE", volume)
			if path == "" {
				path = `\`
			}
			if !strings.HasPrefix(path, `\`) {
				path = `\` + path
			}
			setEnv("HOMEPATH", path)
		} else {
			setEnv("HOMEDRIVE", "")
			setEnv("HOMEPATH", "")
		}
	} else {
		setEnv("USERPROFILE", dir)
		setEnv("HOMEDRIVE", "")
		setEnv("HOMEPATH", "")
	}
}

func setupTempHome(t *testing.T) string {
	t.Helper()
	tmpDir := t.TempDir()
	withHomeDir(t, tmpDir)
	return tmpDir
}

func TestDefault(t *testing.T) {
	cfg := Default()

	if cfg == nil {
		t.Fatal("Default() returned nil")
	}

	if cfg.Timezone != "" {
		t.Errorf("Default timezone = %v, want empty string", cfg.Timezone)
	}

	if cfg.AlwaysExtension != false {
		t.Errorf("Default AlwaysExtension = %v, want false", cfg.AlwaysExtension)
	}

	home, _ := os.UserHomeDir()
	expectedCounterFile := filepath.Join(home, ".stamp", "counters.json")
	if cfg.CounterFile != expectedCounterFile {
		t.Errorf("Default CounterFile = %v, want %v", cfg.CounterFile, expectedCounterFile)
	}
}

func TestLoad_NoConfigFile(t *testing.T) {
	// Create temp dir and set as HOME
	setupTempHome(t)

	// Load should return defaults when no config exists
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v when no config exists", err)
	}

	// Should match defaults
	defaultCfg := Default()
	if cfg.Timezone != defaultCfg.Timezone {
		t.Errorf("Load() timezone = %v, want %v", cfg.Timezone, defaultCfg.Timezone)
	}
	if cfg.AlwaysExtension != defaultCfg.AlwaysExtension {
		t.Errorf("Load() AlwaysExtension = %v, want %v", cfg.AlwaysExtension, defaultCfg.AlwaysExtension)
	}
}

func TestLoad_WithConfigFile(t *testing.T) {
	// Create temp dir and set as HOME
	tmpDir := setupTempHome(t)

	// Create config directory and file
	configDir := filepath.Join(tmpDir, ".stamp")
	os.MkdirAll(configDir, 0o755)

	// Create test config
	testConfig := &Config{
		Timezone:        "Asia/Tokyo",
		AlwaysExtension: true,
		CounterFile:     "~/.stamp/test_counters.json",
	}

	// Write config file
	data, err := yaml.Marshal(testConfig)
	if err != nil {
		t.Fatalf("Failed to marshal test config: %v", err)
	}

	configFile := filepath.Join(configDir, "config.yaml")
	err = os.WriteFile(configFile, data, 0o600)
	if err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	// Load config
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	// Verify loaded values
	if cfg.Timezone != "Asia/Tokyo" {
		t.Errorf("Load() Timezone = %v, want Asia/Tokyo", cfg.Timezone)
	}

	if cfg.AlwaysExtension != true {
		t.Errorf("Load() AlwaysExtension = %v, want true", cfg.AlwaysExtension)
	}

	expectedCounterFile := filepath.Join(tmpDir, ".stamp", "test_counters.json")
	if cfg.CounterFile != expectedCounterFile {
		t.Errorf("Load() CounterFile = %v, want %v", cfg.CounterFile, expectedCounterFile)
	}
}

func TestLoad_InvalidYAML(t *testing.T) {
	// Create temp dir and set as HOME
	tmpDir := setupTempHome(t)

	// Create config directory
	configDir := filepath.Join(tmpDir, ".stamp")
	os.MkdirAll(configDir, 0o755)

	// Write invalid YAML
	configFile := filepath.Join(configDir, "config.yaml")
	err := os.WriteFile(configFile, []byte("invalid: yaml: content:"), 0o600)
	if err != nil {
		t.Fatalf("Failed to write invalid config: %v", err)
	}

	// Load should fail
	cfg, err := Load()
	if err == nil {
		t.Error("Load() should return error for invalid YAML")
	}
	if cfg != nil {
		t.Error("Load() should return nil config for invalid YAML")
	}
}

func TestSave(t *testing.T) {
	// Create temp dir and set as HOME
	tmpDir := setupTempHome(t)

	// Create config
	cfg := &Config{
		Timezone:        "UTC",
		AlwaysExtension: true,
		CounterFile:     "~/.stamp/my_counters.json",
	}

	// Save config
	err := cfg.Save()
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Verify file was created
	configFile := filepath.Join(tmpDir, ".stamp", "config.yaml")
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		t.Error("Config file was not created")
	}

	// Load and verify
	loadedCfg, err := Load()
	if err != nil {
		t.Fatalf("Failed to load saved config: %v", err)
	}

	if loadedCfg.Timezone != cfg.Timezone {
		t.Errorf("Loaded Timezone = %v, want %v", loadedCfg.Timezone, cfg.Timezone)
	}

	if loadedCfg.AlwaysExtension != cfg.AlwaysExtension {
		t.Errorf("Loaded AlwaysExtension = %v, want %v", loadedCfg.AlwaysExtension, cfg.AlwaysExtension)
	}

}

func TestLoad_PartialConfig(t *testing.T) {
	// Create temp dir and set as HOME
	tmpDir := setupTempHome(t)

	// Create config directory
	configDir := filepath.Join(tmpDir, ".stamp")
	os.MkdirAll(configDir, 0o755)

	// Write partial config (only timezone)
	configFile := filepath.Join(configDir, "config.yaml")
	err := os.WriteFile(configFile, []byte("timezone: America/New_York\n"), 0o600)
	if err != nil {
		t.Fatalf("Failed to write partial config: %v", err)
	}

	// Load config
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	// Verify specified value
	if cfg.Timezone != "America/New_York" {
		t.Errorf("Load() Timezone = %v, want America/New_York", cfg.Timezone)
	}

	// Verify defaults for unspecified values
	if cfg.AlwaysExtension != false {
		t.Errorf("Load() AlwaysExtension = %v, want false (default)", cfg.AlwaysExtension)
	}

}
