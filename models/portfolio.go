package models

import (
	"gorm.io/gorm"
)

type Portfolio struct {
	ID          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Title       string         `gorm:"type:varchar(255);not null" json:"title"`
	Count       string         `gorm:"type:varchar(50);not null" json:"count"`
	Description string         `gorm:"type:text;not null" json:"description"`
	CreatedAt   gorm.DeletedAt `json:"created_at"`
	UpdatedAt   gorm.DeletedAt `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	IsDeleted   bool           `json:"is_deleted" gorm:"default:false"`
}