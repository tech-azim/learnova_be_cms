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
	galleryController *controllers.GalleryController,
	videoGalleryController *controllers.VideoGalleryController,
	flyerGalleryController *controllers.FlyerGalleryController,
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
			featureRoute.GET("/active", featureController.FindAllActive)
			featureRoute.GET("/:id", featureController.FindByID)

			// PROTECTED endpoints (perlu authentication)
			featureRoute.POST("", middlewares.AuthMiddleware(), featureController.Create)
			featureRoute.PUT("/:id", middlewares.AuthMiddleware(), featureController.Update)
			featureRoute.DELETE("/:id", middlewares.AuthMiddleware(), featureController.Delete)
		}

		galleryRoute := api.Group("/galleries")
		{
			// PUBLIC endpoints
			galleryRoute.GET("", galleryController.FindAll)
			galleryRoute.GET("/active", galleryController.FindAllActive)
			galleryRoute.GET("/:id", galleryController.FindByID)

			// PROTECTED endpoints (perlu authentication)
			galleryRoute.POST("", middlewares.AuthMiddleware(), galleryController.Create)
			galleryRoute.PUT("/:id", middlewares.AuthMiddleware(), galleryController.Update)
			galleryRoute.DELETE("/:id", middlewares.AuthMiddleware(), galleryController.Delete)
		}

		videoGalleryRoute := api.Group("/video-galleries")
		{
			// PUBLIC endpoints
			videoGalleryRoute.GET("", videoGalleryController.FindAll)
			videoGalleryRoute.GET("/active", videoGalleryController.FindAllActive)
			videoGalleryRoute.GET("/categories", videoGalleryController.FindAllCategories)
			videoGalleryRoute.GET("/by-category", videoGalleryController.FindByCategory) // ?category=Semua&page=1&limit=10
			videoGalleryRoute.GET("/:id", videoGalleryController.FindByID)

			// PROTECTED endpoints (perlu authentication)
			videoGalleryRoute.POST("", middlewares.AuthMiddleware(), videoGalleryController.Create)
			videoGalleryRoute.PUT("/:id", middlewares.AuthMiddleware(), videoGalleryController.Update)
			videoGalleryRoute.DELETE("/:id", middlewares.AuthMiddleware(), videoGalleryController.Delete)
		}

		flyerGalleryRoute := api.Group("/flyer-galleries")
		{
			// PUBLIC endpoints
			flyerGalleryRoute.GET("", flyerGalleryController.FindAll)
			flyerGalleryRoute.GET("/active", flyerGalleryController.FindAllActive)
			flyerGalleryRoute.GET("/:id", flyerGalleryController.FindByID)

			// PROTECTED endpoints (perlu authentication)
			flyerGalleryRoute.POST("", middlewares.AuthMiddleware(), flyerGalleryController.Create)
			flyerGalleryRoute.PUT("/:id", middlewares.AuthMiddleware(), flyerGalleryController.Update)
			flyerGalleryRoute.DELETE("/:id", middlewares.AuthMiddleware(), flyerGalleryController.Delete)
		}
	}
}