package model

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Tool struct {
	gorm.Model

	Name        string         `json:"name" gorm:"uniqueIndex;not null"`
	Description string         `json:"description"`
	InputSchema datatypes.JSON `json:"input_schema" gorm:"type:jsonb"`

	ServerID uint      `json:"-" gorm:"not null"`
	Server   McpServer `json:"-" gorm:"foreignKey:ServerID;references:ID"`
}
