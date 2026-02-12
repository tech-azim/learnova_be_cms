package repositories

import (
	"github.com/tech-azim/be-learnova/models"
	"github.com/tech-azim/be-learnova/utils"
	"gorm.io/gorm"
)

type FeatureRepository interface {
	FindAll(param utils.PaginationParams) ([]models.Feature, int64, error)
	FindByID(id uint) (models.Feature, error)
	Create(feature models.Feature) (models.Feature, error)
	Update(feature models.Feature) (models.Feature, error)
	Delete(id uint) error
	FindAllActive() ([]models.Feature, error) // tambahan untuk get semua feature yang aktif
}

type featureRepository struct {
	db *gorm.DB
}

func NewFeatureRepository(db *gorm.DB) FeatureRepository {
	return &featureRepository{db}
}

// Create implements FeatureRepository.
func (r *featureRepository) Create(feature models.Feature) (models.Feature, error) {
	err := r.db.Create(&feature).Error
	return feature, err
}

// Delete implements FeatureRepository.
func (r *featureRepository) Delete(id uint) error {
	err := r.db.Model(&models.Feature{}).Where("id = ?", id).Update("is_deleted", true).Error
	return err
}

// FindAll implements FeatureRepository.
func (r *featureRepository) FindAll(params utils.PaginationParams) ([]models.Feature, int64, error) {
	offset := (params.Page - 1) * params.Limit

	var features []models.Feature
	var total int64

	query := r.db.Model(&models.Feature{}).Where("is_deleted = ?", false)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Offset(offset).Limit(params.Limit).Find(&features).Error

	return features, total, err
}

// FindByID implements FeatureRepository.
func (r *featureRepository) FindByID(id uint) (models.Feature, error) {
	var feature models.Feature

	err := r.db.Where("id = ? AND is_deleted = ?", id, false).First(&feature).Error

	return feature, err
}

// Update implements FeatureRepository.
func (r *featureRepository) Update(feature models.Feature) (models.Feature, error) {
	err := r.db.Save(&feature).Error

	return feature, err
}

// FindAllActive implements FeatureRepository.
func (r *featureRepository) FindAllActive() ([]models.Feature, error) {
	var features []models.Feature

	err := r.db.Where("is_deleted = ? AND is_active = ?", false, true).
		Order("order ASC").
		Find(&features).Error

	return features, err
}