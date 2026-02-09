package repositories

import (
	"github.com/tech-azim/be-learnova/models"
	"github.com/tech-azim/be-learnova/utils"
	"gorm.io/gorm"
)

type PortfolioRepository interface {
	FindAll(param utils.PaginationParams) ([]models.Portfolio, int64, error)
	FindByID(id uint) (models.Portfolio, error)
	Create(portfolio models.Portfolio) (models.Portfolio, error)
	Update(portfolio models.Portfolio) (models.Portfolio, error)
	Delete(id uint) error
}

type portfolioRepository struct {
	db *gorm.DB
}

func NewPortfolioRepository(db *gorm.DB) PortfolioRepository {
	return &portfolioRepository{db}
}

// Create implements PortfolioRepository.
func (p *portfolioRepository) Create(portfolio models.Portfolio) (models.Portfolio, error) {
	err := p.db.Create(&portfolio).Error
	return portfolio, err
}

// Delete implements PortfolioRepository.
func (p *portfolioRepository) Delete(id uint) error {
	err := p.db.Model(&models.Portfolio{}).Where("id = ?", id).Update("is_deleted", true).Error
	return err
}

// FindAll implements PortfolioRepository.
func (p *portfolioRepository) FindAll(params utils.PaginationParams) ([]models.Portfolio, int64, error) {
	offset := (params.Page - 1) * params.Limit

	var portfolios []models.Portfolio
	var total int64

	query := p.db.Model(&models.Portfolio{}).Where("is_deleted = ?", false)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Offset(offset).Limit(params.Limit).Find(&portfolios).Error

	return portfolios, total, err
}

// FindByID implements PortfolioRepository.
func (p *portfolioRepository) FindByID(id uint) (models.Portfolio, error) {
	var portfolio models.Portfolio

	err := p.db.Where("id = ? AND is_deleted = ?", id, false).First(&portfolio).Error

	return portfolio, err
}

// Update implements PortfolioRepository.
func (p *portfolioRepository) Update(portfolio models.Portfolio) (models.Portfolio, error) {
	err := p.db.Save(&portfolio).Error

	return portfolio, err
}