package api

import (
	"fmt"
	"github.com/duaraghav8/mcpjungle/internal/service"
	"github.com/mark3labs/mcp-go/server"
)

func newMCPProxyServer() (*server.MCPServer, error) {
	mcpProxyServer := server.NewMCPServer(
		"MCPJungle Proxy MCP Server",
		"0.0.1",
		server.WithToolCapabilities(true),
	)
	if err := service.InitMCPProxyServer(mcpProxyServer); err != nil {
		return nil, fmt.Errorf("failed to initialize MCP proxy server: %w", err)
	}
	return mcpProxyServer, nil
}
