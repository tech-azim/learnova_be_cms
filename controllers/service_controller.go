package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tech-azim/be-learnova/models"
	"github.com/tech-azim/be-learnova/services"
	"github.com/tech-azim/be-learnova/utils"
)


type ServiceController struct {
	serviceService services.ServiceService
}

func NewServiceController(serviceService services.ServiceService) *ServiceController {
	return &ServiceController{
		serviceService: serviceService,
	}
}

func (ctrl *ServiceController) Create(c *gin.Context) {
	// Ambil field dari form
	icon := c.PostForm("icon")
	title := c.PostForm("title")
	description := c.PostForm("description")
	color := c.PostForm("color")

	// Validasi field wajib
	if icon == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Icon is required",
		})
		return
	}

	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Title is required",
		})
		return
	}

	if description == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Description is required",
		})
		return
	}

	if color == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Color is required",
		})
		return
	}

	payload := models.Service{
		Icon:        icon,
		Title:       title,
		Description: description,
		Color:       color,
	}

	service, err := ctrl.serviceService.Create(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create service",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data":    service,
		"message": "Service created successfully",
	})
}

func (ctrl *ServiceController) FindAll(c *gin.Context) {
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

	data, total, err := ctrl.serviceService.FindAll(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch services",
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

func (ctrl *ServiceController) FindByID(c *gin.Context) {
	id := c.Param("id")
	uint64Val, err := strconv.ParseUint(id, 10, 0)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid ID format",
			"error":   err.Error(),
		})
		return
	}

	data, err := ctrl.serviceService.FindByID(uint(uint64Val))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Service not found",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}

func (ctrl *ServiceController) Update(c *gin.Context) {
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

	// 2. Cek apakah service exist
	existingService, err := ctrl.serviceService.FindByID(uint(uint64Val))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Service not found",
			"error":   err.Error(),
		})
		return
	}

	// 3. Ambil field dari form
	icon := c.PostForm("icon")
	title := c.PostForm("title")
	description := c.PostForm("description")
	color := c.PostForm("color")

	// Gunakan nilai lama jika tidak ada input baru
	if icon == "" {
		icon = existingService.Icon
	}
	if title == "" {
		title = existingService.Title
	}
	if description == "" {
		description = existingService.Description
	}
	if color == "" {
		color = existingService.Color
	}

	// 4. Buat payload untuk update
	payload := models.Service{
		ID:          uint(uint64Val),
		Icon:        icon,
		Title:       title,
		Description: description,
		Color:       color,
	}

	// 5. Update ke database
	data, err := ctrl.serviceService.Update(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to update service",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    data,
		"message": "Service updated successfully",
	})
}

func (ctrl *ServiceController) Delete(c *gin.Context) {
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

	// 2. Cek apakah service exist
	existingService, err := ctrl.serviceService.FindByID(uint(uint64Val))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Service not found",
			"error":   err.Error(),
		})
		return
	}

	// 3. Soft delete service (set is_deleted = true)
	err = ctrl.serviceService.Delete(uint(uint64Val))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to delete service",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Service deleted successfully",
		"data": gin.H{
			"id":    existingService.ID,
			"title": existingService.Title,
		},
	})
}