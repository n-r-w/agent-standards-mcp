package standards

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/n-r-w/agent-standards-mcp/internal/domain"
)

// FileStandardLoader implements the StandardLoader interface for loading standards from the file system.
type FileStandardLoader struct {
	standardsDir string
}

// NewFileStandardLoader creates a new FileStandardLoader instance.
func NewFileStandardLoader() *FileStandardLoader {
	standardsDir := os.Getenv("AGENT_STANDARDS_MCP_FOLDER")
	if standardsDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			homeDir = "."
		}
		standardsDir = filepath.Join(homeDir, "agent-standards", "standards") // Default directory
	}

	return &FileStandardLoader{
		standardsDir: standardsDir,
	}
}

// ListStandards returns a list of available standard information (name and description).
func (l *FileStandardLoader) ListStandards(_ context.Context) ([]domain.StandardInfo, error) {
	// Find all standard files
	filePaths, err := l.findStandardFiles()
	if err != nil {
		return nil, fmt.Errorf("failed to find standard files: %w", err)
	}

	// Validate all files first
	if err := validateStandardFiles(filePaths, l.standardsDir); err != nil {
		return nil, fmt.Errorf("failed to validate standard files: %w", err)
	}

	// Pre-allocate slice with known capacity
	standardInfos := make([]domain.StandardInfo, 0, len(filePaths))

	for _, filePath := range filePaths {
		// Sanitize file path to prevent path traversal attacks
		cleanPath := filepath.Clean(filePath)

		// Read file content (files already validated by ValidateStandardFiles above)
		content, err := os.ReadFile(cleanPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read file %s: %w", cleanPath, err)
		}

		// Parse frontmatter
		description, _, err := parseFrontmatter(string(content))
		if err != nil {
			return nil, fmt.Errorf("failed to parse frontmatter for %s: %w", filePath, err)
		}

		// Extract standard name from file path
		standardName := extractStandardName(filePath)

		standardInfo := domain.StandardInfo{
			Name:        standardName,
			Description: description,
		}

		standardInfos = append(standardInfos, standardInfo)
	}

	return standardInfos, nil
}

// GetStandards returns the full content of specific standards by their names.
func (l *FileStandardLoader) GetStandards(_ context.Context, standardNames []string) ([]domain.Standard, error) {
	// Pre-allocate slice with known capacity
	standards := make([]domain.Standard, 0, len(standardNames))

	for _, standardName := range standardNames {
		// Construct file path
		filePath := filepath.Join(l.standardsDir, standardName+".md")

		// Validate the file
		if err := validateFile(filePath, l.standardsDir); err != nil {
			// If file doesn't exist, just skip it (don't return error)
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return nil, fmt.Errorf("failed to validate standard file %s: %w", standardName, err)
		}

		// Read file content
		cleanPath := filepath.Clean(filePath)
		content, err := os.ReadFile(cleanPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read standard file %s: %w", standardName, err)
		}

		// Parse frontmatter
		description, standardContent, err := parseFrontmatter(string(content))
		if err != nil {
			return nil, fmt.Errorf("failed to parse frontmatter for standard %s: %w", standardName, err)
		}

		standard := domain.Standard{
			Name:        standardName,
			Description: description,
			Content:     standardContent,
		}

		standards = append(standards, standard)
	}

	return standards, nil
}

// extractStandardName extracts the standard name from a file path by removing the directory and extension.
func extractStandardName(filePath string) string {
	// Get the base filename
	base := filepath.Base(filePath)

	// Remove the extension
	ext := filepath.Ext(base)
	if ext != "" {
		return base[:len(base)-len(ext)]
	}

	return base
}

// findStandardFiles finds all markdown files in the standards directory, excluding hidden files.
func (l *FileStandardLoader) findStandardFiles() ([]string, error) {
	entries, err := os.ReadDir(l.standardsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil // Empty directory is fine
		}
		return nil, fmt.Errorf("failed to read standards directory %s: %w", l.standardsDir, err)
	}

	// Pre-allocate slice with estimated capacity
	files := make([]string, 0, len(entries))

	for _, entry := range entries {
		// Skip hidden files and directories
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		// Only include regular files
		if !entry.Type().IsRegular() {
			continue
		}

		// Only include markdown files
		if filepath.Ext(entry.Name()) != ".md" {
			continue
		}

		filePath := filepath.Join(l.standardsDir, entry.Name())
		files = append(files, filePath)
	}

	return files, nil
}
