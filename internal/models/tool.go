package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Tool struct {
	ID          uuid.UUID `json:"-" gorm:"type:uuid;primaryKey"`
	Name        string    `json:"name" gorm:"uniqueIndex;not null"`
	Description string    `json:"description"`

	ServerID uuid.UUID `json:"-" gorm:"type:uuid"`
	Server   McpServer `json:"-" gorm:"foreignKey:ServerID;references:ID"`
}

func (t *Tool) BeforeCreate(tx *gorm.DB) (err error) {
	t.ID = uuid.New()
	return nil
}
