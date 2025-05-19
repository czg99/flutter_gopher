package bridgegen

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ParsePkgPath returns the package path of the given file path.
// It finds the nearest go.mod file and extracts the module name,
// then calculates the relative package path from the module root.
// Returns:
//   - module: the Go module name from go.mod
//   - pkgPath: the full package import path
//   - err: any error encountered during the process
func ParsePkgPath(path string) (module, pkgPath string, err error) {
	// Convert to absolute path
	path, err = filepath.Abs(path)
	if err != nil {
		err = fmt.Errorf("failed to get absolute path: %w", err)
		return
	}

	// Check if path exists
	fileInfo, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("failed to stat path: %w", err)
		return
	}

	// Determine directory to search from
	var dirToSearch string
	if fileInfo.IsDir() {
		dirToSearch = path
	} else {
		dirToSearch = filepath.Dir(path)
	}

	// Search for go.mod file in current and parent directories
	currentDir := dirToSearch
	for {
		goModPath := filepath.Join(currentDir, "go.mod")
		if _, statErr := os.Stat(goModPath); statErr == nil {
			// Found go.mod, read its content
			data, readErr := os.ReadFile(goModPath)
			if readErr != nil {
				err = fmt.Errorf("failed to read go.mod file: %w", readErr)
				return
			}

			// Extract module name
			lines := strings.Split(string(data), "\n")
			for _, line := range lines {
				trimmedLine := strings.TrimSpace(line)
				if strings.HasPrefix(trimmedLine, "module ") {
					module = strings.TrimSpace(strings.TrimPrefix(trimmedLine, "module "))

					// Calculate relative path from module root to target directory
					relPath, relErr := filepath.Rel(currentDir, dirToSearch)
					if relErr == nil && relPath != "." {
						pkgPath = filepath.ToSlash(filepath.Join(module, relPath))
						return
					}

					// If target is the module root itself
					pkgPath = module
					return
				}
			}

			err = fmt.Errorf("module declaration not found in go.mod")
			return
		}

		// Move to parent directory
		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			// Reached filesystem root without finding go.mod
			err = fmt.Errorf("no go.mod file found in directory hierarchy")
			break
		}
		currentDir = parentDir
	}

	return
}
