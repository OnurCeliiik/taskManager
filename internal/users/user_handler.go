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
