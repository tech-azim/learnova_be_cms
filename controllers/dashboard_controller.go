package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tech-azim/be-learnova/services"
)

type DashboardController struct {
	service services.DashboardService
}

func NewDashboardController(service services.DashboardService) *DashboardController {
	return &DashboardController{service}
}

func (c *DashboardController) GetDashboard(ctx *gin.Context) {
	data, err := c.service.GetDashboardData()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to fetch dashboard data",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Dashboard data fetched successfully",
		"data":    data,
	})
}
