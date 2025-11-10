// Package logging provides structured logging functionality for the agent-standards-mcp server.
package logging

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/n-r-w/agent-standards-mcp/internal/config"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	// maxLogFileSize is the maximum size of a log file before rotation (100MB).
	maxLogFileSize = 100
	// maxLogFiles is the maximum number of old log files to retain.
	maxLogFiles = 7
	// maxLogAge is the maximum number of days to retain old log files.
	maxLogAge = 7
	// dirPermissions is the default permissions for directory creation.
	dirPermissions = 0750
)

// LogRotator provides log rotation functionality using lumberjack.
type LogRotator struct {
	lumberjack *lumberjack.Logger
}

// NewLogRotator creates a new LogRotator with the given configuration.
func NewLogRotator(cfg *config.Config) (*LogRotator, error) {
	if cfg == nil {
		return nil, errors.New("configuration cannot be nil")
	}

	// Create logs directory
	logDir := filepath.Join(cfg.GetFolder(), "logs")
	if err := os.MkdirAll(logDir, dirPermissions); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	// Create lumberjack logger for log rotation
	logFile := filepath.Join(logDir, "agent-standards-mcp.log")
	lumberjackLogger := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    maxLogFileSize, // megabytes
		MaxBackups: maxLogFiles,    // files
		MaxAge:     maxLogAge,      // days
		Compress:   true,           // compress old log files
		LocalTime:  true,           // use local time
	}

	return &LogRotator{
		lumberjack: lumberjackLogger,
	}, nil
}

// Writer returns the underlying writer for the log rotator.
func (lr *LogRotator) Writer() io.Writer {
	return lr.lumberjack
}

// Close closes the log rotator and flushes any pending writes.
func (lr *LogRotator) Close() error {
	// lumberjack.Logger implements io.Closer
	if lr.lumberjack != nil {
		return lr.lumberjack.Close()
	}
	return nil
}
