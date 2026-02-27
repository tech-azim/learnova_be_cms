package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tech-azim/be-learnova/services"
)

type UserController struct {
	service services.UserService
}

// NewUserController membuat instance baru UserController
func NewUserController(service services.UserService) *UserController {
	return &UserController{service}
}

// response helper
func (c *UserController) success(ctx *gin.Context, code int, message string, data any) {
	ctx.JSON(code, gin.H{
		"status":  "success",
		"message": message,
		"data":    data,
	})
}

func (c *UserController) fail(ctx *gin.Context, code int, message string) {
	ctx.JSON(code, gin.H{
		"status":  "error",
		"message": message,
	})
}

// GetAllUsers godoc
// @Summary      List all users
// @Tags         users
// @Produce      json
// @Success      200  {object}  map[string]any
// @Router       /users [get]
func (c *UserController) GetAllUsers(ctx *gin.Context) {
	users, err := c.service.GetAllUsers()
	if err != nil {
		c.fail(ctx, http.StatusInternalServerError, "Failed to fetch users")
		return
	}
	c.success(ctx, http.StatusOK, "Users fetched successfully", users)
}

// GetUserByID godoc
// @Summary      Get user by ID
// @Tags         users
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  map[string]any
// @Failure      400  {object}  map[string]any
// @Failure      404  {object}  map[string]any
// @Router       /users/{id} [get]
func (c *UserController) GetUserByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		c.fail(ctx, http.StatusBadRequest, "Invalid user ID")
		return
	}

	user, err := c.service.GetUserByID(uint(id))
	if err != nil {
		if err.Error() == "user not found" {
			c.fail(ctx, http.StatusNotFound, err.Error())
			return
		}
		c.fail(ctx, http.StatusInternalServerError, "Failed to fetch user")
		return
	}

	c.success(ctx, http.StatusOK, "User fetched successfully", user)
}

// CreateUser godoc
// @Summary      Create new user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        body  body      services.CreateUserInput  true  "User payload"
// @Success      201   {object}  map[string]any
// @Failure      400   {object}  map[string]any
// @Failure      422   {object}  map[string]any
// @Router       /users [post]
func (c *UserController) CreateUser(ctx *gin.Context) {
	var input services.CreateUserInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		c.fail(ctx, http.StatusBadRequest, err.Error())
		return
	}

	user, err := c.service.CreateUser(input)
	if err != nil {
		if err.Error() == "email already registered" {
			c.fail(ctx, http.StatusUnprocessableEntity, err.Error())
			return
		}
		c.fail(ctx, http.StatusInternalServerError, "Failed to create user")
		return
	}

	c.success(ctx, http.StatusCreated, "User created successfully", user)
}

// UpdateUser godoc
// @Summary      Update user by ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id    path      int                       true  "User ID"
// @Param        body  body      services.UpdateUserInput  true  "Update payload"
// @Success      200   {object}  map[string]any
// @Failure      400   {object}  map[string]any
// @Failure      404   {object}  map[string]any
// @Failure      422   {object}  map[string]any
// @Router       /users/{id} [put]
func (c *UserController) UpdateUser(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		c.fail(ctx, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var input services.UpdateUserInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		c.fail(ctx, http.StatusBadRequest, err.Error())
		return
	}

	user, err := c.service.UpdateUser(uint(id), input)
	if err != nil {
		switch err.Error() {
		case "user not found":
			c.fail(ctx, http.StatusNotFound, err.Error())
		case "email already used by another user":
			c.fail(ctx, http.StatusUnprocessableEntity, err.Error())
		default:
			c.fail(ctx, http.StatusInternalServerError, "Failed to update user")
		}
		return
	}

	c.success(ctx, http.StatusOK, "User updated successfully", user)
}

// DeleteUser godoc
// @Summary      Delete user by ID
// @Tags         users
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  map[string]any
// @Failure      400  {object}  map[string]any
// @Failure      404  {object}  map[string]any
// @Router       /users/{id} [delete]
func (c *UserController) DeleteUser(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		c.fail(ctx, http.StatusBadRequest, "Invalid user ID")
		return
	}

	if err := c.service.DeleteUser(uint(id)); err != nil {
		if err.Error() == "user not found" {
			c.fail(ctx, http.StatusNotFound, err.Error())
			return
		}
		c.fail(ctx, http.StatusInternalServerError, "Failed to delete user")
		return
	}

	c.success(ctx, http.StatusOK, "User deleted successfully", nil)
}

// ─── Tambahkan 2 handler baru di UserController ───────────────────────────────
// Letakkan setelah fungsi DeleteUser yang sudah ada

// GetProfile godoc
// @Summary      Get current logged-in user profile
// @Tags         profile
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]any
// @Failure      401  {object}  map[string]any
// @Router       /profile [get]
func (c *UserController) GetProfile(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")

	if !exists {
		c.fail(ctx, http.StatusUnauthorized, "Unauthorized user id not found")
		return
	}

	user, err := c.service.GetProfile(userID.(uint))
	if err != nil {
		if err.Error() == "user not found" {
			c.fail(ctx, http.StatusNotFound, err.Error())
			return
		}
		c.fail(ctx, http.StatusInternalServerError, "Failed to fetch profile")
		return
	}

	c.success(ctx, http.StatusOK, "Profile fetched successfully", user)
}

// UpdateProfile godoc
// @Summary      Update current logged-in user profile
// @Tags         profile
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      services.UpdateProfileInput  true  "Profile payload"
// @Success      200   {object}  map[string]any
// @Failure      400   {object}  map[string]any
// @Failure      401   {object}  map[string]any
// @Failure      422   {object}  map[string]any
// @Router       /profile [put]
func (c *UserController) UpdateProfile(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		c.fail(ctx, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var input services.UpdateProfileInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		c.fail(ctx, http.StatusBadRequest, err.Error())
		return
	}

	user, err := c.service.UpdateProfile(userID.(uint), input)
	if err != nil {
		switch err.Error() {
		case "user not found":
			c.fail(ctx, http.StatusNotFound, err.Error())
		case "email already used by another user":
			c.fail(ctx, http.StatusUnprocessableEntity, err.Error())
		default:
			c.fail(ctx, http.StatusInternalServerError, "Failed to update profile")
		}
		return
	}

	c.success(ctx, http.StatusOK, "Profile updated successfully", user)
}
