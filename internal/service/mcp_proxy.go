package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// InitMCPProxyServer initializes the MCP proxy server.
// It loads all the registered MCP tools from the database into the proxy server.
func InitMCPProxyServer(ps *server.MCPServer) error {
	tools, err := ListTools()
	if err != nil {
		return fmt.Errorf("failed to list tools from DB: %w", err)
	}
	for _, tm := range tools {
		// Add tool to the MCP proxy server
		tool := mcp.NewTool(tm.Name)
		tool.Description = tm.Description

		var inputSchema mcp.ToolInputSchema
		if err := json.Unmarshal(tm.InputSchema, &inputSchema); err != nil {
			return fmt.Errorf(
				"failed to unmarshal input schema %s for tool %s: %w", tm.InputSchema, tm.Name, err,
			)
		}
		tool.InputSchema = inputSchema

		// TODO: Add other attributes to the tool, such as annotations

		ps.AddTool(tool, mcpProxyToolCallHandler)
	}
	return nil
}

// mcpProxyToolCallHandler handles tool calls for the MCP proxy server
// by forwarding the request to the appropriate upstream MCP server and
// relaying the response back.
func mcpProxyToolCallHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name := request.Params.Name
	serverName, toolName, ok := splitServerToolName(name)
	if !ok {
		return nil, fmt.Errorf("invalid input: tool name does not contain a %s separator", serverToolNameSep)
	}

	// get the MCP server details from the database
	server, err := GetMcpServer(serverName)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get details about MCP server %s from DB: %w", serverName, err,
		)
	}

	// connect to the upstream MCP server that actually provides the tool
	mcpClient, err := createMcpServerConn(ctx, server.URL)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to create connection to MCP server %s: %w", serverName, err,
		)
	}
	defer mcpClient.Close()

	// Ensure the tool name is set correctly, ie, without the server name prefix
	request.Params.Name = toolName

	// forward the request to the upstream MCP server and relay the response back
	return mcpClient.CallTool(ctx, request)
}
