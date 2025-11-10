// Package server provides MCP server implementation for the agent-standards-mcp server.
package server

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/n-r-w/agent-standards-mcp/internal/config"
	"github.com/n-r-w/agent-standards-mcp/internal/domain"
	"github.com/n-r-w/agent-standards-mcp/internal/prompt"
	"github.com/n-r-w/agent-standards-mcp/internal/shared"
)

// MCP implements the Server interface using the MCP Go SDK.
type MCP struct {
	cfg            *config.Config
	logger         shared.Logger
	auditLogger    shared.AuditLogger
	standardLoader StandardLoader
	server         *mcp.Server
}

// New creates a new MCP server instance.
func New(
	cfg *config.Config,
	logger shared.Logger,
	auditLogger shared.AuditLogger,
	standardLoader StandardLoader,
) (*MCP, error) {
	if cfg == nil {
		return nil, errors.New("configuration cannot be nil")
	}
	if logger == nil {
		return nil, errors.New("logger cannot be nil")
	}
	if auditLogger == nil {
		return nil, errors.New("audit logger cannot be nil")
	}

	// Create MCP server instance
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "agent-standards-mcp",
		Version: "1.0.0",
		Title:   "Agent Standards MCP Server",
	}, nil)

	return &MCP{
		cfg:            cfg,
		logger:         logger,
		auditLogger:    auditLogger,
		standardLoader: standardLoader,
		server:         server,
	}, nil
}

// Start starts the MCP server with STDIO transport.
func (s *MCP) Start(_ context.Context) error {
	s.logger.Info("Starting MCP server")

	// Create STDIO transport for MCP communication
	transport := &mcp.StdioTransport{}

	// Start serving MCP requests
	return s.server.Run(context.Background(), transport)
}

// Stop gracefully stops the MCP server.
func (s *MCP) Stop(_ context.Context) error {
	s.logger.Info("Stopping MCP server")

	// MCP server doesn't have explicit Close method in this SDK
	// The context cancellation in Run will handle cleanup
	return nil
}

// GetMCPServer returns the underlying MCP server instance for testing purposes.
// This method should only be used in integration tests.
func (s *MCP) GetMCPServer() *mcp.Server {
	return s.server
}

// StartWithTransport starts the MCP server with a custom transport for testing.
// This method should only be used in integration tests.
func (s *MCP) StartWithTransport(ctx context.Context, transport mcp.Transport) error {
	s.logger.Info("Starting MCP server with custom transport")
	return s.server.Run(ctx, transport)
}

// formatStandardInfo formats a single StandardInfo as plain text
func formatStandardInfo(info domain.StandardInfo) string {
	return fmt.Sprintf("%s: %s", info.Name, info.Description)
}

// formatStandard formats a single Standard as plain text with content
func formatStandard(standard domain.Standard) string {
	return fmt.Sprintf("## %s: %s\n```md\n%s\n```", standard.Name, standard.Description, standard.Content)
}

// formatStandardInfos formats multiple StandardInfo objects as plain text
func formatStandardInfos(infos []domain.StandardInfo) string {
	if len(infos) == 0 {
		return "No standards found."
	}

	var builder strings.Builder

	// add prefix
	if len(infos) > 0 {
		builder.WriteString(prompt.LoadRelevantStandardsPrompt() + "\n")
	}

	for i, info := range infos {
		if i > 0 {
			builder.WriteString("\n")
		}
		builder.WriteString(formatStandardInfo(info))
	}

	return builder.String()
}

// formatStandards formats multiple Standard objects as plain text
func formatStandards(standards []domain.Standard) string {
	if len(standards) == 0 {
		return "No standards found."
	}

	var builder strings.Builder

	builder.WriteString(prompt.FollowStandardsPrompt() + "\n\n")

	for i, standard := range standards {
		if i > 0 {
			builder.WriteString("\n\n------\n\n")
		}
		builder.WriteString(formatStandard(standard))
	}
	return builder.String()
}

// RegisterTools registers the list_standards and get_standards tools with the MCP server.
func (s *MCP) RegisterTools() error {
	s.logger.Info("Registering MCP tools")

	// Register list_standards tool
	listStandardsInputSchema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"limit": map[string]any{
				"type":        "integer",
				"description": "Maximum number of standards to return",
			},
		},
	}

	listStandardsOutputSchema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"result": map[string]any{
				"type":        "string",
				"description": "Plain text formatted list of standards with names and descriptions",
			},
		},
	}

	mcp.AddTool(s.server, &mcp.Tool{
		Name:         "list_standards",
		Description:  prompt.ListStandardsPrompt(),
		InputSchema:  listStandardsInputSchema,
		OutputSchema: listStandardsOutputSchema,
		Meta:         mcp.Meta{},
		Annotations:  nil,
		Title:        "List Standards",
	}, func(ctx context.Context, request *mcp.CallToolRequest, input map[string]any) (
		*mcp.CallToolResult, map[string]string, error,
	) {
		result, err := s.handleListStandards(ctx, request, input)
		if err != nil {
			return result, nil, err
		}
		// Extract text content from the result
		var textResult string
		if len(result.Content) > 0 {
			if textContent, ok := result.Content[0].(*mcp.TextContent); ok {
				textResult = textContent.Text
			}
		}
		return result, map[string]string{"result": textResult}, nil
	})

	// Register get_standards tool
	getStandardsInputSchema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"standard_names": map[string]any{
				"type": "array",
				"items": map[string]any{
					"type": "string",
				},
				"description": "List of standard names to retrieve",
			},
		},
		"required": []string{"standard_names"},
	}

	getStandardsOutputSchema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"result": map[string]any{
				"type":        "string",
				"description": "Plain text formatted standards with names, descriptions, and content",
			},
		},
	}

	mcp.AddTool(s.server, &mcp.Tool{
		Name:         "get_standards",
		Description:  prompt.GetStandardsPrompt(),
		InputSchema:  getStandardsInputSchema,
		OutputSchema: getStandardsOutputSchema,
		Meta:         mcp.Meta{},
		Annotations:  nil,
		Title:        "Get Standards",
	}, func(ctx context.Context, request *mcp.CallToolRequest, input map[string]any) (
		*mcp.CallToolResult, map[string]string, error,
	) {
		result, err := s.handleGetStandards(ctx, request, input)
		if err != nil {
			return result, nil, err
		}
		// Extract text content from the result
		var textResult string
		if len(result.Content) > 0 {
			if textContent, ok := result.Content[0].(*mcp.TextContent); ok {
				textResult = textContent.Text
			}
		}
		return result, map[string]string{"result": textResult}, nil
	})

	return nil
}

