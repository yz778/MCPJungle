package server

import (
	"log"
	"os"

	"github.com/duaraghav8/mcpjungle/internal/api"
	"github.com/duaraghav8/mcpjungle/internal/db"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

const ApiV0PathPrefix = "/api/v0"

// Start spins up the registry HTTP server (blocking call).
func Start(port string) {
	_ = godotenv.Load()
	db.Init()

	r := gin.Default()
	r.GET("/healthcheck", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })

	apiV0 := r.Group(ApiV0PathPrefix)
	{
		apiV0.POST("/servers", api.RegisterServerHandler)
		apiV0.DELETE("/servers/:name", api.DeregisterServerHandler)
		apiV0.GET("/servers", api.ListServersHandler)
		apiV0.GET("/tools", api.ListToolsHandler)
		apiV0.POST("/tools/invoke", api.InvokeToolHandler)
	}

	if port == "" {
		port = os.Getenv("PORT")
	}
	if port == "" {
		port = "8080"
	}
	log.Printf("MCP registry listening on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("failed to run: %v", err)
	}
}
