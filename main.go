package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/tech-azim/be-learnova/config"
	"github.com/tech-azim/be-learnova/controllers"
	"github.com/tech-azim/be-learnova/database/seeders"
	"github.com/tech-azim/be-learnova/repositories"
	"github.com/tech-azim/be-learnova/routes"
	"github.com/tech-azim/be-learnova/services"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("env not laod")
	}

	seedFlag := flag.Bool("seed", false , "Run database seeders")
	flag.Parse()
    r := gin.Default()


	config.ConnectDB()

	fmt.Printf("log seedflag %b", *seedFlag)

	if *seedFlag {
		seeders.RunAllSeeder(config.DB)
		return
	}

	userRepo := repositories.NewUserRepository(config.DB)

	authService := services.NewAuthService(userRepo)

	authController := controllers.NewAuthController(authService)


	routes.Router(r, authController)
    r.Run()
}

