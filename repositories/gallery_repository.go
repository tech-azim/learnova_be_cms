package repositories

import (
	"github.com/tech-azim/be-learnova/models"
	"github.com/tech-azim/be-learnova/utils"
	"gorm.io/gorm"
)

type GalleryRepository interface {
	FindAll(param utils.PaginationParams) ([]models.Gallery, int64, error)
	FindByID(id uint) (models.Gallery, error)
	Create(gallery models.Gallery) (models.Gallery, error)
	Update(gallery models.Gallery) (models.Gallery, error)
	Delete(id uint) error
	FindAllActive() ([]models.Gallery, error)
}

type galleryRepository struct {
	db *gorm.DB
}

func NewGalleryRepository(db *gorm.DB) GalleryRepository {
	return &galleryRepository{db}
}

// Create implements GalleryRepository.
func (r *galleryRepository) Create(gallery models.Gallery) (models.Gallery, error) {
	err := r.db.Create(&gallery).Error
	return gallery, err
}

// Delete implements GalleryRepository.
func (r *galleryRepository) Delete(id uint) error {
	err := r.db.Model(&models.Gallery{}).Where("id = ?", id).Update("is_deleted", true).Error
	return err
}

// FindAll implements GalleryRepository.
func (r *galleryRepository) FindAll(params utils.PaginationParams) ([]models.Gallery, int64, error) {
	offset := (params.Page - 1) * params.Limit

	var galleries []models.Gallery
	var total int64

	query := r.db.Model(&models.Gallery{}).Where("is_deleted = ?", false)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("order ASC, date DESC").Offset(offset).Limit(params.Limit).Find(&galleries).Error

	return galleries, total, err
}

// FindByID implements GalleryRepository.
func (r *galleryRepository) FindByID(id uint) (models.Gallery, error) {
	var gallery models.Gallery

	err := r.db.Where("id = ? AND is_deleted = ?", id, false).First(&gallery).Error

	return gallery, err
}

// Update implements GalleryRepository.
func (r *galleryRepository) Update(gallery models.Gallery) (models.Gallery, error) {
	err := r.db.Save(&gallery).Error

	return gallery, err
}

// FindAllActive implements GalleryRepository.
func (r *galleryRepository) FindAllActive() ([]models.Gallery, error) {
	var galleries []models.Gallery

	err := r.db.Where("is_deleted = ? AND is_active = ?", false, true).
		Order("order ASC, date DESC").
		Find(&galleries).Error

	return galleries, err
}