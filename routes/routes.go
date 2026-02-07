package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/tech-azim/be-learnova/controllers"
	"github.com/tech-azim/be-learnova/middlewares"
)

func Router(
    r *gin.Engine,
    authController *controllers.AuthController,
    heroController *controllers.HeroController,
) {
    r.Static("/uploads", "./uploads")
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
            heroRoute.GET("/:id", heroController.FindByID)
            heroRoute.DELETE("/:id", heroController.Delete)
            heroRoute.PUT(":id", heroController.Update)
        }
    }
}