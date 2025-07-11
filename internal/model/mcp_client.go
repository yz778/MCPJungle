package model

import "gorm.io/gorm"

// McpClient represents MCP clients and their access to the MCP Servers provided MCPJungle MCP server
type McpClient struct {
	gorm.Model

	Name        string `json:"name" gorm:"uniqueIndex;not null"`
	Description string `json:"description"`

	AccessToken string `json:"access_token" gorm:"unique; not null"`

	// AllowList contains a list of MCP Server names that this client is allowed to view and call
	AllowList []string `json:"allow_list" gorm:"type:text[]"`
}

// CheckHasServerAccess returns true if this client has access to the specified MCP server.
// If not, it returns false.
func (c *McpClient) CheckHasServerAccess(serverName string) bool {
	for _, allowed := range c.AllowList {
		if allowed == serverName {
			return true
		}
	}
	return false
}
