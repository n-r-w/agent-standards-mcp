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
	}, &mcp.ServerOptions{
		Instructions:                prompt.SystemPrompt(),
		Logger:                      nil,
		PageSize:                    0,
		RootsListChangedHandler:     nil,
		ProgressNotificationHandler: nil,
		CompletionHandler:           nil,
		KeepAlive:                   0,
		SubscribeHandler:            nil,
		UnsubscribeHandler:          nil,
		HasPrompts:                  false,
		HasResources:                false,
		HasTools:                    false,
		GetSessionID:                nil,
		InitializedHandler:          nil,
	})

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

// formatStandardInfos formats multiple StandardInfo objects as plain text
func formatStandardInfos(infos []domain.StandardInfo) string {
	if len(infos) == 0 {
		return "No standards found."
	}

	var builder strings.Builder
	builder.WriteString(prompt.LoadRelevantStandardsPrompt())
	builder.WriteByte('\n')

	for i, info := range infos {
		if i > 0 {
			builder.WriteByte('\n')
		}
		builder.WriteString(formatStandardInfo(info))
	}

	return builder.String()
}

func listStandardsStructuredContent(infos []domain.StandardInfo) map[string]any {
	standards := make([]any, len(infos))
	for i, info := range infos {
		standards[i] = map[string]any{
			"name":        info.Name,
			"description": info.Description,
		}
	}

	return map[string]any{"standards": standards}
}

// getStandardsStructuredContent builds the structured JSON content for get_standards.
// Returns {"standards": [{"name": ..., "description": ..., "content": ...}, ...]}
func getStandardsStructuredContent(domainStandards []domain.Standard) map[string]any {
	standards := make([]map[string]any, len(domainStandards))
	for i, standard := range domainStandards {
		standards[i] = map[string]any{
			"name":        standard.Name,
			"description": standard.Description,
			"content":     standard.Content,
		}
	}

	return map[string]any{"standards": standards}
}

// formatStandards formats multiple Standard objects as plain text
func formatStandards(standards []domain.Standard) string {
	var builder strings.Builder

	builder.WriteString(prompt.FollowStandardsPrompt())
	builder.WriteString("\n\n")

	if len(standards) == 0 {
		builder.WriteString("No standards found.")
		return builder.String()
	}

	builder.WriteString("You MUST follow the loaded standards.\n\n")
	fmt.Fprintf(&builder, "Loaded standards: %d\n", len(standards))

	for _, standard := range standards {
		fmt.Fprintf(&builder, "- %s: %s\n", standard.Name, standard.Description)
	}

	builder.WriteString("\nNote: Full standard bodies are provided in ")
	builder.WriteString("StructuredContent.standards[].content; ")
	builder.WriteString("do not repeat them in plain text unless explicitly requested.")

	return builder.String()
}

// RegisterTools registers the list_standards and get_standards tools with the MCP server.
func (s *MCP) RegisterTools() error {
	s.logger.Info("Registering MCP tools")

	// Register list_standards tool
	listStandardsInputSchema := map[string]any{
		"type":       "object",
		"properties": map[string]any{},
	}

	listStandardsOutputSchema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"standards": map[string]any{
				"type": "array",
				"items": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"name": map[string]any{
							"type": "string",
						},
						"description": map[string]any{
							"type": "string",
						},
					},
					"required": []string{"name", "description"},
				},
			},
		},
		"required": []string{"standards"},
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
		*mcp.CallToolResult, map[string]any, error,
	) {
		result, err := s.handleListStandards(ctx, request, input)
		if err != nil {
			return result, nil, err
		}
		// Return structured content matching the contract: {"standards": [...]}
		structuredOut, ok := result.StructuredContent.(map[string]any)
		if !ok {
			return nil, nil, fmt.Errorf(
				"internal error: StructuredContent is %T, expected map[string]any", result.StructuredContent)
		}
		return result, structuredOut, nil
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

	getStandardSchema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"name":        map[string]any{"type": "string"},
			"description": map[string]any{"type": "string"},
			"content":     map[string]any{"type": "string"},
		},
		"required": []string{"name", "description", "content"},
	}

	getStandardsOutputSchema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"standards": map[string]any{
				"type":  "array",
				"items": getStandardSchema,
			},
		},
		"required": []string{"standards"},
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
		*mcp.CallToolResult, map[string]any, error,
	) {
		result, err := s.handleGetStandards(ctx, request, input)
		if err != nil {
			return result, nil, err
		}
		// Return structured content matching the contract: {"standards": [...]}
		structuredOut, ok := result.StructuredContent.(map[string]any)
		if !ok {
			return nil, nil, fmt.Errorf(
				"internal error: StructuredContent is %T, expected map[string]any", result.StructuredContent)
		}
		return result, structuredOut, nil
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
			StructuredContent: nil,
		}, err
	}

	formattedResult := formatStandardInfos(domainResult)
	structuredContent := listStandardsStructuredContent(domainResult)

	s.auditLogger.LogClientResponse("mcp-client", formattedResult, nil)
	return &mcp.CallToolResult{
		IsError:           false,
		Meta:              mcp.Meta{},
		Content:           []mcp.Content{&mcp.TextContent{Meta: mcp.Meta{}, Annotations: nil, Text: formattedResult}},
		StructuredContent: structuredContent,
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
			StructuredContent: nil,
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
			StructuredContent: nil,
		}, err
	}

	domainResult, err := s.standardLoader.GetStandards(ctx, standardNames)
	if err != nil {
		s.auditLogger.LogClientResponse("mcp-client", nil, err)
		return &mcp.CallToolResult{
			IsError:           true,
			Meta:              mcp.Meta{},
			Content:           []mcp.Content{&mcp.TextContent{Meta: mcp.Meta{}, Annotations: nil, Text: err.Error()}},
			StructuredContent: nil,
		}, err
	}

	formattedResult := formatStandards(domainResult)
	structuredContent := getStandardsStructuredContent(domainResult)

	s.auditLogger.LogClientResponse("mcp-client", formattedResult, nil)
	return &mcp.CallToolResult{
		IsError:           false,
		Meta:              mcp.Meta{},
		Content:           []mcp.Content{&mcp.TextContent{Meta: mcp.Meta{}, Annotations: nil, Text: formattedResult}},
		StructuredContent: structuredContent,
	}, nil
}
