package test

import (
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/n-r-w/agent-standards-mcp/internal/prompt"
	"github.com/stretchr/testify/require"
)

// TestGetStandards_NonExistentStandard tests getting a standard that doesn't exist
func TestGetStandards_NonExistentStandard(t *testing.T) {
	suite := NewTestSuite(t, WithCustomStandardFiles(NonExistentStandardFile()))
	defer suite.Cleanup()

	// Test get_standards for a non-existent standard
	result := AssertToolCallSuccess(t, suite, "get_standards", map[string]any{
		"standard_names": []string{"nonexistent-standard"},
	})

	// Should return empty result with new format (memory 21)
	plainText := AssertPlainTextInput(t, result)
	AssertGetStandardsContentEmpty(t, plainText)
}

// TestGetStandards_MixOfExistentAndNonExistent tests getting a mix of existent and non-existent standards
func TestGetStandards_MixOfExistentAndNonExistent(t *testing.T) {
	suite := NewTestSuite(t, WithCustomStandardFiles(DefaultStandardFiles()))
	defer suite.Cleanup()

	// Test get_standards with a mix of existent and non-existent standards
	result := AssertToolCallSuccess(t, suite, "get_standards", map[string]any{
		"standard_names": []string{"standard1", "nonexistent1", "standard2", "nonexistent2"},
	})

	// Should return only the existent standards (memory 21)
	plainText := AssertPlainTextInput(t, result)
	AssertGetStandardsContentHasHeader(t, plainText)
	AssertGetStandardsContentHasInstructionLine(t, plainText)
	AssertGetStandardsContentHasCount(t, plainText, 2)
	AssertGetStandardsContentListsStandard(t, plainText, "standard1", "A test standard for basic functionality")
	AssertGetStandardsContentListsStandard(t, plainText, "standard2", "Another test standard with different content")
	AssertGetStandardsContentNoBody(t, plainText, "This is the content of standard1")
	AssertGetStandardsContentNoBody(t, plainText, "Standard 2 content here.")

	// StructuredContent should have 2 entries
	response := AssertGetStandardsStructuredContent(t, result)
	require.Len(t, response.Standards, 2)
}

// TestGetStandards_NoFrontmatter tests getting a standard with no frontmatter
func TestGetStandards_NoFrontmatter(t *testing.T) {
	suite := NewTestSuite(t, WithCustomStandardFiles(DefaultStandardFiles()))
	defer suite.Cleanup()

	// Test get_standards for standard with no frontmatter
	result := AssertToolCallSuccess(t, suite, "get_standards", map[string]any{
		"standard_names": []string{"no-description"},
	})

	// Verify that standard with no frontmatter is handled correctly (memory 21)
	plainText := AssertPlainTextInput(t, result)
	AssertGetStandardsContentHasHeader(t, plainText)
	AssertGetStandardsContentHasInstructionLine(t, plainText)
	AssertGetStandardsContentHasCount(t, plainText, 1)
	// Empty description should render as "- no-description:" (no trailing space after colon)
	AssertGetStandardsContentListsStandard(t, plainText, "no-description", "")
	AssertGetStandardsContentNoBody(t, plainText, "This standard has no frontmatter description")

	// StructuredContent should still contain the full body
	response := AssertGetStandardsStructuredContent(t, result)
	require.Len(t, response.Standards, 1)
	require.Contains(t, response.Standards[0].Content, "This standard has no frontmatter description")
}

// TestGetStandards_EmptyStandardList tests calling get_standards with an empty standard list
func TestGetStandards_EmptyStandardList(t *testing.T) {
	suite := NewTestSuite(t, WithCustomStandardFiles(DefaultStandardFiles()))
	defer suite.Cleanup()

	// Test get_standards with empty standard list
	result := AssertToolCallSuccess(t, suite, "get_standards", map[string]any{
		"standard_names": []string{},
	})

	// Should return empty result with new format (memory 21)
	plainText := AssertPlainTextInput(t, result)
	AssertGetStandardsContentEmpty(t, plainText)
}

