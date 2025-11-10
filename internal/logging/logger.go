// Package logging provides structured logging functionality for the agent-standards-mcp server.
package logging

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/n-r-w/agent-standards-mcp/internal/config"
	"github.com/n-r-w/agent-standards-mcp/internal/shared"
)

const (
	// disabledLogLevel is a high level number that effectively disables logging.
	disabledLogLevel = 100
)

// StructuredLogger provides structured logging functionality with slog.
type StructuredLogger struct {
	logger     *slog.Logger
	logRotator *LogRotator
}

var _ shared.Logger = (*StructuredLogger)(nil)

// NewStructuredLogger creates a new StructuredLogger with the given configuration.
func NewStructuredLogger(cfg *config.Config) (*StructuredLogger, error) {
	// Validate configuration
	if cfg == nil {
		return nil, errors.New("configuration cannot be nil")
	}

	// Check if logging is disabled
	if !cfg.IsLoggingEnabled() {
		// Return a logger that writes to /dev/null
		return &StructuredLogger{
			logger:     slog.New(slog.DiscardHandler),
			logRotator: nil,
		}, nil
	}

	// Convert config log level to slog level
	logLevel := cfg.GetLogLevel()
	var slogLevel slog.Level
	switch logLevel {
	case config.LogLevelNone:
		slogLevel = slog.Level(disabledLogLevel) // Disabled logging
	case config.LogLevelDebug:
		slogLevel = slog.LevelDebug
	case config.LogLevelInfo:
		slogLevel = slog.LevelInfo
	case config.LogLevelWarn:
		slogLevel = slog.LevelWarn
	case config.LogLevelError:
		slogLevel = slog.LevelError
	default:
		return nil, fmt.Errorf("invalid log level: %s", logLevel)
	}

	// Create handler with stderr output (MCP compliance)
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level:     slogLevel,
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Only replace source attribute, preserve others for proper structured logging
			if a.Key == "source" {
				return slog.String("source", "agent-standards-mcp")
			}
			return a
		},
	})

	// Create logger
	logger := slog.New(handler)

	var logRotator *LogRotator

	// If logging is enabled, also set up file logging with rotation
	if cfg.IsLoggingEnabled() {
		// Create log rotator for file output
		rotator, err := NewLogRotator(cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to create log rotator: %w", err)
		}
		logRotator = rotator

		// Create multi-writer for both stderr and file
		multiWriter := io.MultiWriter(os.Stderr, rotator.Writer())

		// Create handler with dual output
		handler = slog.NewTextHandler(multiWriter, &slog.HandlerOptions{
			Level:     slogLevel,
			AddSource: true,
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				// Only replace source attribute, preserve others for proper structured logging
				if a.Key == "source" {
					return slog.String("source", "agent-standards-mcp")
				}
				return a
			},
		})

		// Create logger with dual output
		logger = slog.New(handler)
	}

	return &StructuredLogger{
		logger:     logger,
		logRotator: logRotator,
	}, nil
}

// Debug logs a debug message with structured data.
func (s *StructuredLogger) Debug(msg string, args ...any) {
	s.logger.Debug(msg, args...)
}

// Info logs an info message with structured data.
func (s *StructuredLogger) Info(msg string, args ...any) {
	s.logger.Info(msg, args...)
}

// Warn logs a warning message with structured data.
func (s *StructuredLogger) Warn(msg string, args ...any) {
	s.logger.Warn(msg, args...)
}

// Error logs an error message with structured data.
func (s *StructuredLogger) Error(msg string, args ...any) {
	s.logger.Error(msg, args...)
}

// Close closes the structured logger and any underlying resources.
func (s *StructuredLogger) Close() error {
	if s.logRotator != nil {
		return s.logRotator.Close()
	}
	return nil
}

// LoggerFactory creates configured logger instances.
type LoggerFactory struct{}

// NewLoggerFactory creates a new LoggerFactory.
func NewLoggerFactory() *LoggerFactory {
	return &LoggerFactory{}
}

// CreateStructuredLogger creates a new StructuredLogger with the given configuration.
func (lf *LoggerFactory) CreateStructuredLogger(cfg *config.Config) (*StructuredLogger, error) {
	return NewStructuredLogger(cfg)
}

// CreateAudit creates a new Audit logger with the given configuration.
func (lf *LoggerFactory) CreateAudit(cfg *config.Config) (*Audit, error) {
	return NewAudit(cfg)
}
