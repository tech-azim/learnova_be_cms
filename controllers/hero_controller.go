package controllers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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

	// Validasi tipe file (opsional tapi recommended)
	allowedExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}
	
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedExtensions[ext] {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid file type. Only jpg, jpeg, png, gif, webp allowed",
		})
		return
	}

	// Validasi ukuran file (misal max 5MB)
	maxSize := int64(5 * 1024 * 1024) // 5MB
	if file.Size > maxSize {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "File size too large. Maximum 5MB allowed",
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

	// VALIDASI: Cek apakah file sudah ada
	if _, err := os.Stat(filePath); err == nil {
		// File sudah ada, buat nama baru dengan suffix random
		randomSuffix := fmt.Sprintf("_%d", time.Now().UnixNano())
		filenameWithoutExt := strings.TrimSuffix(filename, filepath.Ext(filename))
		filename = fmt.Sprintf("%s%s%s", filenameWithoutExt, randomSuffix, ext)
		filePath = uploadPath + filename
	}

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
	alt := "Hero"
	description := c.PostForm("description")

	// Validasi field wajib
	if title == "" {
		// Hapus file yang sudah diupload
		os.Remove(filePath)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Title is required",
		})
		return
	}

	payload := models.Hero{
		SRC:         filePath, 
		ALT:         alt,
		Description: description,
		Title:       title,
	}

	hero, err := ctrl.heroService.Create(payload)
	if err != nil {
		os.Remove(filePath)
		
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create hero",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"hero":    hero,
		"message": "Hero created successfully",
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


func (ctrl *HeroController) FindById(c *gin.Context) {
	id := c.Param("id")

	uint64Val, err := strconv.ParseUint(id, 10, 0)

	data, err := ctrl.heroService.FindByID(uint(uint64Val))

	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
} 
func (ctrl *HeroController) Update(c *gin.Context) {
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

	// 2. Cek apakah hero exist
	existingHero, err := ctrl.heroService.FindByID(uint(uint64Val))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Hero not found",
			"error":   err.Error(),
		})
		return
	}

	
	filePath := existingHero.SRC 
	oldFilePath := existingHero.SRC 
	newFileUploaded := false 
	
	file, err := c.FormFile("file")
	
	if err == nil { // Kalau ada file baru
		

		uploadPath := "uploads/"
		if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to create upload directory",
				"error":   err.Error(),
			})
			return
		}

		// Nama file baru
		filename := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
		filePath = uploadPath + filename

		// Simpan file baru
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to save file",
				"error":   err.Error(),
			})
			return
		}

		newFileUploaded = true
		fmt.Printf("New file uploaded: %s\n", filePath)
	}

	title := c.PostForm("title")
	alt := "Hero"
	description := c.PostForm("description")

	// Gunakan nilai lama jika tidak ada input baru
	if title == "" {
		title = existingHero.Title
	}
	if alt == "" {
		alt = existingHero.ALT
	}
	if description == "" {
		description = existingHero.Description
	}

	payload := models.Hero{
		ID:          uint(uint64Val),
		SRC:         filePath,
		ALT:         alt,
		Description: description,
		Title:       title,
	}

	data, err := ctrl.heroService.Update(payload)
	if err != nil {
		if newFileUploaded && filePath != oldFilePath {
			if removeErr := os.Remove(filePath); removeErr != nil {
				fmt.Printf("Failed to rollback new file %s: %v\n", filePath, removeErr)
			} else {
				fmt.Printf("Rolled back: Deleted new file %s due to database error\n", filePath)
			}
		}
		
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to update hero",
			"error":   err.Error(),
		})
		return
	}

	if newFileUploaded && oldFilePath != "" && oldFilePath != filePath {
		if err := os.Remove(oldFilePath); err != nil {
			fmt.Printf("Warning: Failed to delete old file %s: %v\n", oldFilePath, err)
		} else {
			fmt.Printf("Successfully deleted old file: %s\n", oldFilePath)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    data,
		"message": "Hero updated successfully",
	})
}

func (ctrl *HeroController) Delete(c *gin.Context) {
	id := c.Param("id")
	uint64Val, err := strconv.ParseUint(id, 10, 0)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid ID format",
			"error":   err.Error(),
		})
		return
	}

	existingHero, err := ctrl.heroService.FindByID(uint(uint64Val))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Hero not found",
			"error":   err.Error(),
		})
		return
	}

	err = ctrl.heroService.Delete(uint(uint64Val))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to delete hero",
			"error":   err.Error(),
		})
		return
	}

	if existingHero.SRC != "" {
		if err := os.Remove(existingHero.SRC); err != nil {
			fmt.Printf("Warning: Failed to delete file %s: %v\n", existingHero.SRC, err)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Hero deleted successfully",
	})
}

func (ctrl *HeroController) FindByID(c *gin.Context){
	id := c.Param("id")
	uint64Val, err := strconv.ParseUint(id, 10, 0)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid ID format",
			"error":   err.Error(),
		})
		return
	}

	data, err := ctrl.heroService.FindByID(uint(uint64Val))

	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
	
}