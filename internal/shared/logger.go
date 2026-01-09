// Package shared provides common interfaces for logging functionality.
package shared //nolint:revive,nolintlint // i like this name :)

//go:generate mockgen -source=logger.go -destination=logger_mock.go -package=shared

// Logger defines the interface for structured logging operations.
type Logger interface {
	// Debug logs a debug message with structured data.
	Debug(msg string, args ...any)

	// Info logs an info message with structured data.
	Info(msg string, args ...any)

	// Warn logs a warning message with structured data.
	Warn(msg string, args ...any)

	// Error logs an error message with structured data.
	Error(msg string, args ...any)
}

// AuditLogger defines the interface for audit logging operations.
type AuditLogger interface {
	// LogClientRequest logs a client request with structured data.
	LogClientRequest(clientID string, method string, params any)

	// LogClientResponse logs a client response with structured data.
	LogClientResponse(clientID string, result any, err error)
}
