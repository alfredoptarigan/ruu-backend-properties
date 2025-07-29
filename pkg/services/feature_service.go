package services

import (
	"alfredo/ruu-properties/pkg/dtos"
	"alfredo/ruu-properties/pkg/repositories"
)

type FeatureService interface {
	Create(request dtos.FeatureRequest) (*dtos.FeatureResponse, error)
}

type featureServiceImpl struct {
	repo repositories.FeatureRepository
}

// Create implements FeatureService.
func (f *featureServiceImpl) Create(request dtos.FeatureRequest) (*dtos.FeatureResponse, error) {
	return f.repo.Create(request)
}

func NewFeatureService(repo repositories.FeatureRepository) FeatureService {
	return &featureServiceImpl{repo: repo}
}
