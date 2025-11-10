// Package config provides configuration management for the agent-standards-mcp server.
package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/caarlos0/env/v11"
)

const (
	// defaultMaxStandards is the default maximum number of standards to load.
	defaultMaxStandards = 100
	// defaultMaxStandardSize is the default maximum size of a single standard file in bytes.
	defaultMaxStandardSize = 10240
)

// Config holds the configuration for the agent-standards-mcp server.
type Config struct {
	LogLevel        string `env:"AGENT_STANDARDS_MCP_LOG_LEVEL" envDefault:"ERROR"`
	Folder          string `env:"AGENT_STANDARDS_MCP_FOLDER" envDefault:"~/agent-standards"`
	MaxStandards    int    `env:"AGENT_STANDARDS_MCP_MAX_STANDARDS" envDefault:"100"`
	MaxStandardSize int    `env:"AGENT_STANDARDS_MCP_MAX_STANDARD_SIZE" envDefault:"10240"`
}

// Load loads configuration from environment variables and validates it.
func Load() (*Config, error) {
	cfg := &Config{
		LogLevel:        "ERROR",
		Folder:          "~/agent-standards",
		MaxStandards:    defaultMaxStandards,
		MaxStandardSize: defaultMaxStandardSize,
	}

	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse environment variables: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return cfg, nil
}

// Validate performs comprehensive validation of the configuration.
func (c *Config) Validate() error {
	if err := c.validateLogLevel(); err != nil {
		return err
	}

	if err := c.validateFolder(); err != nil {
		return err
	}

	if err := c.validateLimits(); err != nil {
		return err
	}

	return nil
}

// validateLogLevel validates the log level setting.
func (c *Config) validateLogLevel() error {
	if c.LogLevel == "" {
		return errors.New("log level cannot be empty")
	}

	return validateLogLevel(c.LogLevel)
}

// validateFolder validates the standards folder path and permissions.
func (c *Config) validateFolder() error {
	if c.Folder == "" {
		return errors.New("folder path cannot be empty")
	}

	// Expand ~ to user home directory and validate
	expandedPath, err := expandPath(c.Folder)
	if err != nil {
		return fmt.Errorf("failed to expand folder path: %w", err)
	}

	// Store the expanded path back to the config
	c.Folder = expandedPath

	return validateDirectoryPath(c.Folder)
}

// validateLimits validates numeric configuration limits.
func (c *Config) validateLimits() error {
	if err := validatePositiveInt(c.MaxStandards, "MaxStandards"); err != nil {
		return err
	}

	if err := validatePositiveInt(c.MaxStandardSize, "MaxStandardSize"); err != nil {
		return err
	}

	return nil
}

// IsLoggingEnabled returns true if logging is enabled (log level is not NONE).
func (c *Config) IsLoggingEnabled() bool {
	return strings.ToUpper(c.LogLevel) != string(LogLevelNone)
}

// GetLogLevel returns the normalized log level.
func (c *Config) GetLogLevel() LogLevel {
	return LogLevel(strings.ToUpper(c.LogLevel))
}

// GetFolder returns the standards folder path.
func (c *Config) GetFolder() string {
	return c.Folder
}

// GetMaxStandards returns the maximum number of standards to load.
func (c *Config) GetMaxStandards() int {
	return c.MaxStandards
}

// GetMaxStandardSize returns the maximum size of a single standard file in bytes.
func (c *Config) GetMaxStandardSize() int {
	return c.MaxStandardSize
}
