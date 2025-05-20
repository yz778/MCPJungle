package service

import (
	"github.com/duaraghav8/mcpjungle/internal/db"
	"github.com/duaraghav8/mcpjungle/internal/models"
)

func RegisterMcpServer(s *models.McpServer) error {
	// validate server name
}

func DeregisterMcpServer(name string) error {}

func ListMcpServers() ([]models.McpServer, error) {}

// GetMcpServer fetches a server from the database by name.
func GetMcpServer(name string) (*models.McpServer, error) {
	var server models.McpServer
	if err := db.DB.Where("name = ?", name).First(&server).Error; err != nil {
		return nil, err
	}
	return &server, nil
}
