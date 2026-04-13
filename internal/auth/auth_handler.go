package auth

import (
	"net/http"
	"task-manager/utils/error"
	"task-manager/utils/jwt"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service AuthService
}

func NewAuthHandler(service AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) RegisterUser(c *gin.Context) {
	var req RegisterUserRequest

	// 1. Bind JSON request to struct
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, error.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request data",
		})
		return
	}

	// 2. Call service
	user, err := h.service.RegisterUser(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, error.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to register user",
		})
		return
	}

	// 3. Generate JWT token for the newly registered user
	token, err := jwt.GenerateToken(user.ID, user.Email, "user") // A user's role at first registration is always "user". Can be modified by an admin later.
	if err != nil {
		c.JSON(http.StatusInternalServerError, error.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to generate token",
		})
		return
	}

	resp := &RegisterResponse{
		Token: token,
		UserResponse: UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role,
		},
	}
	c.JSON(http.StatusCreated, resp)
}

func (h *AuthHandler) LoginUser(c *gin.Context) {
	var req LoginUserRequest

	// 1. Bind JSON request to struct
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, error.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request data",
		})
		return
	}

	// 2. Call service
	user, err := h.service.LoginUser(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, error.ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "Invalid email or password",
		})
		return
	}

	token, err := jwt.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, error.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to generate token",
		})
		return
	}

	resp := &LoginResponse{
		Token: token,
		UserResponse: UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role,
		},
	}

	c.JSON(http.StatusOK, resp)
}
