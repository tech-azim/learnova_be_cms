package repositories

import (
	"github.com/tech-azim/be-learnova/models"
	"github.com/tech-azim/be-learnova/utils"
	"gorm.io/gorm"
)

type HeroRepository interface {
	FindAll(param utils.PaginationParams) ([]models.Hero, error)
	FindByID(id uint) (models.Hero, error)
	Create(hero models.Hero) (models.Hero, error)
	Update(hero models.Hero) (models.Hero, error)
	Delete(id uint) error
}

type heroRepository struct {
	db *gorm.DB
}

func NewHeroRepository(db *gorm.DB) HeroRepository {
	return &heroRepository{db}
}

// Create implements [HeroRepository].
func (h *heroRepository) Create(hero models.Hero) (models.Hero, error) {
	err := h.db.Create(&hero).Error
	return hero, err
}

// Delete implements [HeroRepository].
func (h *heroRepository) Delete(id uint) error {
	err := h.db.Delete(&models.Hero{}, id).Error
	return err
}

// FindAll implements [HeroRepository].
func (h *heroRepository) FindAll(params utils.PaginationParams) ([]models.Hero, error) {
	offset := (params.Page - 1) * params.Limit

	var heroes []models.Hero

	err := h.db.Offset(offset).Limit(params.Limit).Find(&heroes).Error

	return heroes, err
}

// FindByID implements [HeroRepository].
func (h *heroRepository) FindByID(id uint) (models.Hero, error) {
	var hero models.Hero

	err := h.db.First(&hero).Error

	return hero, err
}

// Update implements [HeroRepository].
func (h *heroRepository) Update(hero models.Hero) (models.Hero, error) {
	err := h.db.Save(&hero).Error

	return hero, err
}

