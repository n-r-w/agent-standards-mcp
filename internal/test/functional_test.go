package test

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestListStandards_NoArgs tests list_standards tool with no arguments
func TestListStandards_NoArgs(t *testing.T) {
	suite := NewTestSuite(t, WithCustomStandardFiles(DefaultStandardFiles()))
	defer suite.Cleanup()

	// Test list_standards with no arguments
	result := AssertToolCallSuccess(t, suite, "list_standards", map[string]any{})

	// Verify the result structure - plain text Content
	plainText := AssertPlainTextInput(t, result)
	AssertStandardListCount(t, plainText, 5)
	AssertMultipleStandardsFormat(t, plainText)

	// Verify all standards are present in plain text
	AssertStandardListContains(t, plainText, "standard1")
	AssertStandardListContains(t, plainText, "standard2")
	AssertStandardListContains(t, plainText, "standard3")
	AssertStandardListContains(t, plainText, "no-description")
	AssertStandardListContains(t, plainText, "complex-standard")

	// Validate StructuredContent JSON contract
	// Expected: {"standards": [{"name": "...", "description": "..."}, ...]}
	structuredResponse := AssertListStandardsStructuredContent(t, result)
	require.Len(t, structuredResponse.Standards, 5, "StructuredContent should contain 5 standards")
}

// TestListStandards_EmptyStandardsDir tests list_standards when standards directory is empty
func TestListStandards_EmptyStandardsDir(t *testing.T) {
	suite := NewTestSuite(t, WithCustomStandardFiles(EmptyStandardFiles()))
	defer suite.Cleanup()

	// Test list_standards with empty standards directory
	result := AssertToolCallSuccess(t, suite, "list_standards", map[string]any{})

	// Verify the plain text result
	plainText := AssertPlainTextInput(t, result)
	require.Equal(t, "No standards found.", plainText, "Should return 'No standards found.' for empty directory")

	// Verify the StructuredContent is {"standards": []}
	response := AssertListStandardsStructuredContent(t, result)
	require.Empty(t, response.Standards, "StructuredContent.standards should be empty for empty directory")
}

// TestGetStandards_SingleStandard tests getting a single standard
func TestGetStandards_SingleStandard(t *testing.T) {
	suite := NewTestSuite(t, WithCustomStandardFiles(DefaultStandardFiles()))
	defer suite.Cleanup()

	// Test get_standards for a single standard
	result := AssertToolCallSuccess(t, suite, "get_standards", map[string]any{
		"standard_names": []string{"standard1"},
	})

	// Verify the result structure using new contract assertions (memory 21)
	plainText := AssertPlainTextInput(t, result)

	// NEW CONTRACT: Content is concise summary, no full bodies
	AssertGetStandardsContentHasHeader(t, plainText)
	AssertGetStandardsContentHasInstructionLine(t, plainText)
	AssertGetStandardsContentHasCount(t, plainText, 1)
	AssertGetStandardsContentListsStandard(t, plainText, "standard1", "A test standard for basic functionality")
	AssertGetStandardsContentNote(t, plainText)
	AssertGetStandardsContentNoBody(t, plainText, "This is the content of standard1")

	// StructuredContent should still contain the full body
	response := AssertGetStandardsStructuredContent(t, result)
	require.Len(t, response.Standards, 1)
	require.Contains(t, response.Standards[0].Content, "This is the content of standard1")
}

// TestGetStandards_MultipleStandards tests getting multiple standards
func TestGetStandards_MultipleStandards(t *testing.T) {
	suite := NewTestSuite(t, WithCustomStandardFiles(DefaultStandardFiles()))
	defer suite.Cleanup()

	// Test get_standards for specific standards
	result := AssertToolCallSuccess(t, suite, "get_standards", map[string]any{
		"standard_names": []string{"standard1", "standard2"},
	})

	// Verify the result structure using new contract assertions (memory 21)
	plainText := AssertPlainTextInput(t, result)

	// NEW CONTRACT: Content is concise summary, no full bodies
	AssertGetStandardsContentHasHeader(t, plainText)
	AssertGetStandardsContentHasInstructionLine(t, plainText)
	AssertGetStandardsContentHasCount(t, plainText, 2)
	AssertGetStandardsContentListsStandard(t, plainText, "standard1", "A test standard for basic functionality")
	AssertGetStandardsContentListsStandard(t, plainText, "standard2", "Another test standard with different content")
	AssertGetStandardsContentNote(t, plainText)
	AssertGetStandardsContentNoBody(t, plainText, "This is the content of standard1")
	AssertGetStandardsContentNoBody(t, plainText, "Standard 2 content here.")

	// StructuredContent should still contain the full bodies
	response := AssertGetStandardsStructuredContent(t, result)
	require.Len(t, response.Standards, 2)
	require.Contains(t, response.Standards[0].Content, "This is the content of standard1")
	require.Contains(t, response.Standards[1].Content, "Standard 2 content here.")
}

// TestGetStandards_AllStandards tests getting all standards
func TestGetStandards_AllStandards(t *testing.T) {
	suite := NewTestSuite(t, WithCustomStandardFiles(DefaultStandardFiles()))
	defer suite.Cleanup()

	// Test get_standards for all standards
	result := AssertToolCallSuccess(t, suite, "get_standards", map[string]any{
		"standard_names": []string{"standard1", "standard2", "standard3", "no-description", "complex-standard"},
	})

	// Verify the result structure
	plainText := AssertPlainTextInput(t, result)
	AssertStandardListCount(t, plainText, 5)
	AssertMultipleStandardsFormat(t, plainText)

	// Verify each standard is present
	AssertStandardListContains(t, plainText, "standard1")
	AssertStandardListContains(t, plainText, "standard2")
	AssertStandardListContains(t, plainText, "standard3")
	AssertStandardListContains(t, plainText, "no-description")
	AssertStandardListContains(t, plainText, "complex-standard")
}

// TestGetStandards_CustomStandards tests with custom standard files
func TestGetStandards_CustomStandards(t *testing.T) {
	customStandards := CustomStandardFiles()

	suite := NewTestSuite(t, WithCustomStandardFiles(customStandards))
	defer suite.Cleanup()

	// Test get_standards for custom standards
	result := AssertToolCallSuccess(t, suite, "get_standards", map[string]any{
		"standard_names": []string{"custom1", "custom2"},
	})

	// Verify the result structure using new contract assertions (memory 21)
	plainText := AssertPlainTextInput(t, result)

	// NEW CONTRACT: Content is concise summary, no full bodies
	AssertGetStandardsContentHasHeader(t, plainText)
	AssertGetStandardsContentHasInstructionLine(t, plainText)
	AssertGetStandardsContentHasCount(t, plainText, 2)
	AssertGetStandardsContentListsStandard(t, plainText, "custom1", "Custom standard 1")
	AssertGetStandardsContentListsStandard(t, plainText, "custom2", "Custom standard 2")
	AssertGetStandardsContentNote(t, plainText)
	AssertGetStandardsContentNoBody(t, plainText, "Custom content 1")
	AssertGetStandardsContentNoBody(t, plainText, "Custom content 2")

	// StructuredContent should still contain the full bodies
	response := AssertGetStandardsStructuredContent(t, result)
	require.Len(t, response.Standards, 2)
	require.Contains(t, response.Standards[0].Content, "Custom content 1")
	require.Contains(t, response.Standards[1].Content, "Custom content 2")
}
