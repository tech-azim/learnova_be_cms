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

type GalleryController struct {
	galleryService services.GalleryService
}

func NewGalleryController(galleryService services.GalleryService) *GalleryController {
	return &GalleryController{
		galleryService: galleryService,
	}
}

func (ctrl *GalleryController) Create(c *gin.Context) {
	// Ambil field dari form
	title := c.PostForm("title")
	description := c.PostForm("description")
	url := c.PostForm("url")
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

	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "URL is required",
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

	payload := models.Gallery{
		Title:       title,
		Description: description,
		URL:         url,
		Date:        dateTime,
		Order:       orderInt,
		IsActive:    isActiveBool,
	}

	gallery, err := ctrl.galleryService.Create(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create gallery",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data":    gallery,
		"message": "Gallery created successfully",
	})
}

func (ctrl *GalleryController) FindAll(c *gin.Context) {
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

	data, total, err := ctrl.galleryService.FindAll(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch galleries",
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

func (ctrl *GalleryController) FindAllActive(c *gin.Context) {
	data, err := ctrl.galleryService.FindAllActive()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch active galleries",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}

func (ctrl *GalleryController) FindByID(c *gin.Context) {
	id := c.Param("id")
	uint64Val, err := strconv.ParseUint(id, 10, 0)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid ID format",
			"error":   err.Error(),
		})
		return
	}

	data, err := ctrl.galleryService.FindByID(uint(uint64Val))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Gallery not found",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}

func (ctrl *GalleryController) Update(c *gin.Context) {
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

	// 2. Cek apakah gallery exist
	existingGallery, err := ctrl.galleryService.FindByID(uint(uint64Val))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Gallery not found",
			"error":   err.Error(),
		})
		return
	}

	// 3. Ambil field dari form
	title := c.PostForm("title")
	description := c.PostForm("description")
	url := c.PostForm("url")
	date := c.PostForm("date")
	order := c.PostForm("order")
	isActive := c.PostForm("is_active")

	// Gunakan nilai lama jika tidak ada input baru
	if title == "" {
		title = existingGallery.Title
	}
	if description == "" {
		description = existingGallery.Description
	}
	if url == "" {
		url = existingGallery.URL
	}

	// Parse date
	dateTime := existingGallery.Date
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
	orderInt := existingGallery.Order
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
	isActiveBool := existingGallery.IsActive
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
	payload := models.Gallery{
		ID:          uint(uint64Val),
		Title:       title,
		Description: description,
		URL:         url,
		Date:        dateTime,
		Order:       orderInt,
		IsActive:    isActiveBool,
	}

	// 5. Update ke database
	data, err := ctrl.galleryService.Update(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to update gallery",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    data,
		"message": "Gallery updated successfully",
	})
}

func (ctrl *GalleryController) Delete(c *gin.Context) {
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

	// 2. Cek apakah gallery exist
	existingGallery, err := ctrl.galleryService.FindByID(uint(uint64Val))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Gallery not found",
			"error":   err.Error(),
		})
		return
	}

	// 3. Soft delete gallery (set is_deleted = true)
	err = ctrl.galleryService.Delete(uint(uint64Val))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to delete gallery",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Gallery deleted successfully",
		"data": gin.H{
			"id":    existingGallery.ID,
			"title": existingGallery.Title,
		},
	})
}