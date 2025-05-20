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
	serverName, toolName, ok := strings.Cut(name, "/")
	if !ok {
		// there is no separator "/" in tool name, we cannot extract mcp server name
		// this is invalid input
		return "", errors.New("invalid input: tool name does not contain a '/' separator")
	}

	server, err := GetMcpServer(serverName)
	if err != nil {
		return "", fmt.Errorf(
			"failed to get details about MCP server %s from DB: %w",
			serverName,
			err,
		)
	}

	mcpClient, err := client.NewStreamableHttpClient(server.URL)
	if err != nil {
		return "", fmt.Errorf(
			"failed to create streamable HTTP client for MCP server %s: %w",
			serverName,
			err,
		)
	}

	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "mcpjungle mcp client for " + server.URL,
		Version: "0.1",
	}
	initRequest.Params.Capabilities = mcp.ClientCapabilities{}

	_, err = mcpClient.Initialize(ctx, initRequest)
	if err != nil {
		return "", fmt.Errorf(
			"failed to initialize connection with MCP server %s: %w",
			serverName,
			err,
		)
	}

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
