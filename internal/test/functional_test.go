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

	// Verify the result structure
	plainText := AssertPlainTextInput(t, result)
	AssertStandardListCount(t, plainText, 5)
	AssertMultipleStandardsFormat(t, plainText)

	// Verify all standards are present
	AssertStandardListContains(t, plainText, "standard1")
	AssertStandardListContains(t, plainText, "standard2")
	AssertStandardListContains(t, plainText, "standard3")
	AssertStandardListContains(t, plainText, "no-description")
	AssertStandardListContains(t, plainText, "complex-standard")
}

// TestListStandards_EmptyStandardsDir tests list_standards when standards directory is empty
func TestListStandards_EmptyStandardsDir(t *testing.T) {
	suite := NewTestSuite(t, WithCustomStandardFiles(EmptyStandardFiles()))
	defer suite.Cleanup()

	// Test list_standards with empty standards directory
	result := AssertToolCallSuccess(t, suite, "list_standards", map[string]any{})

	// Verify the result structure
	plainText := AssertPlainTextInput(t, result)
	require.Equal(t, "No standards found.", plainText, "Should return 'No standards found.' for empty directory")
}

// TestGetStandards_SingleStandard tests getting a single standard
func TestGetStandards_SingleStandard(t *testing.T) {
	suite := NewTestSuite(t, WithCustomStandardFiles(DefaultStandardFiles()))
	defer suite.Cleanup()

	// Test get_standards for a single standard
	result := AssertToolCallSuccess(t, suite, "get_standards", map[string]any{
		"standard_names": []string{"standard1"},
	})

	// Verify the result structure
	plainText := AssertPlainTextInput(t, result)
	AssertGetStandardsContainsContent(t, plainText, "standard1", "A test standard for basic functionality", "This is the content of standard1")
}

// TestGetStandards_MultipleStandards tests getting multiple standards
func TestGetStandards_MultipleStandards(t *testing.T) {
	suite := NewTestSuite(t, WithCustomStandardFiles(DefaultStandardFiles()))
	defer suite.Cleanup()

	// Test get_standards for specific standards
	result := AssertToolCallSuccess(t, suite, "get_standards", map[string]any{
		"standard_names": []string{"standard1", "standard2"},
	})

	// Verify the result structure
	plainText := AssertPlainTextInput(t, result)
	AssertStandardListCount(t, plainText, 2)
	AssertMultipleStandardsFormat(t, plainText)

	// Verify we got the expected standards
	AssertStandardListContains(t, plainText, "standard1")
	AssertStandardListContains(t, plainText, "standard2")
	AssertGetStandardsContainsContent(t, plainText, "standard1", "A test standard for basic functionality", "This is the content of standard1")
	AssertGetStandardsContainsContent(t, plainText, "standard2", "Another test standard with different content", "Standard 2 content here.")
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

	// Verify the result structure
	plainText := AssertPlainTextInput(t, result)
	AssertStandardListCount(t, plainText, 2)
	AssertMultipleStandardsFormat(t, plainText)

	// Verify custom standard content
	AssertGetStandardsContainsContent(t, plainText, "custom1", "Custom standard 1", "Custom content 1")
	AssertGetStandardsContainsContent(t, plainText, "custom2", "Custom standard 2", "Custom content 2")
}
