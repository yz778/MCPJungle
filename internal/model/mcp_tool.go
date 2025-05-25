package model

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Tool struct {
	ID          uuid.UUID      `json:"-" gorm:"type:uuid;primaryKey"`
	Name        string         `json:"name" gorm:"uniqueIndex;not null"`
	Description string         `json:"description"`
	InputSchema datatypes.JSON `json:"input_schema" gorm:"type:jsonb"`

	ServerID uuid.UUID `json:"-" gorm:"type:uuid"`
	Server   McpServer `json:"-" gorm:"foreignKey:ServerID;references:ID"`
}

func (t *Tool) BeforeCreate(tx *gorm.DB) (err error) {
	t.ID = uuid.New()
	return nil
}
