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
	registrationController *controllers.RegistrationController,
) {
	r.Static("/uploads", "./uploads")
	api := r.Group("/api/v1")
	{
		authRoute := api.Group("/auth")
		{
			authRoute.POST("/login", authController.Login)
		}

		heroRoute := api.Group("/heros")
		{
			heroRoute.POST("", middlewares.AuthMiddleware(), heroController.Create)
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

		registrationRoute := api.Group("/registrations")
		{
			// PUBLIC (tanpa auth)
			registrationRoute.POST("", registrationController.Create)

			// PROTECTED (pakai auth)
			registrationRoute.GET("", middlewares.AuthMiddleware(), registrationController.FindAll)
			registrationRoute.GET("/:id", middlewares.AuthMiddleware(), registrationController.FindByID)
			registrationRoute.GET("/program/:programId", middlewares.AuthMiddleware(), registrationController.FindByProgramID)
			registrationRoute.GET("/by-email", middlewares.AuthMiddleware(), registrationController.FindByEmail)
			registrationRoute.PUT("/:id", middlewares.AuthMiddleware(), registrationController.Update)
			registrationRoute.DELETE("/:id", middlewares.AuthMiddleware(), registrationController.Delete)
		}
	}
}
