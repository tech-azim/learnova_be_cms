package repositories

import (
	"github.com/tech-azim/be-learnova/models"
	"github.com/tech-azim/be-learnova/utils"
	"gorm.io/gorm"
)

type RegistrationRepository interface {
	FindAll(param utils.PaginationParams) ([]models.Registration, int64, error)
	FindByID(id uint) (models.Registration, error)
	FindByProgramID(programID uint, params utils.PaginationParams) ([]models.Registration, int64, error)
	FindByEmail(email string) (models.Registration, error)
	Create(registration models.Registration) (models.Registration, error)
	Update(registration models.Registration) (models.Registration, error)
	Delete(id uint) error
	CheckEmailExists(email string, programID uint) (bool, error)
}

type registrationRepository struct {
	db *gorm.DB
}

func NewRegistrationRepository(db *gorm.DB) RegistrationRepository {
	return &registrationRepository{db}
}

// Create implements RegistrationRepository.
func (r *registrationRepository) Create(registration models.Registration) (models.Registration, error) {
	err := r.db.Create(&registration).Error
	return registration, err
}

// Delete implements RegistrationRepository.
func (r *registrationRepository) Delete(id uint) error {
	err := r.db.Model(&models.Registration{}).Where("id = ?", id).Update("is_deleted", true).Error
	return err
}

// FindAll implements RegistrationRepository.
func (r *registrationRepository) FindAll(params utils.PaginationParams) ([]models.Registration, int64, error) {
	offset := (params.Page - 1) * params.Limit

	var registrations []models.Registration
	var total int64

	query := r.db.Model(&models.Registration{}).Where("is_deleted = ?", false)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Preload("Program").Offset(offset).Limit(params.Limit).Order("created_at DESC").Find(&registrations).Error

	return registrations, total, err
}

// FindByID implements RegistrationRepository.
func (r *registrationRepository) FindByID(id uint) (models.Registration, error) {
	var registration models.Registration

	err := r.db.Preload("Program").Where("id = ? AND is_deleted = ?", id, false).First(&registration).Error

	return registration, err
}

// FindByProgramID implements RegistrationRepository.
func (r *registrationRepository) FindByProgramID(programID uint, params utils.PaginationParams) ([]models.Registration, int64, error) {
	offset := (params.Page - 1) * params.Limit

	var registrations []models.Registration
	var total int64

	query := r.db.Model(&models.Registration{}).Where("program_id = ? AND is_deleted = ?", programID, false)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Preload("Program").Offset(offset).Limit(params.Limit).Order("created_at DESC").Find(&registrations).Error

	return registrations, total, err
}

// FindByEmail implements RegistrationRepository.
func (r *registrationRepository) FindByEmail(email string) (models.Registration, error) {
	var registration models.Registration

	err := r.db.Preload("Program").Where("email = ? AND is_deleted = ?", email, false).First(&registration).Error

	return registration, err
}

// Update implements RegistrationRepository.
func (r *registrationRepository) Update(registration models.Registration) (models.Registration, error) {
	err := r.db.Save(&registration).Error

	return registration, err
}

// CheckEmailExists implements RegistrationRepository.
func (r *registrationRepository) CheckEmailExists(email string, programID uint) (bool, error) {
	var count int64

	err := r.db.Model(&models.Registration{}).
		Where("email = ? AND program_id = ? AND is_deleted = ?", email, programID, false).
		Count(&count).Error

	return count > 0, err
}
