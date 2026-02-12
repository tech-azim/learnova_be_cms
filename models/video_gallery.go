package models

import (
	"time"

	"gorm.io/gorm"
)

type VideoGallery struct {
	ID          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Title       string         `gorm:"type:varchar(255);not null" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	Thumbnail   string         `gorm:"type:varchar(500);not null" json:"thumbnail"`
	VideoURL    string         `gorm:"type:varchar(500);not null" json:"video_url"`
	Category    string         `gorm:"type:varchar(100);not null" json:"category"`
	Date        time.Time      `gorm:"type:date;not null" json:"date"`
	Order       int            `gorm:"type:int;default:0" json:"order"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	IsDeleted   bool           `gorm:"default:false" json:"is_deleted"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}