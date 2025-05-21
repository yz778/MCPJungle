package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	"strings"

	"github.com/duaraghav8/mcpjungle/internal/db"
	"github.com/duaraghav8/mcpjungle/internal/models"
)

// ListTools returns all tools or filtered by name / tags.
func ListTools() ([]models.Tool, error) {
	var tools []models.Tool
	if err := db.DB.Find(&tools).Error; err != nil {
		return nil, err
	}
	return tools, nil
}

// ListToolsByServer fetches tools from the specified MCP server
func ListToolsByServer(server string) ([]models.Tool, error) {}

// InvokeTool invokes a tool from a registered MCP server and returns its response.
func InvokeTool(ctx context.Context, name string, args map[string]any) (string, error) {
	serverName, toolName, ok := splitServerToolName(name)
	if !ok {
		return "", fmt.Errorf("invalid input: tool name does not contain a %s separator", serverToolNameSep)
	}

	server, err := GetMcpServer(serverName)
	if err != nil {
		return "", fmt.Errorf(
			"failed to get details about MCP server %s from DB: %w",
			serverName,
			err,
		)
	}

	mcpClient, err := createMcpServerConn(ctx, server.URL)
	if err != nil {
		return "", fmt.Errorf(
			"failed to create connection to MCP server %s: %w", serverName, err,
		)
	}
	defer mcpClient.Close()

	callToolReq := mcp.CallToolRequest{}
	callToolReq.Params.Name = toolName
	callToolReq.Params.Arguments = args

	callToolResp, err := mcpClient.CallTool(ctx, callToolReq)
	if err != nil {
		return "", fmt.Errorf(
			"failed to call tool %s on MCP server %s: %w",
			toolName,
			serverName,
			err,
		)
	}

	textContent, ok := callToolResp.Content[0].(mcp.TextContent)
	if !ok {
		return "", errors.New("failed to get text content from tool response")
	}
	return textContent.Text, nil
}

// registerTool registers a tool in the database.
func registerTool(t *models.Tool) error {
	return db.DB.Create(t).Error
}

func deregisterToolsByServer(s *models.McpServer) error {
	if err := db.DB.Where("server_id = ?", s.ID).Delete(&models.Tool{}).Error; err != nil {
		return fmt.Errorf("failed to delete tools for server %s: %w", s.Name, err)
	}
	return nil
}
