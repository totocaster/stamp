package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	Timezone        string `yaml:"timezone"`
	AlwaysExtension bool   `yaml:"always_extension"`
	CounterFile     string `yaml:"counter_file"`
	ProjectStart    int    `yaml:"project_start"`
}

// Default returns the default configuration
func Default() *Config {
	home, _ := os.UserHomeDir()

	return &Config{
		Timezone:        "", // Empty means use system timezone
		AlwaysExtension: false,
		CounterFile:     filepath.Join(home, ".stamp", "counters.json"),
		ProjectStart:    395,
	}
}

// Load loads configuration from ~/.stamp/config.yaml
func Load() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configFile := filepath.Join(home, ".stamp", "config.yaml")

	// Check if config file exists
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		// Use defaults if config doesn't exist
		return Default(), nil
	}

	// Read config file
	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	// Start with defaults
	cfg := Default()

	// Parse YAML and overlay on defaults
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	// Expand counter file path if it starts with ~
	if len(cfg.CounterFile) > 0 && cfg.CounterFile[0] == '~' {
		cfg.CounterFile = filepath.Join(home, cfg.CounterFile[2:])
	}

	return cfg, nil
}

// Save saves the configuration to file
func (c *Config) Save() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configDir := filepath.Join(home, ".stamp")
	configFile := filepath.Join(configDir, "config.yaml")

	// Ensure directory exists
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	// Marshal to YAML
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	// Write to file
	return os.WriteFile(configFile, data, 0644)
}