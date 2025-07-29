package repositories

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"alfredo/ruu-properties/pkg/dtos"
	"alfredo/ruu-properties/pkg/models"
)

type FeatureRepository interface {
	Create(request dtos.FeatureRequest) (*dtos.FeatureResponse, error)
}

type featureRepositoryImpl struct {
	db *gorm.DB
}

// Create implements FeatureRepository.
func (f *featureRepositoryImpl) Create(request dtos.FeatureRequest) (*dtos.FeatureResponse, error) {
	feature := models.Feature{
		UUID:        uuid.New().String(),
		Name:        request.Name,
		Description: request.Description,
	}
	if err := f.db.Debug().Create(&feature).Error; err != nil {
		return &dtos.FeatureResponse{}, fmt.Errorf("%s", "please try again later")
	}
	return &dtos.FeatureResponse{
		UUID:        feature.UUID,
		Name:        feature.Name,
		Description: feature.Description,
		CreatedAt:   feature.CreatedAt.String(),
		UpdatedAt:   feature.UpdatedAt.String(),
	}, nil
}

func NewFeatureRepository(db *gorm.DB) FeatureRepository {
	return &featureRepositoryImpl{db: db}
}
