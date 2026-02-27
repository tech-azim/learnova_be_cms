package services

import (
	"errors"
	"log"

	"github.com/tech-azim/be-learnova/models"
	"github.com/tech-azim/be-learnova/repositories"
	"gorm.io/gorm"
)

type UpdateProfileInput struct {
	Name     string `json:"name"     binding:"omitempty,min=3"`
	Email    string `json:"email"    binding:"omitempty,email"`
	Password string `json:"password" binding:"omitempty,min=6"`
	Phone    string `json:"phone"    binding:"omitempty"`
}

// UserService interface mendefinisikan business logic untuk User
type UserService interface {
	GetAllUsers() ([]models.User, error)
	GetUserByID(id uint) (models.User, error)
	CreateUser(input CreateUserInput) (models.User, error)
	UpdateUser(id uint, input UpdateUserInput) (models.User, error)
	DeleteUser(id uint) error

	GetProfile(id uint) (models.User, error)
	UpdateProfile(id uint, input UpdateProfileInput) (models.User, error)
}

// CreateUserInput DTO untuk membuat user baru
type CreateUserInput struct {
	Name     string `json:"name"     binding:"required,min=3"`
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Phone    string `json:"phone"    binding:"required"`
}

// UpdateUserInput DTO untuk update user (semua field opsional kecuali yang di-tag)
type UpdateUserInput struct {
	Name     string `json:"name"  binding:"omitempty,min=3"`
	Email    string `json:"email" binding:"omitempty,email"`
	Password string `json:"password" binding:"omitempty,min=6"`
	Phone    string `json:"phone"`
}

type userService struct {
	repo repositories.UserRepository
}

// NewUserService membuat instance baru UserService
func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo}
}

// GetAllUsers mengambil semua user
func (s *userService) GetAllUsers() ([]models.User, error) {
	return s.repo.FindAll()
}

// GetUserByID mencari user berdasarkan ID
// Return error "user not found" jika tidak ada
func (s *userService) GetUserByID(id uint) (models.User, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.User{}, errors.New("user not found")
		}
		return models.User{}, err
	}
	return user, nil
}

// CreateUser membuat user baru dengan validasi email unik
func (s *userService) CreateUser(input CreateUserInput) (models.User, error) {
	// Cek apakah email sudah digunakan
	_, err := s.repo.FindByEmail(input.Email)
	if err == nil {
		return models.User{}, errors.New("email already registered")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return models.User{}, err
	}

	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password, // di-hash di repository
		Phone:    input.Phone,
	}

	return s.repo.Create(user)
}

// UpdateUser memperbarui data user berdasarkan ID
// Validasi email unik jika email berubah
func (s *userService) UpdateUser(id uint, input UpdateUserInput) (models.User, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.User{}, errors.New("user not found")
		}
		return models.User{}, err
	}

	// Cek konflik email hanya jika email berubah
	if input.Email != "" && input.Email != user.Email {
		existing, err := s.repo.FindByEmail(input.Email)
		if err == nil && existing.ID != id {
			return models.User{}, errors.New("email already used by another user")
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return models.User{}, err
		}
		user.Email = input.Email
	}

	if input.Name != "" {
		user.Name = input.Name
	}
	if input.Phone != "" {
		user.Phone = input.Phone
	}
	// Password di-hash di repository jika tidak kosong
	if input.Password != "" {
		user.Password = input.Password
	}

	return s.repo.Update(user)
}

// DeleteUser menghapus user berdasarkan ID
func (s *userService) DeleteUser(id uint) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}
	return s.repo.Delete(id)
}

func (s *userService) GetProfile(id uint) (models.User, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.User{}, errors.New("user not found")
		}
		return models.User{}, err
	}
	return user, nil
}

// UpdateProfile memperbarui data user yang sedang login
// Validasi email unik hanya jika email berubah
func (s *userService) UpdateProfile(id uint, input UpdateProfileInput) (models.User, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.User{}, errors.New("user not found")
		}
		return models.User{}, err
	}

	if input.Email != "" && input.Email != user.Email {
		existing, err := s.repo.FindByEmail(input.Email)
		if err == nil && existing.ID != id {
			return models.User{}, errors.New("email already used by another user")
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return models.User{}, err
		}
		user.Email = input.Email
	}

	if input.Name != "" {
		user.Name = input.Name
	}
	if input.Phone != "" {
		user.Phone = input.Phone
	}
	log.Printf("log passsowrd input %s", input.Password)
	log.Printf("log passsowrd user %s", user.Password)

	if input.Password != "" {
		user.Password = input.Password
	}

	return s.repo.Update(user)
}
