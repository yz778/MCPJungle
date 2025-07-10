package api

import (
	"github.com/gin-gonic/gin"
	"github.com/mcpjungle/mcpjungle/internal/model"
	"github.com/mcpjungle/mcpjungle/internal/service/mcp_client"
	"net/http"
)

func listMcpClientsHandler(mcpClientService *mcp_client.McpClientService) gin.HandlerFunc {
	return func(c *gin.Context) {
		clients, err := mcpClientService.ListClients()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, clients)
	}
}

func createMcpClientHandler(mcpClientService *mcp_client.McpClientService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req model.McpClient
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}
		if req.Name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
			return
		}
		client, err := mcpClientService.CreateClient(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, client)
	}
}

func deleteMcpClientHandler(mcpClientService *mcp_client.McpClientService) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("name")
		if name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
			return
		}
		if err := mcpClientService.DeleteClient(name); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusNoContent)
	}
}
