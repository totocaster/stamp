package generator

import (
	"fmt"
	"time"
)

// Generator handles timestamp generation with timezone support
type Generator struct {
	location *time.Location
}

// New creates a new generator with the specified timezone
func New(timezone string) (*Generator, error) {
	loc := time.Local // Default to system timezone

	if timezone != "" {
		var err error
		loc, err = time.LoadLocation(timezone)
		if err != nil {
			return nil, fmt.Errorf("invalid timezone %s: %w", timezone, err)
		}
	}

	return &Generator{
		location: loc,
	}, nil
}

// now returns the current time in the configured timezone
func (g *Generator) now() time.Time {
	return time.Now().In(g.location)
}

// Default generates YYYY-MM-DD-HHMM format
func (g *Generator) Default() string {
	now := g.now()
	return fmt.Sprintf("%04d-%02d-%02d-%02d%02d",
		now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute())
}

// Daily generates YYYY-MM-DD format
func (g *Generator) Daily() string {
	now := g.now()
	return fmt.Sprintf("%04d-%02d-%02d",
		now.Year(), now.Month(), now.Day())
}

// Fleeting generates YYYY-MM-DD-FHHMMSS format
func (g *Generator) Fleeting() string {
	now := g.now()
	return fmt.Sprintf("%04d-%02d-%02d-F%02d%02d%02d",
		now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute(), now.Second())
}

// Voice generates YYYY-MM-DD-VTHHMMSS format
func (g *Generator) Voice() string {
	now := g.now()
	return fmt.Sprintf("%04d-%02d-%02d-VT%02d%02d%02d",
		now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute(), now.Second())
}

// Monthly generates YYYY-MM format
func (g *Generator) Monthly() string {
	now := g.now()
	return fmt.Sprintf("%04d-%02d",
		now.Year(), now.Month())
}

// Yearly generates YYYY format
func (g *Generator) Yearly() string {
	now := g.now()
	return fmt.Sprintf("%04d", now.Year())
}

// GetCurrentDate returns the current date in YYYY-MM-DD format
func (g *Generator) GetCurrentDate() string {
	return g.Daily()
}

// FormatAnalog formats an analog note with given number
func (g *Generator) FormatAnalog(date string, num int) string {
	return fmt.Sprintf("%s-A%d", date, num)
}

// FormatProject formats a project number with optional title
func (g *Generator) FormatProject(num int, title string) string {
	result := fmt.Sprintf("P%04d", num)
	if title != "" {
		result += " " + title
	}
	return result
}