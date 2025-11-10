# Agent Standards MCP Server

A Model Context Protocol (MCP) server that provides agents with access to standards and rules. It enables agents to list available standards and retrieve their content programmatically.

## Why use this MCP server?

Different rules are needed for different cases. Some for coding (and different ones for different languages), others for working with databases, third for working with APIs, etc. If all rules are loaded into the context at once (for example, placed in AGENTS.md), this will lead to:
- Increased cost of LLM requests (more tokens - higher price)
- Slower response (more tokens - longer request processing)
- Loss of LLM focus (too much information - LLM may get confused about what's important and what's not)

Everyone solves this problem differently. For example:
- Claude Code offers a skills mechanism, but unfortunately, it cannot be controlled, because Claude itself decides which skill to use. It is impossible to set rules for when to use one or another skill manually through CLAUDE.md
- In Github Copilot, you can create different sets of rules and specify for which file extensions they should apply. But this will only work at the moment of calling tools that change the content of such files. I.e., at the planning and decision-making stage, the LLM will not have access to these rules and will most likely make incorrect decisions.

This MCP server solves these problems:
- You can explicitly specify to the agent when to load standards by writing this in the rules or commands
- The rules catalog is centralized. There is no binding to the implementation of a specific agent/IDE - can be used with any LLM agent that supports MCP

## Available Tools

The server provides two tools:

- **list_standards**: Lists all available standards with their descriptions
- **get_standards**: Retrieves the full content of specific standards by name

## Logs

By default, the server logs errors only. You can adjust the log level using the `AGENT_STANDARDS_MCP_LOG_LEVEL` environment variable. Available levels are: NONE, DEBUG, INFO, WARN, ERROR. Default location: `~/agent-standards-mcp/logs/`

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

or use Task:

```bash
task build
```

### macOS Installation Notes

macOS may block execution of downloaded binaries by default due to security settings. To allow the executable to run:

1. **First execution attempt**: Run the executable from terminal
   ```bash
   ./agent-standards-mcp
   ```
   This will show a security warning.

2. **Allow execution via System Settings**:
   - Open **System Settings** → **Privacy & Security** → **Security**
   - Find the message about the blocked executable
   - Click **"Allow Anyway"**

3. **Second execution**: Run the executable again
   ```bash
   ./agent-standards-mcp
   ```

4. **Confirm execution**: A dialog will appear asking for confirmation
   - Click **"Open"** and enter your password if prompted
   - The executable will now be allowed to run

After these steps, the executable will be permanently allowed to run on your system.

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

### Configuration

Configure the server with environment variables:

- `AGENT_STANDARDS_MCP_LOG_LEVEL`: Log level (NONE/DEBUG/INFO/WARN/ERROR, default: "ERROR")
- `AGENT_STANDARDS_MCP_FOLDER`: Standards folder path (default: "~/agent-standards")
- `AGENT_STANDARDS_MCP_MAX_STANDARDS`: Maximum number of standards to load (default: 100)
- `AGENT_STANDARDS_MCP_MAX_STANDARD_SIZE`: Maximum size of a standard file in bytes (default: 10240)

## Usage

### Standards Management

Place your standard markdown files in the specified standards folder (default: `~/agent-standards/standards`)

Use the following format:

```markdown
---
description: {A brief description of the standard}
---

## {Standard Title. Start with ##}
{Full content of the standard goes here. Follow ## headings for sections.}
```

LLM Agent will be able to access these standards via the MCP server:
- **List Standards**: Use the `list_standards` tool to get a list of available standard names with descriptions.
- **Get Standard Content**: Use the `get_standards` tool to retrieve the full content of specific standards by name.

### Additional Rules

Some LLMs may require additional rules to properly utilize the standards. You may want to add extra rules in your AGENTS.md/CLAUDE.md etc.
Be careful to prompt injection prevention, because some LLMs (like GPT-5) may stop responding if they decide the rules are unsafe.

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