package counter

import (
	"os"
	"path/filepath"
	"testing"
)

func createTempCounterFile(t *testing.T) string {
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

	// Check that file was created
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

	// Test NextAnalog - should start at 1
	result, err := manager.NextAnalog(date)
	if err != nil {
		t.Fatalf("NextAnalog() error = %v", err)
	}
	if result != "2025-11-12-A1" {
		t.Errorf("NextAnalog() = %v, want 2025-11-12-A1", result)
	}

	// Test NextAnalog again - should increment to 2
	result, err = manager.NextAnalog(date)
	if err != nil {
		t.Fatalf("NextAnalog() error = %v", err)
	}
	if result != "2025-11-12-A2" {
		t.Errorf("NextAnalog() = %v, want 2025-11-12-A2", result)
	}

	// Test CheckAnalog - should show 3 without incrementing
	result, err = manager.CheckAnalog(date)
	if err != nil {
		t.Fatalf("CheckAnalog() error = %v", err)
	}
	if result != "2025-11-12-A3" {
		t.Errorf("CheckAnalog() = %v, want 2025-11-12-A3", result)
	}

	// Verify CheckAnalog didn't increment
	result, err = manager.CheckAnalog(date)
	if err != nil {
		t.Fatalf("CheckAnalog() error = %v", err)
	}
	if result != "2025-11-12-A3" {
		t.Errorf("CheckAnalog() after check = %v, want 2025-11-12-A3", result)
	}

	// Test GetAnalogCounter
	count, err := manager.GetAnalogCounter(date)
	if err != nil {
		t.Fatalf("GetAnalogCounter() error = %v", err)
	}
	if count != 2 {
		t.Errorf("GetAnalogCounter() = %v, want 2", count)
	}

	// Test ResetAnalog
	err = manager.ResetAnalog(date)
	if err != nil {
		t.Fatalf("ResetAnalog() error = %v", err)
	}

	// After reset, should start from 1 again
	result, err = manager.NextAnalog(date)
	if err != nil {
		t.Fatalf("NextAnalog() after reset error = %v", err)
	}
	if result != "2025-11-12-A1" {
		t.Errorf("NextAnalog() after reset = %v, want 2025-11-12-A1", result)
	}
}

//nolint:gocyclo // covers multiple related behaviours in a single flow for readability
func TestManager_ProjectCounters(t *testing.T) {
	counterFile := createTempCounterFile(t)
	manager, err := New(counterFile)
	if err != nil {
		t.Fatalf("Failed to create counter manager: %v", err)
	}

	// Test NextProject - should start at 396 (395 + 1)
	result, err := manager.NextProject("")
	if err != nil {
		t.Fatalf("NextProject() error = %v", err)
	}
	if result != "P0396" {
		t.Errorf("NextProject() = %v, want P0396", result)
	}

	// Test NextProject with title
	result, err = manager.NextProject("Test Project")
	if err != nil {
		t.Fatalf("NextProject() with title error = %v", err)
	}
	if result != "P0397 Test Project" {
		t.Errorf("NextProject() with title = %v, want P0397 Test Project", result)
	}

	// Test CheckProject - should show 398 without incrementing
	result, err = manager.CheckProject()
	if err != nil {
		t.Fatalf("CheckProject() error = %v", err)
	}
	if result != "P0398" {
		t.Errorf("CheckProject() = %v, want P0398", result)
	}

	// Verify CheckProject didn't increment
	result, err = manager.CheckProject()
	if err != nil {
		t.Fatalf("CheckProject() error = %v", err)
	}
	if result != "P0398" {
		t.Errorf("CheckProject() after check = %v, want P0398", result)
	}

	// Test GetProjectCounter
	count, err := manager.GetProjectCounter()
	if err != nil {
		t.Fatalf("GetProjectCounter() error = %v", err)
	}
	if count != 397 {
		t.Errorf("GetProjectCounter() = %v, want 397", count)
	}

	// Test SetProject
	err = manager.SetProject(500)
	if err != nil {
		t.Fatalf("SetProject() error = %v", err)
	}

	result, err = manager.NextProject("")
	if err != nil {
		t.Fatalf("NextProject() after set error = %v", err)
	}
	if result != "P0501" {
		t.Errorf("NextProject() after set = %v, want P0501", result)
	}

	// Test ResetProject
	err = manager.ResetProject()
	if err != nil {
		t.Fatalf("ResetProject() error = %v", err)
	}

	result, err = manager.NextProject("")
	if err != nil {
		t.Fatalf("NextProject() after reset error = %v", err)
	}
	if result != "P0395" {
		t.Errorf("NextProject() after reset = %v, want P0395", result)
	}
}

func TestManager_MultipleDates(t *testing.T) {
	counterFile := createTempCounterFile(t)
	manager, err := New(counterFile)
	if err != nil {
		t.Fatalf("Failed to create counter manager: %v", err)
	}

	// Test that different dates have independent counters
	date1 := "2025-11-12"
	date2 := "2025-11-13"

	// Create some notes for date1
	manager.NextAnalog(date1)
	manager.NextAnalog(date1)

	// Create notes for date2
	result, err := manager.NextAnalog(date2)
	if err != nil {
		t.Fatalf("NextAnalog() for date2 error = %v", err)
	}
	if result != "2025-11-13-A1" {
		t.Errorf("NextAnalog() for different date = %v, want 2025-11-13-A1", result)
	}

	// Verify date1 counter is still correct
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

	// Create first manager and increment counters
	manager1, err := New(counterFile)
	if err != nil {
		t.Fatalf("Failed to create first counter manager: %v", err)
	}

	date := "2025-11-12"
	manager1.NextAnalog(date)
	manager1.NextAnalog(date)
	manager1.NextProject("Test")

	// Create second manager with same file
	manager2, err := New(counterFile)
	if err != nil {
		t.Fatalf("Failed to create second counter manager: %v", err)
	}

	// Verify counters were persisted
	result, err := manager2.NextAnalog(date)
	if err != nil {
		t.Fatalf("NextAnalog() from second manager error = %v", err)
	}
	if result != "2025-11-12-A3" {
		t.Errorf("NextAnalog() from second manager = %v, want 2025-11-12-A3", result)
	}

	result, err = manager2.NextProject("")
	if err != nil {
		t.Fatalf("NextProject() from second manager error = %v", err)
	}
	if result != "P0397" {
		t.Errorf("NextProject() from second manager = %v, want P0397", result)
	}
}

func TestManager_InvalidSetProject(t *testing.T) {
	counterFile := createTempCounterFile(t)
	manager, err := New(counterFile)
	if err != nil {
		t.Fatalf("Failed to create counter manager: %v", err)
	}

	// Test negative value
	err = manager.SetProject(-1)
	if err == nil {
		t.Error("SetProject() with negative value should return error")
	}
}
