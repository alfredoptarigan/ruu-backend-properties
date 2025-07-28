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
	UUID          string    `json:"uuid"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	PhoneNumber   string    `json:"phone_number"`
	Address       string    `json:"address"`
	ContactPerson string    `json:"contact_person"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type ClientGetRequest struct {
	Page      int    `json:"page" query:"page" default:"1"`
	Limit     int    `json:"limit" query:"limit" default:"10"`
	Search    string `json:"search" query:"search"`
	SearchBy  string `json:"search_by" query:"search_by"`
	SortBy    string `json:"sort_by" query:"sort_by" default:"created_at"`
	SortOrder string `json:"sort_order" query:"sort_order" default:"desc"`
}

type ClientUpdateRequest struct {
	UUID          string
	Name          string `json:"name" validate:"omitempty"`
	Email         string `json:"email" validate:"omitempty,email"`
	PhoneNumber   string `json:"phone_number" validate:"omitempty"`
	Address       string `json:"address" validate:"omitempty"`
	ContactPerson string `json:"contact_person" validate:"omitempty"`
}
