package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

// FileExists checks whether a given file or directory exists.
// It returns the absolute path, and an error if any.
func FileExists(path string) (string, error) {
	// Convert to absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return absPath, fmt.Errorf("unable to get absolute path: %w", err)
	}

	// Use os.Stat to get file info
	_, err = os.Stat(absPath)
	if os.IsNotExist(err) {
		return absPath, fmt.Errorf("file does not exist: %s", absPath)
	}

	return absPath, nil
}
