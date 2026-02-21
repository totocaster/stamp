package counter

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Data stores analog counter information keyed by date.
type Data struct {
	Analog map[string]int `json:"analog"` // Date -> counter mapping
}

// Manager handles analog counter persistence and operations.
type Manager struct {
	mu   sync.Mutex
	file string
	data *Data
}

// New creates a new counter manager that stores data in counterFile.
func New(counterFile string) (*Manager, error) {
	// Expand ~ to home directory
	if strings.HasPrefix(counterFile, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		counterFile = filepath.Join(home, counterFile[2:])
	}

	m := &Manager{
		file: counterFile,
		data: &Data{
			Analog: make(map[string]int),
		},
	}

	// Try to load existing data
	if err := m.load(); err != nil {
		// If file doesn't exist or is corrupted, start fresh
		if !os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Warning: Counter file corrupted, starting fresh: %v\n", err)
		}
		m.data = &Data{
			Analog: make(map[string]int),
		}
		// Save initial data
		if err := m.save(); err != nil {
			return nil, err
		}
	}

	return m, nil
}

// load reads counter data from file
func (m *Manager) load() error {
	// Ensure directory exists
	dir := filepath.Dir(m.file)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	data, err := os.ReadFile(m.file)
	if err != nil {
		return err
	}

	var loaded Data
	if err := json.Unmarshal(data, &loaded); err != nil {
		return err
	}
	if loaded.Analog == nil {
		loaded.Analog = make(map[string]int)
	}

	m.data = &loaded
	return nil
}

// save writes counter data to file
func (m *Manager) save() error {
	// Ensure directory exists
	dir := filepath.Dir(m.file)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(m.data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(m.file, data, 0o600)
}

// NextAnalog returns the next analog number for the given date and increments it
func (m *Manager) NextAnalog(date string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Get current counter for the date
	current := m.data.Analog[date]

	// Increment counter
	m.data.Analog[date] = current + 1

	// Save updated data
	if err := m.save(); err != nil {
		// Rollback on save failure
		m.data.Analog[date] = current
		return "", err
	}

	return fmt.Sprintf("%s-A%d", date, current+1), nil
}

// CheckAnalog returns what the next analog number would be without incrementing
func (m *Manager) CheckAnalog(date string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	current := m.data.Analog[date]
	return fmt.Sprintf("%s-A%d", date, current+1), nil
}

// ResetAnalog resets the counter for a specific date
func (m *Manager) ResetAnalog(date string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.data.Analog, date)
	return m.save()
}

// GetAnalogCounter returns the current counter value for a date
func (m *Manager) GetAnalogCounter(date string) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.data.Analog[date], nil
}

// All project counter methods have been removed; sequential IDs now scan the
// filesystem and no longer rely on this storage layer.
