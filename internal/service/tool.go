package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/duaraghav8/mcpjungle/internal/db"
	"github.com/duaraghav8/mcpjungle/internal/models"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

// ListTools returns all tools registered in the registry.
func ListTools() ([]models.Tool, error) {
	var tools []models.Tool
	if err := db.DB.Find(&tools).Error; err != nil {
		return nil, err
	}
	// prepend server name to tool names to ensure we only return the unique names of tools to user
	for i := range tools {
		var s models.McpServer
		if err := db.DB.First(&s, "id = ?", tools[i].ServerID).Error; err != nil {
			return nil, fmt.Errorf("failed to get server for tool %s: %w", tools[i].Name, err)
		}
		tools[i].Name = mergeServerToolNames(s.Name, tools[i].Name)
	}
	return tools, nil
}

// ListToolsByServer fetches tools provided by an MCP server from the registry.
func ListToolsByServer(name string) ([]models.Tool, error) {
	if err := validateServerName(name); err != nil {
		return nil, err
	}

	s, err := GetMcpServer(name)
	if err != nil {
		return nil, fmt.Errorf("failed to get MCP server %s from DB: %w", name, err)
	}

	var tools []models.Tool
	if err := db.DB.Where("server_id = ?", s.ID).Find(&tools).Error; err != nil {
		return nil, fmt.Errorf("failed to get tools for server %s from DB: %w", name, err)
	}

	// prepend server name to tool names to ensure we only return the unique names of tools to user
	for i := range tools {
		tools[i].Name = mergeServerToolNames(s.Name, tools[i].Name)
	}

	return tools, nil
}

func GetTool(name string) (*models.Tool, error) {
	serverName, toolName, ok := splitServerToolName(name)
	if !ok {
		return nil, fmt.Errorf("invalid input: tool name does not contain a %s separator", serverToolNameSep)
	}

	s, err := GetMcpServer(serverName)
	if err != nil {
		return nil, fmt.Errorf("failed to get MCP server %s from DB: %w", serverName, err)
	}

	var tool models.Tool
	if err := db.DB.Where("server_id = ? AND name = ?", s.ID, toolName).First(&tool).Error; err != nil {
		return nil, fmt.Errorf("failed to get tool %s from DB: %w", name, err)
	}
	// set the tool name back to the full name including server name
	tool.Name = name
	return &tool, nil
}

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

	// CallTool() doesn't return an error if the tool is not found.
	// Instead, it returns a "not found" message in the response.
	// TODO: detect this and return an error if the tool is not found.
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

// registerServerTools fetches all tools from an MCP server and registers them in the DB.
func registerServerTools(ctx context.Context, s *models.McpServer, c *client.Client) error {
	// fetch all tools from the server so they can be added to the DB
	resp, err := c.ListTools(ctx, mcp.ListToolsRequest{})
	if err != nil {
		return fmt.Errorf("failed to fetch tools from MCP server %s: %w", s.Name, err)
	}
	for _, tool := range resp.Tools {
		// extracting json schema is currently on best-effort basis
		// if it fails, we log the error and continue with the next tool
		jsonSchema, _ := json.Marshal(tool.InputSchema)

		t := &models.Tool{
			ServerID:    s.ID,
			Name:        tool.GetName(),
			Description: tool.Description,
			InputSchema: jsonSchema,
		}
		if err := db.DB.Create(t).Error; err != nil {
			// TODO: Add error log about this failure
			// If registration of a tool fails, we should not fail the entire server registration.
			// Instead, continue with the next tool.

			//return fmt.Errorf("failed to register tool %s in DB: %w", mergeServerToolNames(s.Name, t.Name), err)
		}
	}
	return nil
}

// deregisterServerTools deletes all tools that belong to an MCP server from the DB.
func deregisterServerTools(s *models.McpServer) error {
	if err := db.DB.Where("server_id = ?", s.ID).Delete(&models.Tool{}).Error; err != nil {
		return fmt.Errorf("failed to delete tools for server %s: %w", s.Name, err)
	}
	return nil
}
