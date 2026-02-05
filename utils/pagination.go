package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type Pagination struct {
	Page int `json:"page"`
	Limit int `json:"limit"`
	TotalRows int `json:"total_rows"`
	TotalPages int `json:"total_pages"`
	Data interface{} `json:"data"`
}

type PaginationParams struct {
	Page int
	Limit int
}

func GetPaginationParams(c *gin.Context) PaginationParams {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// Validasi page minimal 1
	if page < 1 {
		page = 1
	}

	// Validasi limit (min 1, max 100)
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	return PaginationParams{
		Page:  page,
		Limit: limit,
	}
}