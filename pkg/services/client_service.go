package services

import (
	"alfredo/ruu-properties/pkg/dtos"
	"alfredo/ruu-properties/pkg/repositories"
)

type ClientService interface {
	Create(request dtos.ClientRequest) (*dtos.ClientResponse, error)
}

type clientServiceImpl struct {
	clientRepository repositories.ClientRepository
}

func NewClientService(clientRepository repositories.ClientRepository) ClientService {
	return &clientServiceImpl{clientRepository: clientRepository}
}

func (s *clientServiceImpl) Create(request dtos.ClientRequest) (*dtos.ClientResponse, error) {
	return s.clientRepository.Create(request)
}
