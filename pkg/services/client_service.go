package services

import (
	"alfredo/ruu-properties/pkg/dtos"
	"alfredo/ruu-properties/pkg/repositories"
)

type ClientService interface {
	Create(request dtos.ClientRequest) (*dtos.ClientResponse, error)
	GetAll(request dtos.ClientGetRequest) ([]*dtos.ClientResponse, *dtos.PaginationMeta, error)
	GetByID(uuid string) (*dtos.ClientResponse, error)
	Delete(uuid string) error
	Update(request dtos.ClientUpdateRequest) (*dtos.ClientResponse, error)
}

type clientServiceImpl struct {
	clientRepository repositories.ClientRepository
}

func (s *clientServiceImpl) Update(request dtos.ClientUpdateRequest) (*dtos.ClientResponse, error) {
	return s.clientRepository.Update(request)
}

// GetByID implements ClientService.
func (s *clientServiceImpl) GetByID(uuid string) (*dtos.ClientResponse, error) {
	return s.clientRepository.GetByID(uuid)
}

// GetAll implements ClientService.
func (s *clientServiceImpl) GetAll(request dtos.ClientGetRequest) ([]*dtos.ClientResponse, *dtos.PaginationMeta, error) {
	return s.clientRepository.GetAll(request)
}

func NewClientService(clientRepository repositories.ClientRepository) ClientService {
	return &clientServiceImpl{clientRepository: clientRepository}
}

func (s *clientServiceImpl) Create(request dtos.ClientRequest) (*dtos.ClientResponse, error) {
	return s.clientRepository.Create(request)
}

func (s *clientServiceImpl) Delete(uuid string) error {
	return s.clientRepository.Delete(uuid)
}
