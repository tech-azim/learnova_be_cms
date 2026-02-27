package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/tech-azim/be-learnova/config"
	"github.com/tech-azim/be-learnova/controllers"
	"github.com/tech-azim/be-learnova/database/seeders"
	"github.com/tech-azim/be-learnova/repositories"
	"github.com/tech-azim/be-learnova/routes"
	"github.com/tech-azim/be-learnova/services"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		fmt.Println("=== CORS Middleware ===")
		fmt.Println("Origin:", origin)
		fmt.Println("Method:", c.Request.Method)
		fmt.Println("Path:", c.Request.URL.Path)

		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			fmt.Println("OPTIONS request - returning 204")
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("env not load")
	}

	seedFlag := flag.Bool("seed", false, "Run database seeders")
	flag.Parse()

	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(CORSMiddleware())

	r.RedirectTrailingSlash = true

	config.ConnectDB()

	if *seedFlag {
		seeders.RunAllSeeder(config.DB)
		return
	}

	// Initialize Repositories
	userRepo := repositories.NewUserRepository(config.DB)
	heroRepo := repositories.NewHeroRepository(config.DB)
	programRepo := repositories.NewProgramRepository(config.DB)
	registrationRepo := repositories.NewRegistrationRepository(config.DB)
	serviceRepo := repositories.NewServiceRepository(config.DB)
	portolioRepo := repositories.NewPortfolioRepository(config.DB)
	featureRepo := repositories.NewFeatureRepository(config.DB)
	galleryRepo := repositories.NewGalleryRepository(config.DB)
	videoGalleryRepo := repositories.NewVideoGalleryRepository(config.DB)
	flyerGalleryRepo := repositories.NewFlyerGalleryRepository(config.DB)
	dashboardRepo := repositories.NewDashboardRepository(config.DB)

	// Initialize Services
	authService := services.NewAuthService(userRepo)
	userService := services.NewUserService(userRepo) // NEW
	heroService := services.NewHeroService(heroRepo)
	programService := services.NewProgramService(programRepo)
	registrationService := services.RegistrationService(registrationRepo)
	serviceService := services.NewServiceService(serviceRepo)
	portfolioService := services.NewPortfolioService(portolioRepo)
	featureService := services.NewFeatureService(featureRepo)
	galleryService := services.NewGalleryService(galleryRepo)
	videoGalleryService := services.NewVideoGalleryService(videoGalleryRepo)
	flyerGalleryService := services.NewFlyerGalleryService(flyerGalleryRepo)
	dashboardService := services.NewDashboardService(dashboardRepo)

	// Initialize Controllers
	authController := controllers.NewAuthController(authService)
	userController := controllers.NewUserController(userService) // NEW
	heroController := controllers.NewHeroController(heroService)
	programController := controllers.NewProgramController(programService)
	registrationController := controllers.NewRegistrationController(registrationService, programService)
	serviceController := controllers.NewServiceController(serviceService)
	portfolioController := controllers.NewPortfolioController(portfolioService)
	featureController := controllers.NewFeatureController(featureService)
	galleryController := controllers.NewGalleryController(galleryService)
	videoGalleryController := controllers.NewVideoGalleryController(videoGalleryService)
	flyerGalleryController := controllers.NewFlyerGalleryController(flyerGalleryService)
	dashboardController := controllers.NewDashboardController(dashboardService)

	routes.Router(
		r,
		authController,
		heroController,
		programController,
		registrationController,
		serviceController,
		portfolioController,
		featureController,
		galleryController,
		videoGalleryController,
		flyerGalleryController,
		dashboardController,
		userController,
	)

	for _, route := range r.Routes() {
		fmt.Printf("Method: %s | Path: %s\n", route.Method, route.Path)
	}

	fmt.Println("\n🚀 Server starting on port 8080...")
	r.Run(":8080")
}
