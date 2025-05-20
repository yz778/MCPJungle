package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type McpServer struct {
	ID          uuid.UUID `json:"-" gorm:"type:uuid;primaryKey"`
	Name        string    `json:"name" gorm:"uniqueIndex;not null"`
	URL         string    `json:"url" gorm:"not null"`
	Description string    `json:"description"`
}

func (s *McpServer) BeforeCreate(tx *gorm.DB) (err error) {
	s.ID = uuid.New()
	return nil
}
