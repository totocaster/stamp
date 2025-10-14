package generator

import (
	"regexp"
	"strings"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name     string
		timezone string
		wantErr  bool
	}{
		{
			name:     "default timezone",
			timezone: "",
			wantErr:  false,
		},
		{
			name:     "valid timezone Tokyo",
			timezone: "Asia/Tokyo",
			wantErr:  false,
		},
		{
			name:     "valid timezone UTC",
			timezone: "UTC",
			wantErr:  false,
		},
		{
			name:     "invalid timezone",
			timezone: "Invalid/Zone",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen, err := New(tt.timezone)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && gen == nil {
				t.Errorf("New() returned nil generator")
			}
		})
	}
}

func TestGenerator_Default(t *testing.T) {
	gen, err := New("")
	if err != nil {
		t.Fatalf("Failed to create generator: %v", err)
	}

	result := gen.Default()

	// Check format: YYYY-MM-DD-HHMM
	pattern := `^\d{4}-\d{2}-\d{2}-\d{4}$`
	match, _ := regexp.MatchString(pattern, result)
	if !match {
		t.Errorf("Default() = %v, want format YYYY-MM-DD-HHMM", result)
	}
}

func TestGenerator_Daily(t *testing.T) {
	gen, err := New("")
	if err != nil {
		t.Fatalf("Failed to create generator: %v", err)
	}

	result := gen.Daily()

	// Check format: YYYY-MM-DD
	pattern := `^\d{4}-\d{2}-\d{2}$`
	match, _ := regexp.MatchString(pattern, result)
	if !match {
		t.Errorf("Daily() = %v, want format YYYY-MM-DD", result)
	}
}

func TestGenerator_Fleeting(t *testing.T) {
	gen, err := New("")
	if err != nil {
		t.Fatalf("Failed to create generator: %v", err)
	}

	result := gen.Fleeting()

	// Check format: YYYY-MM-DD-FHHMMSS
	pattern := `^\d{4}-\d{2}-\d{2}-F\d{6}$`
	match, _ := regexp.MatchString(pattern, result)
	if !match {
		t.Errorf("Fleeting() = %v, want format YYYY-MM-DD-FHHMMSS", result)
	}
}

func TestGenerator_Voice(t *testing.T) {
	gen, err := New("")
	if err != nil {
		t.Fatalf("Failed to create generator: %v", err)
	}

	result := gen.Voice()

	// Check format: YYYY-MM-DD-VTHHMMSS
	pattern := `^\d{4}-\d{2}-\d{2}-VT\d{6}$`
	match, _ := regexp.MatchString(pattern, result)
	if !match {
		t.Errorf("Voice() = %v, want format YYYY-MM-DD-VTHHMMSS", result)
	}
}

func TestGenerator_Monthly(t *testing.T) {
	gen, err := New("")
	if err != nil {
		t.Fatalf("Failed to create generator: %v", err)
	}

	result := gen.Monthly()

	// Check format: YYYY-MM
	pattern := `^\d{4}-\d{2}$`
	match, _ := regexp.MatchString(pattern, result)
	if !match {
		t.Errorf("Monthly() = %v, want format YYYY-MM", result)
	}
}

func TestGenerator_Yearly(t *testing.T) {
	gen, err := New("")
	if err != nil {
		t.Fatalf("Failed to create generator: %v", err)
	}

	result := gen.Yearly()

	// Check format: YYYY
	pattern := `^\d{4}$`
	match, _ := regexp.MatchString(pattern, result)
	if !match {
		t.Errorf("Yearly() = %v, want format YYYY", result)
	}
}

func TestGenerator_FormatAnalog(t *testing.T) {
	gen, err := New("")
	if err != nil {
		t.Fatalf("Failed to create generator: %v", err)
	}

	tests := []struct {
		name string
		date string
		num  int
		want string
	}{
		{"first note", "2025-11-12", 1, "2025-11-12-A1"},
		{"second note", "2025-11-12", 2, "2025-11-12-A2"},
		{"tenth note", "2025-11-12", 10, "2025-11-12-A10"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := gen.FormatAnalog(tt.date, tt.num); got != tt.want {
				t.Errorf("FormatAnalog() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerator_FormatProject(t *testing.T) {
	gen, err := New("")
	if err != nil {
		t.Fatalf("Failed to create generator: %v", err)
	}

	tests := []struct {
		name  string
		num   int
		title string
		want  string
	}{
		{"without title", 395, "", "P0395"},
		{"with title", 396, "New Project", "P0396 New Project"},
		{"large number", 1234, "", "P1234"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := gen.FormatProject(tt.num, tt.title); got != tt.want {
				t.Errorf("FormatProject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerator_TimezoneConsistency(t *testing.T) {
	// Test that timezone is consistently applied
	gen, err := New("UTC")
	if err != nil {
		t.Fatalf("Failed to create generator with UTC timezone: %v", err)
	}

	// All outputs should be in UTC
	daily := gen.Daily()
	now := time.Now().UTC()
	expectedPrefix := now.Format("2006-01-02")

	if !strings.HasPrefix(daily, expectedPrefix[:10]) {
		t.Errorf("Daily() with UTC timezone = %v, should start with today's date in UTC", daily)
	}
}