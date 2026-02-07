package services

import (
	"github.com/tech-azim/be-learnova/models"
	"github.com/tech-azim/be-learnova/repositories"
	"github.com/tech-azim/be-learnova/utils"
)

type RegistrationService interface {
	Create(registration models.Registration) (models.Registration, error)
	FindAll(params utils.PaginationParams) ([]models.Registration, int64, error)
	FindByID(id uint) (models.Registration, error)
	FindByProgramID(programID uint, params utils.PaginationParams) ([]models.Registration, int64, error)
	FindByEmail(email string) (models.Registration, error)
	Update(registration models.Registration) (models.Registration, error)
	Delete(id uint) error
	CheckEmailExists(email string, programID uint) (bool, error)
}

type registrationService struct {
	registrationRepo repositories.RegistrationRepository
}

func NewRegistrationService(registrationRepo repositories.RegistrationRepository) RegistrationService {
	return &registrationService{
		registrationRepo,
	}
}

// Create implements RegistrationService.
func (s *registrationService) Create(registration models.Registration) (models.Registration, error) {
	result, err := s.registrationRepo.Create(registration)

	if err != nil {
		return models.Registration{}, err
	}

	return result, nil
}

// FindAll implements RegistrationService.
func (s *registrationService) FindAll(params utils.PaginationParams) ([]models.Registration, int64, error) {
	data, total, err := s.registrationRepo.FindAll(params)

	if err != nil {
		return []models.Registration{}, 0, err
	}

	return data, total, nil
}

// FindByID implements RegistrationService.
func (s *registrationService) FindByID(id uint) (models.Registration, error) {
	data, err := s.registrationRepo.FindByID(id)

	if err != nil {
		return models.Registration{}, err
	}

	return data, nil
}

// FindByProgramID implements RegistrationService.
func (s *registrationService) FindByProgramID(programID uint, params utils.PaginationParams) ([]models.Registration, int64, error) {
	data, total, err := s.registrationRepo.FindByProgramID(programID, params)

	if err != nil {
		return []models.Registration{}, 0, err
	}

	return data, total, nil
}

// FindByEmail implements RegistrationService.
func (s *registrationService) FindByEmail(email string) (models.Registration, error) {
	data, err := s.registrationRepo.FindByEmail(email)

	if err != nil {
		return models.Registration{}, err
	}

	return data, nil
}

// Update implements RegistrationService.
func (s *registrationService) Update(registration models.Registration) (models.Registration, error) {
	data, err := s.registrationRepo.Update(registration)

	if err != nil {
		return models.Registration{}, err
	}

	return data, nil
}

// Delete implements RegistrationService.
func (s *registrationService) Delete(id uint) error {
	err := s.registrationRepo.Delete(id)

	if err != nil {
		return err
	}

	return nil
}

// CheckEmailExists implements RegistrationService.
func (s *registrationService) CheckEmailExists(email string, programID uint) (bool, error) {
	exists, err := s.registrationRepo.CheckEmailExists(email, programID)

	if err != nil {
		return false, err
	}

	return exists, nil
}
