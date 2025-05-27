package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mark3labs/mcp-go/server"
)

const V0PathPrefix = "/api/v0"

// NewServer initializes a new Gin server with the MCPJungle MCP server and API endpoints.
func NewServer() (*gin.Engine, error) {
	r := gin.Default()

	r.GET(
		"/health",
		func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		},
	)

	// Set up the proxy MCP server on /mcp
	mcpProxyServer, err := newMCPProxyServer()
	if err != nil {
		return nil, fmt.Errorf("failed to create MCP proxy server: %w", err)
	}
	streamableHttpServer := server.NewStreamableHTTPServer(mcpProxyServer)
	r.Any("/mcp", gin.WrapH(streamableHttpServer))

	// Setup API endpoints
	apiV0 := r.Group(V0PathPrefix)
	{
		apiV0.POST("/servers", registerServerHandler(mcpProxyServer))
		apiV0.DELETE("/servers/:name", deregisterServerHandler(mcpProxyServer))
		apiV0.GET("/servers", listServersHandler())
		apiV0.GET("/tools", listToolsHandler())
		apiV0.POST("/tools/invoke", invokeToolHandler())
		apiV0.GET("/tool", getToolHandler())
	}

	return r, nil
}
