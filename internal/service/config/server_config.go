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

func (s *ServerConfigService) GetConfig() (model.ServerConfig, error) {
	var config model.ServerConfig
	err := s.db.First(&config).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return model.ServerConfig{Initialized: false}, nil
	}
	if err != nil {
		return model.ServerConfig{}, fmt.Errorf("failed to fetch server configuration from db: %v", err)
	}
	return config, nil
}

func (s *ServerConfigService) Init(mode model.ServerMode) error {
	config, err := s.GetConfig()
	if err != nil {
		return err
	}
	if config.Initialized {
		// Config already exists, do nothing
		return nil
	}
	// No config exists, create one
	config = model.ServerConfig{
		Mode:        mode,
		Initialized: true,
	}
	return s.db.Create(&config).Error
}
