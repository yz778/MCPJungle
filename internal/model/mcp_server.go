package model

import "gorm.io/gorm"

type McpServer struct {
	gorm.Model

	Name        string `json:"name" gorm:"uniqueIndex;not null"`
	Description string `json:"description"`

	// URL must be a valid http/https URL.
	// MCPJungle only supports streamable HTTP transport as of now.
	URL string `json:"url" gorm:"not null"`

	// TODO: Store the bearer token in a more secure way, e.g., encrypted in the database.
	// BearerToken is an optional token used for authenticating requests to the MCP server.
	// If present, it will be used to set the Authorization header in all requests to this MCP server.
	BearerToken string `json:"bearer_token,omitempty" gorm:"type:text"`
}
