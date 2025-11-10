// Package main implements the mcp command.
package main

import (
	"context"
	"flag"
	"log/slog"
	"os"

	"github.com/n-r-w/agent-standards-mcp/internal/config"
	"github.com/n-r-w/agent-standards-mcp/internal/logging"
	"github.com/n-r-w/agent-standards-mcp/internal/server"
	"github.com/n-r-w/agent-standards-mcp/internal/standards"
)

// build-time variables that can be set via ldflags
//
//nolint:nolintlint // gochecknoglobals is excluded for this file via .golangci.yml
var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
	builtBy = "unknown"
)

// buildInfo holds build-time information
type buildInfo struct {
	version string
	commit  string
	date    string
	builtBy string
}

// getBuildInfo returns build-time information
func getBuildInfo() buildInfo {
	return buildInfo{
		version: version,
		commit:  commit,
		date:    date,
		builtBy: builtBy,
	}
}

func main() {
	// Add version flag
	showVersion := flag.Bool("version", false, "Show version information")
	flag.Parse()

	if *showVersion {
		info := getBuildInfo()
		logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			AddSource:   false,
			Level:       slog.LevelInfo,
			ReplaceAttr: nil,
		}))
		logger.Info("agent-standards-mcp version info",
			"version", info.version,
			"commit", info.commit,
			"built", info.date,
			"built_by", info.builtBy,
		)
		os.Exit(0)
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	// Create logger factory
	loggerFactory := logging.NewLoggerFactory()

	// Create structured logger
	structuredLogger, err := loggerFactory.CreateStructuredLogger(cfg)
	if err != nil {
		slog.Error("Failed to create structured logger", "error", err)
		os.Exit(1)
	}

	// Create audit logger
	auditLogger, err := loggerFactory.CreateAudit(cfg)
	if err != nil {
		slog.Error("Failed to create audit logger", "error", err)
		os.Exit(1)
	}

	// Test audit logging
	info := getBuildInfo()
	auditLogger.LogClientRequest("test-client", "startup", map[string]any{"version": info.version})

	// Log server startup
	structuredLogger.Info("Starting agent-standards-mcp server",
		"log_level", cfg.GetLogLevel(),
		"standards_folder", cfg.GetFolder(),
		"max_standards", cfg.GetMaxStandards(),
		"max_standard_size", cfg.GetMaxStandardSize(),
	)

	// Create standard loader
	standardLoader := standards.NewFileStandardLoader()

	// Create MCP server
	mcpServer, err := server.New(cfg, structuredLogger, auditLogger, standardLoader)
	if err != nil {
		structuredLogger.Error("Failed to create MCP server", "error", err)
		os.Exit(1)
	}

	// Register MCP tools
	if err := mcpServer.RegisterTools(); err != nil {
		structuredLogger.Error("Failed to register MCP tools", "error", err)
		os.Exit(1)
	}

	// Start server directly (following official MCP SDK pattern)
	ctx := context.Background()
	if err := mcpServer.Start(ctx); err != nil {
		structuredLogger.Error("MCP server failed", "error", err)
		os.Exit(1)
	}
}
