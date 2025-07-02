package api

import (
	"github.com/gin-gonic/gin"
	"github.com/mcpjungle/mcpjungle/internal/model"
	"github.com/mcpjungle/mcpjungle/internal/service/config"
)

func registerInitServerHandler(configService *config.ServerConfigService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Mode string `json:"mode" binding:"required,oneof=development production"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "Invalid request body: " + err.Error()})
			return
		}
		if err := configService.Init(model.ServerMode(req.Mode)); err != nil {
			c.JSON(500, gin.H{"error": "Failed to initialize server: " + err.Error()})
			return
		}
		c.JSON(200, gin.H{"status": "Server initialized successfully", "mode": req.Mode})
	}
}
