# AGENTS.md

This file provides guidance to agents when working with code in this repository.

## Build/Test Commands

- `task lint` - Run golangci-lint (MUST pass before finalizing code)
- `task test` - Run unit tests (tests located alongside implementation files)
- `task test-coverage` - Run tests with coverage report
- `task generate` - Generate mocks using go:generate directives
- `task build` - Build the MCP server binary
- `task run-debug` - Run with DEBUG log level for development

## Critical Project Patterns

- Mock generation: Use `//go:generate mockgen -source=interfaces.go -destination=mocks.go -package=server` pattern
- Standard files must be in `{AGENT_STANDARDS_MCP_FOLDER}/standards/` directory with `.md` extension only
- Frontmatter parsing: Only `description` field is processed from YAML frontmatter, other fields are skipped
- Domain entities are pure (no serialization tags) - separate from transport/data layers
- MCP server uses STDIO transport only - no HTTP or other transports
- Audit logging is mandatory for all client requests/responses via `LogClientRequest`/`LogClientResponse`

## Code Style Requirements

- Error handling: MUST use structured logging with slog, not fmt.Errorf for user-facing errors
- Interface-based design: All major components implement interfaces for testability
- Zero linting errors allowed before code finalization

## Testing Architecture

- Unit test files located alongside implementation files, not in separate test directories
- Integration tests: `internal/test`

## Configuration

- Environment variables only: `AGENT_STANDARDS_MCP_*` prefix
- Default standards folder: `~/agent-standards/standards/`
- Log levels: NONE/DEBUG/INFO/WARN/ERROR (default: ERROR)
- MUST run `task generate` after modifying interfaces to update mocks