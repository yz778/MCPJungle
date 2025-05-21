package service

import (
	"context"
	"fmt"
	"github.com/duaraghav8/mcpjungle/internal/db"
	"github.com/duaraghav8/mcpjungle/internal/models"
	"github.com/mark3labs/mcp-go/mcp"
)

// RegisterMcpServer registers a new MCP server in the database.
// It also registers all the Tools provided by the server.
// Tool registration is on best-effort basis and does not fail the server registration.
func RegisterMcpServer(ctx context.Context, s *models.McpServer) error {
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

	// fetch all tools from the server so they can be added to the DB
	resp, err := c.ListTools(ctx, mcp.ListToolsRequest{})
	if err != nil {
		return fmt.Errorf("failed to fetch tools from MCP server %s: %w", s.Name, err)
	}
	for _, tool := range resp.Tools {
		t := &models.Tool{
			ServerID:    s.ID,
			Name:        tool.GetName(),
			Description: tool.Description,
			// TODO: Also add the tool's input schema, annotation, etc
		}
		if err := registerTool(t); err != nil {
			// TODO: Add error log about this failure
			// If registration of a tool fails, we should not fail the entire server registration.
			// Instead, continue with the next tool.

			//return fmt.Errorf("failed to register tool %s in DB: %w", mergeServerToolNames(s.Name, t.Name), err)
		}
	}

	return nil
}

// DeregisterMcpServer deregisters an MCP server from the database.
// It also deregisters all the tools registered by the server.
// If even a singe tool fails to deregister, the server deregistration fails.
func DeregisterMcpServer(name string) error {
	s, err := GetMcpServer(name)
	if err != nil {
		return fmt.Errorf("failed to get MCP server %s from DB: %w", name, err)
	}
	if err := deregisterToolsByServer(s); err != nil {
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
func ListMcpServers() ([]models.McpServer, error) {
	var servers []models.McpServer
	if err := db.DB.Find(&servers).Error; err != nil {
		return nil, err
	}
	return servers, nil
}

// GetMcpServer fetches a server from the database by name.
func GetMcpServer(name string) (*models.McpServer, error) {
	var server models.McpServer
	if err := db.DB.Where("name = ?", name).First(&server).Error; err != nil {
		return nil, err
	}
	return &server, nil
}
