package services

import (
	"github.com/tech-azim/be-learnova/models"
	"github.com/tech-azim/be-learnova/repositories"
	"github.com/tech-azim/be-learnova/utils"
)

type FlyerGalleryService interface {
	Create(flyerGallery models.FlyerGallery) (models.FlyerGallery, error)
	FindAll(params utils.PaginationParams) ([]models.FlyerGallery, int64, error)
	FindByID(id uint) (models.FlyerGallery, error)
	Update(flyerGallery models.FlyerGallery) (models.FlyerGallery, error)
	Delete(id uint) error
	FindAllActive() ([]models.FlyerGallery, error)
}

type flyerGalleryService struct {
	flyerGalleryRepo repositories.FlyerGalleryRepository
}

func NewFlyerGalleryService(flyerGalleryRepo repositories.FlyerGalleryRepository) FlyerGalleryService {
	return &flyerGalleryService{
		flyerGalleryRepo,
	}
}

// Create implements FlyerGalleryService.
func (s *flyerGalleryService) Create(flyerGallery models.FlyerGallery) (models.FlyerGallery, error) {
	result, err := s.flyerGalleryRepo.Create(flyerGallery)

	if err != nil {
		return models.FlyerGallery{}, err
	}

	return result, nil
}

// FindAll implements FlyerGalleryService.
func (s *flyerGalleryService) FindAll(params utils.PaginationParams) ([]models.FlyerGallery, int64, error) {
	data, total, err := s.flyerGalleryRepo.FindAll(params)

	if err != nil {
		return []models.FlyerGallery{}, 0, err
	}

	return data, total, nil
}

// FindByID implements FlyerGalleryService.
func (s *flyerGalleryService) FindByID(id uint) (models.FlyerGallery, error) {
	data, err := s.flyerGalleryRepo.FindByID(id)

	if err != nil {
		return models.FlyerGallery{}, err
	}

	return data, nil
}

// Update implements FlyerGalleryService.
func (s *flyerGalleryService) Update(flyerGallery models.FlyerGallery) (models.FlyerGallery, error) {
	data, err := s.flyerGalleryRepo.Update(flyerGallery)

	if err != nil {
		return models.FlyerGallery{}, err
	}

	return data, nil
}

// Delete implements FlyerGalleryService.
func (s *flyerGalleryService) Delete(id uint) error {
	err := s.flyerGalleryRepo.Delete(id)

	if err != nil {
		return err
	}

	return nil
}

// FindAllActive implements FlyerGalleryService.
func (s *flyerGalleryService) FindAllActive() ([]models.FlyerGallery, error) {
	data, err := s.flyerGalleryRepo.FindAllActive()

	if err != nil {
		return []models.FlyerGallery{}, err
	}

	return data, nil
}