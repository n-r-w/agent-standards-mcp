// Package test provides testing utilities and assertions for the agent-standards-mcp server.
package test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/n-r-w/agent-standards-mcp/internal/prompt"
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

// ListStandardsResponse represents the expected JSON structure for list_standards StructuredContent.
// This is the contract: {"standards": [{"name": "...", "description": "..."}, ...]}
type ListStandardsResponse struct {
	Standards []StandardItem `json:"standards"`
}

// StandardItem represents a single standard in the list_standards response.
type StandardItem struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// GetStandardsItem represents a single standard item in get_standards StructuredContent response.
type GetStandardsItem struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Content     string `json:"content"`
}

// GetStandardsResponse represents the get_standards StructuredContent response structure.
type GetStandardsResponse struct {
	Standards []GetStandardsItem `json:"standards"`
}

// AssertListStandardsStructuredContent validates that StructuredContent is a JSON object
// with "standards" array field where each item has "name" and "description" string fields.
// This is the expected contract for list_standards tool.
func AssertListStandardsStructuredContent(t *testing.T, result *mcp.CallToolResult) ListStandardsResponse {
	require.NotNil(t, result.StructuredContent, "StructuredContent should not be nil")

	// StructuredContent should be a map (JSON object), not a string
	structuredMap, ok := result.StructuredContent.(map[string]any)
	require.True(t, ok, "StructuredContent should be a JSON object (map[string]any), got %T", result.StructuredContent)

	// Should have "standards" field
	standardsAny, exists := structuredMap["standards"]
	require.True(t, exists, "StructuredContent should have 'standards' field")

	// "standards" should be an array
	standardsArray, ok := standardsAny.([]any)
	require.True(t, ok, "StructuredContent.standards should be an array, got %T", standardsAny)

	// Each item should have "name" and "description" as strings
	var response ListStandardsResponse
	for i, itemAny := range standardsArray {
		itemMap, ok := itemAny.(map[string]any)
		require.True(t, ok, "StructuredContent.standards[%d] should be an object, got %T", i, itemAny)

		nameAny, hasName := itemMap["name"]
		require.True(t, hasName, "StructuredContent.standards[%d] should have 'name' field", i)
		nameStr, ok := nameAny.(string)
		require.True(t, ok, "StructuredContent.standards[%d].name should be string, got %T", i, nameAny)

		descAny, hasDesc := itemMap["description"]
		require.True(t, hasDesc, "StructuredContent.standards[%d] should have 'description' field", i)
		descStr, ok := descAny.(string)
		require.True(t, ok, "StructuredContent.standards[%d].description should be string, got %T", i, descAny)

		response.Standards = append(response.Standards, StandardItem{Name: nameStr, Description: descStr})
	}

	return response
}

// AssertGetStandardsStructuredContent validates the get_standards StructuredContent follows the JSON contract.
// It asserts:
// - StructuredContent is a map (JSON object)
// - Has required field "standards" which is an array
// - Each item has required string fields: "name", "description", "content"
// Returns the parsed response for further assertions.
func AssertGetStandardsStructuredContent(t *testing.T, result *mcp.CallToolResult) GetStandardsResponse {
	require.NotNil(t, result.StructuredContent, "StructuredContent should not be nil")

	// StructuredContent should be a map (JSON object), not a string
	structuredMap, ok := result.StructuredContent.(map[string]any)
	require.True(t, ok, "StructuredContent should be a JSON object (map[string]any), got %T", result.StructuredContent)

	// Should have "standards" field
	standardsAny, exists := structuredMap["standards"]
	require.True(t, exists, "StructuredContent should have 'standards' field")

	// "standards" should be an array
	standardsArray, ok := standardsAny.([]any)
	require.True(t, ok, "StructuredContent.standards should be an array, got %T", standardsAny)

	// Each item should have "name", "description", and "content" as strings
	var response GetStandardsResponse
	for i, itemAny := range standardsArray {
		itemMap, ok := itemAny.(map[string]any)
		require.True(t, ok, "StructuredContent.standards[%d] should be an object, got %T", i, itemAny)

		nameAny, hasName := itemMap["name"]
		require.True(t, hasName, "StructuredContent.standards[%d] should have 'name' field", i)
		nameStr, ok := nameAny.(string)
		require.True(t, ok, "StructuredContent.standards[%d].name should be string, got %T", i, nameAny)

		descAny, hasDesc := itemMap["description"]
		require.True(t, hasDesc, "StructuredContent.standards[%d] should have 'description' field", i)
		descStr, ok := descAny.(string)
		require.True(t, ok, "StructuredContent.standards[%d].description should be string, got %T", i, descAny)

		contentAny, hasContent := itemMap["content"]
		require.True(t, hasContent, "StructuredContent.standards[%d] should have 'content' field", i)
		contentStr, ok := contentAny.(string)
		require.True(t, ok, "StructuredContent.standards[%d].content should be string, got %T", i, contentAny)

		response.Standards = append(response.Standards, GetStandardsItem{
			Name:        nameStr,
			Description: descStr,
			Content:     contentStr,
		})
	}

	return response
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
		if strings.HasPrefix(plainText, prompt.FollowStandardsPrompt()) {
			expectedEmpty := prompt.FollowStandardsPrompt() + "\n\nNo standards found."
			require.Equal(t, expectedEmpty, plainText, "Empty result should return header + 'No standards found.'")
			return
		}

		require.Equal(t, "No standards found.", plainText, "Empty result should return 'No standards found.'")
		return
	}

	countFn := countListStandardsEntries
	if strings.Contains(plainText, "Loaded standards:") {
		countFn = countGetStandardsEntries
	}

	require.Equal(t, expectedCount, countFn(plainText),
		"Plain text should contain exactly %d standards", expectedCount)
}

