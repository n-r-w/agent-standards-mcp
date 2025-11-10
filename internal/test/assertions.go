// Package test provides testing utilities and assertions for the agent-standards-mcp server.
package test

import (
	"context"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/require"
)

// AssertToolsAvailable checks that the expected tools are available in server
func AssertToolsAvailable(t *testing.T, suite *Suite, expectedTools []string) {
	ctx := getContext()

	// Get available tools
	tools, err := suite.ClientSession.ListTools(ctx, &mcp.ListToolsParams{
		Meta:   mcp.Meta{},
		Cursor: "",
	})
	require.NoError(t, err, "Failed to get tools from MCP server")
	require.NotEmpty(t, tools.Tools, "No tools available from MCP server")

	// Create map of available tool names
	availableTools := make(map[string]bool)
	for _, tool := range tools.Tools {
		availableTools[tool.Name] = true
	}

	// Check that all expected tools are available
	for _, expectedTool := range expectedTools {
		require.True(t, availableTools[expectedTool],
			"Expected tool '%s' not found in available tools", expectedTool)
	}
}

// AssertToolCallSuccess calls a tool and verifies successful execution
func AssertToolCallSuccess(t *testing.T, suite *Suite, toolName string, args map[string]any) *mcp.CallToolResult {
	ctx := getContext()

	result, err := suite.ClientSession.CallTool(ctx, &mcp.CallToolParams{
		Meta:      mcp.Meta{},
		Name:      toolName,
		Arguments: args,
	})

	require.NoError(t, err, "Tool call should not return error")
	require.NotNil(t, result, "Tool call result should not be nil")
	require.False(t, result.IsError, "Tool call should not be marked as error")

	return result
}

// AssertToolCallError calls a tool and verifies it returns an error
func AssertToolCallError(t *testing.T, suite *Suite, toolName string, args map[string]any) *mcp.CallToolResult {
	ctx := getContext()

	result, err := suite.ClientSession.CallTool(ctx, &mcp.CallToolParams{
		Meta:      mcp.Meta{},
		Name:      toolName,
		Arguments: args,
	})

	require.NoError(t, err, "Tool call request should succeed even if tool returns error")
	require.NotNil(t, result, "Tool call result should not be nil")
	require.True(t, result.IsError, "Tool call should be marked as error")

	return result
}

// AssertPlainTextInput validates that the result contains plain text content
func AssertPlainTextInput(t *testing.T, result *mcp.CallToolResult) string {
	require.NotNil(t, result.StructuredContent, "Result should contain structured content")
	require.NotEmpty(t, result.Content, "Result should contain content")
	require.Len(t, result.Content, 1, "Result should contain exactly one content item")

	// Extract text content
	textContent, ok := result.Content[0].(*mcp.TextContent)
	require.True(t, ok, "Content should be TextContent")
	require.NotEmpty(t, textContent.Text, "Text content should not be empty")

	return textContent.Text
}

// AssertStandardListContains validates that plain text contains a specific standard by name
func AssertStandardListContains(t *testing.T, plainText string, standardName string) {
	expectedPattern := standardName + ":"
	require.Contains(t, plainText, expectedPattern,
		"Plain text should contain standard '%s' with expected format", standardName)
}

// AssertStandardListCount validates that plain text contains expected number of standards
func AssertStandardListCount(t *testing.T, plainText string, expectedCount int) {
	if expectedCount == 0 {
		require.Equal(t, "No standards found.", plainText, "Empty result should return 'No standards found.'")
		return
	}

	// Check if this is get_standards format (markdown) or list_standards format (plain text)
	if strings.Contains(plainText, "## ") {
		// get_standards format - count standard headers (lines that start with "## " and end with ":")
		lines := strings.Split(plainText, "\n")
		standardCount := 0
		for _, line := range lines {
			line = strings.TrimSpace(line)
			// Count lines that start with "## " and contain ":" (standard headers)
			if strings.HasPrefix(line, "## ") && strings.Contains(line, ":") {
				standardCount++
			}
		}
		require.Equal(t, expectedCount, standardCount,
			"Markdown text should contain exactly %d standards", expectedCount)
	} else {
		// list_standards format - count lines with standard name pattern
		lines := strings.Split(plainText, "\n")
		standardCount := 0
		for _, line := range lines {
			// Count non-empty lines with standard name pattern
			line = strings.TrimSpace(line)
			if line != "" && strings.Contains(line, ":") && !strings.HasPrefix(line, "#") {
				standardCount++
			}
		}
		require.Equal(t, expectedCount, standardCount,
			"Plain text should contain exactly %d standards", expectedCount)
	}
}

