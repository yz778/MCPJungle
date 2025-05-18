package server

import (
	"log"
	"os"

	"github.com/duaraghav8/mcpjungle/internal/api"
	"github.com/duaraghav8/mcpjungle/internal/db"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// Start spins up the registry HTTP server (blocking call).
func Start(port string) {
	_ = godotenv.Load()
	db.Init()

	r := gin.Default()
	r.GET("/healthcheck", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })
	r.POST("/tools", api.RegisterToolHandler)
	r.GET("/tools", api.ListToolsHandler)
	r.DELETE("/tools/:name", api.DeleteToolHandler)
	r.POST("/invoke/:name", api.Invoke)

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
