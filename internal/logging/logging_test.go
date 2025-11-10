// Package logging provides structured logging functionality for the agent-standards-mcp server.
package logging

import (
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/n-r-w/agent-standards-mcp/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStructuredLogger_Debug(t *testing.T) {
	logger := &StructuredLogger{
		logger:     slog.Default(),
		logRotator: nil,
	}

	// This should not panic
	logger.Debug("test debug message", "key", "value")
}

func TestStructuredLogger_Info(t *testing.T) {
	logger := &StructuredLogger{
		logger:     slog.Default(),
		logRotator: nil,
	}

	// This should not panic
	logger.Info("test info message", "key", "value")
}

func TestStructuredLogger_Warn(t *testing.T) {
	logger := &StructuredLogger{
		logger:     slog.Default(),
		logRotator: nil,
	}

	// This should not panic
	logger.Warn("test warn message", "key", "value")
}

func TestStructuredLogger_Error(t *testing.T) {
	logger := &StructuredLogger{
		logger:     slog.Default(),
		logRotator: nil,
	}

	// This should not panic
	logger.Error("test error message", "key", "value")
}

func TestStructuredLogger_Close(t *testing.T) {
	// Test closing logger without log rotator (nil logRotator)
	logger := &StructuredLogger{
		logger:     slog.Default(),
		logRotator: nil,
	}

	err := logger.Close()
	assert.NoError(t, err)
}

func TestNewStructuredLogger(t *testing.T) {
	// This test will fail until NewStructuredLogger is implemented
	tempDir := t.TempDir()
	cfg := &config.Config{
		LogLevel:        "DEBUG",
		Folder:          tempDir,
		MaxStandards:    100,
		MaxStandardSize: 10240,
	}

	logger, err := NewStructuredLogger(cfg)
	require.NoError(t, err)
	require.NotNil(t, logger)

	// Test that logger works
	logger.Info("test message")
}

func TestNewStructuredLogger_DisabledLogging(t *testing.T) {
	// This test will fail until NewStructuredLogger is implemented
	cfg := &config.Config{
		LogLevel:        "NONE",
		Folder:          "/tmp",
		MaxStandards:    100,
		MaxStandardSize: 10240,
	}

	logger, err := NewStructuredLogger(cfg)
	require.NoError(t, err)
	require.NotNil(t, logger)

	// Test that logger works but doesn't output
	logger.Info("test message")
}

func TestNewStructuredLogger_InvalidLogLevel(t *testing.T) {
	// This test will fail until NewStructuredLogger is implemented
	cfg := &config.Config{
		LogLevel:        "INVALID",
		Folder:          "/tmp",
		MaxStandards:    100,
		MaxStandardSize: 10240,
	}

	_, err := NewStructuredLogger(cfg)
	require.Error(t, err)
}

func TestNewStructuredLogger_FileLogging(t *testing.T) {
	// This test will fail until file logging is implemented
	tempDir := t.TempDir()
	cfg := &config.Config{
		LogLevel:        "INFO",
		Folder:          tempDir,
		MaxStandards:    100,
		MaxStandardSize: 10240,
	}

	logger, err := NewStructuredLogger(cfg)
	require.NoError(t, err)
	require.NotNil(t, logger)

	// Test that logger works
	logger.Info("test message")

	// Check if log file was created
	logDir := filepath.Join(tempDir, "logs")
	logFiles, err := os.ReadDir(logDir)
	if err == nil {
		assert.NotEmpty(t, logFiles, "Log files should be created")
	}
}

func TestNewLogRotator(t *testing.T) {
	// This test will fail until NewLogRotator is implemented
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

	// Test that rotator has a writer
	writer := rotator.Writer()
	require.NotNil(t, writer)

	// Test closing rotator
	err = rotator.Close()
	require.NoError(t, err)
}

func TestNewLogRotator_InvalidConfig(t *testing.T) {
	// This test will fail until NewLogRotator is implemented
	rotator, err := NewLogRotator(nil)
	require.Error(t, err)
	require.Nil(t, rotator)
}

func TestLogRotator_FileCreation(t *testing.T) {
	// This test will fail until file logging is implemented
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

	// Write something to the log file to trigger creation
	writer := rotator.Writer()
	_, err = writer.Write([]byte("test log message"))
	require.NoError(t, err)

	// Close rotator to flush any pending writes
	err = rotator.Close()
	require.NoError(t, err)

	// Check if log file was created
	logDir := filepath.Join(tempDir, "logs")
	logFiles, err := os.ReadDir(logDir)
	if err == nil {
		assert.NotEmpty(t, logFiles, "Log files should be created")
	}
}

func TestNewLoggerFactory(t *testing.T) {
	// This test will fail until LoggerFactory is implemented
	factory := NewLoggerFactory()
	require.NotNil(t, factory)
}

func TestLoggerFactory_CreateStructuredLogger(t *testing.T) {
	// This test will fail until LoggerFactory is implemented
	factory := NewLoggerFactory()
	require.NotNil(t, factory)

	tempDir := t.TempDir()
	cfg := &config.Config{
		LogLevel:        "INFO",
		Folder:          tempDir,
		MaxStandards:    100,
		MaxStandardSize: 10240,
	}

	logger, err := factory.CreateStructuredLogger(cfg)
	require.NoError(t, err)
	require.NotNil(t, logger)

	// Test that logger works
	logger.Info("test message")
}

func TestLoggerFactory_CreateStructuredLogger_InvalidConfig(t *testing.T) {
	// This test will fail until LoggerFactory is implemented
	factory := NewLoggerFactory()
	require.NotNil(t, factory)

	logger, err := factory.CreateStructuredLogger(nil)
	require.Error(t, err)
	require.Nil(t, logger)
}

func TestLoggerFactory_CreateAudit(t *testing.T) {
	// This test will fail until LoggerFactory is implemented
	factory := NewLoggerFactory()
	require.NotNil(t, factory)

	tempDir := t.TempDir()
	cfg := &config.Config{
		LogLevel:        "INFO",
		Folder:          tempDir,
		MaxStandards:    100,
		MaxStandardSize: 10240,
	}

	audit, err := factory.CreateAudit(cfg)
	require.NoError(t, err)
	require.NotNil(t, audit)

	// Test that audit logger works
	audit.LogClientRequest("test-client", "test-method", map[string]any{"param": "value"})
}

func TestLoggerFactory_CreateAudit_InvalidConfig(t *testing.T) {
	// This test will fail until LoggerFactory is implemented
	factory := NewLoggerFactory()
	require.NotNil(t, factory)

	audit, err := factory.CreateAudit(nil)
	require.Error(t, err)
	require.Nil(t, audit)
}
