package auth

import (
	"net/http"
	"task-manager/utils/email"
	"task-manager/utils/error"
	"task-manager/utils/jwt"
	"task-manager/utils/logger"

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

	emailService := email.NewEmailService()
	if err := emailService.SendRegistrationEmail(user.Email, user.Name); err != nil {
		logger.Warn("welcome email failed", "email", user.Email, "error", err.Error())
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

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req ForgotPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, error.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "invalid request payload",
		})
		return
	}

	resetToken, err := h.service.ForgotPassword(req.Email)
	if err != nil {
		logger.Warn("password reset requested", "email", req.Email, "error", err.Error())
		c.JSON(http.StatusOK, ForgotPasswordResponse{
			Message: "if email exists, reset link will be sent",
		})
		return
	}

	logger.Info("password reset token generated", "email", req.Email)

	// Send email with reset token
	emailService := email.NewEmailService()
	if err := emailService.SendPasswordResetEmail(req.Email, resetToken); err != nil {
		logger.Error("failed to send reset email", "email", req.Email, "error", err.Error())
		// Still return success to avoid leaking info
		c.JSON(http.StatusOK, ForgotPasswordResponse{
			Message: "if email exists, reset link will be sent",
		})
		return
	}

	c.JSON(http.StatusOK, ForgotPasswordResponse{
		Message: "if email exists, reset link will be sent",
	})
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, error.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "invalid request payload",
		})
		return
	}

	err := h.service.ResetPassword(req.Token, req.NewPassword)
	if err != nil {
		logger.Warn("password reset failed", "error", err.Error())
		c.JSON(http.StatusUnauthorized, error.ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: err.Error(),
		})
		return
	}

	logger.Info("password reset successful")
	c.JSON(http.StatusOK, ResetPasswordResponse{
		Message: "password reset successfully",
	})
}
