package server

import (
	"github.com/duaraghav8/mcpjungle/internal/api"
	"github.com/gin-gonic/gin"
)

const ApiV0PathPrefix = "/api/v0"

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET(
		"/health",
		func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		},
	)

	apiV0 := r.Group(ApiV0PathPrefix)
	{
		apiV0.POST("/servers", api.RegisterServerHandler)
		apiV0.DELETE("/servers/:name", api.DeregisterServerHandler)
		apiV0.GET("/servers", api.ListServersHandler)
		apiV0.GET("/tools", api.ListToolsHandler)
		apiV0.POST("/tools/invoke", api.InvokeToolHandler)
		apiV0.GET("/tool", api.GetToolHandler)
	}

	return r
}
