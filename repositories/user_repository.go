package repositories

import (
	"github.com/tech-azim/be-learnova/models"
	"github.com/tech-azim/be-learnova/utils"
	"gorm.io/gorm"
)

// UserRepository interface untuk operasi database User
// Interface ini memudahkan testing dan mengikuti prinsip SOLID
type UserRepository interface {
	FindAll() ([]models.User, error)
	FindByID(id uint) (models.User, error)
	FindByEmail(email string) (models.User, error)
	Create(user models.User) (models.User, error)
	Update(user models.User) (models.User, error)
	Delete(id uint) error
}

// userRepository implementasi dari UserRepository
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository membuat instance baru UserRepository
// Menggunakan dependency injection untuk koneksi database
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

// FindAll mengambil semua user dari database
func (r *userRepository) FindAll() ([]models.User, error) {
	var users []models.User
	err := r.db.Find(&users).Error
	return users, err
}

// Create menyimpan user baru ke database
// GORM otomatis mengisi ID setelah insert berhasil
func (r *userRepository) Create(user models.User) (models.User, error) {
	hashPassword, errHash := utils.HashPassword(user.Password)
	if errHash != nil {
		user.Password = hashPassword
	}
	err := r.db.Create(&user).Error
	return user, err
}

// Delete menghapus user berdasarkan ID
// Jika model punya DeletedAt, ini soft delete (data tidak benar-benar dihapus)
// SQL: DELETE FROM users WHERE id = ? (atau UPDATE users SET deleted_at = NOW() WHERE id = ?)
func (r *userRepository) Delete(id uint) error {
	err := r.db.Delete(&models.User{}, id).Error
	return err
}

// FindByEmail mencari user berdasarkan email
// Berguna untuk validasi email unique dan proses login
func (r *userRepository) FindByEmail(email string) (models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return user, err
}

// FindByID mencari user berdasarkan ID
// Return error gorm.ErrRecordNotFound jika tidak ditemukan
func (r *userRepository) FindByID(id uint) (models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	return user, err
}

// Update memperbarui data user yang sudah ada
// user.ID harus sudah terisi, Save() akan update semua field
func (r *userRepository) Update(user models.User) (models.User, error) {
	if user.Password != "" {
		hashedPassword, err := utils.HashPassword(user.Password)
		if err != nil {
			return models.User{}, err
		}
		user.Password = hashedPassword
	}
	err := r.db.Save(&user).Error
	return user, err
}