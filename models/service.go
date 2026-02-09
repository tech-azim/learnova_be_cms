package models

import (
	"gorm.io/gorm"
)

type Service struct {
	ID          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Icon        string         `gorm:"type:varchar(255);not null" json:"icon"`        
	Title       string         `gorm:"type:varchar(255);not null" json:"title"`     
	Description string         `gorm:"type:text;not null" json:"description"`       
	Color       string         `gorm:"type:varchar(50);not null" json:"color"`       
	CreatedAt   gorm.DeletedAt `json:"created_at"`                                  
	UpdatedAt   gorm.DeletedAt `json:"updated_at"`                                   
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`          
	IsDeleted bool   `json:"is_deleted" gorm:"default:false"`
}
