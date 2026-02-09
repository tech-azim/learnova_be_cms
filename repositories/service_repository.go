package repositories

import (
	"github.com/tech-azim/be-learnova/models"
	"github.com/tech-azim/be-learnova/utils"
	"gorm.io/gorm"
)

type ServiceRepository interface {
	FindAll(param utils.PaginationParams) ([]models.Service, int64, error)
	FindByID(id uint) (models.Service, error)
	Create(service models.Service) (models.Service, error)
	Update(service models.Service) (models.Service, error)
	Delete(id uint) error
}

type serviceRepository struct {
	db *gorm.DB
}

func NewServiceRepository(db *gorm.DB) ServiceRepository {
	return &serviceRepository{db}
}

// Create implements ServiceRepository.
func (s *serviceRepository) Create(service models.Service) (models.Service, error) {
	err := s.db.Create(&service).Error
	return service, err
}

// Delete implements ServiceRepository.
func (s *serviceRepository) Delete(id uint) error {
	err := s.db.Model(&models.Service{}).Where("id = ?", id).Update("is_deleted", true).Error
	return err
}

// FindAll implements ServiceRepository.
func (s *serviceRepository) FindAll(params utils.PaginationParams) ([]models.Service, int64, error) {
	offset := (params.Page - 1) * params.Limit

	var services []models.Service
	var total int64

	query := s.db.Model(&models.Service{}).Where("is_deleted = ?", false)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Offset(offset).Limit(params.Limit).Find(&services).Error

	return services, total, err
}

// FindByID implements ServiceRepository.
func (s *serviceRepository) FindByID(id uint) (models.Service, error) {
	var service models.Service

	err := s.db.Where("id = ? AND is_deleted = ?", id, false).First(&service).Error

	return service, err
}

// Update implements ServiceRepository.
func (s *serviceRepository) Update(service models.Service) (models.Service, error) {
	err := s.db.Save(&service).Error

	return service, err
}