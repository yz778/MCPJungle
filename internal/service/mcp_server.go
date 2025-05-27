package service

import (
	"context"
	"fmt"
	"github.com/duaraghav8/mcpjungle/internal/db"
	"github.com/duaraghav8/mcpjungle/internal/model"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterMcpServer registers a new MCP server in the database.
// It also registers all the Tools provided by the server.
// Tool registration is on best-effort basis and does not fail the server registration.
// Registered tools are also added to the MCP proxy server.
func RegisterMcpServer(ctx context.Context, s *model.McpServer, mcpProxy *server.MCPServer) error {
	if err := validateServerName(s.Name); err != nil {
		return err
	}

	// test that the server is reachable and is MCP-compliant
	c, err := createMcpServerConn(ctx, s.URL)
	if err != nil {
		return fmt.Errorf("failed to connect to MCP server %s: %w", s.Name, err)
	}
	defer c.Close()

	// register the server in the DB
	if err := db.DB.Create(s).Error; err != nil {
		return fmt.Errorf("failed to register mcp server: %w", err)
	}

	if err = registerServerTools(ctx, s, c, mcpProxy); err != nil {
		return fmt.Errorf("failed to register tools for MCP server %s: %w", s.Name, err)
	}
	return nil
}

// DeregisterMcpServer deregisters an MCP server from the database.
// It also deregisters all the tools registered by the server.
// If even a singe tool fails to deregister, the server deregistration fails.
// A deregistered tool is also removed from the MCP proxy server.
func DeregisterMcpServer(name string, mcpProxy *server.MCPServer) error {
	s, err := GetMcpServer(name)
	if err != nil {
		return fmt.Errorf("failed to get MCP server %s from DB: %w", name, err)
	}
	if err := deregisterServerTools(s, mcpProxy); err != nil {
		return fmt.Errorf(
			"failed to deregister tools for server %s, cannot proceed with server deregistration: %w",
			name,
			err,
		)
	}
	if err := db.DB.Delete(s).Error; err != nil {
		return fmt.Errorf("failed to deregister server %s: %w", name, err)
	}
	return nil
}

// ListMcpServers returns all registered MCP servers.
func ListMcpServers() ([]model.McpServer, error) {
	var servers []model.McpServer
	if err := db.DB.Find(&servers).Error; err != nil {
		return nil, err
	}
	return servers, nil
}

// GetMcpServer fetches a server from the database by name.
func GetMcpServer(name string) (*model.McpServer, error) {
	var server model.McpServer
	if err := db.DB.Where("name = ?", name).First(&server).Error; err != nil {
		return nil, err
	}
	return &server, nil
}
