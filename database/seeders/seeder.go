package seeders

import "gorm.io/gorm"

func RunAllSeeder(db *gorm.DB){
	SeederUsers(db)
}