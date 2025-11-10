// Package logging provides structured logging functionality for the agent-standards-mcp server.
package logging

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/n-r-w/agent-standards-mcp/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogRotator_Configuration(t *testing.T) {
	tempDir := t.TempDir()
	cfg := &config.Config{
		LogLevel:        "INFO",
		Folder:          tempDir,
		MaxStandards:    100,
		MaxStandardSize: 10240,
	}

	rotator, err := NewLogRotator(cfg)
	require.NoError(t, err)
	require.NotNil(t, rotator)

	// Verify lumberjack configuration
	lj := rotator.lumberjack
	assert.Equal(t, filepath.Join(tempDir, "logs", "agent-standards-mcp.log"), lj.Filename)
	assert.Equal(t, maxLogFileSize, lj.MaxSize)    // 100MB
	assert.Equal(t, maxLogFiles, lj.MaxBackups)    // 7 files
	assert.Equal(t, maxLogAge, lj.MaxAge)          // 7 days
	assert.True(t, lj.Compress)                    // compression enabled
	assert.True(t, lj.LocalTime)                   // local time enabled

	// Test closing rotator
	err = rotator.Close()
	require.NoError(t, err)
}

func TestLogRotator_RotationBehavior(t *testing.T) {
	tempDir := t.TempDir()
	cfg := &config.Config{
		LogLevel:        "INFO",
		Folder:          tempDir,
		MaxStandards:    100,
		MaxStandardSize: 10240,
	}

	rotator, err := NewLogRotator(cfg)
	require.NoError(t, err)
	defer rotator.Close()

	// Create a small log file to test basic functionality
	writer := rotator.Writer()
	require.NotNil(t, writer)

	testMessage := "Test log message for rotation testing\n"
	_, err = writer.Write([]byte(testMessage))
	require.NoError(t, err)

	// Flush the write
	err = rotator.Close()
	require.NoError(t, err)

	// Verify file was created
	logFile := filepath.Join(tempDir, "logs", "agent-standards-mcp.log")
	_, err = os.Stat(logFile)
	require.NoError(t, err)

	// Verify file content
	content, err := os.ReadFile(logFile)
	require.NoError(t, err)
	assert.Contains(t, string(content), testMessage)
}

func TestLogRotator_DirectoryCreation(t *testing.T) {
	tempDir := t.TempDir()
	logDir := filepath.Join(tempDir, "subdir", "logs")
	cfg := &config.Config{
		LogLevel:        "INFO",
		Folder:          filepath.Join(tempDir, "subdir"),
		MaxStandards:    100,
		MaxStandardSize: 10240,
	}

	// Verify directory doesn't exist initially
	_, err := os.Stat(logDir)
	assert.True(t, os.IsNotExist(err))

	rotator, err := NewLogRotator(cfg)
	require.NoError(t, err)
	defer rotator.Close()

	// Verify directory was created
	fileInfo, err := os.Stat(logDir)
	require.NoError(t, err)
	assert.True(t, fileInfo.IsDir())

	// Verify directory permissions
	info, err := os.Stat(logDir)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(dirPermissions), info.Mode().Perm())
}

func TestLogRotator_InvalidConfig(t *testing.T) {
	rotator, err := NewLogRotator(nil)
	assert.Error(t, err)
	assert.Nil(t, rotator)
	assert.Contains(t, err.Error(), "configuration cannot be nil")
}

func TestLogRotator_FilePermissions(t *testing.T) {
	tempDir := t.TempDir()
	cfg := &config.Config{
		LogLevel:        "INFO",
		Folder:          tempDir,
		MaxStandards:    100,
		MaxStandardSize: 10240,
	}

	rotator, err := NewLogRotator(cfg)
	require.NoError(t, err)
	defer rotator.Close()

	// Write to log file to create it
	writer := rotator.Writer()
	_, err = writer.Write([]byte("test content"))
	require.NoError(t, err)

	// Check file permissions
	logFile := filepath.Join(tempDir, "logs", "agent-standards-mcp.log")
	info, err := os.Stat(logFile)
	require.NoError(t, err)

	// File should be readable/writable by owner (default file permissions)
	assert.True(t, info.Mode().Perm()&0600 != 0)
}

func TestLogRotator_ConcurrentAccess(t *testing.T) {
	tempDir := t.TempDir()
	cfg := &config.Config{
		LogLevel:        "INFO",
		Folder:          tempDir,
		MaxStandards:    100,
		MaxStandardSize: 10240,
	}

	rotator, err := NewLogRotator(cfg)
	require.NoError(t, err)
	defer rotator.Close()

	// Test concurrent writes
	writer := rotator.Writer()
	require.NotNil(t, writer)

	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(id int) {
			message := fmt.Sprintf("Concurrent test message %d\n", id)
			_, err := writer.Write([]byte(message))
			assert.NoError(t, err)
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			t.Fatal("Concurrent write test timed out")
		}
	}
}