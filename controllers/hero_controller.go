package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tech-azim/be-learnova/models"
	"github.com/tech-azim/be-learnova/services"
	"github.com/tech-azim/be-learnova/utils"
)

type HeroRequest struct  {
	SRC string `json:"src"`
	ALT string `json:"alt"`
	Title string `json:"title"`
	Descipriotn string `json:"description"`
}

type HeroController struct {
	heroService services.HeroService
}

func NewHeroController(heroService services.HeroService) *HeroController {
	return &HeroController{
		heroService: heroService,
	}
}

func (ctrl *HeroController) Create(c *gin.Context){
	var req HeroRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	payload := models.Hero{
		SRC: req.SRC,
		ALT: req.ALT,
		Description: req.Descipriotn,
		Title: req.Title,
	}

	hero, err := ctrl.heroService.Create(payload)

	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H {
			"message":"Failed",
			"error": err,
		})
	}

	c.JSON(http.StatusCreated, gin.H {
		"hero": hero,
		"message":"Sukses",
	})
}

func (ctrl *HeroController) FindAll(c *gin.Context) {
	var params utils.PaginationParams
	
	// Bind query parameters
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	
	// Default values
	if params.Page == 0 {
		params.Page = 1
	}
	if params.Limit == 0 {
		params.Limit = 10
	}
	
	data, err := ctrl.heroService.FindAll(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch heroes",
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"data": data,
		"pagination": gin.H{
			"page":  params.Page,
			"limit": params.Limit,
		},
	})
}