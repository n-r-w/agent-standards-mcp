package test

import (
	"context"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/require"
)

// TestTransport_InMemory tests that tools work correctly with in-memory transport
func TestTransport_InMemory(t *testing.T) {
	suite := NewTestSuite(t, WithCustomStandardFiles(DefaultStandardFiles()))
	defer suite.Cleanup()

	// Verify that expected tools are available
	expectedTools := []string{"list_standards", "get_standards"}
	AssertToolsAvailable(t, suite, expectedTools)

	// Test basic list_standards functionality
	result := AssertToolCallSuccess(t, suite, "list_standards", map[string]any{})

	// Verify the result structure
	plainText := AssertPlainTextInput(t, result)
	AssertStandardListCount(t, plainText, 5)
	AssertMultipleStandardsFormat(t, plainText)

	// Test get_standards functionality
	result = AssertToolCallSuccess(t, suite, "get_standards", map[string]any{
		"standard_names": []string{"standard1"},
	})

	// Verify the result structure
	plainText = AssertPlainTextInput(t, result)
	AssertStandardListCount(t, plainText, 1)
	AssertGetStandardsContainsContent(t, plainText, "standard1", "A test standard for basic functionality", "This is the content of standard1")
}

// TestTransport_ConcurrentConnections tests multiple concurrent connections
func TestTransport_ConcurrentConnections(t *testing.T) {
	// Create multiple server instances to test concurrent handling
	suite1 := NewTestSuite(t, WithCustomStandardFiles(DefaultStandardFiles()))
	defer suite1.Cleanup()

	suite2 := NewTestSuite(t, WithCustomStandardFiles(DefaultStandardFiles()))
	defer suite2.Cleanup()

	// Both should be able to discover tools
	expectedTools := []string{"list_standards", "get_standards"}
	AssertToolsAvailable(t, suite1, expectedTools)
	AssertToolsAvailable(t, suite2, expectedTools)

	// Both should be able to call tools
	result1 := AssertToolCallSuccess(t, suite1, "list_standards", map[string]any{})
	result2 := AssertToolCallSuccess(t, suite2, "list_standards", map[string]any{})

	// Both should return valid results
	plainText1 := AssertPlainTextInput(t, result1)
	plainText2 := AssertPlainTextInput(t, result2)

	AssertStandardListCount(t, plainText1, 5)
	AssertStandardListCount(t, plainText2, 5)
	require.Len(t, plainText1, len(plainText2), "Both suites should return same number of standards")
}

// TestTransport_ServerStartupShutdown tests that server can start and stop properly
func TestTransport_ServerStartupShutdown(t *testing.T) {
	suite := NewTestSuite(t, WithCustomStandardFiles(DefaultStandardFiles()))

	// If we get here, server startup succeeded
	require.NotNil(t, suite.Server, "Server should be created")
	require.NotNil(t, suite.Client, "Client should be created")
	require.NotNil(t, suite.ClientSession, "Client session should be established")

	// The cleanup function will test proper shutdown
	suite.Cleanup()
}

// TestTransport_CustomClientInfo tests with custom client identification
func TestTransport_CustomClientInfo(t *testing.T) {
	suite := NewTestSuite(t,
		WithCustomStandardFiles(DefaultStandardFiles()),
		WithClientInfo("custom-test-client", "2.0.0"),
	)
	defer suite.Cleanup()

	// Verify that expected tools are available
	expectedTools := []string{"list_standards", "get_standards"}
	AssertToolsAvailable(t, suite, expectedTools)

	// Test that tool calls work with custom client
	result := AssertToolCallSuccess(t, suite, "list_standards", map[string]any{})
	plainText := AssertPlainTextInput(t, result)
	AssertStandardListCount(t, plainText, 5)
}

// TestTransport_EmptyStandards tests with empty standards directory
func TestTransport_EmptyStandards(t *testing.T) {
	suite := NewTestSuite(t, WithCustomStandardFiles(EmptyStandardFiles()))
	defer suite.Cleanup()

	// Verify that expected tools are available even with empty standards
	expectedTools := []string{"list_standards", "get_standards"}
	AssertToolsAvailable(t, suite, expectedTools)

	// Test list_standards returns empty result
	result := AssertToolCallSuccess(t, suite, "list_standards", map[string]any{})
	plainText := AssertPlainTextInput(t, result)
	require.Equal(t, "No standards found.", plainText, "Should return 'No standards found.' for empty directory")
}

// TestTransport_ToolDiscoveryVerifiesTools tests that tool discovery properly validates tool availability
func TestTransport_ToolDiscoveryVerifiesTools(t *testing.T) {
	suite := NewTestSuite(t, WithCustomStandardFiles(DefaultStandardFiles()))
	defer suite.Cleanup()

	ctx := context.Background()

	// Get available tools using proper MCP SDK params
	tools, err := suite.ClientSession.ListTools(ctx, &mcp.ListToolsParams{
		Meta:   mcp.Meta{},
		Cursor: "",
	})
	require.NoError(t, err, "Failed to get tools from MCP server")
	require.NotEmpty(t, tools.Tools, "No tools available from MCP server")

	// Verify tool structure
	for _, tool := range tools.Tools {
		require.NotEmpty(t, tool.Name, "Tool should have a name")
		require.NotNil(t, tool.InputSchema, "Tool should have an input schema")

		// Verify that tool is one of the expected tools
		switch tool.Name {
		case "list_standards", "get_standards":
			// Expected tools - OK
		default:
			t.Errorf("Unexpected tool found: %s", tool.Name)
		}
	}
}

// TestTransport_CommandTransport tests basic functionality with command transport (real subprocess)
func TestTransport_CommandTransport(t *testing.T) {
	suite := NewTestSuite(t,
		WithCommandTransport(),
		WithCustomStandardFiles(DefaultStandardFiles()),
	)
	defer suite.Cleanup()

	// Verify that expected tools are available
	expectedTools := []string{"list_standards", "get_standards"}
	AssertToolsAvailable(t, suite, expectedTools)

	// Test basic list_standards functionality
	result := AssertToolCallSuccess(t, suite, "list_standards", map[string]any{})

	// Verify the result structure
	plainText := AssertPlainTextInput(t, result)
	AssertStandardListCount(t, plainText, 5)
	AssertMultipleStandardsFormat(t, plainText)
}
