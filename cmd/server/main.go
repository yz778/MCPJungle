package main

import (
	"log"
	"os"

	"github.com/duaraghav8/mcpjungle/internal/api"
	"github.com/duaraghav8/mcpjungle/internal/db"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	db.Init()

	router := gin.Default()

	router.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	router.POST("/tools", api.RegisterToolHandler)
	router.GET("/tools", api.ListToolsHandler)
	router.DELETE("/tools/:name", api.DeleteToolHandler)
	router.POST("/invoke/:name", api.Invoke)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server listening on :%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("failed to run: %v", err)
	}
}