// AssertStandardContainsDescription validates that a standard in plain text contains expected description
func AssertStandardContainsDescription(t *testing.T, plainText string, standardName, expectedDescription string) {
	// For list_standards, look for "name: description" pattern
	expectedPattern := standardName + ": " + expectedDescription
	require.Contains(t, plainText, expectedPattern,
		"Plain text should contain standard '%s' with description '%s'", standardName, expectedDescription)
}

// AssertStandardContainsContent validates that a standard in plain text contains expected content
func AssertStandardContainsContent(t *testing.T, plainText string, standardName, expectedContent string) {
	// For get_standards, look for markdown format: "## name: description"
	standardHeader := "## " + standardName + ":"
	standardStart := strings.Index(plainText, standardHeader)
	require.GreaterOrEqual(t, standardStart, 0,
		"Plain text should contain standard '%s' with markdown header", standardName)

	// Look for content after the markdown header
	afterHeader := plainText[standardStart:]
	codeBlockStart := strings.Index(afterHeader, "```md\n")
	require.GreaterOrEqual(t, codeBlockStart, 0, "Should find code block after standard header")

	const codeBlockPrefixLength = 6 // len("```md\n")
	contentStart := standardStart + codeBlockStart + codeBlockPrefixLength
	if contentStart >= len(plainText) {
		return
	}

	standardContent := plainText[contentStart:]

	// Find end of code block
	codeBlockEnd := strings.Index(standardContent, "\n```")
	if codeBlockEnd >= 0 {
		standardContent = standardContent[:codeBlockEnd]
	}

	require.Contains(t, standardContent, expectedContent,
		"Standard '%s' should contain expected content", standardName)
}

// AssertGetStandardsContainsContent validates that get_standards result contains full content for a standard
func AssertGetStandardsContainsContent(
	t *testing.T,
	plainText string,
	standardName, expectedDescription, expectedContent string,
) {
	// For get_standards, check for markdown format
	standardHeader := "## " + standardName + ": " + expectedDescription
	require.Contains(t, plainText, standardHeader,
		"Plain text should contain standard '%s' with markdown header and description", standardName)

	// Check for content
	AssertStandardContainsContent(t, plainText, standardName, expectedContent)
}

// AssertMultipleStandardsFormat validates that multiple standards are properly separated
func AssertMultipleStandardsFormat(t *testing.T, plainText string) {
	// Check if this is get_standards format (markdown) or list_standards format (plain text)
	isMarkdownFormat := strings.Contains(plainText, "## ")

	if isMarkdownFormat {
		// get_standards format - check for markdown separators
		if strings.Count(plainText, "## ") > 1 {
			require.Contains(t, plainText, "\n\n------\n\n",
				"Multiple standards should be separated by '------' in markdown format")
		}
		return
	}

	// list_standards format - standards are separated by newlines, not double newlines
	if !strings.Contains(plainText, ":") {
		return
	}

	lines := strings.Split(plainText, "\n")
	standardCount := 0
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && strings.Contains(line, ":") && !strings.HasPrefix(line, "#") {
			standardCount++
		}
	}

	if standardCount > 1 {
		// Check that standards are separated by newlines (single newline is sufficient)
		require.Contains(t, plainText, "\n",
			"Multiple standards should be separated by newlines")
	}
}

// getContext returns a background context for tool calls
func getContext() context.Context {
	return context.Background()
}
