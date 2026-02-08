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

type RegistrationRequest struct {
	Name          string `json:"name" binding:"required"`
	Email         string `json:"email" binding:"required,email"`
	Phone         string `json:"phone" binding:"required"`
	Company       string `json:"company"`
	Position      string `json:"position"`
	ProgramID     uint   `json:"programId" binding:"required"`
	Participants  int    `json:"participants" binding:"required,min=1"`
	PreferredDate string `json:"preferredDate" binding:"required"`
	Message       string `json:"message"`
	Status        string `json:"status"`
}

type RegistrationController struct {
	registrationService services.RegistrationService
	programService      services.ProgramService
}

func NewRegistrationController(registrationService services.RegistrationService, programService services.ProgramService) *RegistrationController {
	return &RegistrationController{
		registrationService: registrationService,
		programService:      programService,
	}
}

func (ctrl *RegistrationController) Create(c *gin.Context) {
	var req RegistrationRequest

	// Bind dan validasi JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request data",
			"error":   err.Error(),
		})
		return
	}

	// Validasi apakah program exists
	_, err := ctrl.programService.FindByID(req.ProgramID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Program not found",
			"error":   err.Error(),
		})
		return
	}

	// Cek apakah email sudah terdaftar di program yang sama
	exists, err := ctrl.registrationService.CheckEmailExists(req.Email, req.ProgramID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to check email existence",
			"error":   err.Error(),
		})
		return
	}

	if exists {
		c.JSON(http.StatusConflict, gin.H{
			"message": "Email already registered for this program",
		})
		return
	}

	// Parse preferred date
	preferredDate, err := time.Parse("2006-01-02", req.PreferredDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid date format. Use YYYY-MM-DD",
			"error":   err.Error(),
		})
		return
	}

	// Validasi: preferred date tidak boleh di masa lalu
	if preferredDate.Before(time.Now().Truncate(24 * time.Hour)) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Preferred date cannot be in the past",
		})
		return
	}

	// Buat payload
	payload := models.Registration{
		Name:          req.Name,
		Email:         req.Email,
		Phone:         req.Phone,
		Company:       req.Company,
		Position:      req.Position,
		ProgramID:     req.ProgramID,
		Participants:  req.Participants,
		PreferredDate: preferredDate,
		Message:       req.Message,
	}

	// Simpan registrasi
	registration, err := ctrl.registrationService.Create(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create registration",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data":    registration,
		"message": "Registration created successfully",
	})
}

func (ctrl *RegistrationController) FindAll(c *gin.Context) {
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

	data, total, err := ctrl.registrationService.FindAll(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch registrations",
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

func (ctrl *RegistrationController) FindByID(c *gin.Context) {
	id := c.Param("id")
	uint64Val, err := strconv.ParseUint(id, 10, 0)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid ID format",
			"error":   err.Error(),
		})
		return
	}

	data, err := ctrl.registrationService.FindByID(uint(uint64Val))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Registration not found",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}

func (ctrl *RegistrationController) FindByProgramID(c *gin.Context) {
	programID := c.Param("programId")
	uint64Val, err := strconv.ParseUint(programID, 10, 0)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid program ID format",
			"error":   err.Error(),
		})
		return
	}

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

	// Validasi apakah program exists
	_, err = ctrl.programService.FindByID(uint(uint64Val))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Program not found",
			"error":   err.Error(),
		})
		return
	}

	data, total, err := ctrl.registrationService.FindByProgramID(uint(uint64Val), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch registrations",
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

func (ctrl *RegistrationController) FindByEmail(c *gin.Context) {
	email := c.Query("email")

	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Email parameter is required",
		})
		return
	}

	data, err := ctrl.registrationService.FindByEmail(email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Registration not found",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}

func (ctrl *RegistrationController) Update(c *gin.Context) {
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

	// 2. Cek apakah registration exist
	existingRegistration, err := ctrl.registrationService.FindByID(uint(uint64Val))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Registration not found",
			"error":   err.Error(),
		})
		return
	}

	// 3. Bind request JSON
	var req RegistrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request data",
			"error":   err.Error(),
		})
		return
	}

	// 4. Validasi apakah program exists (jika program_id berubah)
	if req.ProgramID != existingRegistration.ProgramID {
		_, err := ctrl.programService.FindByID(req.ProgramID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Program not found",
				"error":   err.Error(),
			})
			return
		}

		// Cek email conflict jika pindah program
		if req.Email != existingRegistration.Email || req.ProgramID != existingRegistration.ProgramID {
			exists, err := ctrl.registrationService.CheckEmailExists(req.Email, req.ProgramID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Failed to check email existence",
					"error":   err.Error(),
				})
				return
			}

			if exists {
				c.JSON(http.StatusConflict, gin.H{
					"message": "Email already registered for this program",
				})
				return
			}
		}
	}

	// 5. Parse preferred date
	preferredDate, err := time.Parse("2006-01-02", req.PreferredDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid date format. Use YYYY-MM-DD",
			"error":   err.Error(),
		})
		return
	}

	// 6. Buat payload untuk update
	payload := models.Registration{
		ID:            uint(uint64Val),
		Name:          req.Name,
		Email:         req.Email,
		Phone:         req.Phone,
		Company:       req.Company,
		Position:      req.Position,
		ProgramID:     req.ProgramID,
		Participants:  req.Participants,
		PreferredDate: preferredDate,
		Message:       req.Message,
		Status: 	   req.Status,
	}

	// 7. Update ke database
	data, err := ctrl.registrationService.Update(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to update registration",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    data,
		"message": "Registration updated successfully",
	})
}

func (ctrl *RegistrationController) Delete(c *gin.Context) {
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

	// 2. Cek apakah registration exist
	existingRegistration, err := ctrl.registrationService.FindByID(uint(uint64Val))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Registration not found",
			"error":   err.Error(),
		})
		return
	}

	// 3. Soft delete registration
	err = ctrl.registrationService.Delete(uint(uint64Val))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to delete registration",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Registration deleted successfully",
		"data": gin.H{
			"id":    existingRegistration.ID,
			"name":  existingRegistration.Name,
			"email": existingRegistration.Email,
		},
	})
}
