package dtos

import "time"

type ClientRequest struct {
	Name          string `json:"name" validate:"required"`
	Email         string `json:"email" validate:"required,email"`
	PhoneNumber   string `json:"phone_number" validate:"required"`
	Address       string `json:"address" validate:"required"`
	ContactPerson string `json:"contact_person" validate:"required"`
}

type ClientResponse struct {
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	PhoneNumber   string    `json:"phone_number"`
	Address       string    `json:"address"`
	ContactPerson string    `json:"contact_person"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
