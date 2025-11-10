package standards

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// validateFile validates a single standard file against security and size constraints.
// allowedDir is the base directory that files must be located within.
func validateFile(filePath, allowedDir string) error {
	// Check for path traversal attempts
	if isPathTraversal(filePath, allowedDir) {
		return fmt.Errorf("path traversal detected: %s", filePath)
	}

	// Check if file exists
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file does not exist: %s: %w", filePath, err)
		}
		return fmt.Errorf("failed to stat file %s: %w", filePath, err)
	}

	// Check if path is actually a file (not directory)
	if fileInfo.IsDir() {
		return fmt.Errorf("path is not a file: %s", filePath)
	}

	// Check file size limit
	maxSize, err := getMaxStandardSize()
	if err != nil {
		return fmt.Errorf("failed to get max standard size: %w", err)
	}

	if fileInfo.Size() > maxSize {
		return fmt.Errorf("file size exceeds maximum limit of %d bytes: %d", maxSize, fileInfo.Size())
	}

	return nil
}

// validateStandardFiles validates a list of standard files against count limits.
// allowedDir is the base directory that files must be located within.
func validateStandardFiles(filePaths []string, allowedDir string) error {
	// Check file count limit
	maxStandards, err := getMaxStandards()
	if err != nil {
		return fmt.Errorf("failed to get max standards: %w", err)
	}

	if len(filePaths) > maxStandards {
		return fmt.Errorf("number of files exceeds maximum limit of %d: %d", maxStandards, len(filePaths))
	}

	// Validate each file
	for _, filePath := range filePaths {
		if err := validateFile(filePath, allowedDir); err != nil {
			return fmt.Errorf("validation failed for %s: %w", filePath, err)
		}
	}

	return nil
}

// getMaxStandardSize returns the maximum allowed standard file size in bytes
func getMaxStandardSize() (int64, error) {
	sizeStr := os.Getenv("AGENT_STANDARDS_MCP_MAX_STANDARD_SIZE")
	if sizeStr == "" {
		// Default to 1MB if not set
		return oneMB, nil
	}

	size, err := strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid AGENT_STANDARDS_MCP_MAX_STANDARD_SIZE value: %s", sizeStr)
	}

	return size, nil
}

// getMaxStandards returns the maximum allowed number of standard files
func getMaxStandards() (int, error) {
	countStr := os.Getenv("AGENT_STANDARDS_MCP_MAX_STANDARDS")
	if countStr == "" {
		// Default to 100 if not set
		return defaultMaxStandards, nil
	}

	count, err := strconv.Atoi(countStr)
	if err != nil {
		return 0, fmt.Errorf("invalid AGENT_STANDARDS_MCP_MAX_STANDARDS value: %s", countStr)
	}

	return count, nil
}

// isPathTraversal checks if a path attempts directory traversal
func isPathTraversal(filePath, allowedDir string) bool {
	// Convert to absolute paths for comparison
	absAllowed, err := filepath.Abs(allowedDir)
	if err != nil {
		return true
	}

	absFile, err := filepath.Abs(filePath)
	if err != nil {
		return true
	}

	// Check if the file path is within the allowed directory
	rel, err := filepath.Rel(absAllowed, absFile)
	if err != nil {
		return true
	}

	// Check if the relative path starts with ".." indicating traversal
	return strings.HasPrefix(rel, ".."+string(filepath.Separator)) || rel == ".."
}
