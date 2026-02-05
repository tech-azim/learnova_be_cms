package seeders

import (
	"log"

	"github.com/tech-azim/be-learnova/models"
	"github.com/tech-azim/be-learnova/utils"
	"gorm.io/gorm"
)


func SeederUsers(db *gorm.DB){
	users := []models.User{
      {
		Email: "admin@learnova.com",
		Name: "Email",
		Password: "Password123",
		Phone: "-",
	  },
	}

	log.Print("running seeders user")

	for _, user := range users {
		log.Printf("email %s", user.Email)
		var existingUsers models.User
		hashPassword, errHash := utils.HashPassword(user.Password)
		if errHash != nil {
			log.Printf("error hash password %v", errHash)
			continue
		}
		user.Password = hashPassword
		if err := db.Where("email = ?",user.Email).First(&existingUsers); err != nil {
			if err := db.Create(&user).Error; err != nil {
				log.Printf("Failed to seed user email = %s,%v", user.Email, err)
			}else{
				log.Printf("Success seed email %s", user.Email)
			}
		}else{
			log.Printf("User already exist %s",user.Email )
		}
	}
}