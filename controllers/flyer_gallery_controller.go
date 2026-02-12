package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tech-azim/be-learnova/models"
	"github.com/tech-azim/be-learnova/services"
	"github.com/tech-azim/be-learnova/utils"
)

type FlyerGalleryController struct {
	flyerGalleryService services.FlyerGalleryService
}

func NewFlyerGalleryController(flyerGalleryService services.FlyerGalleryService) *FlyerGalleryController {
	return &FlyerGalleryController{
		flyerGalleryService: flyerGalleryService,
	}
}

func (ctrl *FlyerGalleryController) Create(c *gin.Context) {
	// Ambil field dari form
	title := c.PostForm("title")
	image := c.PostForm("image")
	description := c.PostForm("description")
	isActive := c.PostForm("is_active")

	// Validasi field wajib
	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Title is required",
		})
		return
	}

	if image == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Image is required",
		})
		return
	}

	// Parse is_active (default true jika tidak ada)
	isActiveBool := true
	if isActive != "" {
		var err error
		isActiveBool, err = strconv.ParseBool(isActive)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid is_active format",
				"error":   err.Error(),
			})
			return
		}
	}

	payload := models.FlyerGallery{
		Title:       title,
		Image:       image,
		Description: description,
		IsActive:    isActiveBool,
	}

	flyerGallery, err := ctrl.flyerGalleryService.Create(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create flyer gallery",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data":    flyerGallery,
		"message": "Flyer gallery created successfully",
	})
}

func (ctrl *FlyerGalleryController) FindAll(c *gin.Context) {
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

	data, total, err := ctrl.flyerGalleryService.FindAll(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch flyer galleries",
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

func (ctrl *FlyerGalleryController) FindAllActive(c *gin.Context) {
	data, err := ctrl.flyerGalleryService.FindAllActive()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch active flyer galleries",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}

func (ctrl *FlyerGalleryController) FindByID(c *gin.Context) {
	id := c.Param("id")
	uint64Val, err := strconv.ParseUint(id, 10, 0)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid ID format",
			"error":   err.Error(),
		})
		return
	}

	data, err := ctrl.flyerGalleryService.FindByID(uint(uint64Val))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Flyer gallery not found",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}

func (ctrl *FlyerGalleryController) Update(c *gin.Context) {
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

	// 2. Cek apakah flyer gallery exist
	existingFlyerGallery, err := ctrl.flyerGalleryService.FindByID(uint(uint64Val))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Flyer gallery not found",
			"error":   err.Error(),
		})
		return
	}

	// 3. Ambil field dari form
	title := c.PostForm("title")
	image := c.PostForm("image")
	description := c.PostForm("description")
	isActive := c.PostForm("is_active")

	// Gunakan nilai lama jika tidak ada input baru
	if title == "" {
		title = existingFlyerGallery.Title
	}
	if image == "" {
		image = existingFlyerGallery.Image
	}
	if description == "" {
		description = existingFlyerGallery.Description
	}

	// Parse is_active
	isActiveBool := existingFlyerGallery.IsActive
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
	payload := models.FlyerGallery{
		ID:          uint(uint64Val),
		Title:       title,
		Image:       image,
		Description: description,
		IsActive:    isActiveBool,
	}

	// 5. Update ke database
	data, err := ctrl.flyerGalleryService.Update(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to update flyer gallery",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    data,
		"message": "Flyer gallery updated successfully",
	})
}

func (ctrl *FlyerGalleryController) Delete(c *gin.Context) {
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

	// 2. Cek apakah flyer gallery exist
	existingFlyerGallery, err := ctrl.flyerGalleryService.FindByID(uint(uint64Val))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Flyer gallery not found",
			"error":   err.Error(),
		})
		return
	}

	// 3. Soft delete flyer gallery (set is_deleted = true)
	err = ctrl.flyerGalleryService.Delete(uint(uint64Val))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to delete flyer gallery",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Flyer gallery deleted successfully",
		"data": gin.H{
			"id":    existingFlyerGallery.ID,
			"title": existingFlyerGallery.Title,
		},
	})
}