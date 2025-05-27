package service

import (
	"fmt"
	"github.com/mark3labs/mcp-go/server"
)

// MCPService coordinates operations amongst the registry database, mcp proxy server and upstream MCP servers.
// It is responsible for maintaining data consistency and providing a unified interface for MCP operations.
type MCPService struct {
	mcpProxyServer *server.MCPServer
}

// NewMCPService creates a new instance of MCPService.
// It initializes the MCP proxy server by loading all registered tools from the database.
func NewMCPService(mcpProxyServer *server.MCPServer) (*MCPService, error) {
	s := &MCPService{
		mcpProxyServer: mcpProxyServer,
	}
	if err := s.initMCPProxyServer(); err != nil {
		return nil, fmt.Errorf("failed to initialize MCP proxy server: %w", err)
	}
	return s, nil
}
