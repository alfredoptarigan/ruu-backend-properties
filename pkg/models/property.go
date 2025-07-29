package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Property struct {
	UUID        uuid.UUID      `gorm:"primaryKey;default:uuid_generate_v4()"`
	Name        string         `gorm:"not null" json:"name"`
	Description string         `gorm:"not null" json:"description"`
	CreatedAt   time.Time      `gorm:"default:now()" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"default:now()" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (p *Property) TableName() string {
	return "properties"
}
