package services

import (
	"github.com/tech-azim/be-learnova/models"
	"github.com/tech-azim/be-learnova/repositories"
	"github.com/tech-azim/be-learnova/utils"
)

type ServiceService interface {
	Create(service models.Service) (models.Service, error)
	FindAll(params utils.PaginationParams) ([]models.Service, int64, error)
	FindByID(id uint) (models.Service, error)
	Update(service models.Service) (models.Service, error)
	Delete(id uint) error
}

type serviceService struct {
	serviceRepo repositories.ServiceRepository
}

func NewServiceService(serviceRepo repositories.ServiceRepository) ServiceService {
	return &serviceService{
		serviceRepo,
	}
}

// Create implements ServiceService.
func (s *serviceService) Create(service models.Service) (models.Service, error) {
	result, err := s.serviceRepo.Create(service)

	if err != nil {
		return models.Service{}, err
	}

	return result, nil
}

// FindAll implements ServiceService.
func (s *serviceService) FindAll(params utils.PaginationParams) ([]models.Service, int64, error) {
	data, total, err := s.serviceRepo.FindAll(params)

	if err != nil {
		return []models.Service{}, 0, err
	}

	return data, total, nil
}

// FindByID implements ServiceService.
func (s *serviceService) FindByID(id uint) (models.Service, error) {
	data, err := s.serviceRepo.FindByID(id)

	if err != nil {
		return models.Service{}, err
	}

	return data, nil
}

// Update implements ServiceService.
func (s *serviceService) Update(service models.Service) (models.Service, error) {
	data, err := s.serviceRepo.Update(service)

	if err != nil {
		return models.Service{}, err
	}

	return data, nil
}

// Delete implements ServiceService.
func (s *serviceService) Delete(id uint) error {
	err := s.serviceRepo.Delete(id)

	if err != nil {
		return err
	}

	return nil
}