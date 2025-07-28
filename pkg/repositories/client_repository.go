package repositories

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"alfredo/ruu-properties/pkg/dtos"
	"alfredo/ruu-properties/pkg/models"
)

type ClientRepository interface {
	Create(request dtos.ClientRequest) (*dtos.ClientResponse, error)
	GetAll(request dtos.ClientGetRequest) ([]*dtos.ClientResponse, *dtos.PaginationMeta, error)
	Delete(uuid string) error
	GetByID(uuid string) (*dtos.ClientResponse, error)
	Update(request dtos.ClientUpdateRequest) (*dtos.ClientResponse, error)
}

type clientRepositoryImpl struct {
	db *gorm.DB
}

func (r *clientRepositoryImpl) Update(request dtos.ClientUpdateRequest) (*dtos.ClientResponse, error) {
	var client models.Client
	if err := r.db.Where("uuid = ?", request.UUID).First(&client).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &dtos.ClientResponse{}, fmt.Errorf("%s", "client not found")
		} else {
			return &dtos.ClientResponse{}, fmt.Errorf("%s", err.Error())
		}
	}

	if request.Name != "" {
		client.Name = request.Name
	}
	if request.Email != "" {
		client.Email = request.Email
	}
	if request.PhoneNumber != "" {
		client.PhoneNumber = request.PhoneNumber
	}
	if request.Address != "" {
		client.Address = request.Address
	}
	if request.ContactPerson != "" {
		client.ContactPerson = request.ContactPerson
	}

	if err := r.db.Save(&client).Error; err != nil {
		return &dtos.ClientResponse{}, fmt.Errorf("%s", "please try again later")
	}

	return &dtos.ClientResponse{
		UUID:          client.UUID,
		Name:          client.Name,
		Email:         client.Email,
		PhoneNumber:   client.PhoneNumber,
		Address:       client.Address,
		ContactPerson: client.ContactPerson,
		CreatedAt:     client.CreatedAt,
		UpdatedAt:     client.UpdatedAt,
	}, nil
}

func (r *clientRepositoryImpl) Delete(uuid string) error {
	var client models.Client
	if err := r.db.Where("uuid = ?", uuid).First(&client).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("%s", "client not found")
		} else {
			return fmt.Errorf("%s", "please try again later")
		}
	}
	if err := r.db.Delete(&client).Error; err != nil {
		return fmt.Errorf("%s", "please try again later")
	}
	return nil
}

// GetByID implements ClientRepository.
func (r *clientRepositoryImpl) GetByID(uuid string) (*dtos.ClientResponse, error) {
	var client models.Client

	if err := r.db.Model(&models.Client{}).Where("uuid = ? ", uuid).First(&client).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &dtos.ClientResponse{}, fmt.Errorf("%s", "client not found")
		} else {
			return &dtos.ClientResponse{}, fmt.Errorf("%s", "please try again later")
		}
	}

	return &dtos.ClientResponse{
		UUID:          client.UUID,
		Name:          client.Name,
		Email:         client.Email,
		PhoneNumber:   client.PhoneNumber,
		Address:       client.Address,
		ContactPerson: client.ContactPerson,
		CreatedAt:     client.CreatedAt,
		UpdatedAt:     client.UpdatedAt,
	}, nil
}

// GetAll implements ClientRepository.
func (r *clientRepositoryImpl) GetAll(request dtos.ClientGetRequest) ([]*dtos.ClientResponse, *dtos.PaginationMeta, error) {
	if request.SortBy == "" {
		request.SortBy = "created_at"
	}
	if request.SortOrder == "" {
		request.SortOrder = "desc"
	}

	// Validate sort from frontend if the sort is not asc or desc
	if request.SortOrder != "asc" && request.SortOrder != "desc" {
		request.SortOrder = "desc"
	}

	// Validate sort field (turn on or disabled)
	allowedSortFields := map[string]bool{
		"name":       true,
		"email":      true,
		"created_at": true,
		"updated_at": true,
	}
	if !allowedSortFields[request.SortBy] {
		request.SortBy = "created_at"
	}

	var clients []models.Client
	var total int64

	query := r.db.Model(&models.Client{})

	// Apply search filter if not null
	if request.Search != "" {
		searchPattern := "%" + request.Search + "%"
		switch request.SearchBy {
		case "name":
			query = query.Where("name ILIKE ?", searchPattern)
		case "email":
			query = query.Where("email ILIKE ?", searchPattern)
		case "phone_number":
			query = query.Where("phone_number ILIKE ?", searchPattern)
		case "contact_person":
			query = query.Where("contact_person ILIKE ?", searchPattern)
		default:
			// Global search across multiple fields
			query = query.Where(
				"name ILIKE ? OR email ILIKE ? OR phone_number ILIKE ? OR contact_person ILIKE ? OR address ILIKE ?",
				searchPattern, searchPattern, searchPattern, searchPattern, searchPattern,
			)
		}
	}

	// Count total records
	err := query.Count(&total).Error
	if err != nil {
		return nil, nil, fmt.Errorf("failed to count clients: %w", err)
	}

	// Calculate offset
	offset := (request.Page - 1) * request.Limit

	// Apply pagination and sorting
	sortClause := fmt.Sprintf("%s %s", request.SortBy, request.SortOrder)
	err = query.Order(sortClause).Offset(offset).Limit(request.Limit).Find(&clients).Error
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch clients: %w", err)
	}

	// Convert to response DTOs
	clientResponses := make([]*dtos.ClientResponse, len(clients))
	for i, client := range clients {
		clientResponses[i] = &dtos.ClientResponse{
			UUID:          client.UUID,
			Name:          client.Name,
			Email:         client.Email,
			PhoneNumber:   client.PhoneNumber,
			Address:       client.Address,
			ContactPerson: client.ContactPerson,
			CreatedAt:     client.CreatedAt,
			UpdatedAt:     client.UpdatedAt,
		}
	}

	// Calculate pagination metadata
	totalPages := int((total + int64(request.Page) - 1) / int64(request.Limit))
	paginationMeta := &dtos.PaginationMeta{
		Page:       request.Page,
		Limit:      request.Limit,
		Total:      int(total),
		TotalPages: totalPages,
	}

	return clientResponses, paginationMeta, nil

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
		UUID:          client.UUID,
		Name:          client.Name,
		Email:         client.Email,
		PhoneNumber:   client.PhoneNumber,
		Address:       client.Address,
		ContactPerson: client.ContactPerson,
		CreatedAt:     client.CreatedAt,
		UpdatedAt:     client.UpdatedAt,
	}, nil
}
