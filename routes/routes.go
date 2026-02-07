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
	programController *controllers.ProgramController,
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
			heroRoute.DELETE("/:id", middlewares.AuthMiddleware(), heroController.Delete)
			heroRoute.PUT(":id", middlewares.AuthMiddleware(), heroController.Update)
		}

		programRoute := api.Group("/programs")
		{
			programRoute.POST("", middlewares.AuthMiddleware(), programController.Create)
			programRoute.GET("", middlewares.AuthMiddleware(), programController.FindAll)
			programRoute.GET(":id", programController.FindByID)
			programRoute.DELETE(":id", middlewares.AuthMiddleware(), programController.Delete)
			programRoute.PUT(":id", middlewares.AuthMiddleware(), programController.Update)
		}
	}
}
