//go:build darwin
// +build darwin

package clipboard

import (
	"fmt"

	"golang.design/x/clipboard"
)

// Initialize the clipboard package
func init() {
	// Initialize clipboard access
	err := clipboard.Init()
	if err != nil {
		// Log warning but don't fail
		fmt.Printf("Warning: clipboard initialization failed: %v\n", err)
	}
}

// Copy copies text to the clipboard
func Copy(text string) error {
	// Write to clipboard
	clipboard.Write(clipboard.FmtText, []byte(text))
	return nil
}

// Read reads text from the clipboard
func Read() (string, error) {
	// Read from clipboard
	data := clipboard.Read(clipboard.FmtText)
	return string(data), nil
}