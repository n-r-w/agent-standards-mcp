// Package config provides configuration management for the agent-standards-mcp server.
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// LogLevel represents the allowed log levels for the application.
type LogLevel string

const (
	// LogLevelNone disables logging.
	LogLevelNone LogLevel = "NONE"
	// LogLevelDebug enables debug level logging.
	LogLevelDebug LogLevel = "DEBUG"
	// LogLevelInfo enables info level logging.
	LogLevelInfo LogLevel = "INFO"
	// LogLevelWarn enables warning level logging.
	LogLevelWarn LogLevel = "WARN"
	// LogLevelError enables error level logging.
	LogLevelError LogLevel = "ERROR"
)

const (
	// dirPermissions is the default permissions for directory creation.
	dirPermissions = 0750
)

// validateLogLevel checks if the provided log level is valid.
func validateLogLevel(level string) error {
	normalizedLevel := strings.ToUpper(level)
	switch LogLevel(normalizedLevel) {
	case LogLevelNone, LogLevelDebug, LogLevelInfo, LogLevelWarn, LogLevelError:
		return nil
	default:
		return fmt.Errorf("invalid log level: %s (must be one of: NONE, DEBUG, INFO, WARN, ERROR)", level)
	}
}

// validatePositiveInt checks if the provided integer is positive.
func validatePositiveInt(value int, name string) error {
	if value <= 0 {
		return fmt.Errorf("%s must be positive, got: %d", name, value)
	}
	return nil
}

// expandPath expands ~ to user home directory and resolves the path.
func expandPath(path string) (string, error) {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get user home directory: %w", err)
		}
		return filepath.Join(home, path[2:]), nil
	}
	return path, nil
}

// validateDirectory checks if the directory exists and has appropriate permissions.
func validateDirectory(path string) error {
	// Clean the path to prevent directory traversal
	cleanPath := filepath.Clean(path)

	// Check if directory exists
	fileInfo, err := os.Stat(cleanPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Try to create the directory
			if mkdirErr := os.MkdirAll(cleanPath, dirPermissions); mkdirErr != nil {
				return fmt.Errorf("directory does not exist and failed to create: %s (error: %w)", cleanPath, mkdirErr)
			}
			return nil
		}
		return fmt.Errorf("failed to access directory: %s (error: %w)", cleanPath, err)
	}

	// Check if it's actually a directory
	if !fileInfo.Mode().IsDir() {
		return fmt.Errorf("path is not a directory: %s", cleanPath)
	}

	// Check read permissions
	file, err := os.Open(cleanPath)
	if err != nil {
		return fmt.Errorf("directory lacks read permissions: %s (error: %w)", cleanPath, err)
	}
	if closeErr := file.Close(); closeErr != nil {
		return fmt.Errorf("failed to close directory: %s (error: %w)", cleanPath, closeErr)
	}

	return nil
}

// validateDirectoryPath validates and prepares the directory path.
func validateDirectoryPath(path string) error {
	expandedPath, err := expandPath(path)
	if err != nil {
		return fmt.Errorf("failed to expand path: %w", err)
	}

	return validateDirectory(expandedPath)
}
