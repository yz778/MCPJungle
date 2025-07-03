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

type ServerOptions struct {
	// Port is the HTTP ports to bind the server to
	Port string

	MCPProxyServer *server.MCPServer
	MCPService     *mcp.MCPService
	ConfigService  *config.ServerConfigService
}

// Server represents the MCPJungle registry server that handles MCP proxy and API requests
type Server struct {
	port   string
	router *gin.Engine

	mcpProxyServer *server.MCPServer
	mcpService     *mcp.MCPService

	configService *config.ServerConfigService
}

// NewServer initializes a new Gin server for MCPJungle registry and MCP proxy
func NewServer(opts *ServerOptions) (*Server, error) {
	r, err := newRouter(opts)
	if err != nil {
		return nil, err
	}
	s := &Server{
		port:           opts.Port,
		router:         r,
		mcpProxyServer: opts.MCPProxyServer,
		mcpService:     opts.MCPService,
		configService:  opts.ConfigService,
	}
	return s, nil
}

// IsInitialized returns true if the server is initialized
func (s *Server) IsInitialized() (bool, error) {
	c, err := s.configService.GetConfig()
	if err != nil {
		return false, fmt.Errorf("failed to get server config: %w", err)
	}
	return c.Initialized, nil
}

// GetMode returns the server mode if the server is initialized, otherwise an error
func (s *Server) GetMode() (model.ServerMode, error) {
	ok, err := s.IsInitialized()
	if err != nil {
		return "", fmt.Errorf("failed to check if server is initialized: %w", err)
	}
	if !ok {
		return "", fmt.Errorf("server is not initialized")
	}
	c, err := s.configService.GetConfig()
	if err != nil {
		return "", fmt.Errorf("failed to get server config: %w", err)
	}
	return c.Mode, nil
}

// Init initializes the server configuration in the specified mode
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
func newRouter(opts *ServerOptions) (*gin.Engine, error) {
	r := gin.Default()

	r.GET(
		"/health",
		func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		},
	)

	r.POST("/init", registerInitServerHandler(opts.ConfigService))

	// Set up the MCP proxy server on /mcp
	streamableHttpServer := server.NewStreamableHTTPServer(opts.MCPProxyServer)
	r.Any("/mcp", gin.WrapH(streamableHttpServer))

	// Setup API endpoints
	apiV0 := r.Group(V0PathPrefix)
	{
		apiV0.POST("/servers", registerServerHandler(opts.MCPService))
		apiV0.DELETE("/servers/:name", deregisterServerHandler(opts.MCPService))
		apiV0.GET("/servers", listServersHandler(opts.MCPService))
		apiV0.GET("/tools", listToolsHandler(opts.MCPService))
		apiV0.POST("/tools/invoke", invokeToolHandler(opts.MCPService))
		apiV0.GET("/tool", getToolHandler(opts.MCPService))
	}

	return r, nil
}
