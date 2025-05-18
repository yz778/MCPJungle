package db

import (
	"log"
	"os"

	"github.com/duaraghav8/mcpjungle/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Init initialises the global DB connection using DATABASE_URL env var. Falls back to sqlite if missing.
func Init() {
	dsn := os.Getenv("DATABASE_URL")
	var dialector gorm.Dialector
	if dsn == "" {
		log.Println("[db] DATABASE_URL not set – falling back to embedded SQLite ./mcp.db")
		dialector = sqlite.Open("mcp.db?_busy_timeout=5000&_journal_mode=WAL")
	} else {
		dialector = postgres.Open(dsn)
	}

	var err error
	DB, err = gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Auto‑migrate models
	if err := DB.AutoMigrate(&models.Tool{}); err != nil {
		log.Fatalf("auto‑migration failed: %v", err)
	}
}
