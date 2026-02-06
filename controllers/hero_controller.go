package controllers

import (
	"fmt"
	"net/http"
	"os"
	"time"

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

func (ctrl *HeroController) Create(c *gin.Context) {
	// ambil file
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "File is required",
			"error":   err.Error(),
		})
		return
	}

	// buat folder upload kalau belum ada
	uploadPath := "uploads/"
	if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create upload directory",
			"error":   err.Error(),
		})
		return
	}

	// nama file (biar gak bentrok)
	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
	filePath := uploadPath + filename

	// simpan file
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to save file",
			"error":   err.Error(),
		})
		return
	}

	// ambil field lain dari form
	title := c.PostForm("title")
	alt := c.PostForm("alt")
	description := c.PostForm("description")

	payload := models.Hero{
		SRC:         filePath, // path file
		ALT:         alt,
		Description: description,
		Title:       title,
	}

	hero, err := ctrl.heroService.Create(payload)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"message": "Failed",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"hero":    hero,
		"message": "Sukses",
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
	
	data,total, err := ctrl.heroService.FindAll(params)
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
			"total": total,
		},
	})
}