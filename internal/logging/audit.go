// Package logging provides structured logging functionality for the agent-standards-mcp server.
package logging

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/n-r-w/agent-standards-mcp/internal/config"
	"github.com/n-r-w/agent-standards-mcp/internal/shared"
)

// Audit provides audit logging functionality for client requests.
type Audit struct {
	logger *slog.Logger
}

var _ shared.AuditLogger = (*Audit)(nil)

// NewAudit creates a new Audit logger with the given configuration.
func NewAudit(cfg *config.Config) (*Audit, error) {
	if cfg == nil {
		return nil, errors.New("configuration cannot be nil")
	}

	// Create structured logger for audit logging
	structuredLogger, err := NewStructuredLogger(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create structured logger for audit: %w", err)
	}

	return &Audit{
		logger: structuredLogger.logger,
	}, nil
}

// LogClientRequest logs a client request with structured data.
func (a *Audit) LogClientRequest(clientID string, method string, params any) {
	// Stub implementation - will be implemented later
	a.logger.Info("client_request",
		"client_id", clientID,
		"method", method,
		"params", params,
	)
}

// LogClientResponse logs a client response with structured data.
func (a *Audit) LogClientResponse(clientID string, result any, err error) {
	// Stub implementation - will be implemented later
	if err != nil {
		a.logger.Error("client_response",
			"client_id", clientID,
			"error", err.Error(),
		)
	} else {
		a.logger.Info("client_response",
			"client_id", clientID,
			"result", result,
		)
	}
}
