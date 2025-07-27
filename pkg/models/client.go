package models

import (
	"time"

	"gorm.io/gorm"
)

type Client struct {
	UUID          string         `json:"uuid" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name          string         `json:"name" gorm:"column:name;not null"`
	Email         string         `json:"email" gorm:"column:email;not null;uniqueIndex"`
	PhoneNumber   string         `json:"phone_number" gorm:"column:phone_number;not null;uniqueIndex"`
	Address       string         `json:"address" gorm:"column:address;not null"`
	ContactPerson string         `json:"contact_person" gorm:"column:contact_person;not null"`
	CreatedAt     time.Time      `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt     time.Time      `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at;index"`
}

func (c *Client) TableName() string {
	return "clients"
}
