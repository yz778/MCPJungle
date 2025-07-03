package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mark3labs/mcp-go/server"
	"github.com/mcpjungle/mcpjungle/internal/model"
	"github.com/mcpjungle/mcpjungle/internal/service/config"
	"github.com/mcpjungle/mcpjungle/internal/service/mcp"
)

const V0PathPrefix = "/api/v0"

// Server represents the MCPJungle registry server that handles MCP proxy and API requests
type Server struct {
	port   string
	router *gin.Engine

	mcpProxyServer *server.MCPServer
	mcpService     *mcp.MCPService

	configService *config.ServerConfigService
}

// NewServer initializes a new Gin server for MCPJungle registry and MCP proxy
func NewServer(port string, mcpProxyServer *server.MCPServer, mcpService *mcp.MCPService, configService *config.ServerConfigService) (*Server, error) {
	r, err := newRouter(mcpProxyServer, mcpService, configService)
	if err != nil {
		return nil, err
	}
	s := &Server{
		port:           port,
		router:         r,
		mcpProxyServer: mcpProxyServer,
		mcpService:     mcpService,
		configService:  configService,
	}
	return s, nil
}

func (s *Server) Init(mode model.ServerMode) error {
	if err := s.configService.Init(mode); err != nil {
		return fmt.Errorf("failed to initialize server config in %s mode: %w", mode, err)
	}
	return nil
}

// Start runs the Gin server (blocking call)
func (s *Server) Start() error {
	if err := s.router.Run(":" + s.port); err != nil {
		return fmt.Errorf("failed to run the server: %w", err)
	}
	return nil
}

// newRouter sets up the Gin router with the MCP proxy server and API endpoints.
func newRouter(mcpProxyServer *server.MCPServer, mcpService *mcp.MCPService, configService *config.ServerConfigService) (*gin.Engine, error) {
	r := gin.Default()

	r.GET(
		"/health",
		func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		},
	)

	r.POST("/init", registerInitServerHandler(configService))

	// Set up the MCP proxy server on /mcp
	streamableHttpServer := server.NewStreamableHTTPServer(mcpProxyServer)
	r.Any("/mcp", gin.WrapH(streamableHttpServer))

	// Setup API endpoints
	apiV0 := r.Group(V0PathPrefix)
	{
		apiV0.POST("/servers", registerServerHandler(mcpService))
		apiV0.DELETE("/servers/:name", deregisterServerHandler(mcpService))
		apiV0.GET("/servers", listServersHandler(mcpService))
		apiV0.GET("/tools", listToolsHandler(mcpService))
		apiV0.POST("/tools/invoke", invokeToolHandler(mcpService))
		apiV0.GET("/tool", getToolHandler(mcpService))
	}

	return r, nil
}
