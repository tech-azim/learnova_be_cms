package services

import (
	"github.com/tech-azim/be-learnova/models"
	"github.com/tech-azim/be-learnova/repositories"
	"github.com/tech-azim/be-learnova/utils"
)

type PortfolioService interface {
	Create(portfolio models.Portfolio) (models.Portfolio, error)
	FindAll(params utils.PaginationParams) ([]models.Portfolio, int64, error)
	FindByID(id uint) (models.Portfolio, error)
	Update(portfolio models.Portfolio) (models.Portfolio, error)
	Delete(id uint) error
}

type portfolioService struct {
	portfolioRepo repositories.PortfolioRepository
}

func NewPortfolioService(portfolioRepo repositories.PortfolioRepository) PortfolioService {
	return &portfolioService{
		portfolioRepo,
	}
}

// Create implements PortfolioService.
func (p *portfolioService) Create(portfolio models.Portfolio) (models.Portfolio, error) {
	result, err := p.portfolioRepo.Create(portfolio)

	if err != nil {
		return models.Portfolio{}, err
	}

	return result, nil
}

// FindAll implements PortfolioService.
func (p *portfolioService) FindAll(params utils.PaginationParams) ([]models.Portfolio, int64, error) {
	data, total, err := p.portfolioRepo.FindAll(params)

	if err != nil {
		return []models.Portfolio{}, 0, err
	}

	return data, total, nil
}

// FindByID implements PortfolioService.
func (p *portfolioService) FindByID(id uint) (models.Portfolio, error) {
	data, err := p.portfolioRepo.FindByID(id)

	if err != nil {
		return models.Portfolio{}, err
	}

	return data, nil
}

// Update implements PortfolioService.
func (p *portfolioService) Update(portfolio models.Portfolio) (models.Portfolio, error) {
	data, err := p.portfolioRepo.Update(portfolio)

	if err != nil {
		return models.Portfolio{}, err
	}

	return data, nil
}

// Delete implements PortfolioService.
func (p *portfolioService) Delete(id uint) error {
	err := p.portfolioRepo.Delete(id)

	if err != nil {
		return err
	}

	return nil
}