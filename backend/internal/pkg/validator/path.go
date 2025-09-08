package validator

import (
	"errors"
	"path/filepath"
	"strings"
)

var (
	ErrInvalidPath      = errors.New("invalid path")
	ErrPathTraversal    = errors.New("path traversal detected")
)

// ValidatePath checks if a path is safe
func ValidatePath(path string) error {
	// Check for path traversal
	if strings.Contains(path, "..") {
		return ErrPathTraversal
	}

	// Clean the path
	cleaned := filepath.Clean(path)
	if cleaned != path && path != "" {
		return ErrInvalidPath
	}

	return nil
}