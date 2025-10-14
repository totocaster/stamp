//go:build !darwin
// +build !darwin

package clipboard

import "fmt"

// Copy copies text to the clipboard
func Copy(text string) error {
	return fmt.Errorf("clipboard copy is only supported on macOS")
}

// Read reads text from the clipboard
func Read() (string, error) {
	return "", fmt.Errorf("clipboard read is only supported on macOS")
}