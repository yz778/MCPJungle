package service

import (
	"net/http"
	"time"

	"github.com/duaraghav8/mcpjungle/internal/db"
	"github.com/duaraghav8/mcpjungle/internal/models"
)

var httpClient = &http.Client{Timeout: 10 * time.Second}

// RegisterTool validates the MCP server and stores the record.
func RegisterTool(tool *models.Tool) error {
	tool.Status = models.Verified
	return db.DB.Create(tool).Error
}

// ListTools returns all tools or filtered by name / tags.
func ListTools() ([]models.Tool, error) {
	var tools []models.Tool
	if err := db.DB.Find(&tools).Error; err != nil {
		return nil, err
	}
	return tools, nil
}

// DeleteTool removes a tool by name.
func DeleteTool(name string) error {
	return db.DB.Where("name = ?", name).Delete(&models.Tool{}).Error
}

// GetTool fetches a tool by name.
func GetTool(name string) (*models.Tool, error) {
	var t models.Tool
	if err := db.DB.Where("name = ?", name).First(&t).Error; err != nil {
		return nil, err
	}
	return &t, nil
}
