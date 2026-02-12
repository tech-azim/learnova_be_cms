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

type GalleryController struct {
	galleryService services.GalleryService
}

func NewGalleryController(galleryService services.GalleryService) *GalleryController {
	return &GalleryController{
		galleryService: galleryService,
	}
}

func (ctrl *GalleryController) Create(c *gin.Context) {
	// 1. Ambil file
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "File is required",
			"error":   err.Error(),
		})
		return
	}

	// 2. Validasi tipe file
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

	// 3. Validasi ukuran file (max 5MB)
	maxSize := int64(5 * 1024 * 1024) // 5MB
	if file.Size > maxSize {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "File size too large. Maximum 5MB allowed",
		})
		return
	}

	// 4. Buat folder upload kalau belum ada
	uploadPath := "uploads/"
	if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create upload directory",
			"error":   err.Error(),
		})
		return
	}

	// 5. Generate nama file unik
	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
	filePath := uploadPath + filename

	// 6. Cek apakah file sudah ada
	if _, err := os.Stat(filePath); err == nil {
		randomSuffix := fmt.Sprintf("_%d", time.Now().UnixNano())
		filenameWithoutExt := strings.TrimSuffix(filename, filepath.Ext(filename))
		filename = fmt.Sprintf("%s%s%s", filenameWithoutExt, randomSuffix, ext)
		filePath = uploadPath + filename
	}

	// 7. Simpan file
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to save file",
			"error":   err.Error(),
		})
		return
	}

	// 8. Ambil field dari form
	title := c.PostForm("title")
	description := c.PostForm("description")
	date := c.PostForm("date")
	isActive := c.PostForm("is_active")

	// 9. Validasi field wajib
	if title == "" {
		os.Remove(filePath)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Title is required",
		})
		return
	}

	if date == "" {
		os.Remove(filePath)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Date is required",
		})
		return
	}

	// 10. Parse date
	dateTime, err := time.Parse("2006-01-02", date)
	if err != nil {
		os.Remove(filePath)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid date format. Use YYYY-MM-DD",
			"error":   err.Error(),
		})
		return
	}

	// 11. Parse is_active (default true)
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

	// 12. Buat payload
	payload := models.Gallery{
		Title:       title,
		Description: description,
		URL:         filePath, // Gunakan filePath sebagai URL
		Date:        dateTime,
		IsActive:    isActiveBool,
	}

	// 13. Simpan ke database
	gallery, err := ctrl.galleryService.Create(payload)
	if err != nil {
		os.Remove(filePath)
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

	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

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

	// 3. Handle file upload (opsional untuk update)
	filePath := existingGallery.URL
	oldFilePath := existingGallery.URL
	newFileUploaded := false

	file, err := c.FormFile("file")
	if err == nil { // Ada file baru
		// Validasi tipe file
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

		// Validasi ukuran file
		maxSize := int64(5 * 1024 * 1024) // 5MB
		if file.Size > maxSize {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "File size too large. Maximum 5MB allowed",
			})
			return
		}

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

		// Cek duplikasi
		if _, err := os.Stat(filePath); err == nil {
			randomSuffix := fmt.Sprintf("_%d", time.Now().UnixNano())
			filenameWithoutExt := strings.TrimSuffix(filename, filepath.Ext(filename))
			filename = fmt.Sprintf("%s%s%s", filenameWithoutExt, randomSuffix, ext)
			filePath = uploadPath + filename
		}

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

	// 4. Ambil field dari form
	title := c.PostForm("title")
	description := c.PostForm("description")
	date := c.PostForm("date")
	isActive := c.PostForm("is_active")

	// Gunakan nilai lama jika tidak ada input baru
	if title == "" {
		title = existingGallery.Title
	}
	if description == "" {
		description = existingGallery.Description
	}

	// Parse date
	dateTime := existingGallery.Date
	if date != "" {
		dateTime, err = time.Parse("2006-01-02", date)
		if err != nil {
			// Rollback file baru jika ada
			if newFileUploaded {
				os.Remove(filePath)
			}
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid date format. Use YYYY-MM-DD",
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
			// Rollback file baru jika ada
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

	// 5. Buat payload
	payload := models.Gallery{
		ID:          uint(uint64Val),
		Title:       title,
		Description: description,
		URL:         filePath,
		Date:        dateTime,
		IsActive:    isActiveBool,
	}

	// 6. Update ke database
	data, err := ctrl.galleryService.Update(payload)
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
			"message": "Failed to update gallery",
			"error":   err.Error(),
		})
		return
	}

	// 7. Hapus file lama jika ada file baru
	if newFileUploaded && oldFilePath != "" && oldFilePath != filePath {
		if err := os.Remove(oldFilePath); err != nil {
			fmt.Printf("Warning: Failed to delete old file %s: %v\n", oldFilePath, err)
		} else {
			fmt.Printf("Successfully deleted old file: %s\n", oldFilePath)
		}
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

	// 3. Delete dari database
	err = ctrl.galleryService.Delete(uint(uint64Val))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to delete gallery",
			"error":   err.Error(),
		})
		return
	}

	// 4. Hapus file fisik
	if existingGallery.URL != "" {
		if err := os.Remove(existingGallery.URL); err != nil {
			fmt.Printf("Warning: Failed to delete file %s: %v\n", existingGallery.URL, err)
		} else {
			fmt.Printf("Successfully deleted file: %s\n", existingGallery.URL)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Gallery deleted successfully",
		"data": gin.H{
			"id":    existingGallery.ID,
			"title": existingGallery.Title,
		},
	})
}