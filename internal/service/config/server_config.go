package config

import (
	"errors"
	"fmt"
	"github.com/mcpjungle/mcpjungle/internal/model"
	"gorm.io/gorm"
)

type ServerConfigService struct {
	db *gorm.DB
}

func NewServerConfigService(db *gorm.DB) *ServerConfigService {
	return &ServerConfigService{db: db}
}

func (s *ServerConfigService) Init(mode model.ServerMode) error {
	var config model.ServerConfig
	err := s.db.First(&config).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// No config exists, create one
		config = model.ServerConfig{
			Mode:        mode,
			Initialized: true,
		}
		return s.db.Create(&config).Error
	}
	// If any other error, return it
	if err != nil {
		return fmt.Errorf("failed to fetch server configuration from db: %v", err)
	}
	// Config already exists, do nothing
	return nil
}
