package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/duaraghav8/mcpjungle/internal/db"
	"github.com/duaraghav8/mcpjungle/internal/models"
	"gorm.io/datatypes"
)

var httpClient = &http.Client{Timeout: 10 * time.Second}

// RegisterTool reads .well-known/mcp metadata, validates minimal MCP fields, and stores the record.
func RegisterTool(tool *models.Tool) error {
	metaURL := fmt.Sprintf("%s/.well-known/mcp", tool.URL)
	resp, err := httpClient.Get(metaURL)
	if err != nil {
		return fmt.Errorf("fetching metadata: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status fetching metadata: %d", resp.StatusCode)
	}
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	// Minimal validation – check required keys.
	var meta map[string]any
	if err := json.Unmarshal(raw, &meta); err != nil {
		return fmt.Errorf("invalid json: %w", err)
	}
	if _, ok := meta["schema_version"]; !ok {
		return errors.New("missing required field schema_version in MCP metadata")
	}
	tool.Metadata = datatypes.JSON(raw)
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
