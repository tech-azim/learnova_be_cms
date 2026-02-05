package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/tech-azim/be-learnova/controllers"
)

// Router adalah fungsi untuk setup semua routes aplikasi.
// Function ini menerima parameter:
//   - r: instance Gin Engine
//   - userController: controller untuk handle user endpoints
//   - authController: controller untuk handle auth endpoints
func Router(
	r *gin.Engine,
	authController *controllers.AuthController,
) {
	// Group untuk API versi 1
	// Semua route akan dimulai dengan /api/v1
	api := r.Group("/api/v1")
	{
		// Auth Routes - untuk login, register, dll
		// Base path: /api/v1/auth
		authRoute := api.Group("/auth")
		{
			// POST /api/v1/auth/login
			authRoute.POST("/login", authController.Login)
			
			// POST /api/v1/auth/register
		}

	}
}