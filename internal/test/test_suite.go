package test

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/n-r-w/agent-standards-mcp/internal/config"
	"github.com/n-r-w/agent-standards-mcp/internal/logging"
	"github.com/n-r-w/agent-standards-mcp/internal/server"
	"github.com/n-r-w/agent-standards-mcp/internal/standards"
	"github.com/stretchr/testify/require"
)

const (
	// serverStartupDelay is the time to wait for the server to start
	serverStartupDelay = 10 * time.Millisecond
	// testFilePermissions are the permissions for test files
	testFilePermissions = 0o600
)

// Suite provides a complete MCP client-server setup for integration testing
type Suite struct {
	Client        *mcp.Client
	ClientSession *mcp.ClientSession
	Server        *MCPTestServer
	cleanup       func()
}

// MCPTestServer wraps the MCP server for testing
type MCPTestServer struct {
	Server         *server.MCP
	Config         *config.Config
	StandardLoader server.StandardLoader
}

// SetupOption defines configuration options for test suite
type SetupOption func(*setupConfig)

type setupConfig struct {
	transportType string
	standardFiles map[string]string
	clientName    string
	clientVersion string
}

// WithCustomStandardFiles configures custom standard files
func WithCustomStandardFiles(files map[string]string) SetupOption {
	return func(c *setupConfig) {
		c.standardFiles = files
	}
}

// WithCommandTransport configures the setup to use command transport (real subprocess)
func WithCommandTransport() SetupOption {
	return func(c *setupConfig) {
		c.transportType = "command"
	}
}

// WithClientInfo configures the MCP client identification
func WithClientInfo(name, version string) SetupOption {
	return func(c *setupConfig) {
		c.clientName = name
		c.clientVersion = version
	}
}

// NewTestSuite creates a complete integration test environment
func NewTestSuite(t *testing.T, opts ...SetupOption) *Suite {
	// Default configuration
	config := &setupConfig{
		transportType: "in-memory", // Default to fastest option
		standardFiles: DefaultStandardFiles(),
		clientName:    "test-client",
		clientVersion: "1.0.0",
	}

	// Apply options
	for _, opt := range opts {
		opt(config)
	}

	// Create test server
	testServer := createTestServer(t, config.standardFiles)

	var clientTransport mcp.Transport
	var cleanupFuncs []func()

	ctx := context.Background()

	// Set up transport based on type
	switch config.transportType {
	case "command":
		clientTransport = setupCommandTransport(t, testServer)

	case "in-memory":
		fallthrough
	default:
		// Create in-memory transports - MUST use the SAME pair for client and server
		var serverTransport mcp.Transport
		clientTransport, serverTransport = mcp.NewInMemoryTransports()

		// Start server in background
		serverCtx, cancelServer := context.WithCancel(ctx)
		go func() {
			_ = testServer.Server.GetMCPServer().Run(serverCtx, serverTransport)
		}()

		// Give server time to start
		time.Sleep(serverStartupDelay)

		// Add cleanup for in-memory transport
		cleanupFuncs = append(cleanupFuncs, func() {
			cancelServer()
			time.Sleep(serverStartupDelay)
		})
	}

	// Create MCP client
	client := mcp.NewClient(&mcp.Implementation{
		Name:    config.clientName,
		Version: config.clientVersion,
		Title:   config.clientName,
	}, nil)

	// Connect client to server
	clientSession, err := client.Connect(ctx, clientTransport, nil)
	require.NoError(t, err, "Failed to connect MCP client to server")

	// Create cleanup function
	cleanup := func() {
		if clientSession != nil {
			_ = clientSession.Close()
		}
		// Call cleanup functions in reverse order
		for i := len(cleanupFuncs) - 1; i >= 0; i-- {
			cleanupFuncs[i]()
		}
	}

	return &Suite{
		Client:        client,
		ClientSession: clientSession,
		Server:        testServer,
		cleanup:       cleanup,
	}
}

// Cleanup cleans up test resources
func (ts *Suite) Cleanup() {
	if ts.cleanup != nil {
		ts.cleanup()
	}
}

// createTestServer creates a server instance for testing
func createTestServer(t testing.TB, standardFiles map[string]string) *MCPTestServer {
	// Create temporary standards directory
	tempDir, err := os.MkdirTemp("", "agent-standards-test-*")
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = os.RemoveAll(tempDir)
	})

	// Create test standard files
	for filename, content := range standardFiles {
		filePath := tempDir + "/" + filename
		err := os.WriteFile(filePath, []byte(content), testFilePermissions)
		require.NoError(t, err)
	}

	// Set up environment
	oldStandardsFolder := os.Getenv("AGENT_STANDARDS_MCP_FOLDER")
	_ = os.Setenv("AGENT_STANDARDS_MCP_FOLDER", tempDir)

	t.Cleanup(func() {
		if oldStandardsFolder != "" {
			_ = os.Setenv("AGENT_STANDARDS_MCP_FOLDER", oldStandardsFolder)
		} else {
			_ = os.Unsetenv("AGENT_STANDARDS_MCP_FOLDER")
		}
	})

	// Load configuration
	cfg, err := config.Load()
	require.NoError(t, err)

	// Create logger
	loggerFactory := logging.NewLoggerFactory()
	structuredLogger, err := loggerFactory.CreateStructuredLogger(cfg)
	require.NoError(t, err)

	t.Cleanup(func() {
		if structuredLogger != nil {
			_ = structuredLogger.Close()
		}
	})

	// Create audit logger
	auditLogger, err := loggerFactory.CreateAudit(cfg)
	require.NoError(t, err)

	t.Cleanup(func() {
		if auditLogger != nil {
			// audit logger doesn't have Close method
			_ = auditLogger
		}
	})

	// Create standard loader
	standardLoader := standards.NewFileStandardLoader()

	// Create MCP server
	mcpServer, err := server.New(cfg, structuredLogger, auditLogger, standardLoader)
	require.NoError(t, err)

	// Register tools
	err = mcpServer.RegisterTools()
	require.NoError(t, err)

	return &MCPTestServer{
		Server:         mcpServer,
		Config:         cfg,
		StandardLoader: standardLoader,
	}
}

// setupCommandTransport sets up a command transport for real subprocess testing
func setupCommandTransport(_ *testing.T, server *MCPTestServer) mcp.Transport {
	// Calculate project root directory to ensure subprocess runs from correct location
	projectRoot, err := filepath.Abs(filepath.Join(".", "..", ".."))
	if err != nil {
		// Fall back to safer relative path if absolute path calculation fails
		projectRoot = filepath.Join("..", "..")
	}

	// Prepare server command
	ctx := context.Background()
	serverCmd := exec.CommandContext(ctx, "go", "run", "./cmd/agent-standards-mcp")
	serverCmd.Dir = projectRoot // Ensure we run from project root
	serverCmd.Env = append(os.Environ(),
		"AGENT_STANDARDS_MCP_FOLDER="+server.Config.GetFolder(),
		"AGENT_STANDARDS_MCP_LOG_LEVEL=ERROR", // Reduce log noise in tests
	)

	// Create transport using CommandTransport from MCP SDK
	const terminateDuration = 5 * time.Second // Allow graceful shutdown
	transport := &mcp.CommandTransport{
		Command:           serverCmd,
		TerminateDuration: terminateDuration,
	}

	return transport
}
