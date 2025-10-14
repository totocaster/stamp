package clipboard

import (
	"fmt"
	"runtime"

	"golang.design/x/clipboard"
)

// Initialize the clipboard package
func init() {
	// Initialize clipboard access
	err := clipboard.Init()
	if err != nil && runtime.GOOS == "darwin" {
		// Log warning but don't fail
		fmt.Printf("Warning: clipboard initialization failed: %v\n", err)
	}
}

// Copy copies text to the clipboard
func Copy(text string) error {
	// Only work on macOS for now
	if runtime.GOOS != "darwin" {
		return fmt.Errorf("clipboard copy is only supported on macOS")
	}

	// Write to clipboard
	clipboard.Write(clipboard.FmtText, []byte(text))

	return nil
}

// Read reads text from the clipboard
func Read() (string, error) {
	// Only work on macOS for now
	if runtime.GOOS != "darwin" {
		return "", fmt.Errorf("clipboard read is only supported on macOS")
	}

	// Read from clipboard
	data := clipboard.Read(clipboard.FmtText)

	return string(data), nil
}