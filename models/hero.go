package models

type Hero struct {
	ID uint `json:"id" gorm:"primaryKey;autoIncrement"`
	SRC string `json:"src"`
	ALT string `json:"alt"`
	Title string `json:"title"`
	Description string `json:"description"`
}