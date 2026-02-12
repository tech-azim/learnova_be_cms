package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tech-azim/be-learnova/models"
	"github.com/tech-azim/be-learnova/services"
	"github.com/tech-azim/be-learnova/utils"
)

type FeatureController struct {
	featureService services.FeatureService
}

func NewFeatureController(featureService services.FeatureService) *FeatureController {
	return &FeatureController{
		featureService: featureService,
	}
}

func (ctrl *FeatureController) Create(c *gin.Context) {
	// Ambil field dari form
	icon := c.PostForm("icon")
	title := c.PostForm("title")
	description := c.PostForm("description")
	order := c.PostForm("order")
	isActive := c.PostForm("is_active")

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

	// Parse order (default 0 jika tidak ada)
	orderInt := 0
	if order != "" {
		var err error
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

	payload := models.Feature{
		Icon:        icon,
		Title:       title,
		Description: description,
		Order:       orderInt,
		IsActive:    isActiveBool,
	}

	feature, err := ctrl.featureService.Create(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create feature",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data":    feature,
		"message": "Feature created successfully",
	})
}

func (ctrl *FeatureController) FindAll(c *gin.Context) {
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

	data, total, err := ctrl.featureService.FindAll(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch features",
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

func (ctrl *FeatureController) FindAllActive(c *gin.Context) {
	data, err := ctrl.featureService.FindAllActive()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch active features",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}

func (ctrl *FeatureController) FindByID(c *gin.Context) {
	id := c.Param("id")
	uint64Val, err := strconv.ParseUint(id, 10, 0)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid ID format",
			"error":   err.Error(),
		})
		return
	}

	data, err := ctrl.featureService.FindByID(uint(uint64Val))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Feature not found",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}

func (ctrl *FeatureController) Update(c *gin.Context) {
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

	// 2. Cek apakah feature exist
	existingFeature, err := ctrl.featureService.FindByID(uint(uint64Val))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Feature not found",
			"error":   err.Error(),
		})
		return
	}

	// 3. Ambil field dari form
	icon := c.PostForm("icon")
	title := c.PostForm("title")
	description := c.PostForm("description")
	order := c.PostForm("order")
	isActive := c.PostForm("is_active")

	// Gunakan nilai lama jika tidak ada input baru
	if icon == "" {
		icon = existingFeature.Icon
	}
	if title == "" {
		title = existingFeature.Title
	}
	if description == "" {
		description = existingFeature.Description
	}

	// Parse order
	orderInt := existingFeature.Order
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
	isActiveBool := existingFeature.IsActive
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
	payload := models.Feature{
		ID:          uint(uint64Val),
		Icon:        icon,
		Title:       title,
		Description: description,
		Order:       orderInt,
		IsActive:    isActiveBool,
	}

	// 5. Update ke database
	data, err := ctrl.featureService.Update(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to update feature",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    data,
		"message": "Feature updated successfully",
	})
}

func (ctrl *FeatureController) Delete(c *gin.Context) {
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

	// 2. Cek apakah feature exist
	existingFeature, err := ctrl.featureService.FindByID(uint(uint64Val))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Feature not found",
			"error":   err.Error(),
		})
		return
	}

	// 3. Soft delete feature (set is_deleted = true)
	err = ctrl.featureService.Delete(uint(uint64Val))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to delete feature",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Feature deleted successfully",
		"data": gin.H{
			"id":    existingFeature.ID,
			"title": existingFeature.Title,
		},
	})
}