package sequential

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Spec describes how to detect and format sequential IDs.
type Spec struct {
	Prefix string
	Width  int
	Start  int
}

func (s Spec) normalized() Spec {
	normalized := s
	if normalized.Prefix == "" {
		normalized.Prefix = "P"
	}
	if normalized.Width <= 0 {
		normalized.Width = 4
	}
	if normalized.Start <= 0 {
		normalized.Start = 1
	}
	return normalized
}

// Highest returns the highest numeric component that matches the spec in dir.
func Highest(dir string, spec Spec) (int, error) {
	spec = spec.normalized()

	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0, err
	}

	maxValue := 0
	for _, entry := range entries {
		value, ok := parseName(entry.Name(), spec)
		if !ok {
			continue
		}
		if value > maxValue {
			maxValue = value
		}
	}

	return maxValue, nil
}

// Next returns the next sequential ID formatted according to the spec.
// It also returns the numeric value for callers that need it.
func Next(dir string, spec Spec) (string, int, error) {
	spec = spec.normalized()

	highest, err := Highest(dir, spec)
	if err != nil {
		return "", 0, err
	}

	nextValue := spec.Start
	if highest >= spec.Start {
		nextValue = highest + 1
	}

	return Format(spec, nextValue), nextValue, nil
}

// Format renders a numeric value into the prefixed, zero-padded code.
func Format(spec Spec, value int) string {
	spec = spec.normalized()
	return fmt.Sprintf("%s%0*d", spec.Prefix, spec.Width, value)
}

func parseName(name string, spec Spec) (int, bool) {
	if len(name) < len(spec.Prefix) {
		return 0, false
	}

	prefixChunk := name[:len(spec.Prefix)]
	if !strings.EqualFold(prefixChunk, spec.Prefix) {
		return 0, false
	}

	rest := name[len(spec.Prefix):]
	if len(rest) == 0 {
		return 0, false
	}

	digitsEnd := 0
	for digitsEnd < len(rest) {
		c := rest[digitsEnd]
		if c < '0' || c > '9' {
			break
		}
		digitsEnd++
	}

	if digitsEnd == 0 {
		return 0, false
	}

	value, err := strconv.Atoi(rest[:digitsEnd])
	if err != nil {
		return 0, false
	}

	return value, true
}
