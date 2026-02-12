package config

import (
	"fmt"
	"log"
	"os"

	"github.com/tech-azim/be-learnova/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
	)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect db", err)
	}

	database.AutoMigrate(&models.User{}, &models.Hero{}, &models.Program{}, &models.Registration{}, &models.Service{}, &models.Portfolio{}, &models.Feature{},  &models.Gallery{},  &models.FlyerGallery{}, &models.VideoGallery{},)

	DB = database
	log.Print("Successfully connect database")
}
