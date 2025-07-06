package api

import (
	"github.com/gin-gonic/gin"
	"github.com/mcpjungle/mcpjungle/internal/model"
	"github.com/mcpjungle/mcpjungle/internal/service/mcp"
	"net/http"
)

func registerServerHandler(mcpService *mcp.MCPService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req model.McpServer
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := mcpService.RegisterMcpServer(c, &req); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, req)
	}
}

func deregisterServerHandler(mcpService *mcp.MCPService) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("name")
		if err := mcpService.DeregisterMcpServer(name); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusNoContent)
	}
}

func listServersHandler(mcpService *mcp.MCPService) gin.HandlerFunc {
	return func(c *gin.Context) {
		servers, err := mcpService.ListMcpServers()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, servers)
	}
}
