package services

import (
	"github.com/tech-azim/be-learnova/models"
	"github.com/tech-azim/be-learnova/repositories"
	"github.com/tech-azim/be-learnova/utils"
)

type HeroService interface {
	Create(hero models.Hero) (models.Hero, error)
	FindAll(params utils.PaginationParams) ([]models.Hero,int64, error)
}

type heroService struct {
	heroRepo repositories.HeroRepository
}

func NewHeroService(heroRepo repositories.HeroRepository) HeroService {
	return &heroService{
		heroRepo,
	}
}

// Create implements [HeroService].
func (h *heroService) Create(hero models.Hero) (models.Hero, error) {
	result, err := h.heroRepo.Create(hero)

	if err != nil {
		return models.Hero{}, err
	}

	return result, err
}


func (h *heroService) FindAll(params utils.PaginationParams) ([]models.Hero,int64, error) {
	data,total, err := h.heroRepo.FindAll(params)

	if err != nil {
		return []models.Hero{},0, err
	}
	return data,total, nil
}
