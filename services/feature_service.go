package services

import (
	"github.com/tech-azim/be-learnova/models"
	"github.com/tech-azim/be-learnova/repositories"
	"github.com/tech-azim/be-learnova/utils"
)

type FeatureService interface {
	Create(feature models.Feature) (models.Feature, error)
	FindAll(params utils.PaginationParams) ([]models.Feature, int64, error)
	FindByID(id uint) (models.Feature, error)
	Update(feature models.Feature) (models.Feature, error)
	Delete(id uint) error
	FindAllActive() ([]models.Feature, error)
}

type featureService struct {
	featureRepo repositories.FeatureRepository
}

func NewFeatureService(featureRepo repositories.FeatureRepository) FeatureService {
	return &featureService{
		featureRepo,
	}
}

// Create implements FeatureService.
func (s *featureService) Create(feature models.Feature) (models.Feature, error) {
	result, err := s.featureRepo.Create(feature)

	if err != nil {
		return models.Feature{}, err
	}

	return result, nil
}

// FindAll implements FeatureService.
func (s *featureService) FindAll(params utils.PaginationParams) ([]models.Feature, int64, error) {
	data, total, err := s.featureRepo.FindAll(params)

	if err != nil {
		return []models.Feature{}, 0, err
	}

	return data, total, nil
}

// FindByID implements FeatureService.
func (s *featureService) FindByID(id uint) (models.Feature, error) {
	data, err := s.featureRepo.FindByID(id)

	if err != nil {
		return models.Feature{}, err
	}

	return data, nil
}

// Update implements FeatureService.
func (s *featureService) Update(feature models.Feature) (models.Feature, error) {
	data, err := s.featureRepo.Update(feature)

	if err != nil {
		return models.Feature{}, err
	}

	return data, nil
}

// Delete implements FeatureService.
func (s *featureService) Delete(id uint) error {
	err := s.featureRepo.Delete(id)

	if err != nil {
		return err
	}

	return nil
}

// FindAllActive implements FeatureService.
func (s *featureService) FindAllActive() ([]models.Feature, error) {
	data, err := s.featureRepo.FindAllActive()

	if err != nil {
		return []models.Feature{}, err
	}

	return data, nil
}