package api

import (
	"encoding/json"
	"net/http"

	"github.com/duaraghav8/mcpjungle/internal/service"
	"github.com/gin-gonic/gin"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

// Invoke forwards the JSON body to the tool URL and streams response back.
func Invoke(c *gin.Context) {
	name := c.Param("name")
	tool, err := service.GetTool(name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "tool not found"})
		return
	}

	mcpClient, err := client.NewStreamableHttpClient(tool.URL)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": "failed to create streamable HTTP client for MCP server: " + err.Error(),
			},
		)
		return
	}

	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "MCPJungle client for MCP server " + tool.URL,
		Version: "1.0.0",
	}
	initRequest.Params.Capabilities = mcp.ClientCapabilities{}

	_, err = mcpClient.Initialize(c, initRequest)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": "failed to initialize connection with MCP server: " + err.Error(),
			},
		)
		return
	}

	callToolReq := mcp.CallToolRequest{}
	callToolReq.Params.Name = name

	var args map[string]any
	if err := json.NewDecoder(c.Request.Body).Decode(&args); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "failed to decode request body: " + err.Error(),
			},
		)
		return
	}
	callToolReq.Params.Arguments = args

	callToolResp, err := mcpClient.CallTool(c, callToolReq)
	if err != nil {
		c.JSON(
			http.StatusBadGateway,
			gin.H{
				"error": "failed to call tool: " + err.Error(),
			},
		)
		return
	}
	textContent, ok := callToolResp.Content[0].(mcp.TextContent)
	if !ok {
		c.JSON(
			http.StatusBadGateway,
			gin.H{
				"error": "failed to get text content from tool response",
			},
		)
		return
	}

	c.Status(http.StatusOK)
	c.Writer.WriteString(textContent.Text)
}
