package models

import (
	"time"

	"gorm.io/gorm"
)

type Feature struct {
	// ✅ PERBAIKAN: Tambahkan tag untuk auto-generate UUID
	UUID        string         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"uuid"`
	Name        string         `gorm:"not null" json:"name"`
	Description string         `json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at"`
}

func (f *Feature) TableName() string {
	return "features"
}
