package repositories

import (
	"fmt"

	"gorm.io/gorm"

	"alfredo/ruu-properties/pkg/dtos"
	"alfredo/ruu-properties/pkg/models"
)

type ClientRepository interface {
	Create(request dtos.ClientRequest) (*dtos.ClientResponse, error)
}

type clientRepositoryImpl struct {
	db *gorm.DB
}

func NewClientRepository(db *gorm.DB) ClientRepository {
	return &clientRepositoryImpl{db: db}
}

func (r *clientRepositoryImpl) Create(request dtos.ClientRequest) (*dtos.ClientResponse, error) {
	// Check if client already exists by email
	var existingClient models.Client
	err := r.db.Debug().Where("email = ?", request.Email).First(&existingClient).Error
	if err == nil {
		return nil, fmt.Errorf("client with email %s already exists", request.Email)
	}
	if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// Create new client
	client := models.Client{
		Name:          request.Name,
		Email:         request.Email,
		PhoneNumber:   request.PhoneNumber,
		Address:       request.Address,
		ContactPerson: request.ContactPerson,
	}

	err = r.db.Create(&client).Error
	if err != nil {
		return nil, err
	}

	return &dtos.ClientResponse{
		Name:          client.Name,
		Email:         client.Email,
		PhoneNumber:   client.PhoneNumber,
		Address:       client.Address,
		ContactPerson: client.ContactPerson,
		CreatedAt:     client.CreatedAt,
		UpdatedAt:     client.UpdatedAt,
	}, nil
}
