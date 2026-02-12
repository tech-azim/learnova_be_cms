package models

import (
	"time"

	"gorm.io/gorm"
)

type Feature struct {
    ID          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
    Icon        string         `gorm:"type:varchar(255);not null" json:"icon"`
    Title       string         `gorm:"type:varchar(255);not null" json:"title"`
    Description string         `gorm:"type:text;not null" json:"description"`
    SortOrder   int            `gorm:"type:int;default:0;column:sort_order" json:"order"`
    IsActive    bool           `gorm:"default:true" json:"is_active"`
    IsDeleted   bool           `gorm:"default:false" json:"is_deleted"`
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}