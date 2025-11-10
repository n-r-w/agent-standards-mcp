// Package server provides MCP server implementation for agent-standards-mcp server.
package server

import (
	"context"
	"errors"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/n-r-w/agent-standards-mcp/internal/config"
	"github.com/n-r-w/agent-standards-mcp/internal/domain"
	"github.com/n-r-w/agent-standards-mcp/internal/prompt"
	"github.com/n-r-w/agent-standards-mcp/internal/shared"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestNewServer(t *testing.T) {
	cfg := createTestConfig()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := shared.NewMockLogger(ctrl)
	auditLogger := shared.NewMockAuditLogger(ctrl)
	standardLoader := NewMockStandardLoader(ctrl)

	server, err := New(cfg, logger, auditLogger, standardLoader)
	require.NoError(t, err)
	require.NotNil(t, server)

	// Verify server has correct dependencies
	assert.Equal(t, cfg, server.cfg)
	assert.Equal(t, logger, server.logger)
	assert.Equal(t, auditLogger, server.auditLogger)
	assert.Equal(t, standardLoader, server.standardLoader)
	assert.NotNil(t, server.server)
}

func TestServer_Start(t *testing.T) {
	server, ctrl := createTestServer(t)
	defer ctrl.Finish()

	// Mock logger expectation for Info call

	// Note: We don't actually call Start() as it would block
	// This test verifies server structure is correct for starting
	require.NotNil(t, server.server)
}

func TestServer_Stop(t *testing.T) {
	server, ctrl := createTestServer(t)
	defer ctrl.Finish()

	// Mock logger expectation for Info call
	server.logger.(*shared.MockLogger).EXPECT().
		Info("Stopping MCP server")

	// Test Stop method
	err := server.Stop(context.Background())
	require.NoError(t, err)
}

// Test helper functions

func createTestConfig() *config.Config {
	return &config.Config{
		LogLevel:        "ERROR",
		Folder:          "/tmp",
		MaxStandards:    100,
		MaxStandardSize: 10240,
	}
}

func createTestServer(t *testing.T) (*MCP, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	logger := shared.NewMockLogger(ctrl)
	auditLogger := shared.NewMockAuditLogger(ctrl)
	standardLoader := NewMockStandardLoader(ctrl)

	server, err := New(createTestConfig(), logger, auditLogger, standardLoader)
	require.NoError(t, err)
	require.NotNil(t, server)

	return server, ctrl
}

func createTestStandardInfo(name, description string) domain.StandardInfo {
	return domain.StandardInfo{
		Name:        name,
		Description: description,
	}
}

func createTestStandard(name, description, content string) domain.Standard {
	return domain.Standard{
		Name:        name,
		Description: description,
		Content:     content,
	}
}

// Tests for handleListStandards

func TestMCP_handleListStandards_Success(t *testing.T) {
	server, ctrl := createTestServer(t)
	defer ctrl.Finish()

	ctx := context.Background()
	request := &mcp.CallToolRequest{
		Session: nil,
		Params:  nil,
		Extra:   nil,
	}
	input := map[string]any{"limit": 10}

	expectedStandards := []domain.StandardInfo{
		createTestStandardInfo("test-standard-1", "Test standard 1"),
		createTestStandardInfo("test-standard-2", "Test standard 2"),
	}

	// Set up mock expectations
	server.standardLoader.(*MockStandardLoader).EXPECT().
		ListStandards(ctx).
		Return(expectedStandards, nil)

	server.auditLogger.(*shared.MockAuditLogger).EXPECT().
		LogClientRequest("mcp-client", "list_standards", input)

	server.auditLogger.(*shared.MockAuditLogger).EXPECT().
		LogClientResponse("mcp-client", gomock.Any(), nil)

	// Call handler
	result, err := server.handleListStandards(ctx, request, input)

	// Assertions
	require.NoError(t, err)
	require.NotNil(t, result)
	require.False(t, result.IsError)
	require.Len(t, result.Content, 1)

	// Check that content is plain text
	textContent, ok := result.Content[0].(*mcp.TextContent)
	require.True(t, ok)
	expectedText := prompt.LoadRelevantStandardsPrompt() + "\ntest-standard-1: Test standard 1\ntest-standard-2: Test standard 2"
	assert.Equal(t, expectedText, textContent.Text)
}

func TestMCP_handleListStandards_EmptyResult(t *testing.T) {
	server, ctrl := createTestServer(t)
	defer ctrl.Finish()

	ctx := context.Background()
	request := &mcp.CallToolRequest{
		Session: nil,
		Params:  nil,
		Extra:   nil,
	}
	input := map[string]any{}

	expectedStandards := []domain.StandardInfo{}

	// Set up mock expectations
	server.standardLoader.(*MockStandardLoader).EXPECT().
		ListStandards(ctx).
		Return(expectedStandards, nil)

	server.auditLogger.(*shared.MockAuditLogger).EXPECT().
		LogClientRequest("mcp-client", "list_standards", input)

	server.auditLogger.(*shared.MockAuditLogger).EXPECT().
		LogClientResponse("mcp-client", gomock.Any(), nil)

	// Call handler
	result, err := server.handleListStandards(ctx, request, input)

	// Assertions
	require.NoError(t, err)
	require.NotNil(t, result)
	require.False(t, result.IsError)
	require.Len(t, result.Content, 1)

	// Check that content is plain text
	textContent, ok := result.Content[0].(*mcp.TextContent)
	require.True(t, ok)
	assert.Equal(t, "No standards found.", textContent.Text)
}

func TestMCP_handleListStandards_StandardLoaderError(t *testing.T) {
	server, ctrl := createTestServer(t)
	defer ctrl.Finish()

	ctx := context.Background()
	request := &mcp.CallToolRequest{
		Session: nil,
		Params:  nil,
		Extra:   nil,
	}
	input := map[string]any{}

	expectedError := errors.New("standard loader error")

	// Set up mock expectations
	server.standardLoader.(*MockStandardLoader).EXPECT().
		ListStandards(ctx).
		Return(nil, expectedError)

	server.auditLogger.(*shared.MockAuditLogger).EXPECT().
		LogClientRequest("mcp-client", "list_standards", input)

	server.auditLogger.(*shared.MockAuditLogger).EXPECT().
		LogClientResponse("mcp-client", nil, expectedError)

	// Call handler
	result, err := server.handleListStandards(ctx, request, input)

	// Assertions
	require.Error(t, err)
	require.Equal(t, expectedError, err)
	require.NotNil(t, result)
	require.True(t, result.IsError)
}

// Tests for handleGetStandards

func TestMCP_handleGetStandards_Success(t *testing.T) {
	server, ctrl := createTestServer(t)
	defer ctrl.Finish()

	ctx := context.Background()
	request := &mcp.CallToolRequest{
		Session: nil,
		Params:  nil,
		Extra:   nil,
	}
	input := map[string]any{
		"standard_names": []string{"test-standard-1", "test-standard-2"},
	}

	expectedStandards := []domain.Standard{
		createTestStandard("test-standard-1", "Test standard 1", "Content 1"),
		createTestStandard("test-standard-2", "Test standard 2", "Content 2"),
	}

	// Set up mock expectations
	server.standardLoader.(*MockStandardLoader).EXPECT().
		GetStandards(ctx, []string{"test-standard-1", "test-standard-2"}).
		Return(expectedStandards, nil)

	server.auditLogger.(*shared.MockAuditLogger).EXPECT().
		LogClientRequest("mcp-client", "get_standards", input)

	server.auditLogger.(*shared.MockAuditLogger).EXPECT().
		LogClientResponse("mcp-client", gomock.Any(), nil)

	// Call handler
	result, err := server.handleGetStandards(ctx, request, input)

	// Assertions
	require.NoError(t, err)
	require.NotNil(t, result)
	require.False(t, result.IsError)
	require.Len(t, result.Content, 1)

	// Check that content is plain text
	textContent, ok := result.Content[0].(*mcp.TextContent)
	require.True(t, ok)
	expectedText := prompt.FollowStandardsPrompt() + "\n\n## test-standard-1: Test standard 1\n```md\nContent 1\n```\n\n------\n\n## test-standard-2: Test standard 2\n```md\nContent 2\n```"
	assert.Equal(t, expectedText, textContent.Text)
}

func TestMCP_handleGetStandards_EmptyResult(t *testing.T) {
	server, ctrl := createTestServer(t)
	defer ctrl.Finish()

	ctx := context.Background()
	request := &mcp.CallToolRequest{
		Session: nil,
		Params:  nil,
		Extra:   nil,
	}
	input := map[string]any{
		"standard_names": []string{"nonexistent-standard"},
	}

	expectedStandards := []domain.Standard{}

	// Set up mock expectations
	server.standardLoader.(*MockStandardLoader).EXPECT().
		GetStandards(ctx, []string{"nonexistent-standard"}).
		Return(expectedStandards, nil)

	server.auditLogger.(*shared.MockAuditLogger).EXPECT().
		LogClientRequest("mcp-client", "get_standards", input)

	server.auditLogger.(*shared.MockAuditLogger).EXPECT().
		LogClientResponse("mcp-client", gomock.Any(), nil)

	// Call handler
	result, err := server.handleGetStandards(ctx, request, input)

	// Assertions
	require.NoError(t, err)
	require.NotNil(t, result)
	require.False(t, result.IsError)
	require.Len(t, result.Content, 1)

	// Check that content is plain text
	textContent, ok := result.Content[0].(*mcp.TextContent)
	require.True(t, ok)
	assert.Equal(t, "No standards found.", textContent.Text)
}

// Tests for handleGetStandards input validation

func TestMCP_handleGetStandards_MissingStandardNamesParam(t *testing.T) {
	server, ctrl := createTestServer(t)
	defer ctrl.Finish()

	ctx := context.Background()
	request := &mcp.CallToolRequest{
		Session: nil,
		Params:  nil,
		Extra:   nil,
	}
	input := map[string]any{} // Missing standard_names parameter

	expectedError := errors.New("standard_names parameter is required")

	// Set up mock expectations
	server.auditLogger.(*shared.MockAuditLogger).EXPECT().
		LogClientRequest("mcp-client", "get_standards", input)

	server.auditLogger.(*shared.MockAuditLogger).EXPECT().
		LogClientResponse("mcp-client", nil, expectedError)

	// Call handler
	result, err := server.handleGetStandards(ctx, request, input)

	// Assertions
	require.Error(t, err)
	require.Equal(t, expectedError, err)
	require.NotNil(t, result)
	require.True(t, result.IsError)
}

func TestMCP_handleGetStandards_StandardNamesNotArray(t *testing.T) {
	server, ctrl := createTestServer(t)
	defer ctrl.Finish()

	ctx := context.Background()
	request := &mcp.CallToolRequest{
		Session: nil,
		Params:  nil,
		Extra:   nil,
	}
	input := map[string]any{
		"standard_names": "not-an-array", // Should be array
	}

	expectedError := errors.New("standard_names must be an array of strings")

	// Set up mock expectations
	server.auditLogger.(*shared.MockAuditLogger).EXPECT().
		LogClientRequest("mcp-client", "get_standards", input)

	server.auditLogger.(*shared.MockAuditLogger).EXPECT().
		LogClientResponse("mcp-client", nil, expectedError)

	// Call handler
	result, err := server.handleGetStandards(ctx, request, input)

	// Assertions
	require.Error(t, err)
	require.Equal(t, expectedError, err)
	require.NotNil(t, result)
	require.True(t, result.IsError)
}

func TestMCP_handleGetStandards_StandardNamesArrayWithNonStrings(t *testing.T) {
	server, ctrl := createTestServer(t)
	defer ctrl.Finish()

	ctx := context.Background()
	request := &mcp.CallToolRequest{
		Session: nil,
		Params:  nil,
		Extra:   nil,
	}
	input := map[string]any{
		"standard_names": []any{"valid-string", 123, "another-string"}, // Contains non-string
	}

	expectedError := errors.New("standard_names must be an array of strings")

	// Set up mock expectations
	server.auditLogger.(*shared.MockAuditLogger).EXPECT().
		LogClientRequest("mcp-client", "get_standards", input)

	server.auditLogger.(*shared.MockAuditLogger).EXPECT().
		LogClientResponse("mcp-client", nil, expectedError)

	// Call handler
	result, err := server.handleGetStandards(ctx, request, input)

	// Assertions
	require.Error(t, err)
	require.Equal(t, expectedError, err)
	require.NotNil(t, result)
	require.True(t, result.IsError)
}

// Tests for handleGetStandards error scenarios

func TestMCP_handleGetStandards_StandardLoaderError(t *testing.T) {
	server, ctrl := createTestServer(t)
	defer ctrl.Finish()

	ctx := context.Background()
	request := &mcp.CallToolRequest{
		Session: nil,
		Params:  nil,
		Extra:   nil,
	}
	input := map[string]any{
		"standard_names": []string{"test-standard"},
	}

	expectedError := errors.New("standard loader error")

	// Set up mock expectations
	server.standardLoader.(*MockStandardLoader).EXPECT().
		GetStandards(ctx, []string{"test-standard"}).
		Return(nil, expectedError)

	server.auditLogger.(*shared.MockAuditLogger).EXPECT().
		LogClientRequest("mcp-client", "get_standards", input)

	server.auditLogger.(*shared.MockAuditLogger).EXPECT().
		LogClientResponse("mcp-client", nil, expectedError)

	// Call handler
	result, err := server.handleGetStandards(ctx, request, input)

	// Assertions
	require.Error(t, err)
	require.Equal(t, expectedError, err)
	require.NotNil(t, result)
	require.True(t, result.IsError)
}

// Edge case tests

func TestMCP_handleListStandards_SpecialCharacters(t *testing.T) {
	server, ctrl := createTestServer(t)
	defer ctrl.Finish()

	ctx := context.Background()
	request := &mcp.CallToolRequest{
		Session: nil,
		Params:  nil,
		Extra:   nil,
	}
	input := map[string]any{}

	expectedStandards := []domain.StandardInfo{
		createTestStandardInfo("standard-with-特殊字符", "Standard with special characters: ñáéíóú"),
		createTestStandardInfo("standard-with-emoji", "Standard with emoji: 🚀🔧"),
	}

	// Set up mock expectations
	server.standardLoader.(*MockStandardLoader).EXPECT().
		ListStandards(ctx).
		Return(expectedStandards, nil)

	server.auditLogger.(*shared.MockAuditLogger).EXPECT().
		LogClientRequest("mcp-client", "list_standards", input)

	server.auditLogger.(*shared.MockAuditLogger).EXPECT().
		LogClientResponse("mcp-client", gomock.Any(), nil)

	// Call handler
	result, err := server.handleListStandards(ctx, request, input)

	// Assertions
	require.NoError(t, err)
	require.NotNil(t, result)
	require.False(t, result.IsError)
	require.Len(t, result.Content, 1)

	// Check that content is plain text
	textContent, ok := result.Content[0].(*mcp.TextContent)
	require.True(t, ok)
	expectedText := prompt.LoadRelevantStandardsPrompt() + "\nstandard-with-特殊字符: Standard with special characters: ñáéíóú\nstandard-with-emoji: Standard with emoji: 🚀🔧"
	assert.Equal(t, expectedText, textContent.Text)
}

func TestMCP_handleGetStandards_LargeContent(t *testing.T) {
	server, ctrl := createTestServer(t)
	defer ctrl.Finish()

	ctx := context.Background()
	request := &mcp.CallToolRequest{
		Session: nil,
		Params:  nil,
		Extra:   nil,
	}
	input := map[string]any{
		"standard_names": []string{"large-standard"},
	}

	// Create content that's close to maximum size limit
	largeContent := string(make([]byte, 10200)) // 10KB content
	expectedStandards := []domain.Standard{
		createTestStandard("large-standard", "Large standard", largeContent),
	}

	// Set up mock expectations
	server.standardLoader.(*MockStandardLoader).EXPECT().
		GetStandards(ctx, []string{"large-standard"}).
		Return(expectedStandards, nil)

	server.auditLogger.(*shared.MockAuditLogger).EXPECT().
		LogClientRequest("mcp-client", "get_standards", input)

	server.auditLogger.(*shared.MockAuditLogger).EXPECT().
		LogClientResponse("mcp-client", gomock.Any(), nil)

	// Call handler
	result, err := server.handleGetStandards(ctx, request, input)

	// Assertions
	require.NoError(t, err)
	require.NotNil(t, result)
	require.False(t, result.IsError)
	require.Len(t, result.Content, 1)

	// Check that content is plain text
	textContent, ok := result.Content[0].(*mcp.TextContent)
	require.True(t, ok)
	expectedText := prompt.FollowStandardsPrompt() + "\n\n## large-standard: Large standard\n```md\n" + largeContent + "\n```"
	assert.Equal(t, expectedText, textContent.Text)
}

func TestServer_RegisterTools(t *testing.T) {
	server, ctrl := createTestServer(t)
	defer ctrl.Finish()

	// Mock logger expectation for Info call
	server.logger.(*shared.MockLogger).EXPECT().
		Info("Registering MCP tools")

	// Test RegisterTools method
	err := server.RegisterTools()
	require.NoError(t, err)
}
