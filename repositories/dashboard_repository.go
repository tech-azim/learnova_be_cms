package repositories

import (
	"github.com/tech-azim/be-learnova/models"
	"gorm.io/gorm"
)

type DashboardRepository interface {
	GetTotalProgram() (int64, error)
	GetTotalRegistration() (int64, error)
	GetActiveParticipants() (int64, error)
	GetPendingParticipants() (int64, error)
	GetLatestRegistrations(limit int) ([]models.Registration, error)
	GetRecentActivities(limit int) ([]models.Registration, error)
	GetPopularPrograms(limit int) ([]PopularProgram, error)
}

type PopularProgram struct {
	ID                uint   `json:"id"`
	Title             string `json:"title"`
	Icon              string `json:"icon"`
	Level             string `json:"level"`
	TotalRegistration int64  `json:"total_registration"`
}

type dashboardRepository struct {
	db *gorm.DB
}

func NewDashboardRepository(db *gorm.DB) DashboardRepository {
	return &dashboardRepository{db}
}

func (r *dashboardRepository) GetTotalProgram() (int64, error) {
	var total int64
	err := r.db.Model(&models.Program{}).Where("is_deleted = ?", false).Count(&total).Error
	return total, err
}

func (r *dashboardRepository) GetTotalRegistration() (int64, error) {
	var total int64
	err := r.db.Model(&models.Registration{}).Where("is_deleted = ?", false).Count(&total).Error
	return total, err
}

func (r *dashboardRepository) GetActiveParticipants() (int64, error) {
	var total int64
	err := r.db.Model(&models.Registration{}).
		Where("status = ? AND is_deleted = ?", "active", false).
		Count(&total).Error
	return total, err
}

func (r *dashboardRepository) GetPendingParticipants() (int64, error) {
	var total int64
	err := r.db.Model(&models.Registration{}).
		Where("status = ? AND is_deleted = ?", "pending", false).
		Count(&total).Error
	return total, err
}

func (r *dashboardRepository) GetLatestRegistrations(limit int) ([]models.Registration, error) {
	var registrations []models.Registration
	err := r.db.Preload("Program").
		Where("is_deleted = ?", false).
		Order("created_at DESC").
		Limit(limit).
		Find(&registrations).Error
	return registrations, err
}

func (r *dashboardRepository) GetRecentActivities(limit int) ([]models.Registration, error) {
	var registrations []models.Registration
	err := r.db.Preload("Program").
		Where("is_deleted = ?", false).
		Order("updated_at DESC").
		Limit(limit).
		Find(&registrations).Error
	return registrations, err
}

func (r *dashboardRepository) GetPopularPrograms(limit int) ([]PopularProgram, error) {
	var programs []PopularProgram
	err := r.db.Table("programs").
		Select("programs.id, programs.title, programs.icon, programs.level, COUNT(registrations.id) as total_registration").
		Joins("LEFT JOIN registrations ON registrations.program_id = programs.id AND registrations.is_deleted = false").
		Where("programs.is_deleted = ?", false).
		Group("programs.id, programs.title, programs.icon, programs.level").
		Order("total_registration DESC").
		Limit(limit).
		Scan(&programs).Error
	return programs, err
}
