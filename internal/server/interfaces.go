// Package server provides interfaces for the MCP server implementation.
package server

import (
	"context"

	"github.com/n-r-w/agent-standards-mcp/internal/domain"
)

//go:generate mockgen -source=interfaces.go -destination=mocks.go -package=server

// StandardLoader defines the interface for loading standards from the file system.
type StandardLoader interface {
	// ListStandards returns a list of available standard information (name and description).
	ListStandards(ctx context.Context) ([]domain.StandardInfo, error)

	// GetStandards returns the full content of specific standards by their names.
	GetStandards(ctx context.Context, standardNames []string) ([]domain.Standard, error)
}
