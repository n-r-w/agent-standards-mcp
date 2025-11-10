package test

import (
	"os"
	"path/filepath"
	"testing"
)

// TestUtility_WorkingDirectory verifies that the working directory is set correctly
// This is useful for debugging CI/CD issues and ensuring paths work correctly
func TestUtility_WorkingDirectory(t *testing.T) {
	// Get current working directory
	wd, _ := os.Getwd()
	t.Log("Working directory:", wd)

	// Get absolute path of current directory
	absPath, _ := filepath.Abs(".")
	t.Log("Absolute path of '.':", absPath)

	// Get absolute path of the command directory
	cmdPath, _ := filepath.Abs("./cmd/agent-standards-mcp")
	t.Log("Absolute path of './cmd/agent-standards-mcp':", cmdPath)
}
