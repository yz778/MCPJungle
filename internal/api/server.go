package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mark3labs/mcp-go/server"
	"github.com/mcpjungle/mcpjungle/internal/model"
	"github.com/mcpjungle/mcpjungle/internal/service/config"
	"github.com/mcpjungle/mcpjungle/internal/service/mcp"
	"github.com/mcpjungle/mcpjungle/internal/service/user"
	"net/http"
	"strings"
)

const V0PathPrefix = "/api/v0"

type ServerOptions struct {
	// Port is the HTTP ports to bind the server to
	Port string

	MCPProxyServer *server.MCPServer
	MCPService     *mcp.MCPService
	ConfigService  *config.ServerConfigService
	UserService    *user.UserService
}

// Server represents the MCPJungle registry server that handles MCP proxy and API requests
type Server struct {
	port   string
	router *gin.Engine

	mcpProxyServer *server.MCPServer
	mcpService     *mcp.MCPService

	configService *config.ServerConfigService
	userService   *user.UserService
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
		userService:    opts.UserService,
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

// InitDev initializes the server configuration in the Development mode.
// This method does not create an admin user because that is irrelevant in dev mode.
func (s *Server) InitDev() error {
	_, err := s.configService.Init(model.ModeDev)
	if err != nil {
		return fmt.Errorf("failed to initialize server config in dev mode: %w", err)
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

// requireInitialized is middleware to reject requests to certain routes if the server is not initialized
func requireInitialized(configService *config.ServerConfigService) gin.HandlerFunc {
	return func(c *gin.Context) {
		cfg, err := configService.GetConfig()
		if err != nil || !cfg.Initialized {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "server not initialized"})
			return
		}
		c.Next()
	}
}

// requireAuthIfProd is middleware that checks for a valid admin token if the server is in production mode.
// In development mode, it allows all requests without authentication.
func checkAuthForAPIAccess(configService *config.ServerConfigService, userService *user.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		cfg, err := configService.GetConfig()
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusServiceUnavailable, gin.H{"error": "failed to fetch server config while checking auth"},
			)
			return
		}
		if cfg.Mode == model.ModeDev {
			c.Next()
			return
		}
		authHeader := c.GetHeader("Authorization")
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing access token"})
			return
		}
		_, err = userService.VerifyAdminToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid access token"})
			return
		}
		c.Next()
	}
}

// newRouter sets up the Gin router with the MCP proxy server and API endpoints.
func newRouter(opts *ServerOptions) (*gin.Engine, error) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.GET(
		"/health",
		func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		},
	)

	r.POST("/init", registerInitServerHandler(opts.ConfigService, opts.UserService))

	requireInit := requireInitialized(opts.ConfigService)
	checkAuth := checkAuthForAPIAccess(opts.ConfigService, opts.UserService)

	// Set up the MCP proxy server on /mcp
	streamableHttpServer := server.NewStreamableHTTPServer(opts.MCPProxyServer)
	r.Any("/mcp", requireInit, gin.WrapH(streamableHttpServer))

	// Setup API endpoints
	apiV0 := r.Group(V0PathPrefix, requireInit, checkAuth)
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