// TestListStandards_ComplexStandard tests that complex standards with formatting are handled correctly
func TestListStandards_ComplexStandard(t *testing.T) {
	suite := NewTestSuite(t, WithCustomStandardFiles(DefaultStandardFiles()))
	defer suite.Cleanup()

	// Test list_standards with complex standard
	result := AssertToolCallSuccess(t, suite, "list_standards", map[string]any{})

	// Verify the result structure
	plainText := AssertPlainTextInput(t, result)
	AssertStandardListContains(t, plainText, "complex-standard")
	AssertStandardContainsDescription(t, plainText, "complex-standard", "A more complex standard with advanced features")
	AssertStandardListCount(t, plainText, 5) // All 5 default standards
	AssertMultipleStandardsFormat(t, plainText)
}

// TestGetStandards_DuplicateStandardNames tests requesting the same standard multiple times
func TestGetStandards_DuplicateStandardNames(t *testing.T) {
	suite := NewTestSuite(t, WithCustomStandardFiles(DefaultStandardFiles()))
	defer suite.Cleanup()

	// Test get_standards with duplicate standard names
	result := AssertToolCallSuccess(t, suite, "get_standards", map[string]any{
		"standard_names": []string{"standard1", "standard1", "standard1"},
	})

	// Should return the standard for each occurrence (memory 21: duplicates preserved)
	plainText := AssertPlainTextInput(t, result)
	AssertGetStandardsContentHasHeader(t, plainText)
	AssertGetStandardsContentHasInstructionLine(t, plainText)
	AssertGetStandardsContentHasCount(t, plainText, 3)
	// Count occurrences of the summary line in new format
	standardCount := requireCountOfSubstring(t, plainText, "- standard1:")
	require.Equal(t, 3, standardCount, "Should return summary line for each occurrence")
	// No markdown code blocks in Content
	AssertGetStandardsContentNoBody(t, plainText, "This is the content of standard1")

	// StructuredContent should have 3 entries
	response := AssertGetStandardsStructuredContent(t, result)
	require.Len(t, response.Standards, 3, "StructuredContent should have 3 standard entries")
}

// TestGetStandards_ParameterValidationMissing tests that missing required parameter is caught
func TestGetStandards_ParameterValidationMissing(t *testing.T) {
	suite := NewTestSuite(t, WithCustomStandardFiles(DefaultStandardFiles()))
	defer suite.Cleanup()

	// Test get_standards with missing required parameter
	// This should fail at MCP SDK level due to JSON schema validation
	ctx := getContext()
	_, err := suite.ClientSession.CallTool(ctx, &mcp.CallToolParams{
		Meta:      mcp.Meta{},
		Name:      "get_standards",
		Arguments: map[string]any{},
	})

	// This should fail at MCP SDK level due to JSON schema validation
	require.Error(t, err, "Tool call should fail due to missing required parameter")
	require.Contains(t, err.Error(), "required",
		"Error should indicate missing required parameter")
}

// TestGetStandards_ParameterValidationWrongType tests that wrong parameter type is caught
func TestGetStandards_ParameterValidationWrongType(t *testing.T) {
	suite := NewTestSuite(t, WithCustomStandardFiles(DefaultStandardFiles()))
	defer suite.Cleanup()

	// Test get_standards with wrong parameter type
	// This should fail at MCP SDK level due to JSON schema validation
	ctx := getContext()
	_, err := suite.ClientSession.CallTool(ctx, &mcp.CallToolParams{
		Meta: mcp.Meta{},
		Name: "get_standards",
		Arguments: map[string]any{
			"standard_names": "should-be-array", // Should be array, not string
		},
	})

	// This should fail at MCP SDK level due to JSON schema validation
	require.Error(t, err, "Tool call should fail due to wrong parameter type")
	require.Contains(t, err.Error(), "type",
		"Error should indicate wrong parameter type")
}

// Helper function to count substring occurrences
func requireCountOfSubstring(t *testing.T, text, substring string) int {
	require.NotEmpty(t, substring, "substring must not be empty")
	return strings.Count(text, substring)
}

// Helper function to find substring index
// Tests for get_standards StructuredContent JSON contract validation

