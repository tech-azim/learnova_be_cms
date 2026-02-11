package repositories

import (
	"log"

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

    // ✅ Debug: Log pagination params
    log.Printf("=== FindAll Debug ===")
    log.Printf("Params: Page=%d, Limit=%d, Offset=%d", params.Page, params.Limit, offset)

    query := p.db.Model(&models.Portfolio{}).Where("is_deleted = ?", false)

    // ✅ Debug: Enable SQL logging
    query = query.Debug()

    if err := query.Count(&total).Error; err != nil {
        log.Printf("Error counting portfolios: %v", err)
        return nil, 0, err
    }
    log.Printf("Total count: %d", total)

    err := query.Offset(offset).Limit(params.Limit).Find(&portfolios).Error
    if err != nil {
        log.Printf("Error finding portfolios: %v", err)
        return nil, 0, err
    }

    log.Printf("Found %d portfolios in result", len(portfolios))
    log.Printf("Portfolios data: %+v", portfolios)

    return portfolios, total, nil
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