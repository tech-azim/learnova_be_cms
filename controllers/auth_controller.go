package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tech-azim/be-learnova/services"
)


type AuthController struct {
	authService services.AuthService
}

type LoginRequest struct {
	Email string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func NewAuthController(authService services.AuthService) *AuthController {
	return &AuthController{
		authService: authService, 
	}
}

func (ctrl *AuthController) Login(c *gin.Context){
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H {
			"Mesage": err.Error(),
		})
		return
	}

	token, user, err := ctrl.authService.Login(req.Email,req.Password)
	if  err != nil {
		c.JSON(http.StatusBadRequest, gin.H {
			"Mesage": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H {
		"message":"Sukses",
		"token": token,
		 "user": user,
	})

}