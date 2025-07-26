package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	UUID        string         `json:"uuid" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Email       string         `json:"email" gorm:"type:varchar(255);uniqueIndex"`
	Password    string         `json:"password" gorm:"type:varchar(255)"`
	Name        string         `json:"name" gorm:"type:varchar(255)"`
	PhoneNumber string         `json:"phone_number" gorm:"type:varchar(255)"`
	Image       string         `json:"image" gorm:"type:varchar(255)"`
	Role        string         `json:"role" gorm:"type:varchar(50);default:'user'"`
	CreatedAt   time.Time      `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

func (u *User) TableName() string {
	return "users"
}
