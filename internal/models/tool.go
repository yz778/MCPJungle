package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ToolType string

const (
	MCPServer   ToolType = "mcp_server"
	RESTAPI     ToolType = "rest_api"
	ExternalAPI ToolType = "external_api"
)

type Status string

const (
	Verified   Status = "verified"
	Unverified Status = "unverified"
)

type Tool struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Name        string         `gorm:"uniqueIndex;not null" json:"name"`
	URL         string         `gorm:"not null" json:"url"`
	Type        ToolType       `gorm:"not null" json:"type"`
	Description string         `json:"description"`
	Tags        pq.StringArray `gorm:"type:text[]" json:"tags"`
	Metadata    datatypes.JSON `json:"metadata"`
	Status      Status         `gorm:"not null" json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
}

func (t *Tool) BeforeCreate(tx *gorm.DB) (err error) {
	t.ID = uuid.New()
	t.CreatedAt = time.Now().UTC()
	return nil
}
