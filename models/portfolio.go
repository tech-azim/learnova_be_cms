package models

import (
	"time"

	"gorm.io/gorm"
)

type Portfolio struct {
	ID          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Title       string         `gorm:"type:varchar(255);not null" json:"title"`
	Count       string         `gorm:"type:varchar(50);not null" json:"count"`
	Description string         `gorm:"type:text;not null" json:"description"`
	CreatedAt   time.Time      `json:"created_at"` // ✅ Gunakan time.Time
	UpdatedAt   time.Time      `json:"updated_at"` // ✅ Gunakan time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"` // ✅ Ini sudah benar
	IsDeleted   bool           `json:"is_deleted" gorm:"default:false"`
}