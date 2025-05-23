package db

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Init initialises the global DB connection.
func Init(dsn string) {
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
}