// countGetStandardsEntries counts standards in get_standards format (lines starting with "- " after "Loaded standards:").
func countGetStandardsEntries(plainText string) int {
	lines := strings.Split(plainText, "\n")
	start := -1
	for i, line := range lines {
		if strings.HasPrefix(line, "Loaded standards:") {
			start = i + 1
			break
		}
	}
	if start == -1 {
		return 0
	}

	count := 0
	for _, line := range lines[start:] {
		if strings.HasPrefix(line, "Note:") {
			break
		}
		if strings.HasPrefix(line, "- ") {
			count++
		}
	}

	return count
}

// countListStandardsEntries counts standards in list_standards format (lines with ":" not starting with "#").
func countListStandardsEntries(plainText string) int {
	lines := strings.Split(plainText, "\n")
	count := 0

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && strings.Contains(line, ":") && !strings.HasPrefix(line, "#") {
			count++
		}
	}

	return count
}

// AssertStandardContainsDescription validates that a standard in plain text contains expected description
func AssertStandardContainsDescription(t *testing.T, plainText string, standardName, expectedDescription string) {
	// For list_standards, look for "name: description" pattern
	expectedPattern := standardName + ": " + expectedDescription
	require.Contains(t, plainText, expectedPattern,
		"Plain text should contain standard '%s' with description '%s'", standardName, expectedDescription)
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

// New assertions for get_standards Content summary format (memory 21 contract).
// These assertions enforce the concise summary output (no full standard bodies in Content).

// AssertGetStandardsContentHasHeader verifies Content starts with the follow-standards header.
func AssertGetStandardsContentHasHeader(t *testing.T, plainText string) {
	const expectedHeader = "# MUST FOLLOW STANDARDS BELOW"
	require.True(t, strings.HasPrefix(plainText, expectedHeader),
		"get_standards Content should start with header '%s', got: %q", expectedHeader, plainText[:min(len(plainText), 100)])
}

// AssertGetStandardsContentHasInstructionLine verifies Content contains the follow instruction.
func AssertGetStandardsContentHasInstructionLine(t *testing.T, plainText string) {
	require.Contains(t, plainText, "You MUST follow the loaded standards.",
		"get_standards Content should contain instruction line")
}

// AssertGetStandardsContentHasCount verifies Content contains "Loaded standards: N" line.
func AssertGetStandardsContentHasCount(t *testing.T, plainText string, expectedCount int) {
	expectedLine := fmt.Sprintf("Loaded standards: %d", expectedCount)
	require.Contains(t, plainText, expectedLine,
		"get_standards Content should contain count line '%s'", expectedLine)
}

// AssertGetStandardsContentListsStandard verifies Content contains a summary line for the standard.
// Format: "- name: description".
func AssertGetStandardsContentListsStandard(t *testing.T, plainText, name, description string) {
	var expectedLine string
	if description == "" {
		expectedLine = "- " + name + ":"
	} else {
		expectedLine = "- " + name + ": " + description
	}
	require.Contains(t, plainText, expectedLine,
		"get_standards Content should contain summary line for standard '%s'", name)
}

// AssertGetStandardsContentNoBody verifies Content does NOT contain markdown code blocks
// or known body content snippets. This enforces the new contract where full standard bodies
// are only in StructuredContent, not in plain-text Content.
func AssertGetStandardsContentNoBody(t *testing.T, plainText string, knownBodySnippet string) {
	require.NotContains(t, plainText, "```md",
		"get_standards Content should NOT contain markdown code blocks (full bodies belong in StructuredContent)")
	require.NotContains(t, plainText, "```\n",
		"get_standards Content should NOT contain code block endings (full bodies belong in StructuredContent)")
	if knownBodySnippet != "" {
		require.NotContains(t, plainText, knownBodySnippet,
			"get_standards Content should NOT contain standard body content '%s'", knownBodySnippet)
	}
}

// AssertGetStandardsContentEmpty verifies Content for empty result (N==0).
// Format per memory 21: "# MUST FOLLOW STANDARDS BELOW\n\nNo standards found."
func AssertGetStandardsContentEmpty(t *testing.T, plainText string) {
	const expectedEmpty = "# MUST FOLLOW STANDARDS BELOW\n\nNo standards found."
	require.Equal(t, expectedEmpty, plainText,
		"get_standards Content for empty result should be exactly: %q", expectedEmpty)
}

// AssertGetStandardsContentNote verifies Content contains the note about StructuredContent.
func AssertGetStandardsContentNote(t *testing.T, plainText string) {
	require.Contains(t, plainText, "Note: Full standard bodies are provided in StructuredContent",
		"get_standards Content should contain note about StructuredContent")
}

// getContext returns a background context for tool calls
func getContext() context.Context {
	return context.Background()
}
