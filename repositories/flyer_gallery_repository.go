package repositories

import (
	"github.com/tech-azim/be-learnova/models"
	"github.com/tech-azim/be-learnova/utils"
	"gorm.io/gorm"
)

type FlyerGalleryRepository interface {
	FindAll(param utils.PaginationParams) ([]models.FlyerGallery, int64, error)
	FindByID(id uint) (models.FlyerGallery, error)
	Create(flyerGallery models.FlyerGallery) (models.FlyerGallery, error)
	Update(flyerGallery models.FlyerGallery) (models.FlyerGallery, error)
	Delete(id uint) error
	FindAllActive() ([]models.FlyerGallery, error)
}

type flyerGalleryRepository struct {
	db *gorm.DB
}

func NewFlyerGalleryRepository(db *gorm.DB) FlyerGalleryRepository {
	return &flyerGalleryRepository{db}
}

// Create implements FlyerGalleryRepository.
func (r *flyerGalleryRepository) Create(flyerGallery models.FlyerGallery) (models.FlyerGallery, error) {
	err := r.db.Create(&flyerGallery).Error
	return flyerGallery, err
}

// Delete implements FlyerGalleryRepository.
func (r *flyerGalleryRepository) Delete(id uint) error {
	err := r.db.Model(&models.FlyerGallery{}).Where("id = ?", id).Update("is_deleted", true).Error
	return err
}

// FindAll implements FlyerGalleryRepository.
func (r *flyerGalleryRepository) FindAll(params utils.PaginationParams) ([]models.FlyerGallery, int64, error) {
	offset := (params.Page - 1) * params.Limit

	var flyerGalleries []models.FlyerGallery
	var total int64

	query := r.db.Model(&models.FlyerGallery{}).Where("is_deleted = ?", false)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("created_at DESC").Offset(offset).Limit(params.Limit).Find(&flyerGalleries).Error

	return flyerGalleries, total, err
}

// FindByID implements FlyerGalleryRepository.
func (r *flyerGalleryRepository) FindByID(id uint) (models.FlyerGallery, error) {
	var flyerGallery models.FlyerGallery

	err := r.db.Where("id = ? AND is_deleted = ?", id, false).First(&flyerGallery).Error

	return flyerGallery, err
}

// Update implements FlyerGalleryRepository.
func (r *flyerGalleryRepository) Update(flyerGallery models.FlyerGallery) (models.FlyerGallery, error) {
	err := r.db.Save(&flyerGallery).Error

	return flyerGallery, err
}

// FindAllActive implements FlyerGalleryRepository.
func (r *flyerGalleryRepository) FindAllActive() ([]models.FlyerGallery, error) {
	var flyerGalleries []models.FlyerGallery

	err := r.db.Where("is_deleted = ? AND is_active = ?", false, true).
		Order("created_at DESC").
		Find(&flyerGalleries).Error

	return flyerGalleries, err
}