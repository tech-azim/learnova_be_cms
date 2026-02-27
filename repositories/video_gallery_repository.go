package repositories

import (
	"github.com/tech-azim/be-learnova/models"
	"github.com/tech-azim/be-learnova/utils"
	"gorm.io/gorm"
)

type VideoGalleryRepository interface {
	FindAll(param utils.PaginationParams) ([]models.VideoGallery, int64, error)
	FindByID(id uint) (models.VideoGallery, error)
	FindByCategory(category string, params utils.PaginationParams) ([]models.VideoGallery, int64, error)
	Create(videoGallery models.VideoGallery) (models.VideoGallery, error)
	Update(videoGallery models.VideoGallery) (models.VideoGallery, error)
	Delete(id uint) error
	FindAllActive() ([]models.VideoGallery, error)
	FindAllCategories() ([]string, error)
}

type videoGalleryRepository struct {
	db *gorm.DB
}

func NewVideoGalleryRepository(db *gorm.DB) VideoGalleryRepository {
	return &videoGalleryRepository{db}
}

// Create implements VideoGalleryRepository.
func (r *videoGalleryRepository) Create(videoGallery models.VideoGallery) (models.VideoGallery, error) {
	err := r.db.Create(&videoGallery).Error
	return videoGallery, err
}

// Delete implements VideoGalleryRepository.
func (r *videoGalleryRepository) Delete(id uint) error {
	err := r.db.Model(&models.VideoGallery{}).Where("id = ?", id).Update("is_deleted", true).Error
	return err
}

// FindAll implements VideoGalleryRepository.
func (r *videoGalleryRepository) FindAll(params utils.PaginationParams) ([]models.VideoGallery, int64, error) {
	offset := (params.Page - 1) * params.Limit

	var videoGalleries []models.VideoGallery
	var total int64

	query := r.db.Model(&models.VideoGallery{}).Where("is_deleted = ?", false)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Offset(offset).Limit(params.Limit).Find(&videoGalleries).Error

	return videoGalleries, total, err
}

// FindByID implements VideoGalleryRepository.
func (r *videoGalleryRepository) FindByID(id uint) (models.VideoGallery, error) {
	var videoGallery models.VideoGallery

	err := r.db.Where("id = ? AND is_deleted = ?", id, false).First(&videoGallery).Error

	return videoGallery, err
}

// FindByCategory implements VideoGalleryRepository.
func (r *videoGalleryRepository) FindByCategory(category string, params utils.PaginationParams) ([]models.VideoGallery, int64, error) {
	offset := (params.Page - 1) * params.Limit

	var videoGalleries []models.VideoGallery
	var total int64

	query := r.db.Model(&models.VideoGallery{}).Where("is_deleted = ? AND is_active = ?", false, true)

	// Jika category bukan "Semua", filter berdasarkan category
	if category != "" && category != "Semua" {
		query = query.Where("category = ?", category)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("order ASC, date DESC").Offset(offset).Limit(params.Limit).Find(&videoGalleries).Error

	return videoGalleries, total, err
}

// Update implements VideoGalleryRepository.
func (r *videoGalleryRepository) Update(videoGallery models.VideoGallery) (models.VideoGallery, error) {
	err := r.db.Save(&videoGallery).Error

	return videoGallery, err
}

// FindAllActive implements VideoGalleryRepository.
func (r *videoGalleryRepository) FindAllActive() ([]models.VideoGallery, error) {
	var videoGalleries []models.VideoGallery

	err := r.db.Where("is_deleted = ? AND is_active = ?", false, true).
		Order("order ASC, date DESC").
		Find(&videoGalleries).Error

	return videoGalleries, err
}

// FindAllCategories implements VideoGalleryRepository.
func (r *videoGalleryRepository) FindAllCategories() ([]string, error) {
	var categories []string

	err := r.db.Model(&models.VideoGallery{}).
		Where("is_deleted = ? AND is_active = ?", false, true).
		Distinct("category").
		Pluck("category", &categories).Error

	return categories, err
}
