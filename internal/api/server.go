package api

import (
	"fmt"
	"github.com/duaraghav8/mcpjungle/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/mark3labs/mcp-go/server"
)

const V0PathPrefix = "/api/v0"

type Server struct {
	port           string
	router         *gin.Engine
	mcpProxyServer *server.MCPServer
}

// NewServer initializes a new Gin server for MCPJungle registry and MCP proxy
func NewServer(port string, mcpProxyServer *server.MCPServer) (*Server, error) {
	r, err := newRouter(mcpProxyServer)
	if err != nil {
		return nil, err
	}
	s := &Server{
		port:           port,
		router:         r,
		mcpProxyServer: mcpProxyServer,
	}
	return s, nil
}

// Start runs the Gin server (blocking call)
func (s *Server) Start() error {
	if err := s.router.Run(":" + s.port); err != nil {
		return fmt.Errorf("failed to run the server: %w", err)
	}
	return nil
}

// newRouter sets up the Gin router with the MCP proxy server and API endpoints.
func newRouter(mcpProxyServer *server.MCPServer) (*gin.Engine, error) {
	r := gin.Default()

	r.GET(
		"/health",
		func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		},
	)

	// Set up the MCP proxy server on /mcp
	if err := service.InitMCPProxyServer(mcpProxyServer); err != nil {
		return nil, fmt.Errorf("failed to initialize MCP proxy server: %w", err)
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
