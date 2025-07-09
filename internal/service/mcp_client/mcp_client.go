package mcp_client

import (
	"fmt"
	"github.com/mcpjungle/mcpjungle/internal"
	"github.com/mcpjungle/mcpjungle/internal/model"
	"gorm.io/gorm"
)

// McpClientService provides methods to manage MCP clients in the database.
type McpClientService struct {
	db *gorm.DB
}

func NewMCPClientService(db *gorm.DB) *McpClientService {
	return &McpClientService{db: db}
}

// ListClients retrieves all MCP clients known to mcpjungle from the database
func (m *McpClientService) ListClients() ([]*model.McpClient, error) {
	var clients []*model.McpClient
	if err := m.db.Find(&clients).Error; err != nil {
		return nil, err
	}
	return clients, nil
}

// CreateMcpClient creates a new MCP client in the database and returns it.
func (m *McpClientService) CreateMcpClient(client model.McpClient) (*model.McpClient, error) {
	token, err := internal.GenerateAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}
	client.AccessToken = token
	if err := m.db.Create(&client).Error; err != nil {
		return nil, err
	}
	return &client, nil
}

// DeleteMcpClient removes an MCP client from the database and immediately revokes its access.
// It is an idempotent operation. Deleting a client that does not exist will not return an error.
func (m *McpClientService) DeleteMcpClient(name string) error {
	result := m.db.Where("name = ?", name).Delete(&model.McpClient{})
	return result.Error
}
