package models

import "time"

type Registration struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	Name          string    `json:"name"`
	Email         string    `json:"email" gorm:"unique"`
	Phone         string    `json:"phone"`
	Company       string    `json:"company"`
	Position      string    `json:"position"`

	ProgramID uint    `json:"programId" gorm:"index"`
	Program   Program `json:"program" gorm:"foreignKey:ProgramID"`

	Participants  int       `json:"participants" gorm:"type:int"`
	PreferredDate time.Time `json:"preferredDate" gorm:"type:date"`
	Message       string    `json:"message" gorm:"type:text"`

	Status    string `json:"status" gorm:"type:varchar(20);default:'pending'"`
	IsDeleted bool   `json:"is_deleted" gorm:"default:false"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
