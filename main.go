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

		// Log untuk debugging
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

	// Buat instance Gin tanpa middleware default
	r := gin.New()

	// Tambahkan middleware secara manual dengan urutan yang benar:
	r.Use(gin.Logger())     // 1. Logger
	r.Use(gin.Recovery())   // 2. Recovery
	r.Use(CORSMiddleware()) // 3. CORS - HARUS SEBELUM ROUTES!

	r.RedirectTrailingSlash = true

	config.ConnectDB()

	if *seedFlag {
		seeders.RunAllSeeder(config.DB)
		return
	}

	userRepo := repositories.NewUserRepository(config.DB)
	heroRepo := repositories.NewHeroRepository(config.DB)
	programRepo := repositories.NewProgramRepository(config.DB)
	registrationRepo := repositories.NewRegistrationRepository(config.DB)

	authService := services.NewAuthService(userRepo)
	heroService := services.NewHeroService(heroRepo)
	programService := services.NewProgramService(programRepo)
	registrationService := services.RegistrationService(registrationRepo)

	authController := controllers.NewAuthController(authService)
	heroController := controllers.NewHeroController(heroService)
	programController := controllers.NewProgramController(programService)
	registrationController := controllers.NewRegistrationController(registrationService, programService)

	// Routes dipanggil SETELAH semua middleware global
	routes.Router(r, authController, heroController, programController, registrationController)

	for _, route := range r.Routes() {
		fmt.Printf("Method: %s | Path: %s\n", route.Method, route.Path)
	}

	fmt.Println("\n🚀 Server starting on port 8080...")
	r.Run(":8080")
}
