package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/tech-azim/be-learnova/controllers"
	"github.com/tech-azim/be-learnova/middlewares"
)

// Router adalah fungsi untuk setup semua routes aplikasi.
// Function ini menerima parameter:
//   - r: instance Gin Engine
//   - userController: controller untuk handle user endpoints
//   - authController: controller untuk handle auth endpoints
func Router(
    r *gin.Engine,
    authController *controllers.AuthController,
    heroController *controllers.HeroController,
) {
    api := r.Group("/api/v1")
    {
        authRoute := api.Group("/auth")
        {
            authRoute.POST("/login", authController.Login)
        }

        // Jangan gunakan slash di Group jika ingin path bersih
        heroRoute := api.Group("/heros")
        {
            // POST dengan auth middleware
            heroRoute.POST("", middlewares.AuthMiddleware(), heroController.Create)
            
            // GET tanpa auth
            heroRoute.GET("", heroController.FindAll)
        }
    }
}