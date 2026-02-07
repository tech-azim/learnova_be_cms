package models

import (
	"time"
)

type Program struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	Icon         string         `json:"icon" gorm:"type:varchar(100)"`
	Title        string         `json:"title" gorm:"type:varchar(255)"`
	Duration     string         `json:"duration" gorm:"type:varchar(50)"`
	Participants string         `json:"participants" gorm:"type:varchar(100)"`
	Level        string         `json:"level" gorm:"type:varchar(50)"`
	Description  string         `json:"description" gorm:"type:text"`
	Benefits     []string       `json:"benefits" gorm:"type:text[]"`
	Image        string         `json:"image" gorm:"type:varchar(255)"`
	IsDeleted    bool           `json:"is_deleted" gorm:"default:false"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	Registration []Registration `json:"registration" gorm:"foreignKey:ProgramID"`
}
