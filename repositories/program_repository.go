package repositories

import (
	"github.com/tech-azim/be-learnova/models"
	"github.com/tech-azim/be-learnova/utils"
	"gorm.io/gorm"
)

type ProgramRepository interface {
	FindAll(param utils.PaginationParams) ([]models.Program, int64, error)
	FindByID(id uint) (models.Program, error)
	Create(program models.Program) (models.Program, error)
	Update(program models.Program) (models.Program, error)
	Delete(id uint) error
}

type programRepository struct {
	db *gorm.DB
}

func NewProgramRepository(db *gorm.DB) ProgramRepository {
	return &programRepository{db}
}

// Create implements ProgramRepository.
func (p *programRepository) Create(program models.Program) (models.Program, error) {
	err := p.db.Create(&program).Error
	return program, err
}

// Delete implements ProgramRepository.
func (p *programRepository) Delete(id uint) error {
	err := p.db.Model(&models.Program{}).Where("id = ?", id).Update("is_deleted", true).Error
	return err
}

// FindAll implements ProgramRepository.
func (p *programRepository) FindAll(params utils.PaginationParams) ([]models.Program, int64, error) {
	offset := (params.Page - 1) * params.Limit

	var programs []models.Program
	var total int64

	query := p.db.Model(&models.Program{}).Where("is_deleted = ?", false)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Offset(offset).Limit(params.Limit).Find(&programs).Error

	return programs, total, err
}

// FindByID implements ProgramRepository.
func (p *programRepository) FindByID(id uint) (models.Program, error) {
	var program models.Program

	err := p.db.Where("id = ? AND is_deleted = ?", id, false).First(&program).Error

	return program, err
}

// Update implements ProgramRepository.
func (p *programRepository) Update(program models.Program) (models.Program, error) {
	err := p.db.Save(&program).Error

	return program, err
}
