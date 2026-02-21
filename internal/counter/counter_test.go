package counter

import (
	"os"
	"path/filepath"
	"testing"
)

func createTempCounterFile(t *testing.T) string {
	t.Helper()
	tmpDir := t.TempDir()
	return filepath.Join(tmpDir, "test_counters.json")
}

func TestNew(t *testing.T) {
	counterFile := createTempCounterFile(t)

	manager, err := New(counterFile)
	if err != nil {
		t.Fatalf("Failed to create counter manager: %v", err)
	}

	if manager == nil {
		t.Fatal("New() returned nil manager")
	}

	if _, err := os.Stat(counterFile); os.IsNotExist(err) {
		t.Error("Counter file was not created")
	}
}

func TestManager_AnalogCounters(t *testing.T) {
	counterFile := createTempCounterFile(t)
	manager, err := New(counterFile)
	if err != nil {
		t.Fatalf("Failed to create counter manager: %v", err)
	}

	date := "2025-11-12"

	result, err := manager.NextAnalog(date)
	if err != nil {
		t.Fatalf("NextAnalog() error = %v", err)
	}
	if result != "2025-11-12-A1" {
		t.Errorf("NextAnalog() = %v, want 2025-11-12-A1", result)
	}

	result, err = manager.NextAnalog(date)
	if err != nil {
		t.Fatalf("NextAnalog() error = %v", err)
	}
	if result != "2025-11-12-A2" {
		t.Errorf("NextAnalog() = %v, want 2025-11-12-A2", result)
	}

	result, err = manager.CheckAnalog(date)
	if err != nil {
		t.Fatalf("CheckAnalog() error = %v", err)
	}
	if result != "2025-11-12-A3" {
		t.Errorf("CheckAnalog() = %v, want 2025-11-12-A3", result)
	}

	count, err := manager.GetAnalogCounter(date)
	if err != nil {
		t.Fatalf("GetAnalogCounter() error = %v", err)
	}
	if count != 2 {
		t.Errorf("GetAnalogCounter() = %v, want 2", count)
	}

	if err := manager.ResetAnalog(date); err != nil {
		t.Fatalf("ResetAnalog() error = %v", err)
	}

	result, err = manager.NextAnalog(date)
	if err != nil {
		t.Fatalf("NextAnalog() after reset error = %v", err)
	}
	if result != "2025-11-12-A1" {
		t.Errorf("NextAnalog() after reset = %v, want 2025-11-12-A1", result)
	}
}

func TestManager_MultipleDates(t *testing.T) {
	counterFile := createTempCounterFile(t)
	manager, err := New(counterFile)
	if err != nil {
		t.Fatalf("Failed to create counter manager: %v", err)
	}

	date1 := "2025-11-12"
	date2 := "2025-11-13"

	manager.NextAnalog(date1)
	manager.NextAnalog(date1)

	result, err := manager.NextAnalog(date2)
	if err != nil {
		t.Fatalf("NextAnalog() for date2 error = %v", err)
	}
	if result != "2025-11-13-A1" {
		t.Errorf("NextAnalog() for different date = %v, want 2025-11-13-A1", result)
	}

	result, err = manager.NextAnalog(date1)
	if err != nil {
		t.Fatalf("NextAnalog() for date1 error = %v", err)
	}
	if result != "2025-11-12-A3" {
		t.Errorf("NextAnalog() for date1 after date2 = %v, want 2025-11-12-A3", result)
	}
}

func TestManager_Persistence(t *testing.T) {
	counterFile := createTempCounterFile(t)

	manager1, err := New(counterFile)
	if err != nil {
		t.Fatalf("Failed to create first counter manager: %v", err)
	}

	date := "2025-11-12"
	manager1.NextAnalog(date)
	manager1.NextAnalog(date)

	manager2, err := New(counterFile)
	if err != nil {
		t.Fatalf("Failed to create second counter manager: %v", err)
	}

	result, err := manager2.NextAnalog(date)
	if err != nil {
		t.Fatalf("NextAnalog() from second manager error = %v", err)
	}
	if result != "2025-11-12-A3" {
		t.Errorf("NextAnalog() from second manager = %v, want 2025-11-12-A3", result)
	}
}
