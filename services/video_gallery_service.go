package services

import (
	"github.com/tech-azim/be-learnova/models"
	"github.com/tech-azim/be-learnova/repositories"
	"github.com/tech-azim/be-learnova/utils"
)

type VideoGalleryService interface {
	Create(videoGallery models.VideoGallery) (models.VideoGallery, error)
	FindAll(params utils.PaginationParams) ([]models.VideoGallery, int64, error)
	FindByID(id uint) (models.VideoGallery, error)
	FindByCategory(category string, params utils.PaginationParams) ([]models.VideoGallery, int64, error)
	Update(videoGallery models.VideoGallery) (models.VideoGallery, error)
	Delete(id uint) error
	FindAllActive() ([]models.VideoGallery, error)
	FindAllCategories() ([]string, error)
}

type videoGalleryService struct {
	videoGalleryRepo repositories.VideoGalleryRepository
}

func NewVideoGalleryService(videoGalleryRepo repositories.VideoGalleryRepository) VideoGalleryService {
	return &videoGalleryService{
		videoGalleryRepo,
	}
}

// Create implements VideoGalleryService.
func (s *videoGalleryService) Create(videoGallery models.VideoGallery) (models.VideoGallery, error) {
	result, err := s.videoGalleryRepo.Create(videoGallery)

	if err != nil {
		return models.VideoGallery{}, err
	}

	return result, nil
}

// FindAll implements VideoGalleryService.
func (s *videoGalleryService) FindAll(params utils.PaginationParams) ([]models.VideoGallery, int64, error) {
	data, total, err := s.videoGalleryRepo.FindAll(params)

	if err != nil {
		return []models.VideoGallery{}, 0, err
	}

	return data, total, nil
}

// FindByID implements VideoGalleryService.
func (s *videoGalleryService) FindByID(id uint) (models.VideoGallery, error) {
	data, err := s.videoGalleryRepo.FindByID(id)

	if err != nil {
		return models.VideoGallery{}, err
	}

	return data, nil
}

// FindByCategory implements VideoGalleryService.
func (s *videoGalleryService) FindByCategory(category string, params utils.PaginationParams) ([]models.VideoGallery, int64, error) {
	data, total, err := s.videoGalleryRepo.FindByCategory(category, params)

	if err != nil {
		return []models.VideoGallery{}, 0, err
	}

	return data, total, nil
}

// Update implements VideoGalleryService.
func (s *videoGalleryService) Update(videoGallery models.VideoGallery) (models.VideoGallery, error) {
	data, err := s.videoGalleryRepo.Update(videoGallery)

	if err != nil {
		return models.VideoGallery{}, err
	}

	return data, nil
}

// Delete implements VideoGalleryService.
func (s *videoGalleryService) Delete(id uint) error {
	err := s.videoGalleryRepo.Delete(id)

	if err != nil {
		return err
	}

	return nil
}

// FindAllActive implements VideoGalleryService.
func (s *videoGalleryService) FindAllActive() ([]models.VideoGallery, error) {
	data, err := s.videoGalleryRepo.FindAllActive()

	if err != nil {
		return []models.VideoGallery{}, err
	}

	return data, nil
}

// FindAllCategories implements VideoGalleryService.
func (s *videoGalleryService) FindAllCategories() ([]string, error) {
	data, err := s.videoGalleryRepo.FindAllCategories()

	if err != nil {
		return []string{}, err
	}

	return data, nil
}