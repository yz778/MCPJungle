package api

import (
	"github.com/duaraghav8/mcpjungle/internal/model"
	"github.com/duaraghav8/mcpjungle/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/mark3labs/mcp-go/server"
	"net/http"
)

func registerServerHandler(mcpProxy *server.MCPServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req model.McpServer
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := service.RegisterMcpServer(c, &req, mcpProxy); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, req)
	}
}

func deregisterServerHandler(mcpProxy *server.MCPServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("name")
		if err := service.DeregisterMcpServer(name, mcpProxy); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusNoContent)
	}
}

func listServersHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		servers, err := service.ListMcpServers()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, servers)
	}
}
