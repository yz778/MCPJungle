package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mcpjungle/mcpjungle/internal/model"
)

// initMCPProxyServer initializes the MCP proxy server.
// It loads all the registered MCP tools from the database into the proxy server.
func (m *MCPService) initMCPProxyServer() error {
	tools, err := m.ListTools()
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

		m.mcpProxyServer.AddTool(tool, m.mcpProxyToolCallHandler)
	}
	return nil
}

// mcpProxyToolCallHandler handles tool calls for the MCP proxy server
// by forwarding the request to the appropriate upstream MCP server and
// relaying the response back.
func (m *MCPService) mcpProxyToolCallHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name := request.Params.Name
	serverName, toolName, ok := splitServerToolName(name)
	if !ok {
		return nil, fmt.Errorf("invalid input: tool name does not contain a %s separator", serverToolNameSep)
	}

	serverMode := ctx.Value("mode").(model.ServerMode)
	if serverMode == model.ModeProd {
		// In production mode, we need to check whether the MCP client is authorized to access the MCP server.
		// If not, return error Unauthorized.
		c := ctx.Value("client").(*model.McpClient)
		if !c.CheckHasServerAccess(serverName) {
			return nil, fmt.Errorf(
				"client %s is not authorized to access MCP server %s", c.Name, serverName,
			)
		}
	}

	// get the MCP server details from the database
	server, err := m.GetMcpServer(serverName)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get details about MCP server %s from DB: %w", serverName, err,
		)
	}

	// connect to the upstream MCP server that actually provides the tool
	mcpClient, err := createMcpServerConn(ctx, server)
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
