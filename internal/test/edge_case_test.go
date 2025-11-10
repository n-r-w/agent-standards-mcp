package test

import (
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
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

	// Should return empty result for non-existent standard
	plainText := AssertPlainTextInput(t, result)
	require.Equal(t, "No standards found.", plainText, "Should return 'No standards found.' for non-existent standard")
}

// TestGetStandards_MixOfExistentAndNonExistent tests getting a mix of existent and non-existent standards
func TestGetStandards_MixOfExistentAndNonExistent(t *testing.T) {
	suite := NewTestSuite(t, WithCustomStandardFiles(DefaultStandardFiles()))
	defer suite.Cleanup()

	// Test get_standards with a mix of existent and non-existent standards
	result := AssertToolCallSuccess(t, suite, "get_standards", map[string]any{
		"standard_names": []string{"standard1", "nonexistent1", "standard2", "nonexistent2"},
	})

	// Should return only the existent standards
	plainText := AssertPlainTextInput(t, result)
	AssertStandardListContains(t, plainText, "standard1")
	AssertStandardListContains(t, plainText, "standard2")
	AssertStandardListCount(t, plainText, 2)
	AssertMultipleStandardsFormat(t, plainText)
}

// TestGetStandards_NoFrontmatter tests getting a standard with no frontmatter
func TestGetStandards_NoFrontmatter(t *testing.T) {
	suite := NewTestSuite(t, WithCustomStandardFiles(DefaultStandardFiles()))
	defer suite.Cleanup()

	// Test get_standards for standard with no frontmatter
	result := AssertToolCallSuccess(t, suite, "get_standards", map[string]any{
		"standard_names": []string{"no-description"},
	})

	// Verify that standard with no frontmatter is handled correctly
	plainText := AssertPlainTextInput(t, result)
	AssertStandardListContains(t, plainText, "no-description")
	AssertStandardContainsDescription(t, plainText, "no-description", "")
	AssertStandardContainsContent(t, plainText, "no-description", "This standard has no frontmatter description")
}

// TestGetStandards_EmptyStandardList tests calling get_standards with an empty standard list
func TestGetStandards_EmptyStandardList(t *testing.T) {
	suite := NewTestSuite(t, WithCustomStandardFiles(DefaultStandardFiles()))
	defer suite.Cleanup()

	// Test get_standards with empty standard list
	result := AssertToolCallSuccess(t, suite, "get_standards", map[string]any{
		"standard_names": []string{},
	})

	// Should return empty result
	plainText := AssertPlainTextInput(t, result)
	require.Equal(t, "No standards found.", plainText, "Should return 'No standards found.' for empty standard list")
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

	// Should return the standard for each occurrence
	plainText := AssertPlainTextInput(t, result)
	// Count occurrences of the standard name in markdown format
	standardCount := requireCountOfSubstring(t, plainText, "## standard1:")
	require.Equal(t, 3, standardCount, "Should return standard for each occurrence")
	AssertMultipleStandardsFormat(t, plainText)
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
	count := 0
	start := 0
	for {
		index := requireIndexOfSubstring(text[start:], substring)
		if index == -1 {
			break
		}
		count++
		start += index + len(substring)
	}
	return count
}

// Helper function to find substring index
func requireIndexOfSubstring(text, substring string) int {
	for i := 0; i <= len(text)-len(substring); i++ {
		if text[i:i+len(substring)] == substring {
			return i
		}
	}
	return -1
}
