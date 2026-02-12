package services

import (
	"github.com/tech-azim/be-learnova/models"
	"github.com/tech-azim/be-learnova/repositories"
	"github.com/tech-azim/be-learnova/utils"
)

type GalleryService interface {
	Create(gallery models.Gallery) (models.Gallery, error)
	FindAll(params utils.PaginationParams) ([]models.Gallery, int64, error)
	FindByID(id uint) (models.Gallery, error)
	Update(gallery models.Gallery) (models.Gallery, error)
	Delete(id uint) error
	FindAllActive() ([]models.Gallery, error)
}

type galleryService struct {
	galleryRepo repositories.GalleryRepository
}

func NewGalleryService(galleryRepo repositories.GalleryRepository) GalleryService {
	return &galleryService{
		galleryRepo,
	}
}

// Create implements GalleryService.
func (s *galleryService) Create(gallery models.Gallery) (models.Gallery, error) {
	result, err := s.galleryRepo.Create(gallery)

	if err != nil {
		return models.Gallery{}, err
	}

	return result, nil
}

// FindAll implements GalleryService.
func (s *galleryService) FindAll(params utils.PaginationParams) ([]models.Gallery, int64, error) {
	data, total, err := s.galleryRepo.FindAll(params)

	if err != nil {
		return []models.Gallery{}, 0, err
	}

	return data, total, nil
}

// FindByID implements GalleryService.
func (s *galleryService) FindByID(id uint) (models.Gallery, error) {
	data, err := s.galleryRepo.FindByID(id)

	if err != nil {
		return models.Gallery{}, err
	}

	return data, nil
}

// Update implements GalleryService.
func (s *galleryService) Update(gallery models.Gallery) (models.Gallery, error) {
	data, err := s.galleryRepo.Update(gallery)

	if err != nil {
		return models.Gallery{}, err
	}

	return data, nil
}

// Delete implements GalleryService.
func (s *galleryService) Delete(id uint) error {
	err := s.galleryRepo.Delete(id)

	if err != nil {
		return err
	}

	return nil
}

// FindAllActive implements GalleryService.
func (s *galleryService) FindAllActive() ([]models.Gallery, error) {
	data, err := s.galleryRepo.FindAllActive()

	if err != nil {
		return []models.Gallery{}, err
	}

	return data, nil
}