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
	r.RedirectTrailingSlash = true


	config.ConnectDB()


	if *seedFlag {
		seeders.RunAllSeeder(config.DB)
		return
	}

	userRepo := repositories.NewUserRepository(config.DB)
	heroRepo := repositories.NewHeroRepository(config.DB)

	authService := services.NewAuthService(userRepo)
	heroService := services.NewHeroService(heroRepo)

	authController := controllers.NewAuthController(authService)
	heroController := controllers.NewHeroController(heroService)


	routes.Router(r, authController, heroController)
	for _, route := range r.Routes() {
		fmt.Printf("Method: %s | Path: %s\n", route.Method, route.Path)
	}
	r.Run()
    r.Run()
}

