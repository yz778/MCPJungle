package api

import (
	"github.com/gin-gonic/gin"
	"github.com/mark3labs/mcp-go/server"
)

const V0PathPrefix = "/api/v0"

// NewServer initializes a new Gin server with the MCPJungle MCP server and API endpoints.
func NewServer() *gin.Engine {
	r := gin.Default()

	r.GET(
		"/health",
		func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		},
	)

	// Set up the proxy MCP server on /mcp path
	mcpServer := server.NewMCPServer(
		"MCPJungle Proxy MCP Server",
		"0.0.1",
		server.WithToolCapabilities(true),
	)
	streamableHttpServer := server.NewStreamableHTTPServer(mcpServer)
	r.Any("/mcp", gin.WrapH(streamableHttpServer))

	apiV0 := r.Group(V0PathPrefix)
	{
		apiV0.POST("/servers", registerServerHandler)
		apiV0.DELETE("/servers/:name", deregisterServerHandler)
		apiV0.GET("/servers", listServersHandler)
		apiV0.GET("/tools", listToolsHandler)
		apiV0.POST("/tools/invoke", invokeToolHandler)
		apiV0.GET("/tool", getToolHandler)
	}

	return r
}
