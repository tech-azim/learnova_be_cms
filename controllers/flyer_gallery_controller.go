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

type FlyerGalleryController struct {
	flyerGalleryService services.FlyerGalleryService
}

func NewFlyerGalleryController(flyerGalleryService services.FlyerGalleryService) *FlyerGalleryController {
	return &FlyerGalleryController{
		flyerGalleryService: flyerGalleryService,
	}
}

// validateImageFile validates file extension and size, returns ext and error message
func validateImageFile(filename string, size int64) (string, string) {
	allowedExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}

	ext := strings.ToLower(filepath.Ext(filename))
	if !allowedExtensions[ext] {
		return "", "Invalid file type. Only jpg, jpeg, png, gif, webp allowed"
	}

	maxSize := int64(5 * 1024 * 1024) // 5MB
	if size > maxSize {
		return "", "File size too large. Maximum 5MB allowed"
	}

	return ext, ""
}

// saveUploadedFile saves the uploaded file to uploads/ and returns the file path
func saveUploadedFile(c *gin.Context, filename string, ext string) (string, error) {
	uploadPath := "uploads/"
	if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %w", err)
	}

	generatedName := fmt.Sprintf("%d_%s", time.Now().Unix(), filename)
	filePath := uploadPath + generatedName

	// Handle duplicate filename
	if _, err := os.Stat(filePath); err == nil {
		randomSuffix := fmt.Sprintf("_%d", time.Now().UnixNano())
		nameWithoutExt := strings.TrimSuffix(generatedName, filepath.Ext(generatedName))
		generatedName = fmt.Sprintf("%s%s%s", nameWithoutExt, randomSuffix, ext)
		filePath = uploadPath + generatedName
	}

	return filePath, nil
}

func (ctrl *FlyerGalleryController) Create(c *gin.Context) {
	// 1. Ambil file
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "File is required",
			"error":   err.Error(),
		})
		return
	}

	// 2. Validasi tipe & ukuran file
	ext, errMsg := validateImageFile(file.Filename, file.Size)
	if errMsg != "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": errMsg})
		return
	}

	// 3. Generate path file
	filePath, err := saveUploadedFile(c, file.Filename, ext)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	// 4. Simpan file
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to save file",
			"error":   err.Error(),
		})
		return
	}

	// 5. Ambil field dari form
	title := c.PostForm("title")
	description := c.PostForm("description")
	isActive := c.PostForm("is_active")

	// 6. Validasi field wajib
	if title == "" {
		os.Remove(filePath)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Title is required"})
		return
	}

	// 7. Parse is_active (default true)
	isActiveBool := true
	if isActive != "" {
		isActiveBool, err = strconv.ParseBool(isActive)
		if err != nil {
			os.Remove(filePath)
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid is_active format",
				"error":   err.Error(),
			})
			return
		}
	}

	// 8. Buat payload & simpan ke database
	payload := models.FlyerGallery{
		Title:       title,
		Image:       filePath,
		Description: description,
		IsActive:    isActiveBool,
	}

	flyerGallery, err := ctrl.flyerGalleryService.Create(payload)
	if err != nil {
		os.Remove(filePath)
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

	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if params.Page == 0 {
		params.Page = 1
	}
	if params.Limit == 0 {
		params.Limit = 10
	}

	data, total, err := ctrl.flyerGalleryService.FindAll(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch flyer galleries"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch active flyer galleries"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": data})
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

	c.JSON(http.StatusOK, gin.H{"data": data})
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

	// 3. Handle file upload — gunakan key "image" (konsisten dengan Create & frontend)
	filePath := existingFlyerGallery.Image
	oldFilePath := existingFlyerGallery.Image
	newFileUploaded := false

	file, err := c.FormFile("image") // ← FIX: was "file", sekarang "image"
	if err == nil {                   // Ada file baru diupload
		// Validasi tipe & ukuran
		ext, errMsg := validateImageFile(file.Filename, file.Size)
		if errMsg != "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": errMsg})
			return
		}

		// Generate path file baru
		newFilePath, err := saveUploadedFile(c, file.Filename, ext)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		// Simpan file baru
		if err := c.SaveUploadedFile(file, newFilePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to save file",
				"error":   err.Error(),
			})
			return
		}

		filePath = newFilePath
		newFileUploaded = true
		fmt.Printf("New file uploaded: %s\n", filePath)
	}

	// 4. Ambil field dari form, gunakan nilai lama jika kosong
	title := c.PostForm("title")
	description := c.PostForm("description")
	isActive := c.PostForm("is_active")

	if title == "" {
		title = existingFlyerGallery.Title
	}
	if description == "" {
		description = existingFlyerGallery.Description
	}

	// 5. Parse is_active
	isActiveBool := existingFlyerGallery.IsActive
	if isActive != "" {
		isActiveBool, err = strconv.ParseBool(isActive)
		if err != nil {
			if newFileUploaded {
				os.Remove(filePath)
			}
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid is_active format",
				"error":   err.Error(),
			})
			return
		}
	}

	// 6. Buat payload & update ke database
	payload := models.FlyerGallery{
		ID:          uint(uint64Val),
		Title:       title,
		Image:       filePath,
		Description: description,
		IsActive:    isActiveBool,
	}

	data, err := ctrl.flyerGalleryService.Update(payload)
	if err != nil {
		// Rollback: hapus file baru jika gagal update database
		if newFileUploaded && filePath != oldFilePath {
			if removeErr := os.Remove(filePath); removeErr != nil {
				fmt.Printf("Failed to rollback new file %s: %v\n", filePath, removeErr)
			} else {
				fmt.Printf("Rolled back: Deleted new file %s due to database error\n", filePath)
			}
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to update flyer gallery",
			"error":   err.Error(),
		})
		return
	}

	// 7. Hapus file lama jika berhasil upload file baru
	if newFileUploaded && oldFilePath != "" && oldFilePath != filePath {
		if err := os.Remove(oldFilePath); err != nil {
			fmt.Printf("Warning: Failed to delete old file %s: %v\n", oldFilePath, err)
		} else {
			fmt.Printf("Successfully deleted old file: %s\n", oldFilePath)
		}
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

	// 3. Delete dari database
	if err = ctrl.flyerGalleryService.Delete(uint(uint64Val)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to delete flyer gallery",
			"error":   err.Error(),
		})
		return
	}

	// 4. Hapus file fisik
	if existingFlyerGallery.Image != "" {
		if err := os.Remove(existingFlyerGallery.Image); err != nil {
			fmt.Printf("Warning: Failed to delete file %s: %v\n", existingFlyerGallery.Image, err)
		} else {
			fmt.Printf("Successfully deleted file: %s\n", existingFlyerGallery.Image)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Flyer gallery deleted successfully",
		"data": gin.H{
			"id":    existingFlyerGallery.ID,
			"title": existingFlyerGallery.Title,
		},
	})
}