func TestGetStandards_StructuredContent_Success(t *testing.T) {
	suite := NewTestSuite(t, WithCustomStandardFiles(DefaultStandardFiles()))
	defer suite.Cleanup()

	// Test get_standards with multiple standards
	result := AssertToolCallSuccess(t, suite, "get_standards", map[string]any{
		"standard_names": []string{"standard1", "standard2"},
	})

	// Validate StructuredContent follows the JSON contract (unchanged)
	response := AssertGetStandardsStructuredContent(t, result)
	require.Len(t, response.Standards, 2, "Should return 2 standards")

	// Verify first standard has all required fields populated
	require.Equal(t, "standard1", response.Standards[0].Name)
	require.NotEmpty(t, response.Standards[0].Description)
	require.NotEmpty(t, response.Standards[0].Content)

	// Verify second standard has all required fields populated
	require.Equal(t, "standard2", response.Standards[1].Name)
	require.NotEmpty(t, response.Standards[1].Description)
	require.NotEmpty(t, response.Standards[1].Content)

	// Plain text Content assertions per new contract (memory 21)
	plainText := AssertPlainTextInput(t, result)
	AssertGetStandardsContentHasHeader(t, plainText)
	AssertGetStandardsContentListsStandard(t, plainText, "standard1", "A test standard for basic functionality")
	AssertGetStandardsContentListsStandard(t, plainText, "standard2", "Another test standard with different content")
	AssertGetStandardsContentNoBody(t, plainText, "This is the content of standard1")
}

func TestGetStandards_StructuredContent_Empty(t *testing.T) {
	suite := NewTestSuite(t, WithCustomStandardFiles(DefaultStandardFiles()))
	defer suite.Cleanup()

	// Test get_standards with empty standard list
	result := AssertToolCallSuccess(t, suite, "get_standards", map[string]any{
		"standard_names": []string{},
	})

	// Validate StructuredContent is {"standards": []} for empty result
	response := AssertGetStandardsStructuredContent(t, result)
	require.Empty(t, response.Standards, "Empty result should have empty standards array")

	// Plain text Content should return header + "No standards found." per memory 21 contract
	plainText := AssertPlainTextInput(t, result)
	expectedEmpty := prompt.FollowStandardsPrompt() + "\n\nNo standards found."
	require.Equal(t, expectedEmpty, plainText)
}

func TestGetStandards_StructuredContent_MixOfExistentAndNonExistent(t *testing.T) {
	suite := NewTestSuite(t, WithCustomStandardFiles(DefaultStandardFiles()))
	defer suite.Cleanup()

	// Test get_standards with a mix of existent and non-existent standards
	result := AssertToolCallSuccess(t, suite, "get_standards", map[string]any{
		"standard_names": []string{"standard1", "nonexistent1", "standard2", "nonexistent2"},
	})

	// Validate StructuredContent structure and length (only existent standards)
	response := AssertGetStandardsStructuredContent(t, result)
	require.Len(t, response.Standards, 2, "Should return only 2 existent standards")

	// Each item should have all required fields as strings
	for i, std := range response.Standards {
		require.NotEmpty(t, std.Name, "standards[%d].name should not be empty", i)
		require.IsType(t, "", std.Description, "standards[%d].description should be string", i)
		require.IsType(t, "", std.Content, "standards[%d].content should be string", i)
	}
}

func TestGetStandards_StructuredContent_DuplicateNames(t *testing.T) {
	suite := NewTestSuite(t, WithCustomStandardFiles(DefaultStandardFiles()))
	defer suite.Cleanup()

	// Test get_standards with duplicate standard names
	result := AssertToolCallSuccess(t, suite, "get_standards", map[string]any{
		"standard_names": []string{"standard1", "standard1", "standard1"},
	})

	// Validate StructuredContent structure (duplicates should be preserved)
	response := AssertGetStandardsStructuredContent(t, result)
	require.GreaterOrEqual(t, len(response.Standards), 2, "Duplicate names should result in multiple items")

	// Each item should have all required fields with proper types
	for i, std := range response.Standards {
		require.NotEmpty(t, std.Name, "standards[%d].name should not be empty", i)
		require.Equal(t, "standard1", std.Name, "All items should be standard1")
		require.IsType(t, "", std.Description, "standards[%d].description should be string", i)
		require.IsType(t, "", std.Content, "standards[%d].content should be string", i)
	}
}
