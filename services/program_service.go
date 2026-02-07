package services

import (
	"github.com/tech-azim/be-learnova/models"
	"github.com/tech-azim/be-learnova/repositories"
	"github.com/tech-azim/be-learnova/utils"
)

type ProgramService interface {
	Create(program models.Program) (models.Program, error)
	FindAll(params utils.PaginationParams) ([]models.Program, int64, error)
	FindByID(id uint) (models.Program, error)
	Update(program models.Program) (models.Program, error)
	Delete(id uint) error
}

type programService struct {
	programRepo repositories.ProgramRepository
}

func NewProgramService(programRepo repositories.ProgramRepository) ProgramService {
	return &programService{
		programRepo,
	}
}

// Create implements ProgramService.
func (p *programService) Create(program models.Program) (models.Program, error) {
	result, err := p.programRepo.Create(program)

	if err != nil {
		return models.Program{}, err
	}

	return result, nil
}

// FindAll implements ProgramService.
func (p *programService) FindAll(params utils.PaginationParams) ([]models.Program, int64, error) {
	data, total, err := p.programRepo.FindAll(params)

	if err != nil {
		return []models.Program{}, 0, err
	}

	return data, total, nil
}

// FindByID implements ProgramService.
func (p *programService) FindByID(id uint) (models.Program, error) {
	data, err := p.programRepo.FindByID(id)

	if err != nil {
		return models.Program{}, err
	}

	return data, nil
}

// Update implements ProgramService.
func (p *programService) Update(program models.Program) (models.Program, error) {
	data, err := p.programRepo.Update(program)

	if err != nil {
		return models.Program{}, err
	}

	return data, nil
}

// Delete implements ProgramService.
func (p *programService) Delete(id uint) error {
	err := p.programRepo.Delete(id)

	if err != nil {
		return err
	}

	return nil
}
