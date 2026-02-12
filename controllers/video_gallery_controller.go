package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tech-azim/be-learnova/models"
	"github.com/tech-azim/be-learnova/services"
	"github.com/tech-azim/be-learnova/utils"
)

type VideoGalleryController struct {
	videoGalleryService services.VideoGalleryService
}

func NewVideoGalleryController(videoGalleryService services.VideoGalleryService) *VideoGalleryController {
	return &VideoGalleryController{
		videoGalleryService: videoGalleryService,
	}
}

func (ctrl *VideoGalleryController) Create(c *gin.Context) {
	// Ambil field dari form
	title := c.PostForm("title")
	description := c.PostForm("description")
	thumbnail := c.PostForm("thumbnail")
	videoURL := c.PostForm("video_url")
	category := c.PostForm("category")
	date := c.PostForm("date")
	order := c.PostForm("order")
	isActive := c.PostForm("is_active")

	// Validasi field wajib
	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Title is required",
		})
		return
	}

	if thumbnail == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Thumbnail is required",
		})
		return
	}

	if videoURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Video URL is required",
		})
		return
	}

	if category == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Category is required",
		})
		return
	}

	if date == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Date is required",
		})
		return
	}

	// Parse date (format: YYYY-MM-DD)
	dateTime, err := time.Parse("2006-01-02", date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid date format. Use YYYY-MM-DD",
			"error":   err.Error(),
		})
		return
	}

	// Parse order (default 0 jika tidak ada)
	orderInt := 0
	if order != "" {
		orderInt, err = strconv.Atoi(order)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid order format",
				"error":   err.Error(),
			})
			return
		}
	}

	// Parse is_active (default true jika tidak ada)
	isActiveBool := true
	if isActive != "" {
		isActiveBool, err = strconv.ParseBool(isActive)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid is_active format",
				"error":   err.Error(),
			})
			return
		}
	}

	payload := models.VideoGallery{
		Title:       title,
		Description: description,
		Thumbnail:   thumbnail,
		VideoURL:    videoURL,
		Category:    category,
		Date:        dateTime,
		Order:       orderInt,
		IsActive:    isActiveBool,
	}

	videoGallery, err := ctrl.videoGalleryService.Create(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create video gallery",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data":    videoGallery,
		"message": "Video gallery created successfully",
	})
}

func (ctrl *VideoGalleryController) FindAll(c *gin.Context) {
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

	data, total, err := ctrl.videoGalleryService.FindAll(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch video galleries",
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

func (ctrl *VideoGalleryController) FindAllActive(c *gin.Context) {
	data, err := ctrl.videoGalleryService.FindAllActive()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch active video galleries",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}

func (ctrl *VideoGalleryController) FindByCategory(c *gin.Context) {
	category := c.Query("category")
	
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

	data, total, err := ctrl.videoGalleryService.FindByCategory(category, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch video galleries by category",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": data,
		"pagination": gin.H{
			"page":     params.Page,
			"limit":    params.Limit,
			"total":    total,
			"category": category,
		},
	})
}

func (ctrl *VideoGalleryController) FindAllCategories(c *gin.Context) {
	data, err := ctrl.videoGalleryService.FindAllCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch categories",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}

func (ctrl *VideoGalleryController) FindByID(c *gin.Context) {
	id := c.Param("id")
	uint64Val, err := strconv.ParseUint(id, 10, 0)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid ID format",
			"error":   err.Error(),
		})
		return
	}

	data, err := ctrl.videoGalleryService.FindByID(uint(uint64Val))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Video gallery not found",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}

func (ctrl *VideoGalleryController) Update(c *gin.Context) {
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

	// 2. Cek apakah video gallery exist
	existingVideoGallery, err := ctrl.videoGalleryService.FindByID(uint(uint64Val))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Video gallery not found",
			"error":   err.Error(),
		})
		return
	}

	// 3. Ambil field dari form
	title := c.PostForm("title")
	description := c.PostForm("description")
	thumbnail := c.PostForm("thumbnail")
	videoURL := c.PostForm("video_url")
	category := c.PostForm("category")
	date := c.PostForm("date")
	order := c.PostForm("order")
	isActive := c.PostForm("is_active")

	// Gunakan nilai lama jika tidak ada input baru
	if title == "" {
		title = existingVideoGallery.Title
	}
	if description == "" {
		description = existingVideoGallery.Description
	}
	if thumbnail == "" {
		thumbnail = existingVideoGallery.Thumbnail
	}
	if videoURL == "" {
		videoURL = existingVideoGallery.VideoURL
	}
	if category == "" {
		category = existingVideoGallery.Category
	}

	// Parse date
	dateTime := existingVideoGallery.Date
	if date != "" {
		dateTime, err = time.Parse("2006-01-02", date)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid date format. Use YYYY-MM-DD",
				"error":   err.Error(),
			})
			return
		}
	}

	// Parse order
	orderInt := existingVideoGallery.Order
	if order != "" {
		orderInt, err = strconv.Atoi(order)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid order format",
				"error":   err.Error(),
			})
			return
		}
	}

	// Parse is_active
	isActiveBool := existingVideoGallery.IsActive
	if isActive != "" {
		isActiveBool, err = strconv.ParseBool(isActive)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid is_active format",
				"error":   err.Error(),
			})
			return
		}
	}

	// 4. Buat payload untuk update
	payload := models.VideoGallery{
		ID:          uint(uint64Val),
		Title:       title,
		Description: description,
		Thumbnail:   thumbnail,
		VideoURL:    videoURL,
		Category:    category,
		Date:        dateTime,
		Order:       orderInt,
		IsActive:    isActiveBool,
	}

	// 5. Update ke database
	data, err := ctrl.videoGalleryService.Update(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to update video gallery",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    data,
		"message": "Video gallery updated successfully",
	})
}

func (ctrl *VideoGalleryController) Delete(c *gin.Context) {
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

	// 2. Cek apakah video gallery exist
	existingVideoGallery, err := ctrl.videoGalleryService.FindByID(uint(uint64Val))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Video gallery not found",
			"error":   err.Error(),
		})
		return
	}

	// 3. Soft delete video gallery (set is_deleted = true)
	err = ctrl.videoGalleryService.Delete(uint(uint64Val))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to delete video gallery",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Video gallery deleted successfully",
		"data": gin.H{
			"id":    existingVideoGallery.ID,
			"title": existingVideoGallery.Title,
		},
	})
}