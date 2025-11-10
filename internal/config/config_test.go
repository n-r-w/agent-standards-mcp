package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad_DefaultValues(t *testing.T) {
	// Clear environment variables
	clearEnvVars()

	cfg, err := Load()
	require.NoError(t, err)
	require.NotNil(t, cfg)

	assert.Equal(t, "ERROR", cfg.LogLevel)
	assert.Contains(t, cfg.Folder, "agent-standards")
	assert.Equal(t, 100, cfg.MaxStandards)
	assert.Equal(t, 10240, cfg.MaxStandardSize)
}

func TestLoad_EnvironmentVariables(t *testing.T) {
	// Clear environment variables
	clearEnvVars()

	// Set custom environment variables
	t.Setenv("AGENT_STANDARDS_MCP_LOG_LEVEL", "DEBUG")
	t.Setenv("AGENT_STANDARDS_MCP_FOLDER", "/tmp/custom-standards")
	t.Setenv("AGENT_STANDARDS_MCP_MAX_STANDARDS", "200")
	t.Setenv("AGENT_STANDARDS_MCP_MAX_STANDARD_SIZE", "20480")

	cfg, err := Load()
	require.NoError(t, err)
	require.NotNil(t, cfg)

	assert.Equal(t, "DEBUG", cfg.LogLevel)
	assert.Equal(t, "/tmp/custom-standards", cfg.Folder)
	assert.Equal(t, 200, cfg.MaxStandards)
	assert.Equal(t, 20480, cfg.MaxStandardSize)
}

func TestConfig_ValidateLogLevel(t *testing.T) {
	tests := []struct {
		name        string
		logLevel    string
		expectError bool
	}{
		{"Valid NONE", "NONE", false},
		{"Valid DEBUG", "DEBUG", false},
		{"Valid INFO", "INFO", false},
		{"Valid WARN", "WARN", false},
		{"Valid ERROR", "ERROR", false},
		{"Valid lowercase", "debug", false},
		{"Invalid empty", "", true},
		{"Invalid level", "INVALID", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				LogLevel:        tt.logLevel,
				Folder:          "/tmp",
				MaxStandards:    100,
				MaxStandardSize: 10240,
			}
			err := cfg.validateLogLevel()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestConfig_ValidateFolder(t *testing.T) {
	tests := []struct {
		name        string
		folder      string
		expectError bool
	}{
		{"Valid existing dir", t.TempDir(), false},
		{"Valid creatable dir", filepath.Join(t.TempDir(), "newdir"), false},
		{"Invalid empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				LogLevel:        "ERROR",
				Folder:          tt.folder,
				MaxStandards:    100,
				MaxStandardSize: 10240,
			}
			err := cfg.validateFolder()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestConfig_ValidateLimits(t *testing.T) {
	tests := []struct {
		name            string
		maxStandards    int
		maxStandardSize int
		expectError     bool
	}{
		{"Valid limits", 100, 10240, false},
		{"Zero max standards", 0, 10240, true},
		{"Negative max standards", -1, 10240, true},
		{"Zero max standard size", 100, 0, true},
		{"Negative max standard size", 100, -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				LogLevel:        "ERROR",
				Folder:          "/tmp",
				MaxStandards:    tt.maxStandards,
				MaxStandardSize: tt.maxStandardSize,
			}
			err := cfg.validateLimits()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestConfig_IsLoggingEnabled(t *testing.T) {
	tests := []struct {
		name     string
		logLevel string
		expected bool
	}{
		{"NONE disabled", "NONE", false},
		{"DEBUG enabled", "DEBUG", true},
		{"INFO enabled", "INFO", true},
		{"WARN enabled", "WARN", true},
		{"ERROR enabled", "ERROR", true},
		{"lowercase enabled", "debug", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				LogLevel:        tt.logLevel,
				Folder:          "/tmp",
				MaxStandards:    100,
				MaxStandardSize: 10240,
			}
			assert.Equal(t, tt.expected, cfg.IsLoggingEnabled())
		})
	}
}

func TestConfig_GetLogLevel(t *testing.T) {
	tests := []struct {
		name     string
		logLevel string
		expected LogLevel
	}{
		{"NONE", "NONE", LogLevelNone},
		{"DEBUG", "DEBUG", LogLevelDebug},
		{"INFO", "INFO", LogLevelInfo},
		{"WARN", "WARN", LogLevelWarn},
		{"ERROR", "ERROR", LogLevelError},
		{"lowercase", "debug", LogLevelDebug},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				LogLevel:        tt.logLevel,
				Folder:          "/tmp",
				MaxStandards:    100,
				MaxStandardSize: 10240,
			}
			assert.Equal(t, tt.expected, cfg.GetLogLevel())
		})
	}
}

func TestConfig_Getters(t *testing.T) {
	cfg := &Config{
		LogLevel:        "ERROR",
		Folder:          "/test/folder",
		MaxStandards:    150,
		MaxStandardSize: 15360,
	}

	assert.Equal(t, "/test/folder", cfg.GetFolder())
	assert.Equal(t, 150, cfg.GetMaxStandards())
	assert.Equal(t, 15360, cfg.GetMaxStandardSize())
}

// clearEnvVars clears all relevant environment variables for testing
func clearEnvVars() {
	envVars := []string{
		"AGENT_STANDARDS_MCP_LOG_LEVEL",
		"AGENT_STANDARDS_MCP_FOLDER",
		"AGENT_STANDARDS_MCP_MAX_STANDARDS",
		"AGENT_STANDARDS_MCP_MAX_STANDARD_SIZE",
	}

	for _, envVar := range envVars {
		if err := os.Unsetenv(envVar); err != nil {
			// Log the error but don't fail the test
			// This is a best-effort cleanup operation
		}
	}
}
