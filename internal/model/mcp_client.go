package model

import "gorm.io/gorm"

// McpClient represents MCP clients and their access to the MCP Servers provided MCPJungle MCP server
type McpClient struct {
	gorm.Model

	Name        string `json:"name" gorm:"uniqueIndex;not null"`
	Description string `json:"description"`

	// AllowList contains a list of MCP Server names that this client is allowed to view and call
	AllowList []string `json:"allow_list" gorm:"type:text[];default:array[]::text[]"`
}
