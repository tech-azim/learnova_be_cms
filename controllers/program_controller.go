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

type ProgramRequest struct {
	Icon         string   `json:"icon"`
	Title        string   `json:"title"`
	Duration     string   `json:"duration"`
	Participants string   `json:"participants"`
	Level        string   `json:"level"`
	Description  string   `json:"description"`
	Benefits     []string `json:"benefits"`
	Image        string   `json:"image"`
}

type ProgramController struct {
	programService services.ProgramService
}

func NewProgramController(programService services.ProgramService) *ProgramController {
	return &ProgramController{
		programService: programService,
	}
}

func (ctrl *ProgramController) Create(c *gin.Context) {
	// Ambil file image
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "File is required",
			"error":   err.Error(),
		})
		return
	}

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

	// Validasi ukuran file (max 5MB)
	maxSize := int64(5 * 1024 * 1024) // 5MB
	if file.Size > maxSize {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "File size too large. Maximum 5MB allowed",
		})
		return
	}

	// Buat folder upload kalau belum ada
	uploadPath := "uploads/programs/"
	if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create upload directory",
			"error":   err.Error(),
		})
		return
	}

	// Generate nama file unik
	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
	filePath := uploadPath + filename

	// Validasi: Cek apakah file sudah ada
	if _, err := os.Stat(filePath); err == nil {
		randomSuffix := fmt.Sprintf("_%d", time.Now().UnixNano())
		filenameWithoutExt := strings.TrimSuffix(filename, filepath.Ext(filename))
		filename = fmt.Sprintf("%s%s%s", filenameWithoutExt, randomSuffix, ext)
		filePath = uploadPath + filename
	}

	// Simpan file
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to save file",
			"error":   err.Error(),
		})
		return
	}

	// Ambil field lain dari form
	icon := c.PostForm("icon")
	title := c.PostForm("title")
	duration := c.PostForm("duration")
	participants := c.PostForm("participants")
	level := c.PostForm("level")
	description := c.PostForm("description")

	// Parse benefits (array string)
	benefitsStr := c.PostFormArray("benefits")
	if len(benefitsStr) == 0 {
		// Coba parse dari single string dengan delimiter
		benefitsInput := c.PostForm("benefits")
		if benefitsInput != "" {
			benefitsStr = strings.Split(benefitsInput, ",")
			// Trim whitespace
			for i := range benefitsStr {
				benefitsStr[i] = strings.TrimSpace(benefitsStr[i])
			}
		}
	}

	// Validasi field wajib
	if title == "" {
		os.Remove(filePath)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Title is required",
		})
		return
	}

	if duration == "" {
		os.Remove(filePath)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Duration is required",
		})
		return
	}

	if level == "" {
		os.Remove(filePath)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Level is required",
		})
		return
	}

	payload := models.Program{
		Icon:         icon,
		Title:        title,
		Duration:     duration,
		Participants: participants,
		Level:        level,
		Description:  description,
		Benefits:     benefitsStr,
		Image:        filePath,
	}

	program, err := ctrl.programService.Create(payload)
	if err != nil {
		os.Remove(filePath)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create program",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data":    program,
		"message": "Program created successfully",
	})
}

func (ctrl *ProgramController) FindAll(c *gin.Context) {
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

	data, total, err := ctrl.programService.FindAll(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch programs",
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

func (ctrl *ProgramController) FindByID(c *gin.Context) {
	id := c.Param("id")
	uint64Val, err := strconv.ParseUint(id, 10, 0)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid ID format",
			"error":   err.Error(),
		})
		return
	}

	data, err := ctrl.programService.FindByID(uint(uint64Val))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Program not found",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}

func (ctrl *ProgramController) Update(c *gin.Context) {
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

	// 2. Cek apakah program exist
	existingProgram, err := ctrl.programService.FindByID(uint(uint64Val))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Program not found",
			"error":   err.Error(),
		})
		return
	}

	// 3. Handle file upload (optional)
	filePath := existingProgram.Image
	oldFilePath := existingProgram.Image
	newFileUploaded := false

	file, err := c.FormFile("file")
	if err == nil { // Kalau ada file baru
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

		uploadPath := "uploads/programs/"
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

	// 4. Ambil field dari form
	icon := c.PostForm("icon")
	title := c.PostForm("title")
	duration := c.PostForm("duration")
	participants := c.PostForm("participants")
	level := c.PostForm("level")
	description := c.PostForm("description")

	// Parse benefits
	benefitsStr := c.PostFormArray("benefits")
	if len(benefitsStr) == 0 {
		benefitsInput := c.PostForm("benefits")
		if benefitsInput != "" {
			benefitsStr = strings.Split(benefitsInput, ",")
			for i := range benefitsStr {
				benefitsStr[i] = strings.TrimSpace(benefitsStr[i])
			}
		} else {
			benefitsStr = existingProgram.Benefits
		}
	}

	// Gunakan nilai lama jika tidak ada input baru
	if icon == "" {
		icon = existingProgram.Icon
	}
	if title == "" {
		title = existingProgram.Title
	}
	if duration == "" {
		duration = existingProgram.Duration
	}
	if participants == "" {
		participants = existingProgram.Participants
	}
	if level == "" {
		level = existingProgram.Level
	}
	if description == "" {
		description = existingProgram.Description
	}

	// 5. Buat payload untuk update
	payload := models.Program{
		ID:           uint(uint64Val),
		Icon:         icon,
		Title:        title,
		Duration:     duration,
		Participants: participants,
		Level:        level,
		Description:  description,
		Benefits:     benefitsStr,
		Image:        filePath,
	}

	// 6. Update ke database
	data, err := ctrl.programService.Update(payload)
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
			"message": "Failed to update program",
			"error":   err.Error(),
		})
		return
	}

	// 7. Hapus file lama jika ada file baru yang berhasil diupload
	if newFileUploaded && oldFilePath != "" && oldFilePath != filePath {
		if err := os.Remove(oldFilePath); err != nil {
			fmt.Printf("Warning: Failed to delete old file %s: %v\n", oldFilePath, err)
		} else {
			fmt.Printf("Successfully deleted old file: %s\n", oldFilePath)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    data,
		"message": "Program updated successfully",
	})
}

func (ctrl *ProgramController) Delete(c *gin.Context) {
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

	// 2. Cek apakah program exist
	existingProgram, err := ctrl.programService.FindByID(uint(uint64Val))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Program not found",
			"error":   err.Error(),
		})
		return
	}

	// 3. Soft delete program (set is_deleted = true)
	err = ctrl.programService.Delete(uint(uint64Val))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to delete program",
			"error":   err.Error(),
		})
		return
	}

	// 4. Hapus file image (optional - bisa juga tidak dihapus untuk backup)
	// CATATAN: Karena ini soft delete, mungkin lebih baik file tidak dihapus
	// Tapi jika ingin menghapus file, uncomment code di bawah:
	/*
		if existingProgram.Image != "" {
			if err := os.Remove(existingProgram.Image); err != nil {
				fmt.Printf("Warning: Failed to delete file %s: %v\n", existingProgram.Image, err)
			} else {
				fmt.Printf("Successfully deleted file: %s\n", existingProgram.Image)
			}
		}
	*/

	c.JSON(http.StatusOK, gin.H{
		"message": "Program deleted successfully",
		"data": gin.H{
			"id":    existingProgram.ID,
			"title": existingProgram.Title,
		},
	})
}
