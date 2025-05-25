package migrations

import (
	"fmt"
	"github.com/duaraghav8/mcpjungle/internal/db"
	"github.com/duaraghav8/mcpjungle/internal/model"
)

// Migrate performs the database migration for the application.
func Migrate() error {
	if err := db.DB.AutoMigrate(&model.McpServer{}); err != nil {
		return fmt.Errorf("auto‑migration failed for McpServer model: %v", err)
	}
	if err := db.DB.AutoMigrate(&model.Tool{}); err != nil {
		return fmt.Errorf("auto‑migration failed for Tool model: %v", err)
	}
	return nil
}
