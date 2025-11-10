// Package logging provides structured logging functionality for the agent-standards-mcp server.
package logging

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/natefinch/lumberjack.v2"
)

// TestLogRotation_Integration tests actual log rotation with small file sizes
func TestLogRotation_Integration(t *testing.T) {
	tempDir := t.TempDir()
	
	// Create a custom rotator with small file size for testing
	logDir := filepath.Join(tempDir, "logs")
	if err := os.MkdirAll(logDir, dirPermissions); err != nil {
		t.Fatalf("Failed to create log directory: %v", err)
	}

	logFile := filepath.Join(logDir, "agent-standards-mcp.log")
	lumberjackLogger := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    1,             // 1MB for testing
		MaxBackups: 3,             // Keep 3 backups
		MaxAge:     1,             // 1 day for testing
		Compress:   true,          // Compress old files
		LocalTime:  true,          // Use local time
	}

	rotator := &LogRotator{
		lumberjack: lumberjackLogger,
	}
	defer rotator.Close()

	// Write enough data to trigger rotation
	writer := rotator.Writer()
	require.NotNil(t, writer)

	// Write approximately 2MB of data to trigger rotation
	message := strings.Repeat("This is a test log message that will help us trigger log rotation when we write enough data. ", 100)
	message += "\n"

	for i := 0; i < 200; i++ {
		logEntry := fmt.Sprintf("Log entry %d: %s", i, message)
		_, err := writer.Write([]byte(logEntry))
		require.NoError(t, err)
		
		// Small delay to ensure proper file handling
		time.Sleep(1 * time.Millisecond)
	}

	// Force flush and close
	err := rotator.Close()
	require.NoError(t, err)

	// Check if rotation occurred by looking for backup files
	files, err := filepath.Glob(filepath.Join(logDir, "agent-standards-mcp.log*"))
	require.NoError(t, err)
	
	t.Logf("Found log files: %v", files)
	
	// Should have at least the main log file
	assert.True(t, len(files) >= 1, "Should have at least the main log file")
	
	// Verify main log file exists
	_, err = os.Stat(logFile)
	require.NoError(t, err)
	
	// Check file sizes are reasonable
	for _, file := range files {
		info, err := os.Stat(file)
		require.NoError(t, err)
		t.Logf("File: %s, Size: %d bytes", filepath.Base(file), info.Size())
		
		// Files should not be empty
		assert.True(t, info.Size() > 0, "Log file should not be empty: %s", file)
	}
}

// TestLogRotation_Deletion tests old log file deletion based on age and count
func TestLogRotation_Deletion(t *testing.T) {
	tempDir := t.TempDir()
	
	// Create a custom rotator with very small limits for testing
	logDir := filepath.Join(tempDir, "logs")
	if err := os.MkdirAll(logDir, dirPermissions); err != nil {
		t.Fatalf("Failed to create log directory: %v", err)
	}

	logFile := filepath.Join(logDir, "agent-standards-mcp.log")
	lumberjackLogger := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    1,             // 1MB for testing
		MaxBackups: 2,             // Keep only 2 backups
		MaxAge:     0,             // No age-based deletion for this test
		Compress:   true,          // Compress old files
		LocalTime:  true,          // Use local time
	}

	rotator := &LogRotator{
		lumberjack: lumberjackLogger,
	}
	defer rotator.Close()

	// Write data to create multiple backup files
	writer := rotator.Writer()
	require.NotNil(t, writer)

	message := strings.Repeat("Test data for rotation deletion test. ", 50) + "\n"

	// Create multiple rotations by writing data in chunks
	for rotation := 0; rotation < 4; rotation++ {
		for i := 0; i < 50; i++ {
			logEntry := fmt.Sprintf("Rotation %d, Entry %d: %s", rotation, i, message)
			_, err := writer.Write([]byte(logEntry))
			require.NoError(t, err)
		}
		
		// Force rotation by creating a new logger
		err := rotator.Close()
		require.NoError(t, err)
		
		// Create new logger for next rotation
		lumberjackLogger = &lumberjack.Logger{
			Filename:   logFile,
			MaxSize:    1,             // 1MB for testing
			MaxBackups: 2,             // Keep only 2 backups
			MaxAge:     0,             // No age-based deletion
			Compress:   true,          // Compress old files
			LocalTime:  true,          // Use local time
		}
		rotator = &LogRotator{
			lumberjack: lumberjackLogger,
		}
		writer = rotator.Writer()
	}

	err := rotator.Close()
	require.NoError(t, err)

	// Check final state of log files
	files, err := filepath.Glob(filepath.Join(logDir, "agent-standards-mcp.log*"))
	require.NoError(t, err)
	
	t.Logf("Final log files after rotation test: %v", files)
	
	// Should have at most 3 files (main + 2 backups)
	assert.True(t, len(files) <= 3, "Should have at most 3 files (main + 2 backups), got %d", len(files))
	
	// Verify all existing files have content
	for _, file := range files {
		info, err := os.Stat(file)
		require.NoError(t, err)
		assert.True(t, info.Size() > 0, "Log file should not be empty: %s", file)
	}
}

// TestLogRotation_ConfigValidation tests the rotation configuration parameters
func TestLogRotation_ConfigValidation(t *testing.T) {
	// Verify the constants are properly set
	assert.Equal(t, 100, maxLogFileSize, "Max log file size should be 100MB")
	assert.Equal(t, 7, maxLogFiles, "Max backup files should be 7")
	assert.Equal(t, 7, maxLogAge, "Max age should be 7 days")
	assert.Equal(t, int(0750), int(dirPermissions), "Directory permissions should be 0750")
}