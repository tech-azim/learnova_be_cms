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
	serviceController *controllers.ServiceController,
	portfolioController *controllers.PortfolioController,
	featureController *controllers.FeatureController,
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
			heroRoute.PUT("/:id", middlewares.AuthMiddleware(), heroController.Update)
		}

		programRoute := api.Group("/programs")
		{
			programRoute.POST("", middlewares.AuthMiddleware(), programController.Create)
			programRoute.GET("", middlewares.AuthMiddleware(), programController.FindAll)
			programRoute.GET("/:id", programController.FindByID)
			programRoute.DELETE("/:id", middlewares.AuthMiddleware(), programController.Delete)
			programRoute.PUT("/:id", middlewares.AuthMiddleware(), programController.Update)
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

		serviceRoute := api.Group("/services")
		{
			serviceRoute.POST("", middlewares.AuthMiddleware(), serviceController.Create)
			serviceRoute.GET("", serviceController.FindAll)
			serviceRoute.GET("/:id", serviceController.FindByID)
			serviceRoute.PUT("/:id", middlewares.AuthMiddleware(), serviceController.Update)
			serviceRoute.DELETE("/:id", middlewares.AuthMiddleware(), serviceController.Delete)
		}

		portfolioRoute := api.Group("/portfolios")
		{
			portfolioRoute.POST("", middlewares.AuthMiddleware(), portfolioController.Create)
			portfolioRoute.GET("", portfolioController.FindAll)
			portfolioRoute.GET("/:id", portfolioController.FindByID)
			portfolioRoute.PUT("/:id", middlewares.AuthMiddleware(), portfolioController.Update)
			portfolioRoute.DELETE("/:id", middlewares.AuthMiddleware(), portfolioController.Delete)
		}

		featureRoute := api.Group("/features")
		{
			// PUBLIC endpoints
			featureRoute.GET("", featureController.FindAll)
			featureRoute.GET("/active", featureController.FindAllActive) // endpoint khusus untuk feature yang aktif
			featureRoute.GET("/:id", featureController.FindByID)

			// PROTECTED endpoints (perlu authentication)
			featureRoute.POST("", middlewares.AuthMiddleware(), featureController.Create)
			featureRoute.PUT("/:id", middlewares.AuthMiddleware(), featureController.Update)
			featureRoute.DELETE("/:id", middlewares.AuthMiddleware(), featureController.Delete)
		}
	}
}