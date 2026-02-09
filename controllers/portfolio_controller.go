package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tech-azim/be-learnova/models"
	"github.com/tech-azim/be-learnova/services"
	"github.com/tech-azim/be-learnova/utils"
)

type PortfolioRequest struct {
	Title       string `json:"title"`
	Count       string `json:"count"`
	Description string `json:"description"`
}

type PortfolioController struct {
	portfolioService services.PortfolioService
}

func NewPortfolioController(portfolioService services.PortfolioService) *PortfolioController {
	return &PortfolioController{
		portfolioService: portfolioService,
	}
}

func (ctrl *PortfolioController) Create(c *gin.Context) {
	// Ambil field dari form
	title := c.PostForm("title")
	count := c.PostForm("count")
	description := c.PostForm("description")

	// Validasi field wajib
	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Title is required",
		})
		return
	}

	if count == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Count is required",
		})
		return
	}

	if description == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Description is required",
		})
		return
	}

	payload := models.Portfolio{
		Title:       title,
		Count:       count,
		Description: description,
	}

	portfolio, err := ctrl.portfolioService.Create(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create portfolio",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data":    portfolio,
		"message": "Portfolio created successfully",
	})
}

func (ctrl *PortfolioController) FindAll(c *gin.Context) {
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

	data, total, err := ctrl.portfolioService.FindAll(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch portfolios",
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

func (ctrl *PortfolioController) FindByID(c *gin.Context) {
	id := c.Param("id")
	uint64Val, err := strconv.ParseUint(id, 10, 0)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid ID format",
			"error":   err.Error(),
		})
		return
	}

	data, err := ctrl.portfolioService.FindByID(uint(uint64Val))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Portfolio not found",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}

func (ctrl *PortfolioController) Update(c *gin.Context) {
	// 1. Ambil ID dari URL parameter
	id := c.Param("id")
	uint64Val, err := strconv.ParseUint(id, 10, 0)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid ID format",
			"error":   err.Error(),
		})
		return
	}

	// 2. Cek apakah portfolio exist
	existingPortfolio, err := ctrl.portfolioService.FindByID(uint(uint64Val))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Portfolio not found",
			"error":   err.Error(),
		})
		return
	}

	// 3. Ambil field dari form
	title := c.PostForm("title")
	count := c.PostForm("count")
	description := c.PostForm("description")

	// Gunakan nilai lama jika tidak ada input baru
	if title == "" {
		title = existingPortfolio.Title
	}
	if count == "" {
		count = existingPortfolio.Count
	}
	if description == "" {
		description = existingPortfolio.Description
	}

	// 4. Buat payload untuk update
	payload := models.Portfolio{
		ID:          uint(uint64Val),
		Title:       title,
		Count:       count,
		Description: description,
	}

	// 5. Update ke database
	data, err := ctrl.portfolioService.Update(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to update portfolio",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    data,
		"message": "Portfolio updated successfully",
	})
}

func (ctrl *PortfolioController) Delete(c *gin.Context) {
	// 1. Ambil ID dari URL parameter
	id := c.Param("id")
	uint64Val, err := strconv.ParseUint(id, 10, 0)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid ID format",
			"error":   err.Error(),
		})
		return
	}

	// 2. Cek apakah portfolio exist
	existingPortfolio, err := ctrl.portfolioService.FindByID(uint(uint64Val))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Portfolio not found",
			"error":   err.Error(),
		})
		return
	}

	// 3. Soft delete portfolio (set is_deleted = true)
	err = ctrl.portfolioService.Delete(uint(uint64Val))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to delete portfolio",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Portfolio deleted successfully",
		"data": gin.H{
			"id":    existingPortfolio.ID,
			"title": existingPortfolio.Title,
		},
	})
}