// handleListStandards handles the list_standards tool request.
func (s *MCP) handleListStandards(ctx context.Context, _ *mcp.CallToolRequest, input map[string]any) (
	*mcp.CallToolResult,
	error,
) {
	s.auditLogger.LogClientRequest("mcp-client", "list_standards", input)

	domainResult, err := s.standardLoader.ListStandards(ctx)
	if err != nil {
		s.auditLogger.LogClientResponse("mcp-client", nil, err)
		return &mcp.CallToolResult{
			IsError:           true,
			Meta:              mcp.Meta{},
			Content:           []mcp.Content{&mcp.TextContent{Meta: mcp.Meta{}, Annotations: nil, Text: err.Error()}},
			StructuredContent: err.Error(),
		}, err
	}

	formattedResult := formatStandardInfos(domainResult)

	// Return formatted plain text result
	s.auditLogger.LogClientResponse("mcp-client", formattedResult, nil)
	return &mcp.CallToolResult{
		IsError:           false,
		Meta:              mcp.Meta{},
		Content:           []mcp.Content{&mcp.TextContent{Meta: mcp.Meta{}, Annotations: nil, Text: formattedResult}},
		StructuredContent: formattedResult,
	}, nil
}

// handleGetStandards handles the get_standards tool request.
func (s *MCP) handleGetStandards(ctx context.Context, _ *mcp.CallToolRequest, input map[string]any) (
	*mcp.CallToolResult,
	error,
) {
	s.auditLogger.LogClientRequest("mcp-client", "get_standards", input)

	// Extract standard names from input
	standardNamesRaw, ok := input["standard_names"]
	if !ok {
		err := errors.New("standard_names parameter is required")
		s.auditLogger.LogClientResponse("mcp-client", nil, err)
		return &mcp.CallToolResult{
			IsError:           true,
			Meta:              mcp.Meta{},
			Content:           []mcp.Content{&mcp.TextContent{Meta: mcp.Meta{}, Annotations: nil, Text: err.Error()}},
			StructuredContent: err.Error(),
		}, err
	}

	// Convert standardNamesRaw to []string, handling both []string and []any cases
	var standardNames []string
	var err error

	switch standardNamesTyped := standardNamesRaw.(type) {
	case []string:
		// Direct case (usually from unit tests)
		standardNames = standardNamesTyped
	case []any:
		// JSON unmarshaled case (usually from integration tests)
		standardNames = make([]string, len(standardNamesTyped))
		for i, v := range standardNamesTyped {
			standardName, ok := v.(string)
			if !ok {
				err = errors.New("standard_names must be an array of strings")
				break
			}
			standardNames[i] = standardName
		}
	default:
		err = errors.New("standard_names must be an array of strings")
	}

	if err != nil {
		s.auditLogger.LogClientResponse("mcp-client", nil, err)
		return &mcp.CallToolResult{
			IsError:           true,
			Meta:              mcp.Meta{},
			Content:           []mcp.Content{&mcp.TextContent{Meta: mcp.Meta{}, Annotations: nil, Text: err.Error()}},
			StructuredContent: err.Error(),
		}, err
	}

	domainResult, err := s.standardLoader.GetStandards(ctx, standardNames)
	if err != nil {
		s.auditLogger.LogClientResponse("mcp-client", nil, err)
		return &mcp.CallToolResult{
			IsError:           true,
			Meta:              mcp.Meta{},
			Content:           []mcp.Content{&mcp.TextContent{Meta: mcp.Meta{}, Annotations: nil, Text: err.Error()}},
			StructuredContent: err.Error(),
		}, err
	}

	formattedResult := formatStandards(domainResult)

	// Return formatted plain text result
	s.auditLogger.LogClientResponse("mcp-client", formattedResult, nil)
	return &mcp.CallToolResult{
		IsError:           false,
		Meta:              mcp.Meta{},
		Content:           []mcp.Content{&mcp.TextContent{Meta: mcp.Meta{}, Annotations: nil, Text: formattedResult}},
		StructuredContent: formattedResult,
	}, nil
}
