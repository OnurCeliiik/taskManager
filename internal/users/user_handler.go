package users

import (
	"net/http"

	"task-manager/utils/error"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	service Service
}

func NewUserHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetUser(c *gin.Context) {
	IDStr := c.Param("id")
	ID, err := uuid.Parse(IDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, error.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid user ID",
		})
		return
	}

	user, err := h.service.GetUserByID(ID)
	if err != nil {
		c.JSON(http.StatusNotFound, error.ErrorResponse{
			Code:    http.StatusNotFound,
			Message: "User not found",
		})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *Handler) UpdateUser(c *gin.Context) {
	IDStr := c.Param("id")
	ID, err := uuid.Parse(IDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, error.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "invalid user ID",
		})
		return
	}

	userIDStr := c.GetString("userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, error.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "invalid user ID",
		})
		return
	}

	if ID != userID {
		c.JSON(http.StatusForbidden, error.ErrorResponse{
			Code:    http.StatusForbidden,
			Message: "you are not allowed to update this user",
		})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, error.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "invalid request payload",
		})
		return
	}

	if req.Name == nil && req.Email == nil && req.Password == nil {
		c.JSON(http.StatusBadRequest, error.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "nothing to update",
		})
		return
	}

	resp, err := h.service.UpdateUser(userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, error.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "failed to update user",
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) DeleteUser(c *gin.Context) {
	IDStr := c.Param("id")
	ID, err := uuid.Parse(IDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, error.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "invalid user ID",
		})
		return
	}

	userIDStr := c.GetString("userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, error.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "invalid user ID",
		})
		return
	}

	if ID != userID {
		c.JSON(http.StatusForbidden, error.ErrorResponse{
			Code:    http.StatusForbidden,
			Message: "you are not allowed to delete this user",
		})
		return
	}

	if err := h.service.DeleteUser(userID); err != nil {
		c.JSON(http.StatusInternalServerError, error.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "failed to delete user",
		})
		return
	}

	c.Status(http.StatusNoContent)
}
