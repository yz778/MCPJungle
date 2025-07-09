package api

import (
	"github.com/gin-gonic/gin"
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

func deleteMcpClientHandler(mcpClientService *mcp_client.McpClientService) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("name")
		if name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
			return
		}
		if err := mcpClientService.DeleteMcpClient(name); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusNoContent)
	}
}
