# Agent Standards MCP Server

A Model Context Protocol (MCP) server that provides agents with access to standards and rules. It enables agents to list available standards and retrieve their content programmatically.

## Installation

### Binary Releases

Pre-compiled binaries are available for multiple platforms:

- **Linux (AMD64)**: `agent-standards-mcp-v*-linux-amd64.tar.gz`
- **macOS (Intel)**: `agent-standards-mcp-v*-darwin-amd64.tar.gz`
- **macOS (Apple Silicon)**: `agent-standards-mcp-v*-darwin-arm64.tar.gz`
- **Windows (AMD64)**: `agent-standards-mcp-v*-windows-amd64.zip`

Download the latest release from [GitHub Releases](https://github.com/n-r-w/agent-standards-mcp/releases).

### Build from Source

```bash
go build -o agent-standards-mcp ./cmd/agent-standards-mcp
```

## Usage

The server provides two tools:

- **list_standards**: Lists all available standards with their descriptions
- **get_standards**: Retrieves the full content of specific standards by name

## IDE Integration

### Claude Code
Add the MCP server using the CLI:
```bash
claude-code mcp add agent-standards --command /path/to/agent-standards-mcp
```

### Cursor IDE
Add to your Cursor settings:
```json
{
  "mcp.servers": {
    "agent-standards": {
      "command": "/path/to/agent-standards-mcp"
    }
  }
}
```

## Configuration

Configure the server with environment variables:

- `AGENT_STANDARDS_MCP_LOG_LEVEL`: Log level (NONE/DEBUG/INFO/WARN/ERROR, default: "ERROR")
- `AGENT_STANDARDS_MCP_FOLDER`: Standards folder path (default: "~/agent-standards")
- `AGENT_STANDARDS_MCP_MAX_STANDARDS`: Maximum number of standards to load (default: 100)
- `AGENT_STANDARDS_MCP_MAX_STANDARD_SIZE`: Maximum size of a standard file in bytes (default: 10240)

## Development

### Prerequisites

- Go 1.25.1 or later
- Task (install with `go install github.com/go-task/task/v3/cmd/task@latest`)

### Build Commands

```bash
# Build the application
task build

# Run tests
task test

# Run tests with coverage
task test-coverage

# Run linter
task lint

# Format code
task format

# Clean build artifacts
task clean

# Run the application
task run

# Run with DEBUG log level
task run-debug

# Update dependencies
task deps

# Generate code
task generate

# Install development tools
task install-tools

# Run full CI pipeline
task ci
```

## Release Process

This project uses automated releases with GitHub Actions:

1. **Commit Messages**: Use [conventional commits](https://www.conventionalcommits.org/) format
2. **Version Tags**: Create semantic version tags (e.g., `v1.0.0`)
3. **Automatic Release**: Pushing a tag triggers the release workflow
4. **Generated Assets**: Binaries for all platforms, checksums

### Release Steps

```bash
# 1. Make your changes with conventional commit messages
git commit -m "feat: add new feature"

# 2. Create and push a version tag
git tag v1.0.0
git push origin v1.0.0

# 3. GitHub Actions will automatically create the release
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes with conventional commit messages
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Changelog

See the [GitHub Releases](https://github.com/n-r-w/agent-standards-mcp/releases) page for a detailed changelog.